package recorder

import "testing"

func TestEventPipelineClickAndFill(t *testing.T) {
	click, ok := EventToRecordedStep("click", map[string]string{
		"selector": `[data-testid="submit"]`,
		"text":     "Войти",
	})
	if !ok || click.Action != "click" {
		t.Fatalf("unexpected click step: %+v ok=%v", click, ok)
	}

	fill, ok := EventToRecordedStep("input", map[string]string{
		"selector": "#email",
		"value":    "user@example.com",
	})
	if !ok || fill.Action != "fill" {
		t.Fatalf("unexpected fill step: %+v ok=%v", fill, ok)
	}

	lines := RecordedStepsToLines(NormalizeSteps([]RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		click,
		fill,
	}))
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %v", lines)
	}
	if lines[0] != `открыт "https://example.com"` {
		t.Fatalf("unexpected goto line: %q", lines[0])
	}
}

func TestEventPipelineSignature(t *testing.T) {
	step, ok := EventToRecordedStep("draw-signature", map[string]string{
		"selector": "canvas",
	})
	if !ok || step.Action != "draw-signature" {
		t.Fatalf("unexpected signature step: %+v", step)
	}
	line, ok := RecordedStepToLine(step)
	if !ok || line != `рисую подпись в "canvas"` {
		t.Fatalf("unexpected line: %q ok=%v", line, ok)
	}
}
