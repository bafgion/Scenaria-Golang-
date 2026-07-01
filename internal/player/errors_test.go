package player

import (
	"errors"
	"testing"
)

func TestExecutionFailurePartialResult(t *testing.T) {
	partial := ExecutionResult{Scenarios: 2, ScenarioResults: []ScenarioResult{{Scenario: "A"}}}
	err := NewExecutionFailure(errors.New("boom"), partial)
	got, ok := PartialExecutionResult(err)
	if !ok || got.Scenarios != 2 {
		t.Fatalf("PartialExecutionResult = %+v %v", got, ok)
	}
}

func TestExecutionFailureWithoutResults(t *testing.T) {
	err := executionFailure(errors.New("boom"), ExecutionResult{})
	if _, ok := PartialExecutionResult(err); ok {
		t.Fatal("expected plain error without partial results")
	}
}
