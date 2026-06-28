package selector

import (
	"fmt"
	"strings"
)

func Normalize(raw string) string {
	return strings.TrimSpace(raw)
}

func ValidateSyntax(selector string) error {
	trimmed := Normalize(selector)
	if trimmed == "" {
		return fmt.Errorf("selector is empty")
	}
	return nil
}

func Describe(selector string) string {
	trimmed := Normalize(selector)
	if strings.HasPrefix(trimmed, "#") {
		return "id"
	}
	if strings.HasPrefix(trimmed, ".") {
		return "class"
	}
	if strings.HasPrefix(trimmed, "text=") {
		return "text"
	}
	if strings.Contains(trimmed, ">>") {
		return "chain"
	}
	return "css"
}
