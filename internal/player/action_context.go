package player

import (
	"context"
	"fmt"
	"strings"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

const (
	gotoMaxAttempts = 3
	gotoRetryDelay  = 500 * time.Millisecond
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

func pageGoto(ctx context.Context, page playwright.Page, url string, waitUntil *playwright.WaitUntilState) error {
	if page == nil {
		return fmt.Errorf("browser page is not available")
	}
	if UrlsMatch(page.URL(), url) {
		return nil
	}
	var lastErr error
	for attempt := 0; attempt < gotoMaxAttempts; attempt++ {
		if attempt > 0 {
			if err := sleepContext(ctx, gotoRetryDelay); err != nil {
				return err
			}
		}
		lastErr = pageGotoOnce(ctx, page, url, waitUntil)
		if lastErr == nil {
			return nil
		}
		if ctx != nil && ctx.Err() != nil {
			return lastErr
		}
		if !isRetryableGotoError(lastErr) {
			return lastErr
		}
	}
	return lastErr
}

func isRetryableGotoError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "net::") ||
		strings.Contains(msg, "connection") ||
		strings.Contains(msg, "econnreset") ||
		strings.Contains(msg, "err_connection")
}

func pageGotoOnce(ctx context.Context, page playwright.Page, url string, waitUntil *playwright.WaitUntilState) error {
	opts := playwright.PageGotoOptions{
		WaitUntil: waitUntil,
		Timeout:   timeoutMs(ctx, NavTimeoutMs),
	}
	if ctx == nil {
		if _, err := page.Goto(url, opts); err != nil {
			return fmt.Errorf("goto failed: %w", err)
		}
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		_, err := page.Goto(url, opts)
		errCh <- err
	}()
	select {
	case <-ctx.Done():
		drainAsync(errCh)
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("goto failed: %w", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return nil
	}
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
		drainAsync(errCh)
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
		drainAsync(errCh)
		return ctx.Err()
	case err := <-errCh:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
}

func expectDownload(ctx context.Context, page playwright.Page, trigger func() error) (playwright.Download, error) {
	if page == nil {
		return nil, fmt.Errorf("browser page is not available")
	}
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
		drainChan(ch)
		return nil, ctx.Err()
	case r := <-ch:
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return r.download, r.err
	}
}
