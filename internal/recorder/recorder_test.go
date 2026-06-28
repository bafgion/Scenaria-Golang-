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
