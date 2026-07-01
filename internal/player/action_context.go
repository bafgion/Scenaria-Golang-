package player

import (
	"context"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

const (
	// ActionTimeoutMs is the default per-action timeout (click, fill, etc.).
	ActionTimeoutMs = 9000
	// LocatorWaitTimeoutMs is the default WaitFor timeout for asserts and wait-* steps.
	LocatorWaitTimeoutMs = 30000
)

// timeoutMs returns a Playwright timeout capped by the context deadline when present.
func timeoutMs(ctx context.Context, defaultMs float64) *float64 {
	if ctx == nil {
		return playwright.Float(defaultMs)
	}
	if err := ctx.Err(); err != nil {
		return playwright.Float(1)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return playwright.Float(defaultMs)
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return playwright.Float(1)
	}
	ms := float64(remaining.Milliseconds())
	if ms > defaultMs {
		ms = defaultMs
	}
	if ms < 1 {
		ms = 1
	}
	return playwright.Float(ms)
}

func capWaitDuration(ctx context.Context, duration time.Duration) time.Duration {
	if ctx == nil {
		return duration
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return duration
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return 0
	}
	if duration > remaining {
		return remaining
	}
	return duration
}

func waitForLocator(ctx context.Context, locator playwright.Locator, opts playwright.LocatorWaitForOptions) error {
	if opts.Timeout == nil {
		opts.Timeout = timeoutMs(ctx, LocatorWaitTimeoutMs)
	}
	err := locator.WaitFor(opts)
	if ctx != nil && ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
