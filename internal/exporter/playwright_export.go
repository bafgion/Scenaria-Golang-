package exporter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func WritePlaywrightTS(path string, feature *gherkin.Feature, baseURL string) error {
	content, err := renderPlaywrightTS(feature, baseURL)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write ts export %q: %w", path, err)
	}
	return nil
}

func WritePlaywrightPython(path string, feature *gherkin.Feature, baseURL string) error {
	content, err := renderPlaywrightPython(feature, baseURL)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write python export %q: %w", path, err)
	}
	return nil
}

func renderPlaywrightTS(feature *gherkin.Feature, baseURL string) (string, error) {
	var b strings.Builder
	b.WriteString("import { test, expect } from '@playwright/test';\n\n")
	b.WriteString("test.describe(" + strconv.Quote(feature.Title) + ", () => {\n")

	for _, scenario := range feature.Scenarios {
		b.WriteString("  test(" + strconv.Quote(scenario.Title) + ", async ({ page }) => {\n")
		steps := mergeSteps(feature, scenario)
		for _, step := range steps {
			line, err := renderTSAction(step, baseURL)
			if err != nil {
				return "", err
			}
			b.WriteString("    " + line + "\n")
		}
		b.WriteString("  });\n\n")
	}
	b.WriteString("});\n")
	return b.String(), nil
}

func renderPlaywrightPython(feature *gherkin.Feature, baseURL string) (string, error) {
	var b strings.Builder
	b.WriteString("from playwright.sync_api import sync_playwright\n\n")
	b.WriteString("def run():\n")
	b.WriteString("    with sync_playwright() as p:\n")
	b.WriteString("        browser = p.chromium.launch(headless=True)\n")
	b.WriteString("        page = browser.new_page()\n\n")

	for _, scenario := range feature.Scenarios {
		b.WriteString("        # Scenario: " + scenario.Title + "\n")
		steps := mergeSteps(feature, scenario)
		for _, step := range steps {
			line, err := renderPythonAction(step, baseURL)
			if err != nil {
				return "", err
			}
			b.WriteString("        " + line + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("        browser.close()\n\n")
	b.WriteString("if __name__ == '__main__':\n")
	b.WriteString("    run()\n")
	return b.String(), nil
}

func mergeSteps(feature *gherkin.Feature, scenario gherkin.Scenario) []gherkin.Step {
	out := make([]gherkin.Step, 0, len(feature.Background)+len(scenario.Steps))
	out = append(out, feature.Background...)
	out = append(out, scenario.Steps...)
	return out
}

func renderTSAction(step gherkin.Step, baseURL string) (string, error) {
	action, err := stepdsl.Parse(step)
	if err != nil {
		return "", err
	}
	switch action.Kind {
	case "goto":
		return "await page.goto(" + strconv.Quote(stepdsl.ResolveURL(action.Value1, baseURL)) + ");", nil
	case "click":
		return "await page.click(" + strconv.Quote(action.Value1) + ");", nil
	case "fill":
		return "await page.fill(" + strconv.Quote(action.Value2) + ", " + strconv.Quote(action.Value1) + ");", nil
	case "assert-text":
		return "await expect(page.locator('body')).toContainText(" + strconv.Quote(action.Value1) + ");", nil
	case "assert-url-contains":
		return "expect(page.url()).toContain(" + strconv.Quote(action.Value1) + ");", nil
	case "press":
		return "await page.keyboard.press(" + strconv.Quote(action.Value1) + ");", nil
	case "wait":
		ms, err := durationToMS(action.Value1)
		if err != nil {
			return "", err
		}
		return "await page.waitForTimeout(" + strconv.Itoa(ms) + ");", nil
	default:
		return "", fmt.Errorf("unsupported action kind for TS export: %s", action.Kind)
	}
}

func renderPythonAction(step gherkin.Step, baseURL string) (string, error) {
	action, err := stepdsl.Parse(step)
	if err != nil {
		return "", err
	}
	switch action.Kind {
	case "goto":
		return "page.goto(" + strconv.Quote(stepdsl.ResolveURL(action.Value1, baseURL)) + ")", nil
	case "click":
		return "page.click(" + strconv.Quote(action.Value1) + ")", nil
	case "fill":
		return "page.fill(" + strconv.Quote(action.Value2) + ", " + strconv.Quote(action.Value1) + ")", nil
	case "assert-text":
		return "assert " + strconv.Quote(action.Value1) + " in page.content()", nil
	case "assert-url-contains":
		return "assert " + strconv.Quote(action.Value1) + " in page.url", nil
	case "press":
		return "page.keyboard.press(" + strconv.Quote(action.Value1) + ")", nil
	case "wait":
		ms, err := durationToMS(action.Value1)
		if err != nil {
			return "", err
		}
		return "page.wait_for_timeout(" + strconv.Itoa(ms) + ")", nil
	default:
		return "", fmt.Errorf("unsupported action kind for Python export: %s", action.Kind)
	}
}

func durationToMS(raw string) (int, error) {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid wait duration %q: %w", raw, err)
	}
	return int(d.Milliseconds()), nil
}
