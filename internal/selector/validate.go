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
	for _, runnable := range gherkin.ExpandFeature(feature) {
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
	return nil
}

func selectorsFromAction(action stepdsl.Action) []string {
	switch action.Kind {
	case "click", "double-click", "hover", "clear", "check", "uncheck", "scroll-to",
		"assert-visible", "assert-hidden", "wait-visible", "wait-hidden", "press-in",
		"download-click", "draw-signature", "remember-field":
		return []string{action.Value1}
	case "fill", "select", "upload", "assert-text", "fill-generated", "prompt-email-code":
		if action.Value2 != "" {
			return []string{action.Value2}
		}
	}
	return nil
}
