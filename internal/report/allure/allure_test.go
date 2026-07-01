package allure

import (
	"os"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWritePlaceholder(t *testing.T) {
	dir := t.TempDir()
	if err := WritePlaceholder(Options{OutputDir: dir}); err != nil {
		t.Fatalf("WritePlaceholder: %v", err)
	}
}

func TestWriteResults(t *testing.T) {
	dir := t.TempDir()
	result := player.ExecutionResult{
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "login.feature", Scenario: "Успешный вход", Status: "passed"},
			{
				FeaturePath:   "login.feature",
				Scenario:      "Неверный пароль",
				Status:        "failed",
				Message:       "assertion failed",
				ScreenshotPNG: []byte{0x89, 0x50, 0x4e, 0x47},
			},
		},
	}
	if err := WriteResults(dir, result); err != nil {
		t.Fatalf("WriteResults: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 2 results + 1 attachment, got %d files", len(entries))
	}
}

func TestWriteResultsWithTrace(t *testing.T) {
	dir := t.TempDir()
	result := player.ExecutionResult{
		ScenarioResults: []player.ScenarioResult{
			{
				FeaturePath: "x.feature",
				Scenario:    "fail",
				Status:      "failed",
				Message:     "boom",
				TraceZIP:    []byte("PK\x03\x04"),
			},
		},
	}
	if err := WriteResults(dir, result); err != nil {
		t.Fatalf("WriteResults: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 2 {
		t.Fatalf("expected result + trace attachment, got %d", len(entries))
	}
}

func TestWriteResultsUniqueTimestamps(t *testing.T) {
	dir := t.TempDir()
	result := player.ExecutionResult{
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "a.feature", Scenario: "One", Status: "passed"},
			{FeaturePath: "a.feature", Scenario: "Two", Status: "passed"},
		},
	}
	if err := WriteResults(dir, result); err != nil {
		t.Fatalf("WriteResults: %v", err)
	}
	// Re-run should replace prior results, not accumulate stale files.
	if err := WriteResults(dir, result); err != nil {
		t.Fatalf("WriteResults second run: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 result files after clean, got %d", len(entries))
	}
}

func TestMapStatus(t *testing.T) {
	if mapStatus("passed") != "passed" || mapStatus("failed") != "failed" || mapStatus("skipped") != "skipped" {
		t.Fatal("unexpected status mapping")
	}
}
