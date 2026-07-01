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

func sleepContext(ctx context.Context, duration time.Duration) error {
	if ctx == nil {
		time.Sleep(duration)
		return nil
	}
	duration = capWaitDuration(ctx, duration)
	if duration <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func waitForLocator(ctx context.Context, locator playwright.Locator, opts playwright.LocatorWaitForOptions) error {
	if opts.Timeout == nil {
		opts.Timeout = timeoutMs(ctx, LocatorWaitTimeoutMs)
	}
	if ctx == nil {
		return locator.WaitFor(opts)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- locator.WaitFor(opts)
	}()
	select {
	case <-ctx.Done():
		go func() { _ = <-errCh }()
		return ctx.Err()
	case err := <-errCh:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
}

func pressKey(ctx context.Context, page playwright.Page, key string) error {
	if ctx == nil {
		return page.Keyboard().Press(key)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- page.Keyboard().Press(key)
	}()
	select {
	case <-ctx.Done():
		go func() { _ = <-errCh }()
		return ctx.Err()
	case err := <-errCh:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
}

func expectDownload(ctx context.Context, page playwright.Page, trigger func() error) (playwright.Download, error) {
	if ctx == nil {
		return page.ExpectDownload(trigger)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	type result struct {
		download playwright.Download
		err      error
	}
	ch := make(chan result, 1)
	go func() {
		download, err := page.ExpectDownload(trigger)
		ch <- result{download, err}
	}()
	select {
	case <-ctx.Done():
		go func() { _ = <-ch }()
		return nil, ctx.Err()
	case r := <-ch:
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return r.download, r.err
	}
}
