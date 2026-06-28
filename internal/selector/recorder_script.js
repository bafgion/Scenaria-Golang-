(() => {
  if (window.__scenariaRecorder) return;

  function cssEscape(value) {
    if (window.CSS && CSS.escape) return CSS.escape(value);
    return String(value).replace(/["\\]/g, '\\$&');
  }

  function visibleText(el) {
    if (!el || el.nodeType !== 1) return '';
    const clone = el.cloneNode(true);
    clone.querySelectorAll('input, textarea, select, script, style, noscript, svg').forEach((n) => n.remove());
    return (clone.innerText || clone.textContent || '').trim().replace(/\s+/g, ' ');
  }

  function labelTextForControl(el) {
    if (!el || el.nodeType !== 1) return '';
    const id = el.id;
    if (id) {
      const linked = el.ownerDocument.querySelector(`label[for="${cssEscape(id)}"]`);
      if (linked) return visibleText(linked);
    }
    const parentLabel = el.closest('label');
    if (parentLabel) return visibleText(parentLabel);
    return '';
  }

  function hasTextSelector(el, text) {
    const normalized = String(text || visibleText(el) || '').trim();
    if (!normalized || normalized.length < 2 || normalized.length > 80) return null;
    const escaped = normalized.replace(/"/g, '\\"');
    const tag = (el.tagName || '').toUpperCase();
    if (tag === 'BUTTON') return `button:has-text("${escaped}")`;
    if (tag === 'A') return `a:has-text("${escaped}")`;
    const role = el.getAttribute('role');
    if (role === 'button') return `button:has-text("${escaped}")`;
    if (role === 'link') return `a:has-text("${escaped}")`;
    return `button:has-text("${escaped}")`;
  }

  function clickableAncestor(el) {
    if (!el || el.nodeType !== 1) return null;
    const interactive = el.closest('button, a, [role="button"], [role="link"], [role="menuitem"], [role="tab"]');
    if (interactive) return interactive;
    let node = el;
    for (let depth = 0; node && depth < 8; depth++) {
      if (node.tagName === 'BUTTON' || node.tagName === 'A') return node;
      const role = node.getAttribute('role');
      if (role && ['button', 'link', 'menuitem', 'tab'].includes(role)) return node;
      node = node.parentElement;
    }
    return el;
  }

  function findCanvas(el) {
    if (!el || el.nodeType !== 1) return null;
    if (el.tagName === 'CANVAS') return el;
    return el.closest('canvas');
  }

  function buildCanvasSelector(canvas) {
    if (!canvas) return null;
    const testId = canvas.getAttribute('data-testid');
    if (testId) return `[data-testid="${cssEscape(testId)}"]`;
    if (canvas.id) return `#${cssEscape(canvas.id)}`;
    const aria = canvas.getAttribute('aria-label');
    if (aria && aria.trim()) return `canvas[aria-label="${cssEscape(aria.trim())}"]`;
    const wrap = canvas.closest('[data-testid]');
    if (wrap) {
      const wrapId = wrap.getAttribute('data-testid');
      if (wrapId) return `[data-testid="${cssEscape(wrapId)}"] canvas`;
    }
    const sig = canvas.closest('[class*="sign"], [class*="signature"], [data-signature]');
    if (sig && sig.id) return `#${cssEscape(sig.id)} canvas`;
    return 'canvas';
  }

  function isSignatureCanvas(canvas) {
    if (!canvas) return false;
    const cls = (canvas.className || '').toLowerCase();
    if (cls.includes('sign') || cls.includes('signature')) return true;
    const parent = canvas.parentElement;
    if (parent) {
      const pcls = (parent.className || '').toLowerCase();
      if (pcls.includes('sign') || pcls.includes('signature')) return true;
      if (parent.getAttribute('data-signature')) return true;
    }
    return canvas.getAttribute('role') === 'img' && !!canvas.getAttribute('aria-label');
  }

    if (!el || !['INPUT', 'TEXTAREA'].includes(el.tagName)) return null;
    const tag = el.tagName.toLowerCase();
    const testId = el.getAttribute('data-testid');
    if (testId) return `[data-testid="${cssEscape(testId)}"]`;
    if (el.id) return `#${cssEscape(el.id)}`;
    const label = labelTextForControl(el);
    if (label && label.length >= 2) return `label:has-text("${label.slice(0, 60).replace(/"/g, '\\"')}")`;
    const placeholder = el.getAttribute('placeholder');
    if (placeholder) return `${tag}[placeholder="${cssEscape(placeholder)}"]`;
    const aria = el.getAttribute('aria-label');
    if (aria) return `[aria-label="${cssEscape(aria.trim())}"]`;
    const name = el.getAttribute('name');
    if (name) return `${tag}[name="${cssEscape(name)}"]`;
    return null;
  }

  function countMatchingClickables(doc, label) {
    if (!doc || !label) return 0;
    let count = 0;
    const nodes = doc.querySelectorAll('button, a, [role="button"], [role="link"]');
    for (const node of nodes) {
      if (visibleText(node).trim() === label) count++;
    }
    return count;
  }

  function buildContextualClickSelector(target) {
    if (!target || target.nodeType !== 1) return null;
    const label = visibleText(target).trim();
    if (!label || label.length > 40) return null;
    if (countMatchingClickables(target.ownerDocument, label) <= 1) return null;
    let node = target.parentElement;
    for (let depth = 0; node && depth < 8; depth++) {
      const caption = visibleText(node).trim();
      if (caption.length >= 6 && caption !== label && caption.length <= 80) {
        const escapedCaption = caption.replace(/"/g, '\\"');
        const escapedLabel = label.replace(/"/g, '\\"');
        const tag = (target.tagName || 'BUTTON').toLowerCase();
        const btnTag = tag === 'a' ? 'a' : 'button';
        return `div:has-text("${escapedCaption}") >> ${btnTag}:has-text("${escapedLabel}")`;
      }
      node = node.parentElement;
    }
    return null;
  }

  function buildClickSelector(el) {
    if (!el || el.nodeType !== 1) return null;
    const target = clickableAncestor(el) || el;
    const testId = target.getAttribute('data-testid');
    if (testId) return `[data-testid="${cssEscape(testId)}"]`;
    if (target.id) return `#${cssEscape(target.id)}`;
    const aria = target.getAttribute('aria-label');
    if (aria && aria.trim()) return `[aria-label="${cssEscape(aria.trim())}"]`;
    const contextual = buildContextualClickSelector(target);
    if (contextual) return contextual;
    return hasTextSelector(target, '');
  }

  function buildSelector(el) {
    if (!el || el.nodeType !== 1) return '';
    if (el.tagName === 'CANVAS' || findCanvas(el) === el) {
      const canvasSel = buildCanvasSelector(el);
      if (canvasSel) return canvasSel;
    }
    const click = buildClickSelector(el);
    if (click) return click;
    const input = buildInputSelector(el);
    if (input) return input;
    if (el.id) return `#${cssEscape(el.id)}`;
    return el.tagName.toLowerCase();
  }

  function clickContextCaption(clickEl) {
    let node = clickEl && clickEl.parentElement;
    for (let depth = 0; node && depth < 8; depth++) {
      const text = visibleText(node);
      if (text.length >= 8 && text.length <= 120) return text.slice(0, 80);
      node = node.parentElement;
    }
    return '';
  }

  function collect(el, type) {
    const isField = el && ['INPUT', 'TEXTAREA', 'SELECT'].includes(el.tagName);
    const target = type === 'click' ? (clickableAncestor(el) || el) : el;
    if (!target) return {};
    return {
      tag: (target.tagName || '').toUpperCase(),
      id: target.id || '',
      name: target.getAttribute('name') || '',
      text: visibleText(target).slice(0, 120),
      testid: target.getAttribute('data-testid') || '',
      selector: buildSelector(el) || buildSelector(target),
      value: target.value || '',
      inputtype: (target.type || 'text').toLowerCase(),
      captiontext: isField ? labelTextForControl(target).slice(0, 120) : '',
      contexttext: type === 'click' ? clickContextCaption(target) : '',
      placeholder: target.getAttribute('placeholder') || '',
      arialabel: (target.getAttribute('aria-label') || '').trim(),
      role: target.getAttribute('role') || ''
    };
  }

  window.__scenariaRecorder = { events: [], paused: false };
  const push = (type, el) => {
    if (!el || window.__scenariaRecorder.paused) return;
    window.__scenariaRecorder.events.push({ type, detail: collect(el, type), ts: Date.now() });
  };

  let lastClickAt = 0;
  let lastClickKey = '';

  function resolveClickTarget(event) {
    const x = event.clientX;
    const y = event.clientY;
    let el = event.target;
    if (typeof document.elementFromPoint === 'function') {
      const top = document.elementFromPoint(x, y);
      if (top && top.nodeType === 1) {
        el = top;
      }
    }
    return el;
  }

  function shouldSkipDuplicateClick(target) {
    const key = (target && (target.id || target.getAttribute('data-testid') || visibleText(target).slice(0, 40))) || '';
    const now = Date.now();
    if (key && key === lastClickKey && now-lastClickAt < 400) {
      return true;
    }
    lastClickAt = now;
    lastClickKey = key;
    return false;
  }

  document.addEventListener('click', (e) => {
    const el = resolveClickTarget(e);
    if (!el || shouldSkipDuplicateClick(el)) return;
    const canvas = findCanvas(el);
    if (canvas && isSignatureCanvas(canvas)) {
      push('draw-signature', canvas);
      return;
    }
    push('click', el);
  }, true);
  document.addEventListener('input', (e) => push('input', e.target), true);
  document.addEventListener('change', (e) => push('change', e.target), true);
})();
