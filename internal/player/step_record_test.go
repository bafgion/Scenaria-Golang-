package player

import (
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestStepRecordLifecycle(t *testing.T) {
	ctx := NewRunContext(nil, 1, t.TempDir())
	step := gherkin.Step{Keyword: "Когда", Text: `кликаю "#x"`, Line: 5}
	idx := ctx.beginLeafStep(step)
	if idx != 0 {
		t.Fatalf("expected index 0, got %d", idx)
	}
	ctx.completeLeafStep(idx, "#x", time.Now().Add(-12*time.Millisecond), nil, nil)
	recs := ctx.StepRecords()
	if len(recs) != 1 || recs[0].Status != "passed" || recs[0].Selector != "#x" {
		t.Fatalf("unexpected records: %+v", recs)
	}
	if recs[0].DurationMS < 10 {
		t.Fatalf("expected duration >= 10ms, got %d", recs[0].DurationMS)
	}
}

func TestActionSelector(t *testing.T) {
	if actionSelector(stepdsl.Action{Kind: "click", Value1: "#a"}) != "#a" {
		t.Fatal("click selector")
	}
	if actionSelector(stepdsl.Action{Kind: "fill", Value1: "value", Value2: "#name"}) != "#name" {
		t.Fatal("fill selector")
	}
	if actionSelector(stepdsl.Action{Kind: "assert-text", Value1: "Hello", Value2: "#msg"}) != "#msg" {
		t.Fatal("assert-text selector")
	}
	if actionSelector(stepdsl.Action{Kind: "remember-url", Value1: "u"}) != "u" {
		t.Fatal("goto selector")
	}
}
