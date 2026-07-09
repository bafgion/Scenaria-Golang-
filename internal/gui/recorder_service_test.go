package gui

import "testing"

func TestRecorderServiceNilSessionGuards(t *testing.T) {
	svc := NewRecorderService()

	got := svc.PollBrowserSession(nil, "browser-1")
	if got.BrowserOpen || got.Recording || got.Paused || got.StepCount != 0 {
		t.Fatalf("expected zero browser session dto, got %+v", got)
	}

	if err := svc.FocusBrowser(nil); err == nil {
		t.Fatal("expected focus error for nil session")
	}
	if err := svc.UpdateRecordingOptions(nil, true, true, true, false, false, 500, true); err == nil {
		t.Fatal("expected update options error for nil session")
	}

	pick := svc.PickSelector(nil, nil)
	if pick.Error == "" {
		t.Fatal("expected picker error for nil session")
	}

	if _, _, err := svc.BeginCapture(nil); err == nil {
		t.Fatal("expected begin capture error for nil session")
	}
	if _, err := svc.StopCapture(nil); err == nil {
		t.Fatal("expected stop capture error for nil session")
	}
	if _, ok := svc.UndoLastRecordedStep(nil); ok {
		t.Fatal("expected undo false for nil session")
	}
}
