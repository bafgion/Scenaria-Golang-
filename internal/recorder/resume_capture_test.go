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
		t.Fatal("after stop the next capture is a fresh segment")
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

	var notifies []int
	notify := func(index int, _ string) {
		notifies = append(notifies, index)
	}

	syncSteps := func() {
		if ShouldSyncRecordedStepsOnCaptureStart(s) {
			for i, st := range recorded {
				if line, ok := RecordedStepToLine(st); ok {
					notify(i, line)
				}
			}
		}
		_ = s.BeginCapture()
	}

	syncSteps()
	if len(notifies) != 3 {
		t.Fatalf("first capture: expected 3 notifies, got %v", notifies)
	}

	s.EndCapture()
	notifies = nil
	syncSteps()
	if len(notifies) != 0 {
		t.Fatalf("after stop buffer is empty; expected no replay, got %v", notifies)
	}
	if len(recorded) != 0 {
		t.Fatalf("expected recorded cleared after stop, got %d", len(recorded))
	}
}
