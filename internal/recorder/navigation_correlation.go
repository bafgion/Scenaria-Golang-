package recorder

import (
	"strings"
	"time"
)

const navCorrelationWindow = 15 * time.Second

type navCorrelation struct {
	pending bool
	since   time.Time
}

func (n *navCorrelation) noteNavCausingStep(at time.Time) {
	if n == nil {
		return
	}
	n.pending = true
	n.since = at
}

func (n *navCorrelation) expire(now time.Time) {
	if n == nil || !n.pending {
		return
	}
	if now.Sub(n.since) > navCorrelationWindow {
		n.pending = false
	}
}

// classifyURLNavigation decides how to record a URL change after browser events
// were drained for the current poll tick.
func classifyURLNavigation(correlation *navCorrelation, now time.Time, _ bool) (kind string) {
	if correlation == nil || !correlation.pending {
		return "goto"
	}
	if now.Sub(correlation.since) > navCorrelationWindow {
		correlation.pending = false
		return "goto"
	}
	correlation.pending = false
	return "wait-url"
}

func isNavCausingStep(step RecordedStep) bool {
	switch step.Action {
	case "click", "double-click", "download-click", "press-in", "check", "uncheck":
		return true
	case "press":
		key := strings.ToLower(strings.TrimSpace(step.Value))
		return key == "enter" || key == "return"
	default:
		return false
	}
}
