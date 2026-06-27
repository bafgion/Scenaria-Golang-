package player

import (
	"context"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
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

func (e *PlaywrightExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
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
		action, parseErr := parseStepAction(step)
		if parseErr != nil {
			result.Status = "failed"
			result.Message = parseErr.Error()
			return result, nil
		}
		if err := executeAction(page, action, e.options.BaseURL); err != nil {
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

type stepAction struct {
	Kind   string
	Value1 string
	Value2 string
}

func parseStepAction(step gherkin.Step) (stepAction, error) {
	text := strings.TrimSpace(step.Text)
	lower := strings.ToLower(text)
	values := extractQuotedValues(text)

	switch {
	case strings.HasPrefix(lower, "открываю ") || strings.HasPrefix(lower, "перехожу на "):
		if len(values) < 1 {
			return stepAction{}, fmt.Errorf("line %d: navigation step must contain URL in quotes", step.Line)
		}
		return stepAction{Kind: "goto", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "нажимаю ") || strings.HasPrefix(lower, "кликаю "):
		if len(values) < 1 {
			return stepAction{}, fmt.Errorf("line %d: click step must contain selector in quotes", step.Line)
		}
		return stepAction{Kind: "click", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "ввожу "):
		if len(values) < 2 {
			return stepAction{}, fmt.Errorf("line %d: fill step must contain value and selector in quotes", step.Line)
		}
		return stepAction{Kind: "fill", Value1: values[0], Value2: values[1]}, nil
	case strings.HasPrefix(lower, "вижу "):
		if len(values) < 1 {
			return stepAction{}, fmt.Errorf("line %d: assertion step must contain expected text in quotes", step.Line)
		}
		return stepAction{Kind: "assert-text", Value1: values[0]}, nil
	default:
		return stepAction{}, fmt.Errorf("line %d: unsupported step text %q for playwright engine", step.Line, step.Text)
	}
}

func executeAction(page playwright.Page, action stepAction, baseURL string) error {
	switch action.Kind {
	case "goto":
		_, err := page.Goto(resolveURL(action.Value1, baseURL))
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
	default:
		return fmt.Errorf("unsupported action kind %q", action.Kind)
	}
}

func resolveURL(raw string, baseURL string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	if strings.TrimSpace(baseURL) == "" {
		return trimmed
	}
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasPrefix(trimmed, "/") {
		return base + trimmed
	}
	return base + "/" + strings.TrimLeft(trimmed, "/")
}

func extractQuotedValues(s string) []string {
	values := make([]string, 0)
	for {
		start := strings.Index(s, `"`)
		if start < 0 {
			break
		}
		s = s[start+1:]
		end := strings.Index(s, `"`)
		if end < 0 {
			break
		}
		values = append(values, s[:end])
		s = s[end+1:]
	}
	return values
}
