package report

import (
	"strings"
	"testing"
)

func TestTrimPayloadArtifacts(t *testing.T) {
	big := strings.Repeat("x", 5<<20)
	payload := &htmlReportPayload{
		Scenarios: []htmlScenario{{
			Screenshot: "data:image/jpeg;base64," + big,
			Steps:      []htmlStep{{Screenshot: big}},
		}},
	}
	if !trimPayloadArtifacts(payload, 1<<20) {
		t.Fatal("expected trim")
	}
	if payload.Scenarios[0].Screenshot != "" || payload.Scenarios[0].Steps[0].Screenshot != "" {
		t.Fatal("screenshots not cleared")
	}
	if !payload.ArtifactsTrimmed {
		t.Fatal("flag not set")
	}
}

func TestTrimPayloadClearsTraceEvents(t *testing.T) {
	big := strings.Repeat("x", 1<<20)
	events := make([]htmlTraceEvent, 8)
	for i := range events {
		events[i] = htmlTraceEvent{OffsetMS: int64(i), Title: big}
	}
	payload := &htmlReportPayload{
		Scenarios: []htmlScenario{{TraceEvents: events}},
	}
	if !trimPayloadArtifacts(payload, 1<<20) {
		t.Fatal("expected trim")
	}
	if len(payload.Scenarios[0].TraceEvents) != 0 {
		t.Fatal("trace events should be cleared")
	}
}

func TestComputeRegressionsNewFailure(t *testing.T) {
	sc := htmlScenario{
		Status:      "failed",
		FailedStep:  intPtr(1),
		HistoryRuns: []htmlHistoryEntry{{Status: "passed", At: "t0"}},
		Steps: []htmlStep{
			{Index: 0, Status: "passed"},
			{Index: 1, Status: "failed", Text: "click"},
		},
	}
	regs := computeRegressions(sc)
	if len(regs) != 1 || regs[0].Kind != "new_failure" {
		t.Fatalf("unexpected: %+v", regs)
	}
}

func intPtr(v int) *int { return &v }
