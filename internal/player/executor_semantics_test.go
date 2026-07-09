package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestEvaluateConditionRequiresPage(t *testing.T) {
	ctx := NewRunContext(nil, 1, "")
	ok, err := ctx.EvaluateConditionResult(&gherkin.Condition{Type: "visible", Selector: "#missing"})
	if err == nil {
		t.Fatal("expected error when page is not available")
	}
	if ok {
		t.Fatal("expected false on error")
	}
}

func TestWhileLoopPropagatesConditionError(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{MaxLoopIterations: 5})
	steps := []gherkin.Step{{
		Block:     gherkin.BlockWhile,
		Condition: &gherkin.Condition{Type: "visible", Selector: "#missing"},
		Line:      1,
		Children:  []gherkin.Step{{Text: `жду 1 мс`, Line: 2}},
	}}
	err := exec.ExecuteSteps(t.Context(), &browserSession{}, steps, NewRunContext(nil, 1, t.TempDir()))
	if err == nil {
		t.Fatal("expected while condition error")
	}
}
