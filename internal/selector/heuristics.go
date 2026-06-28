package selector

import (
	"fmt"
	"regexp"
	"strings"
)

var cssEscapeRE = regexp.MustCompile(`([!"#$%&'()*+,./:;<=>?@[\\\]^` + "`" + `{|}~])`)

// ElementInfo describes a DOM element for selector building (recording).
type ElementInfo struct {
	Tag         string
	ID          string
	Classes     []string
	Name        string
	Type        string
	Placeholder string
	Role        string
	Label       string
	Text        string
	TestID      string
	AriaLabel   string
}

func BuildFromElement(el ElementInfo) string {
	tag := strings.ToLower(strings.TrimSpace(el.Tag))
	if tag == "" {
		tag = "*"
	}
	if id := strings.TrimSpace(el.ID); id != "" {
		return "#" + cssEscape(id)
	}
	if testID := strings.TrimSpace(el.TestID); testID != "" {
		return fmt.Sprintf(`[data-testid=%q]`, testID)
	}
	if name := strings.TrimSpace(el.Name); name != "" {
		return fmt.Sprintf(`%s[name=%q]`, tag, name)
	}
	if aria := strings.TrimSpace(el.AriaLabel); aria != "" {
		return fmt.Sprintf(`[aria-label=%q]`, aria)
	}
	if label := strings.TrimSpace(el.Label); label != "" {
		return fmt.Sprintf(`text=%q`, label)
	}
	if text := strings.TrimSpace(el.Text); text != "" && len([]rune(text)) <= 80 {
		return fmt.Sprintf(`text=%q`, text)
	}
	if placeholder := strings.TrimSpace(el.Placeholder); placeholder != "" {
		return fmt.Sprintf(`%s[placeholder=%q]`, tag, placeholder)
	}
	if role := strings.TrimSpace(el.Role); role != "" {
		return fmt.Sprintf(`role=%s`, role)
	}
	for _, className := range el.Classes {
		className = strings.TrimSpace(className)
		if className != "" && !strings.HasPrefix(className, "ng-") {
			return tag + "." + cssEscape(className)
		}
	}
	if typ := strings.TrimSpace(el.Type); typ != "" {
		return fmt.Sprintf(`%s[type=%q]`, tag, typ)
	}
	return tag
}

func cssEscape(value string) string {
	return cssEscapeRE.ReplaceAllStringFunc(value, func(s string) string {
		return `\` + s
	})
}

// RecorderHeuristicsJS is injected into the browser during live recording.
const RecorderHeuristicsJS = `(() => {
  if (window.__scenariaRecorder) return;
  const visibleText = (el) => {
    if (!el) return '';
    const clone = el.cloneNode(true);
    clone.querySelectorAll('script,style,input,textarea').forEach(n => n.remove());
    return (clone.textContent || '').replace(/\s+/g, ' ').trim();
  };
  const clickableAncestor = (el) => {
    let node = el;
    while (node && node.nodeType === 1) {
      const tag = (node.tagName || '').toLowerCase();
      if (tag === 'button' || tag === 'a' || node.getAttribute('role') === 'button') return node;
      if (node.onclick) return node;
      node = node.parentElement;
    }
    return el;
  };
  const buildSelector = (el) => {
    if (!el) return '';
    const target = clickableAncestor(el);
    if (target.id) return '#' + CSS.escape(target.id);
    const testId = target.getAttribute('data-testid');
    if (testId) return '[data-testid="' + testId + '"]';
    const name = target.getAttribute('name');
    if (name) return target.tagName.toLowerCase() + '[name="' + name + '"]';
    const text = visibleText(target);
    if (text && text.length <= 80) return 'text="' + text + '"';
    return target.tagName.toLowerCase();
  };
  const collect = (el, type) => {
    const target = type === 'click' ? clickableAncestor(el) : el;
    return {
      tag: (target.tagName || '').toUpperCase(),
      id: target.id || '',
      name: target.getAttribute('name') || '',
      text: visibleText(target).slice(0, 80),
      testid: target.getAttribute('data-testid') || '',
      selector: buildSelector(target),
      value: target.value || ''
    };
  };
  window.__scenariaRecorder = { events: [] };
  const push = (type, el) => window.__scenariaRecorder.events.push({ type, detail: collect(el, type), ts: Date.now() });
  document.addEventListener('click', (e) => { if (e.target) push('click', e.target); }, true);
  document.addEventListener('input', (e) => { if (e.target) push('input', e.target); }, true);
})();`
