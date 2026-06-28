package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func testFeature() *gherkin.Feature {
	return &gherkin.Feature{
		Title: "Login",
		Background: []gherkin.Step{
			{Keyword: "Когда", Text: `открываю "/login"`},
		},
		Scenarios: []gherkin.Scenario{
			{
				Title: "Success",
				Steps: []gherkin.Step{
					{Keyword: "Когда", Text: `ввожу "admin" в "#username"`},
					{Keyword: "И", Text: `ввожу "secret" в "#password"`},
					{Keyword: "И", Text: `нажимаю "#submit"`},
					{Keyword: "Тогда", Text: `вижу ".dashboard"`},
					{Keyword: "И", Text: `проверяю текст "Панель" в ".dashboard"`},
					{Keyword: "И", Text: `url содержит "dashboard"`},
				},
			},
		},
	}
}

func TestWritePlaywrightTS(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "out.spec.ts")
	if err := WritePlaywrightTS(path, testFeature(), "https://example.com"); err != nil {
		t.Fatalf("WritePlaywrightTS returned error: %v", err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read ts export: %v", err)
	}
	text := string(payload)
	mustContain := []string{
		"import { test, expect } from '@playwright/test';",
		`await page.goto("https://example.com/login");`,
		`await page.fill("#username", "admin");`,
		`await expect(page.locator(".dashboard")).toBeVisible();`,
		`await expect(page.locator(".dashboard")).toContainText("Панель");`,
		`expect(page.url()).toContain("dashboard");`,
	}
	for _, expected := range mustContain {
		if !strings.Contains(text, expected) {
			t.Fatalf("ts export missing %q\n%s", expected, text)
		}
	}
}

func TestWritePlaywrightPython(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "out.py")
	if err := WritePlaywrightPython(path, testFeature(), "https://example.com"); err != nil {
		t.Fatalf("WritePlaywrightPython returned error: %v", err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read python export: %v", err)
	}
	text := string(payload)
	mustContain := []string{
		"from playwright.sync_api import sync_playwright",
		`page.goto("https://example.com/login")`,
		`page.fill("#username", "admin")`,
		`assert page.locator(".dashboard").is_visible()`,
		`assert "Панель" in page.locator(".dashboard").inner_text()`,
		`assert "dashboard" in page.url`,
	}
	for _, expected := range mustContain {
		if !strings.Contains(text, expected) {
			t.Fatalf("python export missing %q\n%s", expected, text)
		}
	}
}
