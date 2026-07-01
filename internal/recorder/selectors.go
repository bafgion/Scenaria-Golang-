package recorder

import (
	"regexp"
	"strings"
)

var recordedHasTextRE = regexp.MustCompile(`(?:button|a|div):has-text\("((?:[^"\\]|\\.)*)"\)`)

// canonicalizeRecordedSelector maps legacy :has-text() selectors to Playwright text= engine.
func canonicalizeRecordedSelector(sel string) string {
	sel = strings.TrimSpace(sel)
	if sel == "" {
		return sel
	}
	return recordedHasTextRE.ReplaceAllStringFunc(sel, func(match string) string {
		sub := recordedHasTextRE.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		text := strings.ReplaceAll(sub[1], `\"`, `"`)
		text = strings.ReplaceAll(text, `\\`, `\`)
		return `text="` + text + `"`
	})
}
