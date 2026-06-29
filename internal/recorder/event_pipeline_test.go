package recorder

import (
	"strings"
	"testing"
)

func TestDecodeEventsFromPlaywrightPayload(t *testing.T) {
	raw := []map[string]any{
		{
			"type": "click",
			"detail": map[string]string{
				"selector": "#submit",
				"text":     "Войти",
			},
			"ts": int64(42),
		},
	}
	events, err := decodeEvents(raw)
	if err != nil {
		t.Fatalf("decodeEvents: %v", err)
	}
	if len(events) != 1 || events[0].Type != "click" {
		t.Fatalf("got %+v", events)
	}
	if events[0].Detail["selector"] != "#submit" {
		t.Fatalf("detail: %+v", events[0].Detail)
	}
}

func TestEventPipelineDeterministicOrder(t *testing.T) {
	recorded := []RecordedStep{{Action: "goto", Value: "https://example.com"}}
	batch := []recorderEvent{
		{Type: "input", Detail: map[string]string{"selector": "#email", "value": "a"}},
		{Type: "input", Detail: map[string]string{"selector": "#email", "value": "user@example.com"}},
		{Type: "click", Detail: map[string]string{"selector": "#submit", "text": "Войти"}},
	}
	for _, event := range batch {
		detail := normalizeDetail(event.Detail)
		step, ok := EventToRecordedStep(event.Type, detail)
		if !ok {
			t.Fatalf("unexpected event %+v", event)
		}
		appendCoalescedStep(&recorded, step)
	}
	if len(recorded) != 3 {
		t.Fatalf("expected 3 steps, got %+v", recorded)
	}
	if recorded[1].Action != "fill" || recorded[1].Value != "user@example.com" {
		t.Fatalf("coalesced fill: %+v", recorded[1])
	}
	if recorded[2].Action != "click" || recorded[2].Selector != "#submit" {
		t.Fatalf("click step: %+v", recorded[2])
	}

	lines := RecordedStepsToLines(NormalizeSteps(recorded))
	if len(lines) != 3 {
		t.Fatalf("lines: %v", lines)
	}
	if !strings.Contains(lines[1], "user@example.com") {
		t.Fatalf("fill line: %q", lines[1])
	}
}

func TestAppendCoalescedStepAddsHoverBeforeMenuClick(t *testing.T) {
	var recorded []RecordedStep
	appendCoalescedStep(&recorded, RecordedStep{
		Action:        "click",
		Selector:      "a.item",
		HoverSelector: "nav.menu",
		HoverText:     "Каталог",
	})
	if len(recorded) != 2 {
		t.Fatalf("expected hover+click, got %+v", recorded)
	}
	if recorded[0].Action != "hover" || recorded[0].Selector != "nav.menu" {
		t.Fatalf("hover: %+v", recorded[0])
	}
	if recorded[1].Action != "click" || recorded[1].Selector != "a.item" {
		t.Fatalf("click: %+v", recorded[1])
	}
}

func TestAppendCoalescedStepSkipsRedundantHoverBeforeClick(t *testing.T) {
	var recorded []RecordedStep
	appendCoalescedStep(&recorded, RecordedStep{Action: "hover", Selector: "nav.menu"})
	appendCoalescedStep(&recorded, RecordedStep{
		Action:        "click",
		Selector:      "a.item",
		HoverSelector: "nav.menu",
	})
	if len(recorded) != 2 {
		t.Fatalf("expected single hover before click, got %+v", recorded)
	}
}
