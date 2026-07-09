package player

import (
	"context"
	"errors"
	"testing"
)

func TestRunProgressFromContext(t *testing.T) {
	var seen int
	ctx := WithRunProgress(context.Background(), func(ev RunProgressEvent) {
		if ev.Phase == ProgressScenarioStart {
			seen++
		}
	})
	emitRunProgress(ctx, RunProgressEvent{Phase: ProgressScenarioStart, Index: 1, Total: 2})
	if seen != 1 {
		t.Fatalf("expected progress callback, got %d", seen)
	}
}

func TestFinalizeParallelScenarioUsesExplicitStatusForSuccess(t *testing.T) {
	var events []RunProgressEvent
	ctx := WithRunProgress(context.Background(), func(ev RunProgressEvent) {
		events = append(events, ev)
	})
	results := make([]ScenarioResult, 1)
	runCase := RunCase{CaseID: "case-1", FeaturePath: "a.feature", Name: "A"}
	finalizeParallelScenario(ctx, nil, results, 0, 1, runCase, ScenarioResult{Status: "canceled", Message: "stop"})

	if len(events) != 1 {
		t.Fatalf("expected one progress event, got %d", len(events))
	}
	if events[0].Success {
		t.Fatal("canceled scenario must not be reported as successful")
	}
	if results[0].Status != "canceled" || results[0].CaseID != "case-1" {
		t.Fatalf("unexpected finalized result: %+v", results[0])
	}
}

func TestTerminalScenarioResultForAcquireError(t *testing.T) {
	runCase := RunCase{CaseID: "case-1", FeaturePath: "a.feature", Name: "A"}
	result := terminalScenarioResult(runCase, ScenarioResult{}, errors.New("browser pool is closed"))
	if result.Status != "failed" {
		t.Fatalf("status = %q, want failed", result.Status)
	}
	if result.CaseID != runCase.CaseID || result.FeaturePath != runCase.FeaturePath || result.Scenario != runCase.Name {
		t.Fatalf("result did not preserve case identity: %+v", result)
	}
}

func TestFillCanceledResultsEmitsTerminalProgress(t *testing.T) {
	base, cancel := context.WithCancel(context.Background())
	var events []RunProgressEvent
	ctx := WithRunProgress(base, func(ev RunProgressEvent) {
		events = append(events, ev)
	})
	cancel()

	results := make([]ScenarioResult, 1)
	fillCanceledResults(ctx, results, []RunCase{{CaseID: "case-1", FeaturePath: "a.feature", Name: "A"}})
	if results[0].Status != "canceled" {
		t.Fatalf("status = %q, want canceled", results[0].Status)
	}
	if len(events) != 1 || events[0].Phase != ProgressScenarioDone || events[0].Success {
		t.Fatalf("unexpected progress events: %+v", events)
	}
}
