package player

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/bafgion/scenaria-golang/internal/selector"
	playwright "github.com/mxschmitt/playwright-go"
)

func locatorPollTimeout(ctx context.Context) time.Duration {
	if ctx == nil {
		return LocatorWaitTimeoutMs * time.Millisecond
	}
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return 0
		}
		if remaining < LocatorWaitTimeoutMs*time.Millisecond {
			return remaining
		}
	}
	return LocatorWaitTimeoutMs * time.Millisecond
}

func waitForLocatorEnabledState(ctx context.Context, page playwright.Page, sel string, wantEnabled bool) error {
	locator := selector.ResolveChainedLocator(page, sel)
	deadline := time.Now().Add(locatorPollTimeout(ctx))
	backoff := 100 * time.Millisecond
	state := "enabled"
	if !wantEnabled {
		state = "disabled"
	}
	for {
		visible, visErr := locator.IsVisible()
		if visErr == nil && visible {
			enabled, err := locator.IsEnabled()
			if err == nil && enabled == wantEnabled {
				return nil
			}
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return fmt.Errorf("timed out waiting for element %q to be %s", sel, state)
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

func assertLocatorEnabledState(ctx context.Context, page playwright.Page, sel string, wantEnabled bool) error {
	locator := selector.ResolveChainedLocator(page, sel)
	if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		if wantEnabled {
			return fmt.Errorf("expected element %q to be enabled: %w", sel, err)
		}
		return fmt.Errorf("expected element %q to be disabled: %w", sel, err)
	}
	enabled, err := locator.IsEnabled()
	if err != nil {
		return fmt.Errorf("check enabled state for %q: %w", sel, err)
	}
	if wantEnabled && !enabled {
		return fmt.Errorf("expected element %q to be enabled", sel)
	}
	if !wantEnabled && enabled {
		return fmt.Errorf("expected element %q to be disabled", sel)
	}
	return nil
}

func assertLocatorTextMatchesRegex(ctx context.Context, page playwright.Page, pattern, sel string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex %q: %w", pattern, err)
	}
	locator := selector.ResolveChainedLocator(page, sel)
	if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		return fmt.Errorf("assert text regex failed: %w", err)
	}
	text, err := locator.InnerText()
	if err != nil {
		return fmt.Errorf("read text failed: %w", err)
	}
	if !re.MatchString(text) {
		return fmt.Errorf("expected text matching %q in %q, got %q", pattern, sel, text)
	}
	return nil
}

func assertLocatorSelected(ctx context.Context, page playwright.Page, sel string) error {
	locator := selector.ResolveChainedLocator(page, sel)
	if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		return fmt.Errorf("expected element %q to be selected: %w", sel, err)
	}
	selected, err := locator.Evaluate(`el => {
		if (el.getAttribute('aria-selected') === 'true') return true;
		if (el.getAttribute('aria-checked') === 'true') return true;
		if (el.matches('[aria-selected="true"], [aria-checked="true"]')) return true;
		if (el.classList.contains('active') || el.classList.contains('selected')) return true;
		if (el.tagName === 'OPTION' && el.selected) return true;
		return false;
	}`, nil)
	if err != nil {
		return fmt.Errorf("check selected state for %q: %w", sel, err)
	}
	ok, _ := selected.(bool)
	if !ok {
		return fmt.Errorf("expected element %q to be selected", sel)
	}
	return nil
}
