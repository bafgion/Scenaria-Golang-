package report

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestBuildCICompare(t *testing.T) {
	prev := &RunSummaryDetailed{
		GeneratedAt: "2026-01-01T00:00:00Z",
		Items: []ScenarioSummary{
			{Path: "a.feature", Scenario: "OK", Status: "passed", DurationMS: 100},
			{Path: "b.feature", Scenario: "Fail", Status: "failed", DurationMS: 500},
		},
	}
	scenarios := []htmlScenario{
		{FeaturePath: "a.feature", Scenario: "OK", Status: "failed", DurationMS: 120},
		{FeaturePath: "b.feature", Scenario: "Fail", Status: "passed", DurationMS: 400},
		{FeaturePath: "c.feature", Scenario: "Slow", Status: "passed", DurationMS: 2000},
	}
	ci := buildCICompare(scenarios, prev)
	if ci == nil {
		t.Fatal("expected ci compare")
	}
	if len(ci.NewFailures) != 1 || ci.NewFailures[0].Scenario != "OK" {
		t.Fatalf("new failures: %+v", ci.NewFailures)
	}
	if len(ci.Fixed) != 1 || ci.Fixed[0].Scenario != "Fail" {
		t.Fatalf("fixed: %+v", ci.Fixed)
	}
}

func TestBuildDurationSparkline(t *testing.T) {
	sc := htmlScenario{
		DurationMS: 300,
		HistoryRuns: []htmlHistoryEntry{
			{DurationMS: 200},
			{DurationMS: 250},
		},
	}
	spark := buildDurationSparkline(sc)
	if len(spark) != 3 || spark[0] != 250 || spark[2] != 300 {
		t.Fatalf("sparkline: %v", spark)
	}
}

func TestBuildHTMLStepsTraceOffset(t *testing.T) {
	sr := player.ScenarioResult{
		StepRecords: []player.StepRecord{
			{Index: 0, Status: "passed", DurationMS: 100},
			{Index: 1, Status: "failed", DurationMS: 50, PageContext: "Title\nhttps://x"},
		},
	}
	steps := buildHTMLSteps(sr, nil, nil, true, true, "", "test")
	if steps[1].TraceOffsetMS != 100 {
		t.Fatalf("offset want 100 got %d", steps[1].TraceOffsetMS)
	}
	if steps[1].PageContext == "" {
		t.Fatal("expected page context")
	}
}
