package player

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var digitsOnlyRE = regexp.MustCompile(`\d`)

// ParseNumberFromText extracts a normalized integer string from UI text like "7 180 ₽" or "1 795".
func ParseNumberFromText(text string) (string, error) {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = strings.TrimSpace(text)
	for _, cur := range []string{"₽", "руб.", "руб", "$", "€", "EUR", "USD"} {
		text = strings.ReplaceAll(text, cur, "")
	}
	text = strings.TrimSpace(text)
	digits := digitsOnlyRE.FindAllString(text, -1)
	if len(digits) == 0 {
		return "", fmt.Errorf("no number found in %q", text)
	}
	out := strings.Join(digits, "")
	if out == "" {
		return "", fmt.Errorf("no number found in %q", text)
	}
	// Reject strings that were only punctuation after stripping.
	if !strings.ContainsFunc(text, unicode.IsDigit) {
		return "", fmt.Errorf("no number found in %q", text)
	}
	return out, nil
}
