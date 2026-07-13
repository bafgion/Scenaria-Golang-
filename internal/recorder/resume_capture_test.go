package recorder

import "testing"

func TestShouldSyncRecordedStepsOnCaptureStart(t *testing.T) {
	s := NewLiveSession()
	s.InitBrowseMode()
	if !ShouldSyncRecordedStepsOnCaptureStart(s) {
		t.Fatal("browse mode should sync steps on first capture")
	}
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	s.EndCapture()
	if !ShouldSyncRecordedStepsOnCaptureStart(s) {
		t.Fatal("after full EndCapture the next capture is a fresh segment")
	}
}

func TestEndCaptureClearsRecordedSteps(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: "#one"},
	}
	s.Bind(nil, &steps)
	s.InitBrowseMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	if len(steps) == 0 {
		t.Fatal("expected steps during capture")
	}
	s.EndCapture()
	if len(steps) != 0 {
		t.Fatalf("expected steps cleared after stop, got %d", len(steps))
	}
	if s.RecordedStepCount() != 0 {
		t.Fatalf("expected step count 0 after stop, got %d", s.RecordedStepCount())
	}
}

func TestSyncRecordedStepsOnCaptureStartCallbackCount(t *testing.T) {
	s := NewLiveSession()
	s.InitBrowseMode()
	recorded := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: "#one"},
		{Action: "click", Selector: "#two"},
	}
	s.Bind(nil, &recorded)

	var events []RecordStepEvent
	notify := func(event RecordStepEvent) {
		events = append(events, event)
	}

	syncSteps := func() {
		if ShouldSyncRecordedStepsOnCaptureStart(s) {
			notifySnapshot(notify, recorded)
		}
		_ = s.BeginCapture()
	}

	syncSteps()
	if len(events) != 1 || events[0].Op != RecordStepSnapshot || len(events[0].Lines) != 3 {
		t.Fatalf("first capture: expected snapshot with 3 lines, got %+v", events)
	}

	s.StopCapturePreserveBuffer()
	s.ResetCaptureSegment()
	events = nil
	syncSteps()
	if len(events) != 0 {
		t.Fatalf("after preserve-buffer stop expected no replay, got %+v", events)
	}
	if len(recorded) != 0 {
		t.Fatalf("expected recorded cleared after segment reset, got %d", len(recorded))
	}

	s.EndCapture()
	events = nil
	recorded = []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: "#one"},
	}
	syncSteps()
	if len(events) != 1 || events[0].Op != RecordStepSnapshot {
		t.Fatalf("after full EndCapture expected replay snapshot, got %+v", events)
	}
}
