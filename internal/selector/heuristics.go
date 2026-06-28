package selector

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed recorder_script.js
var recorderScriptJS string

// RecorderHeuristicsJS is injected into the browser during live recording.
var RecorderHeuristicsJS = recorderScriptJS

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
