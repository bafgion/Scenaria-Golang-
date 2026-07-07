package player

import (
	"context"
	"fmt"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

func needsSequentialFill(selector, generatorKind string) bool {
	if canonical, ok := NormalizeGeneratorName(generatorKind); ok && canonical == "phone" {
		return true
	}
	lower := strings.ToLower(selector)
	return strings.Contains(lower, "type=tel") || strings.Contains(lower, "[tel]")
}

func fillLocatorInput(ctx context.Context, locator playwright.Locator, selector, value, generatorKind string) error {
	timeout := timeoutMs(ctx, ActionTimeoutMs)
	if needsSequentialFill(selector, generatorKind) {
		if err := locator.Click(playwright.LocatorClickOptions{Timeout: timeout}); err != nil {
			return fmt.Errorf("focus for sequential fill failed: %w", err)
		}
		if err := locator.PressSequentially(value, playwright.LocatorPressSequentiallyOptions{Timeout: timeout}); err != nil {
			return fmt.Errorf("sequential fill failed: %w", err)
		}
		return nil
	}
	if err := locator.Fill(value, playwright.LocatorFillOptions{Timeout: timeout}); err != nil {
		return fmt.Errorf("fill failed: %w", err)
	}
	return nil
}
