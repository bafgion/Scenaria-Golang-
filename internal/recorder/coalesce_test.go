package recorder

import (
	"strings"
	"testing"
)

func TestApplyCoalescedStepEmpty(t *testing.T) {
	step := RecordedStep{Action: "click", Selector: "#a"}
	out, emitted := ApplyCoalescedStep(nil, step)
	if len(out) != 1 || emitted == nil || emitted.Selector != "#a" {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestApplyCoalescedStepCoalescesFill(t *testing.T) {
	steps := []RecordedStep{{Action: "fill", Selector: "#email", Value: "a"}}
	next := RecordedStep{Action: "fill", Selector: "#email", Value: "b@x.com"}
	out, emitted := ApplyCoalescedStep(steps, next)
	if len(out) != 1 || out[0].Value != "b@x.com" || emitted == nil {
		t.Fatalf("got %+v", out)
	}
}

func TestApplyCoalescedStepReplacesClickBeforeFill(t *testing.T) {
	fragile := "div > label:nth-of-type(1) > div:nth-of-type(2) > input"
	steps := []RecordedStep{{Action: "click", Selector: fragile}}
	next := RecordedStep{Action: "fill", Selector: fragile, Value: "Ivan", Text: "Имя"}
	out, emitted := ApplyCoalescedStep(steps, next)
	if len(out) != 1 || out[0].Action != "fill" {
		t.Fatalf("got %+v", out)
	}
	if !strings.HasPrefix(out[0].Selector, `label:has-text("Имя`) {
		t.Fatalf("selector: %q", out[0].Selector)
	}
	if emitted == nil || emitted.Value != "Ivan" {
		t.Fatalf("emitted=%+v", emitted)
	}
}

func TestApplyCoalescedStepSkipsTabAfterFill(t *testing.T) {
	steps := []RecordedStep{{Action: "fill", Selector: "#name", Value: "Ivan"}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "press", Selector: "#name", Value: "Tab"})
	if len(out) != 1 || emitted != nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestApplyCoalescedStepReplacesClickWithCheck(t *testing.T) {
	steps := []RecordedStep{{Action: "click", Selector: `label:has-text("Согласен")`}}
	check := RecordedStep{Action: "check", Selector: `label:has-text("Согласен")`}
	out, emitted := ApplyCoalescedStep(steps, check)
	if len(out) != 1 || out[0].Action != "check" || emitted == nil {
		t.Fatalf("got %+v", out)
	}
}

func TestApplyCoalescedStepCheckboxFill(t *testing.T) {
	steps := []RecordedStep{{Action: "click", Selector: "input"}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{
		Action: "fill", Selector: "input", Value: "on", InputType: "checkbox",
	})
	if len(out) != 1 || out[0].Action != "check" || emitted == nil {
		t.Fatalf("got %+v", out)
	}
}

func TestApplyCoalescedStepSkipsDuplicateHover(t *testing.T) {
	steps := []RecordedStep{{Action: "hover", Selector: "nav"}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "hover", Selector: "nav"})
	if len(out) != 1 || emitted != nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}

func TestApplyCoalescedStepSkipsDuplicateClick(t *testing.T) {
	steps := []RecordedStep{{Action: "click", Selector: "#ok"}}
	out, emitted := ApplyCoalescedStep(steps, RecordedStep{Action: "click", Selector: "#ok"})
	if len(out) != 1 || emitted != nil {
		t.Fatalf("got %+v emitted=%v", out, emitted)
	}
}
