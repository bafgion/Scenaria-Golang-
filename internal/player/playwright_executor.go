package player

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type PlaywrightExecutorOptions struct {
	BrowserName string
	Headless    bool
	BaseURL     string
	AutoInstall bool
}

type PlaywrightExecutor struct {
	options PlaywrightExecutorOptions
}

func NewPlaywrightExecutor(options PlaywrightExecutorOptions) *PlaywrightExecutor {
	return &PlaywrightExecutor{options: options}
}

func (e *PlaywrightExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	result := ScenarioResult{
		FeaturePath: input.FeaturePath,
		Scenario:    input.Scenario.Title,
		Status:      "passed",
	}

	if e.options.AutoInstall {
		if err := playwright.Install(); err != nil {
			return ScenarioResult{}, fmt.Errorf("playwright install failed: %w", err)
		}
	}

	pw, err := playwright.Run()
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("start playwright: %w", err)
	}
	defer func() { _ = pw.Stop() }()

	browser, err := launchBrowser(pw, e.options)
	if err != nil {
		return ScenarioResult{}, err
	}
	defer func() { _ = browser.Close() }()

	page, err := browser.NewPage()
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("create browser page: %w", err)
	}

	steps := make([]gherkin.Step, 0, len(input.Feature.Background)+len(input.Scenario.Steps))
	steps = append(steps, input.Feature.Background...)
	steps = append(steps, input.Scenario.Steps...)

	for _, step := range steps {
		action, parseErr := stepdsl.Parse(step)
		if parseErr != nil {
			result.Status = "failed"
			result.Message = parseErr.Error()
			return result, nil
		}
		if err := executeAction(ctx, page, action, e.options.BaseURL); err != nil {
			result.Status = "failed"
			result.Message = fmt.Sprintf("line %d: %v", step.Line, err)
			return result, nil
		}
	}

	return result, nil
}

func launchBrowser(pw *playwright.Playwright, options PlaywrightExecutorOptions) (playwright.Browser, error) {
	name := strings.ToLower(strings.TrimSpace(options.BrowserName))
	if name == "" {
		name = "chromium"
	}

	launchOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(options.Headless),
	}

	switch name {
	case "chromium":
		return pw.Chromium.Launch(launchOpts)
	case "firefox":
		return pw.Firefox.Launch(launchOpts)
	case "webkit":
		return pw.WebKit.Launch(launchOpts)
	default:
		return nil, fmt.Errorf("unsupported browser %q (supported: chromium, firefox, webkit)", name)
	}
}

func executeAction(ctx context.Context, page playwright.Page, action stepdsl.Action, baseURL string) error {
	switch action.Kind {
	case "goto":
		_, err := page.Goto(stepdsl.ResolveURL(action.Value1, baseURL))
		if err != nil {
			return fmt.Errorf("goto failed: %w", err)
		}
		return nil
	case "click":
		if err := page.Click(action.Value1); err != nil {
			return fmt.Errorf("click failed: %w", err)
		}
		return nil
	case "fill":
		if err := page.Fill(action.Value2, action.Value1); err != nil {
			return fmt.Errorf("fill failed: %w", err)
		}
		return nil
	case "assert-text":
		content, err := page.Content()
		if err != nil {
			return fmt.Errorf("read page content failed: %w", err)
		}
		if !strings.Contains(content, action.Value1) {
			return fmt.Errorf("expected text %q not found", action.Value1)
		}
		return nil
	case "assert-url-contains":
		if !strings.Contains(page.URL(), action.Value1) {
			return fmt.Errorf("expected URL to contain %q, got %q", action.Value1, page.URL())
		}
		return nil
	case "press":
		if err := page.Keyboard().Press(action.Value1); err != nil {
			return fmt.Errorf("key press failed: %w", err)
		}
		return nil
	case "wait":
		duration, err := time.ParseDuration(action.Value1)
		if err != nil {
			return fmt.Errorf("invalid wait duration %q: %w", action.Value1, err)
		}
		timer := time.NewTimer(duration)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	default:
		return fmt.Errorf("unsupported action kind %q", action.Kind)
	}
}
