package recorder

import (
	"regexp"
	"strings"
)

var unstableSelectorValueRE = regexp.MustCompile(`(?i)(radix|:r[0-9]+|headlessui|react-aria|ember[0-9]|^:r[0-9])`)

// ClickHasTextSelector builds a Playwright-style tag:has-text selector (Python recorder parity).
func ClickHasTextSelector(label, tag string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}
	tag = normalizeClickTag(tag)
	return tag + `:has-text("` + escapeSelectorText(label) + `")`
}

// ContextualClickSelector builds div >> button style chains like the legacy Python recorder.
func ContextualClickSelector(context, label string) string {
	context = strings.TrimSpace(context)
	label = strings.TrimSpace(label)
	if context == "" || label == "" {
		return ""
	}
	if len(context) > 80 {
		context = context[:60]
	}
	if len(label) > 60 {
		label = label[:40]
	}
	return `div:has-text("` + escapeSelectorText(context) + `") >> button:has-text("` + escapeSelectorText(label) + `")`
}

func normalizeClickTag(tag string) string {
	switch strings.ToLower(strings.TrimSpace(tag)) {
	case "a", "link":
		return "a"
	case "button", "menuitem", "tab":
		return "button"
	case "div", "span", "li":
		return strings.ToLower(strings.TrimSpace(tag))
	default:
		return "button"
	}
}

func isUnstableSelectorValue(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	lower := strings.ToLower(value)
	if unstableSelectorValueRE.MatchString(lower) {
		return true
	}
	if strings.Contains(value, ":") && len(value) > 6 {
		return true
	}
	if len(value) > 28 {
		return true
	}
	if strings.Count(value, "-") >= 5 {
		return true
	}
	if regexp.MustCompile(`^[a-f0-9]{16,}$`).MatchString(lower) {
		return true
	}
	return false
}
