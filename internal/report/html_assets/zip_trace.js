/* Playwright trace.zip reader for offline HTML report (deflate + stored entries). */
(function (root) {
  async function inflateRaw(data) {
    if (typeof DecompressionStream === 'undefined') {
      throw new Error('DecompressionStream unavailable');
    }
    const ds = new DecompressionStream('deflate-raw');
    const stream = new Blob([data]).stream().pipeThrough(ds);
    return new Uint8Array(await new Response(stream).arrayBuffer());
  }

  async function unzipEntries(arrayBuffer) {
    const view = new DataView(arrayBuffer);
    const dec = new TextDecoder();
    const out = {};
    let offset = 0;
    const total = view.byteLength;
    while (offset + 30 <= total) {
      if (view.getUint32(offset, true) !== 0x04034b50) break;
      const method = view.getUint16(offset + 8, true);
      const compSize = view.getUint32(offset + 18, true);
      const nameLen = view.getUint16(offset + 26, true);
      const extraLen = view.getUint16(offset + 28, true);
      const nameStart = offset + 30;
      if (nameStart + nameLen > total) break;
      const name = dec.decode(new Uint8Array(arrayBuffer, nameStart, nameLen));
      const dataStart = nameStart + nameLen + extraLen;
      if (dataStart + compSize > total) break;
      const comp = new Uint8Array(arrayBuffer, dataStart, compSize);
      let raw;
      if (method === 0) {
        raw = comp;
      } else if (method === 8) {
        raw = await inflateRaw(comp);
      } else {
        offset = dataStart + compSize;
        continue;
      }
      out[name.replace(/\\/g, '/')] = raw;
      offset = dataStart + compSize;
    }
    return out;
  }

  function traceActionTitle(row) {
    const method = String(row.method || '').trim();
    if (!method) return '';
    const cls = String(row.class || '').trim();
    let title = cls ? cls + '.' + method : method;
    const params = row.params;
    if (params && typeof params === 'object') {
      for (const key of ['url', 'selector', 'text', 'name']) {
        if (params[key]) {
          let s = String(params[key]).trim();
          if (s.length > 60) s = s.slice(0, 60) + '…';
          title += ' ' + s;
          break;
        }
      }
    }
    return title;
  }

  function parseTraceActions(bytes, limit) {
    const text = new TextDecoder().decode(bytes);
    const events = [];
    let base = null;
    for (const line of text.split('\n')) {
      if (!line.trim()) continue;
      let row;
      try { row = JSON.parse(line); } catch (_) { continue; }
      if (row.type !== 'action') continue;
      if (base == null) base = row.startTime || 0;
      const title = traceActionTitle(row);
      if (!title) continue;
      events.push({
        offset_ms: Math.round((row.startTime || 0) - base),
        title: title,
        kind: row.method || '',
      });
      if (events.length >= limit) return events;
    }
    return events;
  }

  function parseNetworkLine(line) {
    let raw;
    try { raw = JSON.parse(line); } catch (_) { return null; }
    const snap = raw.snapshot;
    if (!snap || !snap.request) return null;
    const req = snap.request;
    const url = String(req.url || '').trim();
    if (!url) return null;
    const method = String(req.method || 'GET').trim() || 'GET';
    const status = snap.response && snap.response.status ? Number(snap.response.status) : 0;
    const fail = String(snap._failureText || snap.failureText || '').trim();
    if (status > 0 && status < 400 && !fail) return null;
    const at = snap._monotonicTime || 0;
    let snippet;
    if (fail) snippet = method + ' ' + url + ' — ' + fail;
    else if (status >= 400) snippet = method + ' ' + url + ' — HTTP ' + status;
    else snippet = method + ' ' + url;
    return { at: at, snippet: snippet };
  }

  function parseNetworkFailures(bytes, limit) {
    const text = new TextDecoder().decode(bytes);
    const out = [];
    let base = null;
    for (const line of text.split('\n')) {
      if (!line.trim()) continue;
      const row = parseNetworkLine(line);
      if (!row) continue;
      if (base == null) base = row.at;
      out.push({ offset_ms: Math.round(row.at - base), snippet: row.snippet });
      if (out.length >= limit) break;
    }
    return out;
  }

  async function parsePlaywrightTraceZip(arrayBuffer) {
    const files = await unzipEntries(arrayBuffer);
    const traceName = Object.keys(files).find((k) => k.endsWith('.trace'));
    if (!traceName) throw new Error('no .trace in zip');
    const events = parseTraceActions(files[traceName], 120);
    const netName = Object.keys(files).find((k) => k.endsWith('.network'));
    const network = netName ? parseNetworkFailures(files[netName], 32) : [];
    return { events: events, network: network };
  }

  root.ScenariaTraceZip = { parsePlaywrightTraceZip: parsePlaywrightTraceZip };
})(typeof window !== 'undefined' ? window : globalThis);
