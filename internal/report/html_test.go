package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteHTML(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/report.html"
	plan := player.ExecutionPlan{
		Cases: []player.RunCase{{
			FeaturePath: "demo.feature",
			Name:        "S1",
			Tags:        []string{"@smoke"},
			Steps: []gherkin.Step{{
				Keyword: "Когда",
				Text:    `открыт "https://example.com"`,
				Line:    3,
			}},
		}},
	}
	result := player.ExecutionResult{
		Mode:      "dry-run",
		Files:     1,
		Scenarios: 1,
		Steps:     1,
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "demo.feature", Scenario: "S1", Status: "dry-run", Message: "dry-run"},
		},
	}
	if err := WriteHTML(path, result, HTMLOptions{Plan: plan}); err != nil {
		t.Fatalf("WriteHTML failed: %v", err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read html: %v", err)
	}
	text := string(payload)
	for _, want := range []string{
		"scenaria-report-data",
		"demo.feature",
		"Timeline",
		"Inspector",
		`example.com`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in html", want)
		}
	}
}

func TestBuildHTMLStepsFromRecords(t *testing.T) {
	sr := player.ScenarioResult{
		Status: "failed",
		StepRecords: []player.StepRecord{{
			Index: 0, Text: "click", Selector: "#btn", Status: "passed", DurationMS: 10,
		}, {
			Index: 1, Text: "assert", Status: "failed", Error: "boom", DurationMS: 5,
		}},
	}
	steps := buildHTMLSteps(sr, nil, nil, false, true, "", "test")
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if len(steps[0].Tips) == 0 {
		t.Fatal("expected selector tip for non-testid selector")
	}
}

func TestWriteHTMLStoresScreenshotsAsFiles(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	if err := WriteHTML(path, player.ExecutionResult{
		Mode: "browser",
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath:   "a.feature",
			Scenario:      "S",
			Status:        "failed",
			ScreenshotPNG: png,
			StepRecords: []player.StepRecord{{
				Index: 0, Text: "click", Status: "failed", ScreenshotPNG: png,
			}},
		}},
	}, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if strings.Contains(text, "data:image/") {
		t.Fatal("screenshot was embedded as data URL")
	}
	if !strings.Contains(text, "screenshots/000_a__S__scenario.png") {
		t.Fatal("scenario screenshot path not embedded")
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", "000_a__S__scenario.png")); err != nil {
		t.Fatalf("scenario screenshot not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", "000_a__S__step_000.png")); err != nil {
		t.Fatalf("step screenshot not written: %v", err)
	}
}

func TestWriteHTMLRefusesUnmarkedArtifactDir(t *testing.T) {
	tmp := t.TempDir()
	screenshots := filepath.Join(tmp, "screenshots")
	if err := os.MkdirAll(screenshots, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(screenshots, "foreign.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := WriteHTML(filepath.Join(tmp, "report.html"), player.ExecutionResult{}, HTMLOptions{})
	if err == nil || !strings.Contains(err.Error(), "unmarked html artifact directory") {
		t.Fatalf("expected unmarked artifact dir error, got %v", err)
	}
}

func TestWriteHTMLCleansMarkedArtifacts(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	result := player.ExecutionResult{ScenarioResults: []player.ScenarioResult{{
		FeaturePath:   "a.feature",
		Scenario:      "S",
		Status:        "failed",
		ScreenshotPNG: png,
	}}}
	if err := WriteHTML(path, result, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(tmp, "screenshots", "stale.png")
	if err := os.WriteFile(stale, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteHTML(path, result, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale artifact should be removed, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", htmlArtifactMarker)); err != nil {
		t.Fatalf("marker missing: %v", err)
	}
}

func TestWriteHTMLWithBridgeURL(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/report.html"
	if err := WriteHTML(path, player.ExecutionResult{
		Mode: "dry-run",
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "a.feature", Scenario: "S", Status: "dry-run"},
		},
	}, HTMLOptions{BridgeURL: "http://127.0.0.1:19999"}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "http://127.0.0.1:19999") {
		t.Fatal("bridge url not embedded")
	}
}

func TestInferStepStatus(t *testing.T) {
	fs := 1
	if inferStepStatus("failed", 0, &fs) != "passed" {
		t.Fatal("step 0 should be passed")
	}
	if inferStepStatus("failed", 1, &fs) != "failed" {
		t.Fatal("step 1 should be failed")
	}
	if inferStepStatus("failed", 2, &fs) != "skipped" {
		t.Fatal("step 2 should be skipped")
	}
}
