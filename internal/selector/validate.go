package selector

import (
	"context"
	"fmt"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type ValidationIssue struct {
	FeaturePath string
	Line        int
	Selector    string
	Message     string
}

type Validator struct {
	Headless bool
}

func (v Validator) ValidateFeature(path string, feature *gherkin.Feature) ([]ValidationIssue, error) {
	issues := make([]ValidationIssue, 0)
	for _, runnable := range gherkin.ExpandFeatureAtPath(feature, path) {
		for _, step := range gherkin.FlattenSteps(runnable.Steps) {
			if step.Block != "" {
				continue
			}
			action, err := stepdsl.Parse(step)
			if err != nil {
				continue
			}
			selectors := selectorsFromAction(action)
			for _, sel := range selectors {
				if err := ValidateSyntax(sel); err != nil {
					issues = append(issues, ValidationIssue{
						FeaturePath: path,
						Line:        step.Line,
						Selector:    sel,
						Message:     err.Error(),
					})
				}
			}
		}
	}
	return issues, nil
}

func (v Validator) ValidateVisible(ctx context.Context, page playwright.Page, selector string, timeout time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := ValidateSyntax(selector); err != nil {
		return err
	}
	locator := page.Locator(selector)
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	}); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func selectorsFromAction(action stepdsl.Action) []string {
	switch action.Kind {
	case "click", "double-click", "hover", "clear", "check", "uncheck", "scroll-to",
		"assert-visible", "assert-hidden", "assert-enabled", "assert-disabled", "assert-selected",
		"wait-visible", "wait-hidden", "wait-enabled", "wait-disabled", "press-in",
		"download-click", "draw-signature", "remember-field", "remember-number":
		return []string{action.Value1}
	case "fill", "select", "upload", "assert-text", "assert-text-regex", "fill-generated", "prompt-email-code":
		if action.Value2 != "" {
			return []string{action.Value2}
		}
	}
	return nil
}
