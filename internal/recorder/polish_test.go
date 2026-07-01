package recorder

import "testing"

func TestPolishIncomingStepSkipsGenericHover(t *testing.T) {
	_, ok := polishIncomingStep(RecordedStep{Action: "hover", Selector: "img"})
	if ok {
		t.Fatal("expected generic hover to be skipped")
	}
}

func TestApplyCoalescedStepDropsDuplicateGoto(t *testing.T) {
	steps := []RecordedStep{{Action: "goto", Value: "https://example.com"}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "goto", Value: "https://example.com"})
	if len(out) != 1 || emitted != nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestApplyCoalescedStepDropsHoverBeforeSameClick(t *testing.T) {
	selector := `nav >> text="Novinki"`
	steps := []RecordedStep{{Action: "hover", Selector: selector}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "click", Selector: selector})
	if len(out) != 1 || out[0].Action != "click" || emitted == nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestApplyCoalescedStepSkipsRepeatedHoverSinceGoto(t *testing.T) {
	steps := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "hover", Selector: "nav.menu"},
		{Action: "click", Selector: "a.item"},
	}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "hover", Selector: "nav.menu"})
	if len(out) != 3 || emitted != nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestNormalizeDropsRedundantHoverBeforeClick(t *testing.T) {
	selector := `text="Novinki"`
	steps := []RecordedStep{
		{Action: "hover", Selector: selector},
		{Action: "click", Selector: selector},
		{Action: "hover", Selector: selector},
		{Action: "hover", Selector: selector},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 || out[0].Action != "click" {
		t.Fatalf("got %+v", out)
	}
}

func TestAppendGotoStepDedupes(t *testing.T) {
	var steps []RecordedStep
	steps = append(steps, RecordedStep{Action: "goto", Value: "https://example.com"})
	appendGotoStep(&steps, "https://example.com", nil)
	if len(steps) != 1 {
		t.Fatalf("got %+v", steps)
	}
}
