package selector

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed heuristics_script.js
var heuristicsScriptJS string

//go:embed recorder_script.js
var recorderListenersJS string

//go:embed browser_toolbar.js
var browserToolbarJS string

// HeuristicsJS exposes selector helpers for picker and recorder (no event listeners).
var HeuristicsJS = heuristicsScriptJS

// RecorderListenersJS attaches DOM listeners; requires HeuristicsJS on the page first.
var RecorderListenersJS = recorderListenersJS

// RecorderScript installs heuristics + recorder listeners on the current document.
var RecorderScript = heuristicsScriptJS + "\n" + recorderListenersJS

// RecorderHeuristicsJS is the full in-browser recorder bundle (legacy name).
var RecorderHeuristicsJS = RecorderScript

// BrowserToolbarJS is injected as a fixed panel inside the Playwright browser window.
var BrowserToolbarJS = browserToolbarJS

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
	Title       string
}

func BuildFromElement(el ElementInfo) string {
	tag := strings.ToLower(strings.TrimSpace(el.Tag))
	if tag == "" {
		tag = "*"
	}
	if tag == "canvas" {
		if testID := strings.TrimSpace(el.TestID); testID != "" {
			return fmt.Sprintf(`[data-testid=%q]`, testID)
		}
		if aria := strings.TrimSpace(el.AriaLabel); aria != "" {
			return fmt.Sprintf(`canvas[aria-label=%q]`, aria)
		}
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
	if title := strings.TrimSpace(el.Title); title != "" {
		return fmt.Sprintf(`[title=%q]`, title)
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
