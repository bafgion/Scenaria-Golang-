package player

import (
	"context"
	"time"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

const (
	// DefaultMaxActionRetries is the number of retries after the first attempt (3 tries total).
	DefaultMaxActionRetries = 2
	defaultRetryBackoff     = 200 * time.Millisecond
)

func isRetryableAction(kind string) bool {
	switch kind {
	case "click", "double-click", "hover", "fill", "check", "uncheck",
		"assert-visible", "assert-hidden", "assert-text",
		"wait-visible", "wait-hidden", "scroll-to", "drag-drop":
		return true
	default:
		return false
	}
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
	attempts := e.maxActionRetries() + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt > 0 {
			if err := sleepRetryBackoff(ctx, attempt, e.retryBackoff()); err != nil {
				return err
			}
		}
		lastErr = executeAction(ctx, session, action, e.options.BaseURL, runCtx)
		if lastErr == nil {
			return nil
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
