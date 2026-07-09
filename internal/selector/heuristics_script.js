(() => {
  if (window.__scenariaHeuristics) return;

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

  function findControlForLabel(labelEl) {
    if (!labelEl || labelEl.tagName !== 'LABEL') return null;
    const forId = labelEl.getAttribute('for');
    if (forId) {
      const doc = labelEl.ownerDocument;
      const byId = doc.getElementById(forId);
      if (byId && isInputLikeElement(byId)) return byId;
    }
    const nested = labelEl.querySelector('input:not([type="checkbox"]):not([type="radio"]), textarea, select');
    if (nested && isInputLikeElement(nested)) return nested;
    return null;
  }

  function isInputLikeElement(el) {
    if (!el || el.nodeType !== 1) return false;
    const tag = el.tagName;
    if (tag === 'INPUT') {
      const type = (el.type || 'text').toLowerCase();
      return type !== 'checkbox' && type !== 'radio' && type !== 'button' && type !== 'submit' && type !== 'reset';
    }
    if (tag === 'TEXTAREA' || tag === 'SELECT') return true;
    if (el.isContentEditable || el.getAttribute('contenteditable') === 'true') return true;
    const role = (el.getAttribute('role') || '').toLowerCase();
    return ['textbox', 'combobox', 'searchbox', 'spinbutton'].includes(role);
  }

  function inputElementTag(el) {
    if (!el) return 'input';
    if (el.tagName === 'TEXTAREA') return 'textarea';
    if (el.tagName === 'SELECT') return 'select';
    if (el.isContentEditable || el.getAttribute('contenteditable') === 'true') {
      return '[contenteditable="true"]';
    }
    const role = (el.getAttribute('role') || '').toLowerCase();
    if (['textbox', 'combobox', 'searchbox', 'spinbutton'].includes(role)) {
      const tag = (el.tagName || 'input').toLowerCase();
      if (tag !== 'input' && tag !== 'textarea' && tag !== 'select') {
        return `[role="${cssEscape(role)}"]`;
      }
    }
    return (el.tagName || 'input').toLowerCase();
  }

  function resolveInputFromPick(rawEl) {
    if (!rawEl || rawEl.nodeType !== 1) return null;
    if (isInputLikeElement(rawEl)) return rawEl;
    if (rawEl.tagName === 'LABEL') return findControlForLabel(rawEl);
    const nested = rawEl.closest('label')?.querySelector(
      'input:not([type="checkbox"]):not([type="radio"]), textarea, select, [contenteditable="true"], [role="textbox"], [role="combobox"], [role="searchbox"], [role="spinbutton"]',
    );
    if (nested && isInputLikeElement(nested)) return nested;
    return null;
  }

  function buildAdjacentLabelSelector(el) {
    if (!el || !isInputLikeElement(el)) return null;
    const id = el.id;
    if (id) {
      const label = el.ownerDocument.querySelector(`label[for="${cssEscape(id)}"]`);
      if (label && !el.closest('label')) {
        const labelText = visibleText(label).trim();
        if (labelText.length >= 2) {
          const escaped = labelText.slice(0, 60).replace(/"/g, '\\"');
          return `label:has-text("${escaped}") >> ${inputElementTag(el)}`;
        }
      }
    }
    let prev = el.previousElementSibling;
    while (prev) {
      if (prev.tagName === 'LABEL') {
        const labelText = visibleText(prev).trim();
        if (labelText.length >= 2) {
          const escaped = labelText.slice(0, 60).replace(/"/g, '\\"');
          return `label:has-text("${escaped}") >> ${inputElementTag(el)}`;
        }
      }
      prev = prev.previousElementSibling;
    }
    return null;
  }

  function hasTextSelector(el, text) {
    return clickHasTextSelector(el, text);
  }

  function clickTagFor(target) {
    if (!target || target.nodeType !== 1) return 'button';
    const tag = (target.tagName || '').toLowerCase();
    const role = (target.getAttribute('role') || '').toLowerCase();
    if (tag === 'a' || role === 'link') return 'a';
    if (tag === 'button' || ['button', 'menuitem', 'tab'].includes(role)) return 'button';
    if (tag === 'div' || tag === 'span' || tag === 'li') return tag;
    return tag || 'button';
  }

  function clickHasTextSelector(el, textOverride) {
    const normalized = String(textOverride || visibleText(el) || '').trim();
    if (!normalized || normalized.length < 2 || normalized.length > 80) return null;
    const escaped = normalized.replace(/"/g, '\\"');
    const target = clickableAncestor(el) || el;
    const tag = clickTagFor(target);
    return `${tag}:has-text("${escaped}")`;
  }

  function clickableAncestor(el) {
    if (!el || el.nodeType !== 1) return null;
    if (el.namespaceURI === 'http://www.w3.org/2000/svg' || (el.tagName && el.tagName.toLowerCase() === 'svg')) {
      const parentInteractive = el.closest('button, a, [role="button"], [role="link"], [role="menuitem"], [role="tab"]');
      if (parentInteractive) return parentInteractive;
    }
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

  function strategyOrder(kind) {
    const cfg = window.__scenariaSelectorOrder || {};
    const defaults = kind === 'input'
      ? ['testid', 'id', 'name', 'aria', 'label', 'adjacent', 'placeholder']
      : ['testid', 'aria', 'title', 'id', 'contextual', 'text'];
    const order = cfg[kind];
    return Array.isArray(order) && order.length ? order : defaults;
  }

  function clickStrategyBuilders(target) {
    return {
      testid() {
        const testId = target.getAttribute('data-testid');
        return testId ? `[data-testid="${cssEscape(testId)}"]` : null;
      },
      id() {
        return target.id ? `#${cssEscape(target.id)}` : null;
      },
      aria() {
        const aria = target.getAttribute('aria-label');
        return aria && aria.trim() ? `[aria-label="${cssEscape(aria.trim())}"]` : null;
      },
      title() {
        const title = target.getAttribute('title');
        return title && title.trim() ? `[title="${cssEscape(title.trim())}"]` : null;
      },
      contextual() {
        return buildContextualClickSelector(target);
      },
      text() {
        return clickHasTextSelector(target, '');
      },
    };
  }

  function inputStrategyBuilders(el, tag) {
    return {
      testid() {
        const testId = el.getAttribute('data-testid');
        return testId ? `[data-testid="${cssEscape(testId)}"]` : null;
      },
      id() {
        return el.id ? `#${cssEscape(el.id)}` : null;
      },
      label() {
        const label = labelTextForControl(el);
        if (!label || label.length < 2) return null;
        const escaped = label.slice(0, 60).replace(/"/g, '\\"');
        return `label:has-text("${escaped}") >> ${tag}`;
      },
      adjacent() {
        return buildAdjacentLabelSelector(el);
      },
      placeholder() {
        const placeholder = el.getAttribute('placeholder');
        return placeholder ? `${tag}[placeholder="${cssEscape(placeholder)}"]` : null;
      },
      aria() {
        const aria = el.getAttribute('aria-label');
        return aria ? `[aria-label="${cssEscape(aria.trim())}"]` : null;
      },
      name() {
        const name = el.getAttribute('name');
        return name ? `${tag}[name="${cssEscape(name)}"]` : null;
      },
    };
  }

  function buildInputSelector(el) {
    if (!isInputLikeElement(el)) return null;
    const tag = inputElementTag(el);
    const builders = inputStrategyBuilders(el, tag);
    for (const key of strategyOrder('input')) {
      const sel = builders[key] && builders[key]();
      if (sel) return sel;
    }
    return null;
  }

  function navScopeTag(el) {
    if (!el || el.nodeType !== 1) return null;
    const nav = el.closest('nav, header, [role="navigation"], [role="menubar"]');
    if (!nav) return null;
    const tag = (nav.tagName || 'nav').toUpperCase();
    return tag === 'HEADER' ? 'header' : 'nav';
  }

  function scopedTextSelector(scopeTag, text, el) {
    const normalized = String(text || '').trim();
    if (!normalized || normalized.length < 2 || normalized.length > 80) return null;
    const escaped = normalized.replace(/"/g, '\\"');
    const innerTag = el ? clickTagFor(clickableAncestor(el) || el) : 'button';
    if (scopeTag) {
      return `${scopeTag} >> ${innerTag}:has-text("${escaped}")`;
    }
    return `${innerTag}:has-text("${escaped}")`;
  }

  function buildMenuTriggerSelector(trigger) {
    if (!trigger || trigger.nodeType !== 1) return null;
    const text = visibleText(trigger).trim();
    if (!text) return buildClickSelector(trigger);
    if (countMatchingClickables(trigger.ownerDocument, text) <= 1) {
      return clickHasTextSelector(trigger, text) || buildClickSelector(trigger);
    }
    const scopeTag = navScopeTag(trigger);
    if (scopeTag) {
      return scopedTextSelector(scopeTag, text, trigger);
    }
    return buildClickSelector(trigger);
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
    const scopeTag = navScopeTag(target);
    if (scopeTag) {
      return scopedTextSelector(scopeTag, label, target);
    }
    let node = target.parentElement;
    for (let depth = 0; node && depth < 8; depth++) {
      const caption = visibleText(node).trim();
      if (caption.length >= 6 && caption.length <= 40 && caption !== label) {
        const escapedCaption = caption.replace(/"/g, '\\"');
        const escapedLabel = label.replace(/"/g, '\\"');
        const parentTag = (node.tagName || 'div').toLowerCase();
        const clickTag = clickTagFor(target);
        return `${parentTag}:has-text("${escapedCaption}") >> ${clickTag}:has-text("${escapedLabel}")`;
      }
      node = node.parentElement;
    }
    return null;
  }

  function buildClickSelector(el) {
    if (!el || el.nodeType !== 1) return null;
    const target = clickableAncestor(el) || el;
    const builders = clickStrategyBuilders(target);
    for (const key of strategyOrder('click')) {
      const sel = builders[key] && builders[key]();
      if (sel) return sel;
    }
    return null;
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
    return '';
  }

  function cssPathSegment(el) {
    if (!el || el.nodeType !== 1) return '';
    const tag = (el.tagName || '').toLowerCase();
    if (!tag || tag === 'html') return tag;
    if (el.id) return `#${cssEscape(el.id)}`;
    const testId = el.getAttribute('data-testid');
    if (testId) return `[data-testid="${cssEscape(testId)}"]`;
    const parent = el.parentElement;
    if (!parent) return tag;
    let sameTagIndex = 0;
    let sameTagCount = 0;
    for (const child of Array.from(parent.children)) {
      if ((child.tagName || '').toLowerCase() !== tag) continue;
      sameTagCount++;
      if (child === el) sameTagIndex = sameTagCount;
    }
    return sameTagCount > 1 ? `${tag}:nth-of-type(${sameTagIndex})` : tag;
  }

  function buildCssPathSelector(el) {
    if (!el || el.nodeType !== 1) return '';
    const root = el.getRootNode();
    const parts = [];
    let node = el;
    for (let depth = 0; node && node.nodeType === 1 && depth < 8; depth++) {
      const segment = cssPathSegment(node);
      if (!segment || segment === 'html') break;
      parts.unshift(segment);
      if (segment.startsWith('#') || segment.startsWith('[data-testid=')) break;
      if ((node.tagName || '').toLowerCase() === 'body') break;
      node = node.parentElement;
    }
    const selector = parts.join(' > ');
    if (!selector) return '';
    const scope = (typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot) ? root : (el.ownerDocument || document);
    const matches = querySelectorAllApprox(scope, selector);
    if (matches && matches.length === 1 && matches[0] === el) return selector;
    return '';
  }

  function buildRecorderSelector(el, kind) {
    if (!el || el.nodeType !== 1) return '';
    const resolvedKind = kind || pickKindForElement(el);
    const candidates = generateCandidates(el, resolvedKind);
    const exact = candidates.find((cand) => cand.valid && cand.selector);
    if (exact) return prefixShadowChain(el, exact.selector);
    const usable = candidates.find((cand) =>
      cand.selector &&
      cand.matches_count === 1 &&
      cand.matches_picked &&
      cand.visible &&
      cand.actionable &&
      !(cand.warnings || []).includes('wrong-target')
    );
    if (usable) return prefixShadowChain(el, usable.selector);
    const cssPath = buildCssPathSelector(resolvedKind === 'input' ? el : (clickableAncestor(el) || el));
    if (cssPath) return prefixShadowChain(el, cssPath);
    const fallback = resolvedKind === 'input' ? buildInputSelector(el) : buildSelector(el);
    return prefixShadowChain(el, fallback);
  }

  function isElementVisible(el) {
    if (!el || el.nodeType !== 1) return false;
    const style = window.getComputedStyle(el);
    if (style.display === 'none' || style.visibility === 'hidden' || Number(style.opacity) === 0) {
      return false;
    }
    const rect = el.getBoundingClientRect();
    return rect.width > 0 || rect.height > 0;
  }

  function isActionable(el, action) {
    if (!el || el.nodeType !== 1) return false;
    if (!isElementVisible(el)) return false;
    const style = window.getComputedStyle(el);
    if (style.pointerEvents === 'none') return false;
    if (el.disabled) return false;
    if (el.getAttribute('aria-disabled') === 'true') return false;
    const tag = (el.tagName || '').toUpperCase();
    if (action === 'fill' || action === 'select') {
      if (el.readOnly) return false;
      if (el.getAttribute('aria-readonly') === 'true') return false;
      if (action === 'select' && tag !== 'SELECT' && (el.getAttribute('role') || '').toLowerCase() !== 'combobox') {
        return false;
      }
    }
    if (action === 'check' || action === 'uncheck') {
      return tag === 'INPUT' && (el.type || '').toLowerCase() === 'checkbox';
    }
    if (action === 'hover' || action === 'click') return true;
    if (action === 'fill') return isInputLikeElement(el);
    return true;
  }

  function parseHasTextSelector(selector) {
    const match = String(selector || '').match(/:has-text\("((?:[^"\\]|\\.)*)"\)/);
    if (!match) return null;
    const text = match[1].replace(/\\"/g, '"').replace(/\\\\/g, '\\');
    const css = selector.replace(match[0], '') || '*';
    return { css, text };
  }

  function querySelectorAllApprox(root, selector) {
    if (!root || !selector) return null;
    const trimmed = String(selector).trim();
    if (!trimmed) return [];
    const chainIdx = trimmed.indexOf('>>');
    if (chainIdx >= 0) {
      const left = trimmed.slice(0, chainIdx).trim();
      const right = trimmed.slice(chainIdx + 2).trim();
      const containers = querySelectorAllApprox(root, left);
      if (!containers) return null;
      const out = [];
      for (const container of containers) {
        const inner = querySelectorAllApprox(container, right);
        if (!inner) return null;
        out.push(...inner);
      }
      return out;
    }
    const hasText = parseHasTextSelector(trimmed);
    try {
      if (hasText) {
        return Array.from(root.querySelectorAll(hasText.css)).filter((el) => visibleText(el).includes(hasText.text));
      }
      return Array.from(root.querySelectorAll(trimmed));
    } catch (_) {
      return null;
    }
  }

  function countMatchesInDoc(doc, selector) {
    const matches = querySelectorAllApprox(doc, selector);
    return matches ? matches.length : -1;
  }

  function matchesPickedElement(selector, pickedEl, actionTarget, matchesCount) {
    if (!selector || !pickedEl) return false;
    if (matchesCount < 0) return true;
    try {
      const doc = pickedEl.ownerDocument || document;
      const nodes = querySelectorAllApprox(doc, selector);
      if (!nodes) return true;
      if (nodes.length !== 1) return false;
      const matched = nodes[0];
      if (matched === pickedEl || matched === actionTarget) return true;
      if (matched.contains(pickedEl) || pickedEl.contains(matched)) return true;
      return false;
    } catch (_) {
      return false;
    }
  }

  function strategyScore(kind, key) {
    const idx = strategyOrder(kind).indexOf(key);
    return idx >= 0 ? idx : strategyOrder(kind).length;
  }

  function candidateWarnings(kind, key, selector, matchesCount, visible, matchesPicked, actionable) {
    const warnings = [];
    if (key === 'text' || (kind === 'input' && selector.startsWith('label:has-text(') && !selector.includes('>>'))) {
      warnings.push('text-only');
    }
    if (matchesCount > 1) warnings.push('not-unique');
    if (matchesCount === 0) warnings.push('no-matches');
    if (!visible) warnings.push('not-visible');
    if (!matchesPicked && matchesCount === 1) warnings.push('wrong-target');
    if (!actionable) warnings.push('not-actionable');
    return warnings;
  }

  function matchesSelectorOnElement(el, selector) {
    if (!el || !selector) return false;
    if (selector.includes('>>') || selector.includes(':has-text(')) return false;
    try {
      const root = el.getRootNode();
      const scope = (typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot) ? root : (el.ownerDocument || document);
      const found = scope.querySelector(selector);
      return found === el;
    } catch (_) {
      return false;
    }
  }

  const LIBRARY_SCORE_BONUS = -100;
  const TEXT_PENALTY_WHEN_LIBRARY = 50;

  function libraryPacksEnabled() {
    const cfg = window.__scenariaLibraryHeuristics || {};
    return { mui: cfg.mui !== false, ant: cfg.ant !== false };
  }

  function closestLibraryMatch(el, selectors) {
    if (!el || !el.closest) return null;
    try {
      return el.closest(selectors);
    } catch (_) {
      return null;
    }
  }

  function isPortalContainer(el) {
    return !!closestLibraryMatch(el, [
      '[role="presentation"]',
      '.MuiPopover-root', '.MuiModal-root', '.MuiMenu-root', '.MuiPopper-root',
      '.ant-dropdown', '.ant-select-dropdown', '.ant-picker-dropdown', '.ant-menu-submenu',
    ]);
  }

  function muiRoot(el) {
    return closestLibraryMatch(el, '.MuiButton-root, .MuiIconButton-root, .MuiInputBase-root, .MuiMenuItem-root, .MuiCheckbox-root, .MuiSwitch-root, [class*="MuiButton-root"]');
  }

  function antRoot(el) {
    return closestLibraryMatch(el, '.ant-btn, .ant-input, .ant-select, .ant-select-selector, .ant-menu-item, [class^="ant-"], [class*=" ant-"]');
  }

  function buildMuiCandidate(el, kind, actionTarget) {
    const root = muiRoot(el) || muiRoot(actionTarget);
    if (!root) return null;
    const testId = root.getAttribute('data-testid');
    if (testId) return { selector: `[data-testid="${cssEscape(testId)}"]`, strategy: 'mui-testid' };
    if (kind === 'input') {
      const input = el.matches('input, textarea') ? el : root.querySelector('input, textarea, .MuiInputBase-input');
      if (!input) return null;
      const name = input.getAttribute('name');
      if (name) return { selector: `.MuiInputBase-input[name="${cssEscape(name)}"]`, strategy: 'mui-input' };
      return { selector: '.MuiInputBase-input', strategy: 'mui-input' };
    }
    if (root.classList.contains('MuiMenuItem-root')) {
      const text = visibleText(root).trim();
      if (text.length >= 2 && text.length <= 60) {
        const escaped = text.replace(/"/g, '\\"');
        return { selector: `.MuiMenuItem-root:has-text("${escaped}")`, strategy: 'mui-menu-item' };
      }
      return { selector: '.MuiMenuItem-root', strategy: 'mui-menu-item' };
    }
    if (root.classList.contains('MuiButton-root') || (root.className && String(root.className).includes('MuiButton-root'))) {
      const text = visibleText(root).trim();
      if (text.length >= 2 && text.length <= 60) {
        const escaped = text.replace(/"/g, '\\"');
        return { selector: `.MuiButton-root:has-text("${escaped}")`, strategy: 'mui-button' };
      }
      return { selector: '.MuiButton-root', strategy: 'mui-button' };
    }
    if (root.classList.contains('MuiIconButton-root')) {
      return { selector: '.MuiIconButton-root', strategy: 'mui-icon-button' };
    }
    return null;
  }

  function buildAntCandidate(el, kind, actionTarget) {
    const root = antRoot(el) || antRoot(actionTarget);
    if (!root) return null;
    const testId = root.getAttribute('data-testid');
    if (testId) return { selector: `[data-testid="${cssEscape(testId)}"]`, strategy: 'ant-testid' };
    if (kind === 'input') {
      const input = el.matches('input, textarea') ? el : root.querySelector('input, textarea, .ant-input');
      if (!input) return null;
      const placeholder = input.getAttribute('placeholder');
      if (placeholder) return { selector: `.ant-input[placeholder="${cssEscape(placeholder)}"]`, strategy: 'ant-input' };
      return { selector: '.ant-input', strategy: 'ant-input' };
    }
    const btn = root.closest('.ant-btn') || (root.classList.contains('ant-btn') ? root : null);
    if (btn) {
      const text = visibleText(btn).trim();
      if (text.length >= 2 && text.length <= 60) {
        const escaped = text.replace(/"/g, '\\"');
        return { selector: `.ant-btn:has-text("${escaped}")`, strategy: 'ant-button' };
      }
      return { selector: '.ant-btn', strategy: 'ant-button' };
    }
    const menuItem = root.closest('.ant-dropdown-menu-item, .ant-select-item, .ant-menu-item')
      || (root.classList.contains('ant-menu-item') ? root : null);
    if (menuItem) {
      const text = visibleText(menuItem).trim();
      if (text.length >= 2 && text.length <= 60) {
        const escaped = text.replace(/"/g, '\\"');
        return { selector: `.ant-menu-item:has-text("${escaped}")`, strategy: 'ant-menu-item' };
      }
    }
    const select = root.closest('.ant-select');
    if (select) return { selector: '.ant-select', strategy: 'ant-select' };
    return null;
  }

  function buildCandidateEntry(kind, key, sel, el, actionTarget, suggested, inShadow, score, extraWarnings) {
    const doc = el.ownerDocument || document;
    let matchesCount = countMatchesInDoc(doc, sel);
    const visible = isElementVisible(actionTarget);
    let matchesPicked = matchesPickedElement(sel, el, actionTarget, matchesCount);
    if (inShadow && matchesSelectorOnElement(actionTarget, sel)) {
      matchesCount = 1;
      matchesPicked = true;
    }
    const actionable = isActionable(actionTarget, suggested);
    const warnings = candidateWarnings(kind, key, sel, matchesCount, visible, matchesPicked, actionable);
    if (extraWarnings && extraWarnings.length) {
      for (const w of extraWarnings) {
        if (!warnings.includes(w)) warnings.push(w);
      }
    }
    const unique = matchesCount === 1 || matchesCount < 0;
    const valid = unique && matchesPicked && visible && actionable && matchesCount !== 0;
    return {
      selector: sel,
      strategy: key,
      score,
      matches_count: matchesCount,
      unique,
      visible,
      matches_picked: matchesPicked,
      warnings,
      valid,
      actionable,
    };
  }

  function generateLibraryCandidates(el, kind) {
    const packs = libraryPacksEnabled();
    if (!packs.mui && !packs.ant) return [];
    const actionTarget = kind === 'input' ? el : (clickableAncestor(el) || el);
    const suggested = suggestAction(el);
    const inShadow = (() => {
      const root = el.getRootNode();
      return typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot;
    })();
    const portalWarnings = isPortalContainer(el) ? ['portal-menu'] : [];
    const out = [];
    const seen = new Set();
    const add = (built, strategyKey) => {
      if (!built || !built.selector || seen.has(built.selector)) return;
      seen.add(built.selector);
      out.push(buildCandidateEntry(
        kind, strategyKey, built.selector, el, actionTarget, suggested, inShadow, LIBRARY_SCORE_BONUS, portalWarnings,
      ));
    };
    if (packs.mui) add(buildMuiCandidate(el, kind, actionTarget), 'mui');
    if (packs.ant) add(buildAntCandidate(el, kind, actionTarget), 'ant');
    return out;
  }

  function shouldPenalizeGenericText(kind, key, selector) {
    if (key === 'text') return true;
    if (kind === 'input' && key === 'label') return true;
    if (kind === 'input' && key === 'placeholder' && selector && selector.startsWith('label:has-text(')) return true;
    return false;
  }

  function generateCandidates(el, kind) {
    if (!el || el.nodeType !== 1) return [];
    const actionTarget = kind === 'input' ? el : (clickableAncestor(el) || el);
    const tag = inputElementTag(el);
    const builders = kind === 'input' ? inputStrategyBuilders(el, tag) : clickStrategyBuilders(actionTarget);
    const suggested = suggestAction(el);
    const inShadow = (() => {
      const root = el.getRootNode();
      return typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot;
    })();
    const libraryCandidates = generateLibraryCandidates(el, kind);
    const hasValidLibrary = libraryCandidates.some((c) => c.valid);
    const seen = new Set(libraryCandidates.map((c) => c.selector));
    const candidates = [...libraryCandidates];
    for (const key of strategyOrder(kind)) {
      if (!builders[key]) continue;
      const sel = builders[key]();
      if (!sel || seen.has(sel)) continue;
      seen.add(sel);
      let score = strategyScore(kind, key);
      if (hasValidLibrary && shouldPenalizeGenericText(kind, key, sel)) {
        score += TEXT_PENALTY_WHEN_LIBRARY;
      }
      candidates.push(buildCandidateEntry(
        kind, key, sel, el, actionTarget, suggested, inShadow, score, [],
      ));
    }
    candidates.sort((a, b) => {
      if (a.valid !== b.valid) return a.valid ? -1 : 1;
      if (a.score !== b.score) return a.score - b.score;
      return (a.matches_count < 0 ? 99 : a.matches_count) - (b.matches_count < 0 ? 99 : b.matches_count);
    });
    return candidates;
  }

  function suggestAction(el) {
    if (!el || el.nodeType !== 1) return 'click';
    const tag = (el.tagName || '').toUpperCase();
    if (tag === 'INPUT') {
      const type = (el.type || 'text').toLowerCase();
      if (type === 'checkbox') return el.checked ? 'uncheck' : 'check';
      if (type === 'radio') return 'click';
      return 'fill';
    }
    if (tag === 'TEXTAREA') return 'fill';
    if (tag === 'SELECT') return 'select';
    if (el.isContentEditable || el.getAttribute('contenteditable') === 'true') return 'fill';
    const role = (el.getAttribute('role') || '').toLowerCase();
    if (role === 'combobox') return 'select';
    if (['textbox', 'searchbox', 'spinbutton'].includes(role)) return 'fill';
    return 'click';
  }

  function pickKindForElement(el) {
    return isInputLikeElement(el) ? 'input' : 'click';
  }

  function prefixSelectorChain(prefix, selector) {
    if (!prefix || !selector) return selector || prefix || '';
    return `${prefix} >> ${selector}`;
  }

  function prefixPickerResult(prefix, result) {
    if (!prefix || !result) return result;
    result.selector = prefixSelectorChain(prefix, result.selector);
    if (Array.isArray(result.candidates)) {
      result.candidates = result.candidates.map((cand) => ({
        ...cand,
        selector: prefixSelectorChain(prefix, cand.selector),
      }));
    }
    return result;
  }

  function shadowHostPrefix(el) {
    const root = el.getRootNode();
    if (!root || typeof ShadowRoot === 'undefined' || !(root instanceof ShadowRoot)) return '';
    const host = root.host;
    if (!host || host.nodeType !== 1) return '';
    if (host.id) return `#${cssEscape(host.id)}`;
    const testId = host.getAttribute('data-testid');
    if (testId) return `[data-testid="${cssEscape(testId)}"]`;
    return '';
  }

  function prefixShadowChain(el, selector) {
    const hostPrefix = shadowHostPrefix(el);
    if (!hostPrefix || !selector) return selector;
    if (selector.startsWith(hostPrefix + ' >> ')) return selector;
    return `${hostPrefix} >> ${selector}`;
  }

  function prefixPickerResultChains(el, result) {
    if (!result) return result;
    result.selector = prefixShadowChain(el, result.selector);
    if (Array.isArray(result.candidates)) {
      result.candidates = result.candidates.map((cand) => ({
        ...cand,
        selector: prefixShadowChain(el, cand.selector),
      }));
    }
    return result;
  }

  function buildPickerResult(el, kind) {
    if (!el || el.nodeType !== 1) return null;
    const resolvedKind = kind || pickKindForElement(el);
    const candidates = generateCandidates(el, resolvedKind);
    if (!candidates.length) {
      const fallback = resolvedKind === 'input' ? buildInputSelector(el) : buildSelector(el);
      if (!fallback) return null;
      const result = {
        selector: prefixShadowChain(el, fallback),
        candidates: [],
        warnings: ['fallback'],
        suggested_action: suggestAction(el),
      };
      return result;
    }
    const best = candidates.find((c) => c.valid) || candidates[0];
    const warnings = [...(best.warnings || [])];
    if (!best.valid) warnings.push('low-confidence');
    const result = {
      selector: prefixShadowChain(el, best.selector),
      candidates: candidates.slice(0, 8).map((cand) => ({
        ...cand,
        selector: prefixShadowChain(el, cand.selector),
      })),
      warnings,
      suggested_action: suggestAction(el),
    };
    return result;
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
    const isField = isInputLikeElement(el);
    const target = type === 'click' ? (clickableAncestor(el) || el) : (resolveInputFromPick(el) || el);
    if (!target) return {};
    const selectorKind = isInputLikeElement(target) ? 'input' : 'click';
    const detail = {
      tag: (target.tagName || '').toUpperCase(),
      id: target.id || '',
      name: target.getAttribute('name') || '',
      text: visibleText(target).slice(0, 120),
      testid: target.getAttribute('data-testid') || '',
      selector: buildRecorderSelector(target, selectorKind) || buildSelector(el) || buildSelector(target),
      value: target.value || '',
      inputtype: (target.type || 'text').toLowerCase(),
      captiontext: isField ? labelTextForControl(target).slice(0, 120) : '',
      contexttext: type === 'click' ? clickContextCaption(target) : '',
      placeholder: target.getAttribute('placeholder') || '',
      arialabel: (target.getAttribute('aria-label') || '').trim(),
      title: (target.getAttribute('title') || '').trim(),
      role: target.getAttribute('role') || '',
      checked: target.checked ? 'true' : 'false',
    };
    return detail;
  }

  window.__scenariaHeuristics = {
    cssEscape,
    visibleText,
    labelTextForControl,
    findControlForLabel,
    isInputLikeElement,
    inputElementTag,
    resolveInputFromPick,
    buildAdjacentLabelSelector,
    hasTextSelector,
    clickHasTextSelector,
    clickTagFor,
    clickableAncestor,
    findCanvas,
    buildCanvasSelector,
    isSignatureCanvas,
    buildInputSelector,
    buildSelector,
    buildMenuTriggerSelector,
    navScopeTag,
    scopedTextSelector,
    collect,
    generateCandidates,
    buildPickerResult,
    suggestAction,
    isElementVisible,
    isActionable,
    pickKindForElement,
    prefixPickerResult,
    prefixSelectorChain,
    prefixShadowChain,
    shadowHostPrefix,
    buildCssPathSelector,
    buildRecorderSelector,
  };
})();
