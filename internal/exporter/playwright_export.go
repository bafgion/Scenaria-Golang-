package exporter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	case "assert-visible":
		return "await expect(page.locator(" + strconv.Quote(action.Value1) + ")).toBeVisible();", nil
	case "assert-hidden":
		return "await expect(page.locator(" + strconv.Quote(action.Value1) + ")).toBeHidden();", nil
	case "assert-text":
		return "await expect(page.locator(" + strconv.Quote(action.Value2) + ")).toContainText(" + strconv.Quote(action.Value1) + ");", nil
	case "assert-url":
		return "expect(page.url()).toBe(" + strconv.Quote(action.Value1) + ");", nil
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
	case "wait-visible":
		return "await page.locator(" + strconv.Quote(action.Value1) + ").waitFor({ state: 'visible' });", nil
	case "wait-hidden":
		return "await page.locator(" + strconv.Quote(action.Value1) + ").waitFor({ state: 'hidden' });", nil
	case "reload":
		return "await page.reload();", nil
	case "go-back":
		return "await page.goBack();", nil
	case "close-browser":
		return "await page.context().browser()?.close();", nil
	case "check":
		return "await page.check(" + strconv.Quote(action.Value1) + ");", nil
	case "uncheck":
		return "await page.uncheck(" + strconv.Quote(action.Value1) + ");", nil
	case "select":
		return "await page.selectOption(" + strconv.Quote(action.Value2) + ", " + strconv.Quote(action.Value1) + ");", nil
	case "hover":
		return "await page.hover(" + strconv.Quote(action.Value1) + ");", nil
	case "double-click":
		return "await page.dblclick(" + strconv.Quote(action.Value1) + ");", nil
	case "clear":
		return "await page.fill(" + strconv.Quote(action.Value1) + ", '');", nil
	case "scroll-to":
		return "await page.locator(" + strconv.Quote(action.Value1) + ").scrollIntoViewIfNeeded();", nil
	case "remember-text", "remember-field", "remember-url":
		return "// " + action.Kind + " (runtime variable)", nil
	case "fill-generated":
		return "await page.fill(" + strconv.Quote(action.Value2) + ", /* generated: " + action.Value1 + " */ '');", nil
	case "switch-tab":
		return "// switch tab: " + action.Mode, nil
	case "close-tab":
		return "// close current tab", nil
	case "assert-tab-count":
		return fmt.Sprintf("expect(context.pages()).toHaveLength(%d);", action.IntVal), nil
	case "download-click":
		return "const download = await page.waitForEvent('download'); await page.click(" + strconv.Quote(action.Value1) + ");", nil
	case "assert-download-contains":
		return "// assert downloaded file contains " + strconv.Quote(action.Value1), nil
	case "draw-signature":
		return "// draw signature in " + strconv.Quote(action.Value1), nil
	case "prompt-email-code":
		return "await page.fill(" + strconv.Quote(action.Value2) + ", process.env.SCENARIA_EMAIL_CODE || '');", nil
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
	case "assert-visible":
		return "assert page.locator(" + strconv.Quote(action.Value1) + ").is_visible()", nil
	case "assert-hidden":
		return "assert not page.locator(" + strconv.Quote(action.Value1) + ").is_visible()", nil
	case "assert-text":
		return "assert " + strconv.Quote(action.Value1) + " in page.locator(" + strconv.Quote(action.Value2) + ").inner_text()", nil
	case "assert-url":
		return "assert page.url() == " + strconv.Quote(action.Value1), nil
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
	case "wait-visible":
		return "page.locator(" + strconv.Quote(action.Value1) + ").wait_for(state='visible')", nil
	case "wait-hidden":
		return "page.locator(" + strconv.Quote(action.Value1) + ").wait_for(state='hidden')", nil
	case "reload":
		return "page.reload()", nil
	case "go-back":
		return "page.go_back()", nil
	case "close-browser":
		return "browser.close()", nil
	case "check":
		return "page.check(" + strconv.Quote(action.Value1) + ")", nil
	case "uncheck":
		return "page.uncheck(" + strconv.Quote(action.Value1) + ")", nil
	case "select":
		return "page.select_option(" + strconv.Quote(action.Value2) + ", " + strconv.Quote(action.Value1) + ")", nil
	case "hover":
		return "page.hover(" + strconv.Quote(action.Value1) + ")", nil
	case "double-click":
		return "page.dblclick(" + strconv.Quote(action.Value1) + ")", nil
	case "clear":
		return "page.fill(" + strconv.Quote(action.Value1) + ", '')", nil
	case "scroll-to":
		return "page.locator(" + strconv.Quote(action.Value1) + ").scroll_into_view_if_needed()", nil
	case "remember-text", "remember-field", "remember-url":
		return "# " + action.Kind + " (runtime variable)", nil
	case "fill-generated":
		return "page.fill(" + strconv.Quote(action.Value2) + ", generate_" + action.Value1 + "())", nil
	case "switch-tab", "close-tab", "draw-signature", "assert-download-contains":
		return "# " + action.Kind, nil
	case "assert-tab-count":
		return fmt.Sprintf("assert len(context.pages) == %d", action.IntVal), nil
	case "download-click":
		return "with page.expect_download() as download_info:\n    page.click(" + strconv.Quote(action.Value1) + ")", nil
	case "prompt-email-code":
		return "page.fill(" + strconv.Quote(action.Value2) + ", os.environ.get('SCENARIA_EMAIL_CODE', ''))", nil
	default:
		return "", fmt.Errorf("unsupported action kind for Python export: %s", action.Kind)
	}
}

func durationToMS(raw string) (int, error) {
	d, err := stepdsl.ParseWaitDuration(raw)
	if err != nil {
		return 0, err
	}
	return int(d.Milliseconds()), nil
}
