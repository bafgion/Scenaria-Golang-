package gui

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
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

	got, _, err := svc.finalizeGUIReports(tmp, RunRequest{HTMLPath: htmlPath}, plan, result, runErr, false)
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

func TestResolvePreviousSummaryPathUsesLatestPointer(t *testing.T) {
	root := t.TempDir()
	latestSummary := filepath.Join(root, ".scenaria", "runs", "run-prev", "summary.json")
	if err := os.MkdirAll(filepath.Dir(latestSummary), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := report.WriteRunSummaryDetailed(latestSummary, report.RunSummaryDetailed{Mode: "browser", Scenarios: 1}); err != nil {
		t.Fatal(err)
	}
	if err := report.WriteLatestRunPointer(root, report.RunArtifactLayout{
		RunID:       "run-prev",
		SummaryJSON: latestSummary,
	}); err != nil {
		t.Fatal(err)
	}
	currentSummary := filepath.Join(root, ".scenaria", "runs", "run-new", "summary.json")
	got := resolvePreviousSummaryPath(root, currentSummary, "run-new")
	if got != latestSummary {
		t.Fatalf("previous summary path = %q, want %q", got, latestSummary)
	}
}
