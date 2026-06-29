package recorder

import "testing"

func TestEventsToStepClick(t *testing.T) {
	step, ok := EventsToStep("click", map[string]string{"tag": "BUTTON", "id": "save"})
	if !ok {
		t.Fatal("expected click step")
	}
	if step != `нажимаю "#save"` {
		t.Fatalf("unexpected step: %q", step)
	}
}

func TestEventsToStepSelect(t *testing.T) {
	step, ok := EventsToStep("change", map[string]string{"tag": "SELECT", "id": "lang", "value": "ru"})
	if !ok {
		t.Fatal("expected select step")
	}
	if step != `выбираю "ru" в "#lang"` {
		t.Fatalf("unexpected step: %q", step)
	}
}

func TestEventsToStepDrawSignature(t *testing.T) {
	step, ok := EventsToStep("draw-signature", map[string]string{"tag": "CANVAS", "id": "sign"})
	if !ok {
		t.Fatal("expected draw-signature step")
	}
	if step != `рисую подпись в "#sign"` {
		t.Fatalf("unexpected step: %q", step)
	}
}

func TestEventToRecordedStepDrawSignature(t *testing.T) {
	rec, ok := EventToRecordedStep("draw-signature", map[string]string{"selector": "#canvas-sign"})
	if !ok {
		t.Fatal("expected recorded step")
	}
	line, ok := RecordedStepToLine(rec)
	if !ok || line != `рисую подпись в "#canvas-sign"` {
		t.Fatalf("unexpected line: %q ok=%v", line, ok)
	}
}
