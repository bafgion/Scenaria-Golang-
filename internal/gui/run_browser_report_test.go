package gui

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestFinalizeGUIReportsOnFailedRun(t *testing.T) {
	tmp := t.TempDir()
	scenaria := filepath.Join(tmp, ".scenaria")
	if err := os.MkdirAll(scenaria, 0o755); err != nil {
		t.Fatal(err)
	}
	htmlPath := filepath.Join(scenaria, "report.html")

	svc := NewService()
	if _, err := svc.OpenProject(tmp); err != nil {
		t.Fatal(err)
	}

	plan := player.ExecutionPlan{
		Cases: []player.RunCase{{FeaturePath: "a.feature", Name: "Fail"}},
	}
	result := player.ExecutionResult{
		Mode: "browser",
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "a.feature",
			Scenario:    "Fail",
			Status:      "failed",
			Message:     "assertion failed",
		}},
	}
	runErr := errors.New(`scenario "Fail" failed: assertion failed`)

	got, err := svc.finalizeGUIReports(tmp, RunRequest{HTMLPath: htmlPath}, plan, result, runErr, false)
	if err == nil {
		t.Fatal("expected run error after failed scenario")
	}
	if len(got.ScenarioResults) != 1 {
		t.Fatalf("unexpected result: %+v", got)
	}
	if _, statErr := os.Stat(htmlPath); statErr != nil {
		t.Fatalf("HTML report not written on failure: %v", statErr)
	}
}
