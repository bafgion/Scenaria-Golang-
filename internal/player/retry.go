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

func (e *StepExecutor) runAction(ctx context.Context, session *browserSession, action stepdsl.Action, runCtx *RunContext, stepIdx int) error {
	if !e.shouldRetryAction(action.Kind) {
		return executeAction(ctx, session, action, e.options.BaseURL, runCtx)
	}
	attempts, err := e.runWithRetries(ctx, action.Kind, func() error {
		return executeAction(ctx, session, action, e.options.BaseURL, runCtx)
	})
	if runCtx != nil && stepIdx >= 0 && attempts > 1 {
		runCtx.setStepRetryAttempts(stepIdx, attempts-1)
	}
	return err
}

func (e *StepExecutor) runWithRetries(ctx context.Context, label string, fn func() error) (int, error) {
	attempts := e.maxActionRetries() + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return attempt + 1, err
		}
		if attempt > 0 {
			logx.Debug("step retry", "kind", label, "attempt", attempt+1, "max", attempts, "error", lastErr)
			if err := sleepRetryBackoff(ctx, attempt, e.retryBackoff()); err != nil {
				return attempt + 1, err
			}
		}
		lastErr = fn()
		if lastErr == nil {
			return attempt + 1, nil
		}
		if !isRetryableStepError(lastErr) {
			return attempt + 1, lastErr
		}
	}
	return attempts, lastErr
}
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
