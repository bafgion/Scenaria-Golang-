(() => {
  if (window.__scenariaRecorder) return;
  const H = window.__scenariaHeuristics;
  if (!H) return;
  const {
    visibleText,
    clickableAncestor,
    findCanvas,
    buildSelector,
    isSignatureCanvas,
    collect,
  } = H;

  window.__scenariaRecorder = { events: [], paused: false, filterImportant: false, navOnly: false, hoverRecord: false, scrollBeforeClick: false, hoverRecordMinMs: 600 };
  const cfg = () => window.__scenariaRecorder || {};

  function pushDetail(type, detail) {
    if (!detail || cfg().paused) return;
    if (cfg().navOnly && !['goto', 'scroll-to'].includes(type)) return;
    window.__scenariaRecorder.events.push({ type, detail, ts: Date.now() });
  }

  function isNavTarget(el) {
    if (!el || el.nodeType !== 1) return false;
    const tag = (el.tagName || '').toUpperCase();
    if (tag === 'A') return true;
    const role = el.getAttribute('role');
    return role === 'link';
  }

  function isImportantTarget(el) {
    if (!el || el.nodeType !== 1) return false;
    const tag = (el.tagName || '').toUpperCase();
    if (['BUTTON', 'A', 'INPUT', 'TEXTAREA', 'SELECT', 'LABEL'].includes(tag)) return true;
    if (findCanvas(el)) return true;
    const role = el.getAttribute('role');
    if (role && ['button', 'link', 'menuitem', 'tab', 'checkbox', 'radio', 'combobox'].includes(role)) return true;
    if (el.getAttribute('data-testid') || el.id) return true;
    return false;
  }

  const push = (type, el) => {
    if (!el || cfg().paused) return;
    const c = cfg();
    if (c.navOnly) {
      if (type !== 'click' || !isNavTarget(clickableAncestor(el) || el)) return;
    } else if (c.filterImportant && type === 'click' && !isImportantTarget(clickableAncestor(el) || el)) {
      return;
    }
    pushDetail(type, collect(el, type));
  };

  const RECORD_KEYS = new Set([
    'Enter', 'Tab', 'Escape', 'Backspace', 'Delete',
    'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight',
    'Home', 'End', 'PageUp', 'PageDown',
  ]);

  function pushKey(key, el) {
    if (!key || cfg().paused || cfg().navOnly) return;
    const field = el && el.closest
      ? el.closest('input, textarea, [contenteditable="true"]')
      : null;
    const isTextField = field && ['INPUT', 'TEXTAREA'].includes(field.tagName);
    if (isTextField) {
      const detail = collect(field, 'press-in');
      detail.value = key;
      pushDetail('press-in', detail);
      return;
    }
    pushDetail('press', { value: key, key });
  }

  let lastClickAt = 0;
  let lastClickKey = '';
  let lastClickSelector = { selector: '', at: 0 };
  let lastHoverTrigger = null;

  function isSubmenuContainer(node) {
    if (!node || node.nodeType !== 1) return false;
    if (node.matches('ul, ol, nav, [role="menu"], [class*="sub"], [class*="drop"], [class*="mega"], [class*="menu"]')) {
      return true;
    }
    return node.children.length > 2 && !!node.querySelector('a, button');
  }

  function findMenuHoverTrigger(el) {
    if (!el || el.nodeType !== 1) return null;

    let item = el.closest('li');
    while (item) {
      const trigger = item.querySelector(':scope > a, :scope > button, :scope > [role="button"]');
      if (trigger) {
        for (const child of Array.from(item.children)) {
          if (child === trigger) continue;
          if (isSubmenuContainer(child) && child.contains(el) && trigger !== el) {
            const selector = buildSelector(trigger);
            if (!selector) return null;
            return {
              selector,
              text: (trigger.innerText || trigger.textContent || '').trim().slice(0, 120),
            };
          }
        }
      }
      item = item.parentElement ? item.parentElement.closest('li') : null;
    }

    let node = el.parentElement;
    for (let depth = 0; node && depth < 8; depth++) {
      const directTriggers = node.querySelectorAll(':scope > a, :scope > button, :scope > [role="button"]');
      for (const trigger of directTriggers) {
        if (trigger === el || trigger.contains(el)) continue;
        for (const child of node.children) {
          if (child === trigger) continue;
          if (isSubmenuContainer(child) && child.contains(el)) {
            const selector = buildSelector(trigger);
            if (selector) {
              return {
                selector,
                text: (trigger.innerText || trigger.textContent || '').trim().slice(0, 120),
              };
            }
          }
        }
      }
      node = node.parentElement;
    }

    const header = el.closest('header, nav');
    if (header) {
      const triggers = header.querySelectorAll('a, button');
      for (const trigger of triggers) {
        let sibling = trigger.nextElementSibling;
        while (sibling) {
          if (isSubmenuContainer(sibling) && sibling.contains(el)) {
            const selector = buildSelector(trigger);
            if (selector) {
              return {
                selector,
                text: (trigger.innerText || trigger.textContent || '').trim().slice(0, 120),
              };
            }
          }
          sibling = sibling.nextElementSibling;
        }
      }
    }

    return null;
  }

  function sameNavContext(a, b) {
    const navA = a && a.closest ? a.closest('header, nav') : null;
    const navB = b && b.closest ? b.closest('header, nav') : null;
    return !!navA && navA === navB;
  }

  function rememberHoverTarget(el) {
    if (!el || el.nodeType !== 1) return;
    const trigger = el.closest('a,button,[role="button"]');
    if (!trigger) return;
    const selector = buildSelector(trigger);
    if (!selector) return;
    const text = (trigger.innerText || trigger.textContent || '').trim().slice(0, 120);
    lastHoverTrigger = { element: trigger, selector, text, at: Date.now() };
  }

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

  function scrollIntoViewIfNeeded(el) {
    if (!cfg().scrollBeforeClick || !el || el.nodeType !== 1) return;
    try {
      el.scrollIntoView({ block: 'center', inline: 'nearest' });
    } catch (_) {
      /* ignore */
    }
  }

  function onDocumentClick(e) {
    const el = resolveClickTarget(e);
    if (!el || shouldSkipDuplicateClick(el)) return;
    scrollIntoViewIfNeeded(clickableAncestor(el) || el);
    const tag = (el.tagName || '').toUpperCase();
    const inputType = (el.type || '').toLowerCase();
    if (tag === 'INPUT' && ['checkbox', 'radio', 'file'].includes(inputType)) return;
    const canvas = findCanvas(el);
    if (canvas && isSignatureCanvas(canvas)) {
      push('draw-signature', canvas);
      return;
    }

    const clickRoot = clickableAncestor(el) || el;
    const c = cfg();
    if (c.navOnly && !isNavTarget(clickRoot)) return;
    if (c.filterImportant && !isImportantTarget(clickRoot)) return;
    const detail = collect(el, 'click');
    if (!detail.selector) return;

    const now = Date.now();
    if (detail.selector === lastClickSelector.selector && now - lastClickSelector.at < 600) return;
    lastClickSelector = { selector: detail.selector, at: now };

    let hover = findMenuHoverTrigger(el);
    if (!hover && lastHoverTrigger && now - lastHoverTrigger.at < 12000) {
      if (
        lastHoverTrigger.selector !== detail.selector &&
        sameNavContext(el, lastHoverTrigger.element) &&
        !lastHoverTrigger.element.contains(el)
      ) {
        hover = { selector: lastHoverTrigger.selector, text: lastHoverTrigger.text };
      }
    }
    if (hover) {
      detail.hoverselector = hover.selector;
      detail.hovertext = hover.text || '';
    }
    pushDetail('click', detail);
  }

  function onDocumentInput(e) {
    if (cfg().navOnly) return;
    const el = e.target;
    if (!el || el.nodeType !== 1) return;
    if (el.tagName === 'INPUT' && (el.type || '').toLowerCase() === 'file') return;
    push('input', el);
  }

  function onDocumentChange(e) {
    if (cfg().navOnly) return;
    const el = e.target;
    if (!el || el.nodeType !== 1) return;
    if (el.tagName === 'INPUT' && (el.type || '').toLowerCase() === 'file') {
      const files = el.files;
      if (!files || !files.length) return;
      const detail = collect(el, 'upload');
      detail.value = files[0].name;
      pushDetail('upload', detail);
      return;
    }
    if (el.tagName === 'INPUT' && (el.type || '').toLowerCase() === 'checkbox') {
      push(el.checked ? 'check' : 'uncheck', el);
      return;
    }
    if (el.tagName === 'INPUT' && (el.type || '').toLowerCase() === 'radio') {
      if (el.checked) push('check', el);
      return;
    }
    push('change', el);
  }

  function onDocumentKeyDown(e) {
    if (cfg().paused || cfg().navOnly) return;
    if (e.repeat) return;
    const key = e.key;
    if (e.ctrlKey || e.metaKey || e.altKey) {
      if (key.length === 1) {
        const parts = [];
        if (e.ctrlKey || e.metaKey) parts.push('Control');
        if (e.altKey) parts.push('Alt');
        if (e.shiftKey) parts.push('Shift');
        parts.push(key.toUpperCase());
        pushKey(parts.join('+'), e.target);
      }
      return;
    }
    if (RECORD_KEYS.has(key)) {
      pushKey(key, e.target);
    }
  }

  let scrollTimer = null;
  let lastScrollKey = '';
  function elementAtViewportCenter(doc) {
    const cx = Math.floor((doc.documentElement.clientWidth || window.innerWidth) / 2);
    const cy = Math.floor((doc.documentElement.clientHeight || window.innerHeight) / 2);
    const el = doc.elementFromPoint(cx, cy);
    return el && el.nodeType === 1 ? el : null;
  }

  function onScroll(e) {
    if (cfg().paused || cfg().navOnly) return;
    if (scrollTimer) clearTimeout(scrollTimer);
    scrollTimer = setTimeout(() => {
      const root = e.target;
      const doc = root && root.ownerDocument ? root.ownerDocument : document;
      const el = elementAtViewportCenter(doc);
      if (!el) return;
      const detail = collect(el, 'scroll-to');
      if (!detail.selector) return;
      if (detail.selector === lastScrollKey) return;
      lastScrollKey = detail.selector;
      pushDetail('scroll-to', detail);
    }, 400);
  }

  let dragSource = null;
  function onDragStart(e) {
    if (cfg().paused) return;
    dragSource = e.target;
  }

  function onDrop(e) {
    if (cfg().paused || !dragSource) return;
    const target = e.target;
    if (!target || target.nodeType !== 1) return;
    const srcSel = buildSelector(dragSource);
    const dstSel = buildSelector(target);
    dragSource = null;
    if (!srcSel || !dstSel || srcSel === dstSel) return;
    pushDetail('drag-drop', { selector: srcSel, target: dstSel });
  }

  let lastHoverAt = 0;
  let lastHoverKey = '';
  let hoverTimer = null;
  let hoverTimerEl = null;
  function onDocumentMouseOver(e) {
    const el = e.target;
    if (!el || el.nodeType !== 1) return;
    rememberHoverTarget(el);
    if (!cfg().hoverRecord || cfg().paused) return;
    if (hoverTimerEl === el) return;
    if (hoverTimer) {
      clearTimeout(hoverTimer);
      hoverTimer = null;
    }
    hoverTimerEl = el;
    const minMs = cfg().hoverRecordMinMs || 600;
    hoverTimer = setTimeout(() => {
      hoverTimer = null;
      if (!cfg().hoverRecord || cfg().paused) return;
      const key = (el.id || el.getAttribute('data-testid') || visibleText(el).slice(0, 40)) || '';
      const now = Date.now();
      if (key && key === lastHoverKey && now - lastHoverAt < minMs) return;
      lastHoverAt = now;
      lastHoverKey = key;
      push('hover', el);
    }, minMs);
  }

  const attachedRoots = new WeakSet();

  function attachRoot(root) {
    if (!root || attachedRoots.has(root)) return;
    attachedRoots.add(root);
    root.addEventListener('click', onDocumentClick, true);
    root.addEventListener('input', onDocumentInput, true);
    root.addEventListener('change', onDocumentChange, true);
    root.addEventListener('mouseover', onDocumentMouseOver, true);
    root.addEventListener('keydown', onDocumentKeyDown, true);
    root.addEventListener('scroll', onScroll, true);
    root.addEventListener('dragstart', onDragStart, true);
    root.addEventListener('drop', onDrop, true);
    if (root === document) {
      window.addEventListener('scroll', onScroll, true);
    }
  }

  function scanNode(node) {
    if (!node || node.nodeType !== 1) return;
    if (node.shadowRoot) attachRoot(node.shadowRoot);
    if (node.tagName === 'IFRAME') {
      const hookIframe = () => {
        try {
          const doc = node.contentDocument;
          if (doc) {
            attachRoot(doc);
            observeRoot(doc);
          }
        } catch (_) {
          /* cross-origin iframe */
        }
      };
      hookIframe();
      node.addEventListener('load', hookIframe);
    }
    if (node.querySelectorAll) {
      node.querySelectorAll('*').forEach((child) => {
        if (child.shadowRoot) attachRoot(child.shadowRoot);
      });
    }
  }

  function observeRoot(root) {
    attachRoot(root);
    if (typeof MutationObserver === 'undefined') return;
    const target = root === document ? root.documentElement : root;
    const observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        mutation.addedNodes.forEach((node) => scanNode(node));
      }
    });
    observer.observe(target, { childList: true, subtree: true });
  }

  observeRoot(document);
})();
