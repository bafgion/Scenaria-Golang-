package gherkin

import "strings"

func lineIndentLevel(raw string) int {
	if strings.TrimSpace(raw) == "" {
		return -1
	}
	prefix := raw[:len(raw)-len(strings.TrimLeft(raw, " \t"))]
	tabs := strings.Count(prefix, "\t")
	if tabs > 0 {
		return tabs
	}
	spaces := len(prefix) - strings.Count(prefix, "\t")
	if spaces == 0 {
		return 0
	}
	if spaces >= 2 {
		return spaces / 2
	}
	return 1
}
