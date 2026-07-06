package report

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestComputeRunDiff(t *testing.T) {
	fs := 1
	sc := htmlScenario{
		Status:     "failed",
		FailedStep: &fs,
		HistoryRuns: []htmlHistoryEntry{
			{Status: "passed", At: "2026-06-02T10:00:00Z"},
			{Status: "failed", At: "2026-06-01T10:00:00Z", FailedStep: func() *int { v := 0; return &v }()},
		},
		Steps: []htmlStep{
			{Index: 0, Text: "open", Status: "passed"},
			{Index: 1, Text: "click", Status: "failed"},
		},
	}
	diff := computeRunDiff(sc)
	if diff == nil || len(diff.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %+v", diff)
	}
	if len(diff.Rows) != 2 {
		t.Fatalf("expected 2 rows")
	}
	if diff.Rows[1].Cells[0] != "skipped" || diff.Rows[1].Cells[1] != "passed" || diff.Rows[1].Cells[2] != "failed" {
		t.Fatalf("unexpected cells: %+v", diff.Rows[1].Cells)
	}
}

func TestBuildStepDurationSparkline(t *testing.T) {
	sc := htmlScenario{
		HistoryRuns: []htmlHistoryEntry{
			{StepDurations: []int{180, 90}},
			{StepDurations: []int{160, 70}},
		},
		Steps: []htmlStep{{Index: 0, DurationMS: 200}},
	}
	spark := buildStepDurationSparkline(sc, 0, 200)
	if len(spark) != 3 || spark[0] != 160 || spark[2] != 200 {
		t.Fatalf("spark: %v", spark)
	}
}

func TestStepDurationsFromRecordsPlayer(t *testing.T) {
	durs := player.RunstatusEntry(player.ScenarioResult{
		FeaturePath: "a.feature", Scenario: "S", Status: "passed",
		StepRecords: []player.StepRecord{
			{Index: 0, DurationMS: 10},
			{Index: 1, DurationMS: 20},
		},
	}, "playwright").StepDurations
	if len(durs) != 2 || durs[1] != 20 {
		t.Fatalf("durs: %v", durs)
	}
}
