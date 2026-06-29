package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteJUnitGolden(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "junit.xml")

	result := player.ExecutionResult{
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "login.feature", Scenario: "Успешный вход", Status: "passed"},
			{FeaturePath: "login.feature", Scenario: "Неверный пароль", Status: "failed", Message: "текст не совпал"},
			{FeaturePath: "smoke.feature", Scenario: "Dry-run", Status: "skipped", Message: "dry-run mode"},
		},
	}
	if err := WriteJUnit(path, result); err != nil {
		t.Fatalf("WriteJUnit: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	gotText := strings.ReplaceAll(string(got), "\r\n", "\n")
	wantPath := filepath.Join("testdata", "junit_golden.xml")
	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read golden %q: %v", wantPath, err)
	}
	wantText := strings.ReplaceAll(string(want), "\r\n", "\n")

	if strings.TrimSpace(gotText) != strings.TrimSpace(wantText) {
		t.Fatalf("junit output differs from golden %q\n--- got ---\n%s\n--- want ---\n%s", wantPath, gotText, wantText)
	}
}
