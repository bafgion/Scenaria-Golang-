package recorder

import "testing"

func TestNormalizeStepsCoalescesFill(t *testing.T) {
	steps := []RecordedStep{
		{Action: "fill", Selector: "#email", Value: "a"},
		{Action: "fill", Selector: "#email", Value: "ab@x.com"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
	if out[0].Value != "ab@x.com" {
		t.Fatalf("unexpected value: %q", out[0].Value)
	}
}

func TestNormalizeStepsCoalescesSelect(t *testing.T) {
	steps := []RecordedStep{
		{Action: "select", Selector: "#lang", Value: "en"},
		{Action: "select", Selector: "#lang", Value: "ru"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
	if out[0].Value != "ru" {
		t.Fatalf("unexpected value: %q", out[0].Value)
	}
}

func TestNormalizeStepsDropsDuplicateGoto(t *testing.T) {
	steps := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "goto", Value: "https://example.com"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
}

func TestNormalizeStepsDropsDuplicateClick(t *testing.T) {
	steps := []RecordedStep{
		{Action: "click", Selector: `button:has-text("OK")`},
		{Action: "click", Selector: `button:has-text("OK")`},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
}
