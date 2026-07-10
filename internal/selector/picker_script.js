(() => {
  if (typeof window.__shopPickerCleanup === 'function') {
    window.__shopPickerCleanup();
  }
  if (window.__shopPickerActive) return;
  if (window !== window.top) return;
  const H = window.__scenariaHeuristics;
  if (!H) return;

  window.__shopPickerActive = true;

  const SHIELD_ID = '__shopPickerShield';
  const OVERLAY_ID = '__shopPickerOverlay';
  const HINT_ID = '__shopPickerHint';
  const SKIP_IDS = new Set([SHIELD_ID, OVERLAY_ID, HINT_ID]);

  function isCrossOriginIframe(iframeEl) {
    try {
      return !iframeEl.contentDocument;
    } catch (_) {
      return true;
    }
  }

  function buildIframeSelector(el) {
    if (!el || el.tagName !== 'IFRAME') return null;
    const src = el.getAttribute('src') || '';
    if (src.includes('telegram.org')) {
      return 'iframe[src*="telegram.org"]';
    }
    if (el.id) return `#${H.cssEscape(el.id)}`;
    const title = el.getAttribute('title');
    if (title) return `iframe[title="${H.cssEscape(title)}"]`;
    const name = el.getAttribute('name');
    if (name) return `iframe[name="${H.cssEscape(name)}"]`;
    try {
      const url = new URL(src, window.location.href);
      if (url.host) {
        return `iframe[src*="${H.cssEscape(url.host)}"]`;
      }
    } catch (_) {
      /* ignore */
    }
    return H.buildSelector(el);
  }

  function pickHitAt(doc, x, y, iframeEl) {
    const base = doc.nodeType === 11 ? doc : doc;
    let elements = [];
    try {
      elements = base.elementsFromPoint(x, y);
    } catch (_) {
      elements = (doc.nodeType === 9 ? doc : document).elementsFromPoint(x, y);
    }
    for (const el of elements) {
      if (!el || el.nodeType !== 1) continue;
      if (SKIP_IDS.has(el.id)) continue;
      const tag = (el.tagName || '').toUpperCase();
      if (tag === 'HTML' || tag === 'BODY') continue;
      if (el.shadowRoot) {
        const inner = pickHitAt(el.shadowRoot, x, y, iframeEl);
        if (inner && inner.el) return inner;
      }
      if (el.tagName === 'IFRAME' && !iframeEl) {
        if (isCrossOriginIframe(el)) {
          return { el, iframe: null, crossOriginIframe: true };
        }
        try {
          const frameDoc = el.contentDocument || (el.contentWindow && el.contentWindow.document);
          if (frameDoc) {
            const rect = el.getBoundingClientRect();
            const inner = pickHitAt(frameDoc, x - rect.left, y - rect.top, el);
            if (inner && inner.el) return inner;
          }
        } catch (_) {
          return { el, iframe: null, crossOriginIframe: true };
        }
      }
      return { el, iframe: iframeEl, crossOriginIframe: false };
    }
    return null;
  }

  function elementUnderPointer(x, y) {
    return pickHitAt(document, x, y, null);
  }

  function normalizeSvgTarget(rawEl) {
    if (!rawEl || rawEl.nodeType !== 1) return rawEl;
    if (rawEl.namespaceURI === 'http://www.w3.org/2000/svg' || (rawEl.tagName && rawEl.tagName.toLowerCase() === 'svg')) {
      const clickable = H.clickableAncestor(rawEl);
      if (clickable && clickable !== rawEl) return clickable;
    }
    return rawEl;
  }

  function resolvePickTarget(hit) {
    if (!hit || !hit.el) return { el: null, kind: null, iframePrefix: null };
    if (hit.crossOriginIframe) {
      return { el: hit.el, kind: 'iframe', iframePrefix: null };
    }
    let rawEl = normalizeSvgTarget(hit.el);
    if (SKIP_IDS.has(rawEl.id)) return { el: null, kind: null, iframePrefix: null };
    const iframePrefix = hit.iframe ? buildIframeSelector(hit.iframe) : null;

    if (rawEl.tagName === 'IFRAME') {
      return { el: rawEl, kind: 'iframe', iframePrefix: null };
    }

    const canvas = H.findCanvas(rawEl);
    if (canvas) {
      return { el: canvas, kind: 'click', iframePrefix };
    }

    const inputControl = H.resolveInputFromPick(rawEl);
    if (inputControl) {
      return { el: inputControl, kind: 'input', iframePrefix };
    }

    return { el: rawEl, kind: 'click', iframePrefix };
  }

  function buildPickPayload(target) {
    if (!target.el) return null;
    if (target.kind === 'iframe') {
      const selector = buildIframeSelector(target.el);
      if (!selector) return null;
      const crossOrigin = isCrossOriginIframe(target.el);
      const warnings = crossOrigin ? ['cross-origin-iframe'] : ['iframe'];
      return {
        selector,
        candidates: [{ selector, strategy: 'iframe', score: 0, matches_count: -1, unique: true, visible: true, matches_picked: true, warnings }],
        warnings,
        suggested_action: 'click',
      };
    }
    const canvasSel = target.el.tagName === 'CANVAS' ? H.buildCanvasSelector(target.el) : null;
    if (canvasSel && target.el.tagName === 'CANVAS') {
      const warnings = canvasSel === 'canvas' ? ['generic-canvas'] : [];
      const payload = {
        selector: canvasSel,
        candidates: [{ selector: canvasSel, strategy: 'canvas', score: 0, matches_count: -1, unique: true, visible: H.isElementVisible(target.el), matches_picked: true, warnings }],
        warnings,
        suggested_action: 'click',
      };
      return target.iframePrefix ? H.prefixPickerResult(target.iframePrefix, payload) : payload;
    }
    const result = H.buildPickerResult(target.el, target.kind);
    if (!result) return null;
    if (target.iframePrefix) {
      const prefixed = H.prefixPickerResult(target.iframePrefix, result);
      const warnings = prefixed.warnings || (prefixed.warnings = []);
      if (!warnings.includes('iframe')) warnings.push('iframe');
      if (Array.isArray(prefixed.candidates)) {
        for (const cand of prefixed.candidates) {
          if (!cand.warnings) cand.warnings = [];
          if (!cand.warnings.includes('iframe')) cand.warnings.push('iframe');
        }
      }
      return prefixed;
    }
    return result;
  }

  function removeOverlay() {
    document.getElementById(OVERLAY_ID)?.remove();
  }

  function showOverlay(el) {
    removeOverlay();
    const rect = el.getBoundingClientRect();
    if (!rect.width && !rect.height) return;
    const box = document.createElement('div');
    box.id = OVERLAY_ID;
    box.style.cssText = [
      'position:fixed',
      'pointer-events:none',
      'z-index:2147483646',
      `left:${rect.left}px`,
      `top:${rect.top}px`,
      `width:${rect.width}px`,
      `height:${rect.height}px`,
      'border:2px solid #5ec8f2',
      'background:rgba(79,195,247,0.12)',
      'border-radius:3px',
    ].join(';');
    document.body.appendChild(box);
  }

  const hint = document.createElement('div');
  hint.id = HINT_ID;
  hint.textContent = 'Кликните по элементу · Esc — отмена';
  hint.style.cssText = [
    'position:fixed',
    'top:8px',
    'left:50%',
    'transform:translateX(-50%)',
    'z-index:2147483647',
    'background:#094771',
    'color:#fff',
    'padding:6px 12px',
    'border-radius:4px',
    'font:12px sans-serif',
    'pointer-events:none',
  ].join(';');
  document.body.appendChild(hint);

  function finishPick(payload) {
    const done = window.pickSelectorDone;
    if (typeof done === 'function') {
      Promise.resolve(done(payload)).catch(() => {});
    }
    window.__shopPickerCleanup && window.__shopPickerCleanup();
  }

  function onMove(event) {
    const hit = elementUnderPointer(event.clientX, event.clientY);
    const target = resolvePickTarget(hit);
    if (!target.el) return;
    showOverlay(target.el);
  }

  function onPick(event) {
    event.preventDefault();
    event.stopPropagation();
    event.stopImmediatePropagation();
    const hit = elementUnderPointer(event.clientX, event.clientY);
    const target = resolvePickTarget(hit);
    const payload = buildPickPayload(target);
    if (!payload || !payload.selector) return;
    finishPick(payload);
  }

  function blockPointerDown(event) {
    event.preventDefault();
    event.stopPropagation();
    event.stopImmediatePropagation();
  }

  function onKey(event) {
    if (event.key === 'Escape') {
      const cancel = window.pickSelectorCancel;
      if (typeof cancel === 'function') {
        Promise.resolve(cancel()).catch(() => {});
      }
      window.__shopPickerCleanup && window.__shopPickerCleanup();
    }
  }

  let shield = document.createElement('div');
  shield.id = SHIELD_ID;
  shield.style.cssText = [
    'position:fixed',
    'inset:0',
    'z-index:2147483645',
    'cursor:crosshair',
    'background:transparent',
  ].join(';');
  document.body.appendChild(shield);
  shield.addEventListener('pointerdown', blockPointerDown, true);
  shield.addEventListener('mousedown', blockPointerDown, true);
  shield.addEventListener('mousemove', onMove, true);
  shield.addEventListener('click', onPick, true);

  window.__shopPickerCleanup = () => {
    shield.removeEventListener('pointerdown', blockPointerDown, true);
    shield.removeEventListener('mousedown', blockPointerDown, true);
    shield.removeEventListener('mousemove', onMove, true);
    shield.removeEventListener('click', onPick, true);
    shield.remove();
    shield = null;
    document.removeEventListener('keydown', onKey, true);
    removeOverlay();
    hint.remove();
    window.__shopPickerActive = false;
    delete window.__shopPickerCleanup;
  };

  document.addEventListener('keydown', onKey, true);
})();
