package player

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

const (
	// DefaultMaxActionRetries is the number of retries after the first attempt (3 tries total).
	DefaultMaxActionRetries = 2
	defaultRetryBackoff     = 200 * time.Millisecond
)

func isRetryableAction(kind string) bool {
	switch kind {
	case "click", "double-click", "hover", "fill", "check", "uncheck", "clear",
		"select", "select-option", "press-in", "scroll-to", "scroll-into-view",
		"drag-drop", "upload", "download-click",
		"assert-visible", "assert-hidden", "assert-enabled", "assert-disabled", "assert-selected",
		"assert-text", "assert-text-regex", "assert-value", "assert-count",
		"assert-url", "assert-url-contains",
		"wait-visible", "wait-hidden", "wait-enabled", "wait-disabled", "wait-url":
		return true
	default:
		return false
	}
}

// isRetryableStepError reports transient Playwright / network failures worth retrying.
func isRetryableStepError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if isRetryableGotoError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	permanent := []string{
		"text mismatch",
		"expected text",
		"expected \"",
		"does not match",
		"strict mode violation",
		"not a valid",
		"invalid selector",
		"unsupported action kind",
		"requires run context",
		"no downloaded file recorded",
	}
	for _, needle := range permanent {
		if strings.Contains(msg, needle) {
			return false
		}
	}
	transient := []string{
		"timeout",
		"timed out",
		"interrupted",
		"detached",
		"stale",
		"not attached",
		"execution context",
		"target closed",
		"element is not visible",
		"element is not enabled",
		"element is outside of the viewport",
		"waiting for",
		"not stable",
		"navigation",
		"frame was detached",
	}
	for _, needle := range transient {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}

func (e *StepExecutor) maxActionRetries() int {
	if e == nil || e.options.MaxActionRetries == 0 {
		return DefaultMaxActionRetries
	}
	if e.options.MaxActionRetries < 0 {
		return 0
	}
	return e.options.MaxActionRetries
}

func (e *StepExecutor) retryBackoff() time.Duration {
	if e != nil && e.options.RetryBackoff > 0 {
		return e.options.RetryBackoff
	}
	return defaultRetryBackoff
}

func (e *StepExecutor) runAction(ctx context.Context, session *browserSession, action stepdsl.Action, runCtx *RunContext) error {
	if !isRetryableAction(action.Kind) {
		return executeAction(ctx, session, action, e.options.BaseURL, runCtx)
	}
	return e.runWithRetries(ctx, action.Kind, func() error {
		return executeAction(ctx, session, action, e.options.BaseURL, runCtx)
	})
}

func (e *StepExecutor) runWithRetries(ctx context.Context, label string, fn func() error) error {
	attempts := e.maxActionRetries() + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt > 0 {
			logx.Debug("step retry", "kind", label, "attempt", attempt+1, "max", attempts, "error", lastErr)
			if err := sleepRetryBackoff(ctx, attempt, e.retryBackoff()); err != nil {
				return err
			}
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !isRetryableStepError(lastErr) {
			return lastErr
		}
	}
	return lastErr
}

func sleepRetryBackoff(ctx context.Context, attempt int, base time.Duration) error {
	if attempt < 1 {
		attempt = 1
	}
	delay := base
	for i := 1; i < attempt; i++ {
		delay *= 2
	}
	if jitter := retryJitter(delay); jitter > 0 {
		delay += jitter
	}
	delay = capWaitDuration(ctx, delay)
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func retryJitter(delay time.Duration) time.Duration {
	if delay <= 0 {
		return 0
	}
	max := delay / 4
	if max <= 0 {
		return 0
	}
	return time.Duration(time.Now().UnixNano() % int64(max))
}
