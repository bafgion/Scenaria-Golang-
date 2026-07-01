package player

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const defaultURLAssertTimeout = 5 * time.Second

func urlAssertTimeout(ctx context.Context) time.Duration {
	if ctx == nil {
		return defaultURLAssertTimeout
	}
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return 0
		}
		if remaining < defaultURLAssertTimeout {
			return remaining
		}
	}
	return defaultURLAssertTimeout
}

func waitForURL(ctx context.Context, page interface{ URL() string }, expected string, contains bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	deadline := time.Now().Add(urlAssertTimeout(ctx))
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	backoff := 100 * time.Millisecond
	for {
		current := page.URL()
		matched := false
		if contains {
			matched = strings.Contains(current, expected)
		} else {
			matched = UrlsMatch(current, expected)
		}
		if matched {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			if contains {
				return fmt.Errorf("expected URL to contain %q, got %q", expected, current)
			}
			return fmt.Errorf("expected URL %q, got %q", expected, current)
		}
		delay := backoff
		if delay > remaining {
			delay = remaining
		}
		if err := sleepContext(ctx, delay); err != nil {
			return err
		}
		if backoff < 500*time.Millisecond {
			backoff *= 2
		}
	}
}
