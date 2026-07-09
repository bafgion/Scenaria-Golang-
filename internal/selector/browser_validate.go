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
	"github.com/bafgion/scenaria-golang/internal/playwrightrt"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type BrowserValidateOptions struct {
	BrowserName string
	Headless    bool
	BaseURL     string
	Timeout     time.Duration
	Mode        ValidationMode
}

type StepValidation struct {
	Line       int
	StepText   string
	Selector   string
	Status     string
	Message    string
	Mode       string
	ActionKind string
	MatchCount int
	Limitation string
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
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	startURL := firstGotoURL(feature, opts.BaseURL, path)
	if startURL == "" {
		return nil, nil
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 8 * time.Second
	}
	opts.Mode = NormalizeValidationMode(string(opts.Mode))
	limitation := validationLimitation(opts.Mode)

	if err := paths.EnsurePlaywrightEngine(opts.BrowserName); err != nil {
		return nil, fmt.Errorf("install playwright: %w", err)
	}
	pw, release, err := playwrightrt.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

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
	if _, err := page.Goto(startURL, playwright.PageGotoOptions{Timeout: playwright.Float(float64(opts.Timeout.Milliseconds()))}); err != nil {
		return nil, fmt.Errorf("goto %q: %w", startURL, err)
	}

	issues := make([]StepValidation, 0)
	for _, runnable := range gherkin.ExpandFeatureAtPath(feature, path) {
		steps := collectBrowserValidateSteps(runnable.Steps)
		flowReplayCursor := 0
		for i, item := range steps {
			if err := ctx.Err(); err != nil {
				return issues, err
			}
			if opts.Mode == ValidationModeFlow {
				for j := flowReplayCursor; j < i; j++ {
					if err := replaySafeStep(ctx, page, steps[j].action, opts.BaseURL, opts.Timeout); err != nil {
						issues = append(issues, StepValidation{
							Line:       steps[j].step.Line,
							StepText:   steps[j].text,
							Status:     "warning",
							Message:    fmt.Sprintf("flow replay: %v", err),
							Mode:       string(opts.Mode),
							ActionKind: steps[j].action.Kind,
							Limitation: limitation,
						})
					}
				}
				flowReplayCursor = i
			}
			step := item.step
			action := item.action
			stepText := item.text
			if action.Kind == "goto" {
				url := stepdsl.ResolveURL(action.Value1, opts.BaseURL)
				if _, err := page.Goto(url, playwright.PageGotoOptions{Timeout: playwright.Float(float64(opts.Timeout.Milliseconds()))}); err != nil {
					issues = append(issues, StepValidation{
						Line:       step.Line,
						StepText:   stepText,
						Status:     "missing",
						Message:    fmt.Sprintf("goto failed: %v", err),
						Mode:       string(opts.Mode),
						ActionKind: action.Kind,
						Limitation: limitation,
					})
				} else {
					if opts.Mode == ValidationModeFlow {
						flowReplayCursor = i + 1
					}
					issues = append(issues, StepValidation{
						Line:       step.Line,
						StepText:   stepText,
						Status:     "found",
						Message:    "страница открыта",
						Mode:       string(opts.Mode),
						ActionKind: action.Kind,
						Limitation: limitation,
					})
				}
				continue
			}
			selectors := selectorsFromAction(action)
			if len(selectors) == 0 {
				issues = append(issues, StepValidation{
					Line:       step.Line,
					StepText:   stepText,
					Status:     "skipped",
					Message:    "без селектора",
					Mode:       string(opts.Mode),
					ActionKind: action.Kind,
					Limitation: limitation,
				})
				continue
			}
			priorFlow := countPriorFlowSteps(steps, i)
			for _, sel := range selectors {
				result := StepValidation{
					Line:       step.Line,
					StepText:   stepText,
					Selector:   sel,
					Mode:       string(opts.Mode),
					ActionKind: action.Kind,
					Limitation: limitation,
				}
				validation := v.ValidateActionTarget(ctx, page, action, sel, opts.Timeout)
				result.MatchCount = validation.MatchCount
				if validation.OK {
					result.Status = "found"
					result.Message = validation.Message
					if len(validation.Warnings) > 0 {
						result.Status = "warning"
						result.Message = validation.Message + "; " + strings.Join(validation.Warnings, "; ")
					}
				} else if opts.Mode == ValidationModeStatic && priorFlow > 0 {
					result.Status = "warning"
					result.Message = validation.Message + "; возможно требуются предыдущие шаги сценария (dynamic UI)"
				} else {
					result.Status = "missing"
					result.Message = validation.Message
				}
				issues = append(issues, result)
			}
		}
	}
	return issues, nil
}

type browserValidateStep struct {
	step   gherkin.Step
	action stepdsl.Action
	text   string
}

func collectBrowserValidateSteps(steps []gherkin.Step) []browserValidateStep {
	out := make([]browserValidateStep, 0)
	for _, step := range gherkin.FlattenSteps(steps) {
		if step.Block != "" {
			continue
		}
		action, err := stepdsl.Parse(step)
		if err != nil {
			continue
		}
		out = append(out, browserValidateStep{
			step:   step,
			action: action,
			text:   strings.TrimSpace(step.Keyword + " " + step.Text),
		})
	}
	return out
}

func countPriorFlowSteps(steps []browserValidateStep, index int) int {
	count := 0
	for i := 0; i < index && i < len(steps); i++ {
		if isPriorFlowStep(steps[i].action) {
			count++
		}
	}
	return count
}

func isPriorFlowStep(action stepdsl.Action) bool {
	switch action.Kind {
	case "click", "double-click", "hover", "fill", "fill-generated", "select", "check", "uncheck",
		"scroll-to", "clear", "upload", "press-in", "download-click":
		return true
	default:
		return false
	}
}

func (v Validator) validateHidden(ctx context.Context, page playwright.Page, selector string, timeout time.Duration) error {
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
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	}); err != nil {
		return fmt.Errorf("selector %q is still visible: %w", selector, err)
	}
	if err := ctx.Err(); err != nil {
		return err
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
