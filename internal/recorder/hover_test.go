package recorder

import (
	"strings"
	"testing"
)

func TestAppendRecordedStepInsertsHoverBeforeMenuClick(t *testing.T) {
	steps := make([]RecordedStep, 0)
	appendRecordedStep(&steps, RecordedStep{
		Action: "click", Selector: "a.item", HoverSelector: "nav.menu", HoverText: "Menu",
	}, nil)
	if len(steps) != 2 {
		t.Fatalf("got %+v", steps)
	}
	if steps[0].Action != "hover" || steps[0].Selector != "nav.menu" {
		t.Fatalf("hover step: %+v", steps[0])
	}
	if steps[1].Action != "click" || steps[1].Selector != "a.item" {
		t.Fatalf("click step: %+v", steps[1])
	}
}

func TestAppendRecordedStepSkipsDuplicateHover(t *testing.T) {
	steps := []RecordedStep{{Action: "hover", Selector: "nav.menu"}}
	appendRecordedStep(&steps, RecordedStep{
		Action: "click", Selector: "a.item", HoverSelector: "nav.menu",
	}, nil)
	if len(steps) != 2 {
		t.Fatalf("got %+v", steps)
	}
}

func TestRecordedStepsToLinesEmitsHoverBeforeClick(t *testing.T) {
	lines := RecordedStepsToLines([]RecordedStep{{
		Action: "click", Selector: "a.item", HoverSelector: "nav.menu",
	}})
	if len(lines) != 2 {
		t.Fatalf("got %+v", lines)
	}
	if !strings.Contains(lines[0], `навожу "nav.menu"`) {
		t.Fatalf("hover line: %q", lines[0])
	}
	if !strings.Contains(lines[1], `нажимаю "a.item"`) {
		t.Fatalf("click line: %q", lines[1])
	}
}

func TestEventToRecordedStepClickHoverSelector(t *testing.T) {
	step, ok := EventToRecordedStep("click", map[string]string{
		"selector": "a.item", "hoverselector": "nav.menu", "hovertext": "Menu",
	})
	if !ok || step.HoverSelector != "nav.menu" || step.HoverText != "Menu" {
		t.Fatalf("got %+v ok=%v", step, ok)
	}
}
