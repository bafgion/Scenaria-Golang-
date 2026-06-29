package recorder

import "testing"

func TestNormalizeStepsSkipsClickAfterCheck(t *testing.T) {
	steps := []RecordedStep{
		{Action: "check", Selector: "#agree"},
		{Action: "click", Selector: "#agree"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 || out[0].Action != "check" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestNormalizeStepsDedupesScroll(t *testing.T) {
	steps := []RecordedStep{
		{Action: "scroll-to", Selector: "#footer"},
		{Action: "scroll-to", Selector: "#footer"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 scroll, got %d", len(out))
	}
}

func TestEventToRecordedStepExtendedActions(t *testing.T) {
	upload, ok := EventToRecordedStep("upload", map[string]string{"selector": "#file", "value": "data.csv"})
	if !ok || upload.Action != "upload" {
		t.Fatalf("upload: %+v ok=%v", upload, ok)
	}
	line, ok := RecordedStepToLine(upload)
	if !ok || line != `загружаю файл "data.csv" в "#file"` {
		t.Fatalf("upload line: %q", line)
	}

	drag, ok := EventToRecordedStep("drag-drop", map[string]string{"selector": "#a", "target": "#b"})
	if !ok || drag.Action != "drag-drop" {
		t.Fatalf("drag: %+v", drag)
	}
	line, ok = RecordedStepToLine(drag)
	if !ok || line != `перетаскиваю "#a" в "#b"` {
		t.Fatalf("drag line: %q", line)
	}
}
