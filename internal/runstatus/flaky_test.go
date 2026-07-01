package runstatus

import "testing"

func TestFlakyStats(t *testing.T) {
	entries := []Entry{
		{Path: "a.feature::S1", Success: true, At: "2026-01-03T10:00:00Z"},
		{Path: "a.feature::S1", Success: false, At: "2026-01-02T10:00:00Z", FailedStep: intPtr(2)},
		{Path: "a.feature::S1", Success: true, At: "2026-01-01T10:00:00Z"},
		{Path: "b.feature::S2", Success: false, At: "2026-01-01T09:00:00Z", FailedStep: intPtr(1)},
		{Path: "b.feature::S2", Success: false, At: "2026-01-01T08:00:00Z", FailedStep: intPtr(1)},
	}
	scenarios, steps := FlakyStats(entries)
	if len(scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(scenarios))
	}
	var s1 *ScenarioFlakyStat
	for i := range scenarios {
		if scenarios[i].Path == "a.feature::S1" {
			s1 = &scenarios[i]
			break
		}
	}
	if s1 == nil {
		t.Fatal("missing a.feature::S1")
	}
	if !s1.Flaky || s1.Failures != 1 || s1.Passes != 2 {
		t.Fatalf("unexpected s1 stats: %+v", s1)
	}
	if len(steps) != 1 || steps[0].Step != 1 || steps[0].Failures != 2 {
		t.Fatalf("unexpected step stats: %+v", steps)
	}
}

func intPtr(v int) *int { return &v }
