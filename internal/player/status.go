package player

import "strings"

// ResultStatusIsSuccessful is the UI/history success flag, not a failure counter.
// Only a scenario that actually passed should be rendered as successful.
func ResultStatusIsSuccessful(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "passed", "pass", "ok", "success", "dry-run", "dryrun":
		return true
	default:
		return false
	}
}
