package selector

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type BrowserValidateOptions struct {
	Headless bool
	BaseURL  string
	Timeout  time.Duration
}

func (v Validator) ValidateFeatureInBrowser(ctx context.Context, path string, feature *gherkin.Feature, opts BrowserValidateOptions) ([]ValidationIssue, error) {
	startURL := firstGotoURL(feature, opts.BaseURL)
	if startURL == "" {
		return nil, fmt.Errorf("feature has no goto step for browser validation")
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 8 * time.Second
	}

	if err := playwright.Install(); err != nil {
		return nil, fmt.Errorf("install playwright: %w", err)
	}
	paths.ConfigurePlaywrightBrowsers()
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(opts.Headless),
	})
	if err != nil {
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	if _, err := page.Goto(startURL); err != nil {
		return nil, fmt.Errorf("goto %q: %w", startURL, err)
	}

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
			if action.Kind == "goto" {
				url := stepdsl.ResolveURL(action.Value1, opts.BaseURL)
				if _, err := page.Goto(url); err != nil {
					issues = append(issues, ValidationIssue{
						FeaturePath: path,
						Line:        step.Line,
						Message:     fmt.Sprintf("goto failed: %v", err),
					})
				}
				continue
			}
			for _, sel := range selectorsFromAction(action) {
				if err := v.ValidateVisible(ctx, page, sel, opts.Timeout); err != nil {
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

func firstGotoURL(feature *gherkin.Feature, baseURL string) string {
	for _, runnable := range gherkin.ExpandFeature(feature) {
		for _, step := range gherkin.FlattenSteps(runnable.Steps) {
			action, err := stepdsl.Parse(step)
			if err == nil && action.Kind == "goto" {
				return stepdsl.ResolveURL(action.Value1, baseURL)
			}
		}
	}
	if strings.TrimSpace(baseURL) != "" {
		return baseURL
	}
	return ""
}
