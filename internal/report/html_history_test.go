package report

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

func TestLookupHistoryRuns(t *testing.T) {
	fs := 1
	entries := []runstatus.Entry{
		{Path: "a.feature::S1", Success: false, At: "2026-06-01T10:00:00Z", Message: "boom", FailedStep: &fs},
		{Path: "a.feature::S1", Success: true, At: "2026-05-31T10:00:00Z"},
		{Path: "b.feature::S2", Success: true, At: "2026-05-30T10:00:00Z"},
	}
	runs := lookupHistoryRuns(entries, player.ScenarioResult{FeaturePath: "a.feature", Scenario: "S1"}, 3)
	if len(runs) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(runs))
	}
	if runs[0].Status != "failed" || runs[1].Status != "passed" {
		t.Fatalf("unexpected order: %+v", runs)
	}
}

func TestLookupHistoryRunsLimit(t *testing.T) {
	entries := make([]runstatus.Entry, 0, 6)
	for i := 0; i < 6; i++ {
		entries = append(entries, runstatus.Entry{
			Path: "x.feature::Y", Success: true, At: "2026-06-01T10:00:00Z",
		})
	}
	runs := lookupHistoryRuns(entries, player.ScenarioResult{FeaturePath: "x.feature", Scenario: "Y"}, 5)
	if len(runs) != 5 {
		t.Fatalf("expected limit 5, got %d", len(runs))
	}
}
