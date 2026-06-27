package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestRunRecord(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "recorded.feature")

	err := RunRecord([]string{
		"--output", out,
		"--feature", "Логин",
		"--scenario", "Успех",
		"--step", `открываю "https://example.com"`,
		"--step", `нажимаю "#login"`,
	})
	if err != nil {
		t.Fatalf("RunRecord returned error: %v", err)
	}

	feature, err := gherkin.ParseFeatureFile(out)
	if err != nil {
		t.Fatalf("failed to parse recorded feature: %v", err)
	}
	if feature.Title != "Логин" {
		t.Fatalf("unexpected recorded feature title: %q", feature.Title)
	}
	if len(feature.Scenarios) != 1 || len(feature.Scenarios[0].Steps) != 2 {
		t.Fatalf("unexpected recorded scenarios: %+v", feature.Scenarios)
	}
}

func TestParseRecordOptionsErrors(t *testing.T) {
	if _, err := parseRecordOptions(nil); err == nil {
		t.Fatal("expected usage error")
	}
	if _, err := parseRecordOptions([]string{"--output"}); err == nil {
		t.Fatal("expected missing output value error")
	}
	if _, err := parseRecordOptions([]string{"--output", "a.feature", "--step"}); err == nil {
		t.Fatal("expected missing step value error")
	}
}

func TestRunRecordDefaultStep(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "recorded.feature")
	if err := RunRecord([]string{"--output", out}); err != nil {
		t.Fatalf("RunRecord returned error: %v", err)
	}
	payload, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read recorded feature: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("recorded feature file must not be empty")
	}
}
