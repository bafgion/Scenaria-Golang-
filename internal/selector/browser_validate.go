package selector

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/browserconfig"
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type BrowserValidateOptions struct {
	BrowserName string
	Headless    bool
	BaseURL     string
	Timeout     time.Duration
}

type StepValidation struct {
	Line     int
	StepText string
	Selector string
	Status   string
	Message  string
}

func (v Validator) ValidateFeatureInBrowser(ctx context.Context, path string, feature *gherkin.Feature, opts BrowserValidateOptions) ([]ValidationIssue, error) {
	detailed, err := v.ValidateFeatureInBrowserDetailed(ctx, path, feature, opts)
	if err != nil {
		return nil, err
	}
	issues := make([]ValidationIssue, 0)
	for _, step := range detailed {
		if step.Status == "missing" || step.Status == "warning" {
			issues = append(issues, ValidationIssue{
				FeaturePath: path,
				Line:        step.Line,
				Selector:    step.Selector,
				Message:     step.Message,
			})
		}
	}
	return issues, nil
}

func (v Validator) ValidateFeatureInBrowserDetailed(ctx context.Context, path string, feature *gherkin.Feature, opts BrowserValidateOptions) ([]StepValidation, error) {
	startURL := firstGotoURL(feature, opts.BaseURL, path)
	if startURL == "" {
		return nil, nil
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 8 * time.Second
	}

	if err := paths.EnsurePlaywrightEngine(opts.BrowserName); err != nil {
		return nil, fmt.Errorf("install playwright: %w", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	defer pw.Stop()

	name := strings.ToLower(strings.TrimSpace(opts.BrowserName))
	if name == "" {
		name = "chromium"
	}
	launchOpts := browserconfig.LaunchOptions(name, opts.Headless, 0)
	var browser playwright.Browser
	switch name {
	case "chromium":
		browser, err = pw.Chromium.Launch(launchOpts)
	case "firefox":
		browser, err = pw.Firefox.Launch(launchOpts)
	case "webkit":
		browser, err = pw.WebKit.Launch(launchOpts)
	default:
		return nil, fmt.Errorf("unsupported browser %q (supported: chromium, firefox, webkit)", name)
	}
	if err != nil {
		return nil, err
	}
	defer browser.Close()

	appCfg, _ := settings.LoadDefaultAppSettings()
	httpCreds := httpauth.PlaywrightHTTPCredentials(startURL, appCfg)
	ctxOpts := browserconfig.NewContextOptions(opts.Headless, httpCreds)
	browserCtx, err := browser.NewContext(ctxOpts)
	if err != nil {
		return nil, err
	}
	defer browserCtx.Close()

	page, err := browserCtx.NewPage()
	if err != nil {
		return nil, err
	}
	if _, err := page.Goto(startURL); err != nil {
		return nil, fmt.Errorf("goto %q: %w", startURL, err)
	}

	issues := make([]StepValidation, 0)
	for _, runnable := range gherkin.ExpandFeatureAtPath(feature, path) {
		for _, step := range gherkin.FlattenSteps(runnable.Steps) {
			if step.Block != "" {
				continue
			}
			action, err := stepdsl.Parse(step)
			if err != nil {
				continue
			}
			stepText := strings.TrimSpace(step.Keyword + " " + step.Text)
			if action.Kind == "goto" {
				url := stepdsl.ResolveURL(action.Value1, opts.BaseURL)
				if _, err := page.Goto(url); err != nil {
					issues = append(issues, StepValidation{
						Line:     step.Line,
						StepText: stepText,
						Status:   "missing",
						Message:  fmt.Sprintf("goto failed: %v", err),
					})
				} else {
					issues = append(issues, StepValidation{
						Line:     step.Line,
						StepText: stepText,
						Status:   "found",
						Message:  "страница открыта",
					})
				}
				continue
			}
			selectors := selectorsFromAction(action)
			if len(selectors) == 0 {
				issues = append(issues, StepValidation{
					Line:     step.Line,
					StepText: stepText,
					Status:   "skipped",
					Message:  "без селектора",
				})
				continue
			}
			expectHidden := action.Kind == "assert-hidden" || action.Kind == "wait-hidden"
			for _, sel := range selectors {
				result := StepValidation{
					Line:     step.Line,
					StepText: stepText,
					Selector: sel,
				}
				if expectHidden {
					if err := v.validateHidden(ctx, page, sel, opts.Timeout); err != nil {
						result.Status = "warning"
						result.Message = err.Error()
					} else {
						result.Status = "found"
						result.Message = "элемент скрыт"
					}
				} else if err := v.ValidateVisible(ctx, page, sel, opts.Timeout); err != nil {
					result.Status = "missing"
					result.Message = err.Error()
				} else {
					result.Status = "found"
					result.Message = "элемент найден"
				}
				issues = append(issues, result)
			}
		}
	}
	return issues, nil
}

func (v Validator) validateHidden(ctx context.Context, page playwright.Page, selector string, timeout time.Duration) error {
	if err := ValidateSyntax(selector); err != nil {
		return err
	}
	locator := page.Locator(selector)
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	}); err != nil {
		return fmt.Errorf("selector %q is still visible: %w", selector, err)
	}
	return nil
}

func firstGotoURL(feature *gherkin.Feature, baseURL, featurePath string) string {
	for _, runnable := range gherkin.ExpandFeatureAtPath(feature, featurePath) {
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
