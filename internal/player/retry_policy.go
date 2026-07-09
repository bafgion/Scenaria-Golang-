package player

// RetryPolicy controls which action categories may be retried.
// Nil pointers use defaults: waits and assertions retry, risky actions do not.
type RetryPolicy struct {
	Waits       *bool
	Assertions  *bool
	Actions     *bool
}

func retryEnabled(flag *bool, defaultValue bool) bool {
	if flag == nil {
		return defaultValue
	}
	return *flag
}

func (p RetryPolicy) waitsEnabled() bool {
	return retryEnabled(p.Waits, true)
}

func (p RetryPolicy) assertionsEnabled() bool {
	return retryEnabled(p.Assertions, true)
}

func (p RetryPolicy) actionsEnabled() bool {
	return retryEnabled(p.Actions, false)
}

func isWaitRetryAction(kind string) bool {
	switch kind {
	case "wait-visible", "wait-hidden", "wait-enabled", "wait-disabled", "wait-url":
		return true
	default:
		return false
	}
}

func isAssertionRetryAction(kind string) bool {
	switch kind {
	case "assert-visible", "assert-hidden", "assert-enabled", "assert-disabled", "assert-selected",
		"assert-text", "assert-text-regex", "assert-value", "assert-count",
		"assert-url", "assert-url-contains":
		return true
	default:
		return false
	}
}

func isRiskyRetryAction(kind string) bool {
	switch kind {
	case "click", "double-click", "hover", "fill", "check", "uncheck", "clear",
		"select", "select-option", "press-in", "scroll-to", "scroll-into-view",
		"drag-drop", "upload", "download-click":
		return true
	default:
		return false
	}
}

func (e *StepExecutor) shouldRetryAction(kind string) bool {
	if e == nil || e.maxActionRetries() == 0 {
		return false
	}
	policy := e.options.RetryPolicy
	switch {
	case isWaitRetryAction(kind):
		return policy.waitsEnabled()
	case isAssertionRetryAction(kind):
		return policy.assertionsEnabled()
	case isRiskyRetryAction(kind):
		return policy.actionsEnabled()
	default:
		return false
	}
}
