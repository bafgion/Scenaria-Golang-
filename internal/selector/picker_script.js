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

  function resolvePickTarget(rawEl) {
    if (!rawEl || rawEl.nodeType !== 1) return { el: null, selector: null };
    if (SKIP_IDS.has(rawEl.id)) return { el: null, selector: null };
    if (rawEl.tagName === 'IFRAME') {
      return { el: rawEl, selector: buildIframeSelector(rawEl) };
    }
    const canvas = H.findCanvas(rawEl);
    if (canvas) {
      return { el: canvas, selector: H.buildCanvasSelector(canvas) || H.buildSelector(canvas) };
    }
    const textInput =
      rawEl.tagName === 'INPUT' || rawEl.tagName === 'TEXTAREA'
        ? rawEl
        : rawEl.closest('label')?.querySelector('input:not([type="checkbox"]):not([type="radio"]), textarea');
    if (textInput) {
      return { el: textInput, selector: H.buildInputSelector(textInput) || H.buildSelector(textInput) };
    }
    return { el: rawEl, selector: H.buildSelector(rawEl) };
  }

  function elementUnderPointer(x, y) {
    const stack = document.elementsFromPoint(x, y);
    for (const el of stack) {
      if (!el || el.nodeType !== 1) continue;
      if (SKIP_IDS.has(el.id)) continue;
      return el;
    }
    return null;
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

  function finishPick(selector) {
    const done = window.pickSelectorDone;
    if (typeof done === 'function') {
      Promise.resolve(done(selector)).catch(() => {});
    }
    window.__shopPickerCleanup && window.__shopPickerCleanup();
  }

  function onMove(event) {
    const raw = elementUnderPointer(event.clientX, event.clientY);
    const target = resolvePickTarget(raw);
    if (!target.el) return;
    showOverlay(target.el);
  }

  function onPick(event) {
    event.preventDefault();
    event.stopPropagation();
    event.stopImmediatePropagation();
    const raw = elementUnderPointer(event.clientX, event.clientY);
    const target = resolvePickTarget(raw);
    if (!target.el || !target.selector) return;
    finishPick(target.selector);
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
