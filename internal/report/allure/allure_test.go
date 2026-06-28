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

func TestMapStatus(t *testing.T) {
	if mapStatus("passed") != "passed" || mapStatus("failed") != "failed" || mapStatus("skipped") != "skipped" {
		t.Fatal("unexpected status mapping")
	}
}
