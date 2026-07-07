package player

import (
	"fmt"
	"strings"
)

const maxSnapshotRunes = 6000

func captureDOMSnapshot(session *browserSession) string {
	page, err := session.currentPage()
	if err != nil {
		return ""
	}
	html, err := page.Content()
	if err != nil {
		return ""
	}
	html = strings.TrimSpace(html)
	return truncateSnapshot(html, maxSnapshotRunes)
}

func captureA11ySnapshot(session *browserSession) string {
	page, err := session.currentPage()
	if err != nil {
		return ""
	}
	raw, err := page.Evaluate(a11ySnapshotJS)
	if err != nil || raw == nil {
		return ""
	}
	text, _ := raw.(string)
	return truncateSnapshot(strings.TrimSpace(text), maxSnapshotRunes)
}

func truncateSnapshot(s string, maxRunes int) string {
	if maxRunes <= 0 || s == "" {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "\n… (truncated)"
}

const a11ySnapshotJS = `() => {
  const lines = [];
  const interesting = new Set(['button','link','input','heading','img','navigation','main','form','textbox','checkbox','radio','combobox','list','listitem','tab','dialog','alert']);
  const walk = (el, depth) => {
    if (!el || depth > 14 || lines.length > 90) return;
    const role = (el.getAttribute && el.getAttribute('role')) || (el.tagName || '').toLowerCase();
    const label = (el.getAttribute && (el.getAttribute('aria-label') || el.getAttribute('title'))) || '';
    let text = '';
    if (el.childNodes && el.childNodes.length === 1 && el.childNodes[0].nodeType === 3) {
      text = (el.textContent || '').trim().slice(0, 48);
    }
    if (interesting.has(role) || label || (text && text.length > 1)) {
      const pad = '  '.repeat(Math.min(depth, 8));
      let line = pad + role;
      if (label) line += ' "' + label + '"';
      if (text) line += ' ' + text;
      lines.push(line);
    }
    const kids = el.children || [];
    for (let i = 0; i < kids.length; i++) walk(kids[i], depth + 1);
  };
  if (document.body) walk(document.body, 0);
  return lines.join('\n');
}`

func formatSnapshotPreview(kind string, n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s snapshot (%d chars)", kind, n)
}
