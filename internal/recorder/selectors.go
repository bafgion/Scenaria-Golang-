package recorder

import (
	"strings"
)

// canonicalizeRecordedSelector keeps Playwright :has-text selectors as recorded (Python parity).
func canonicalizeRecordedSelector(sel string) string {
	return strings.TrimSpace(sel)
}
