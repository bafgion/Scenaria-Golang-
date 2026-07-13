package recorder

import "testing"

func TestStopCapturePreserveBufferKeepsCaptureEver(t *testing.T) {
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
	s.StopCapturePreserveBuffer()
	s.ResetCaptureSegment()
	if s.CaptureEnabled() {
		t.Fatal("expected capture disabled")
	}
	if !s.CaptureEverEnabled() {
		t.Fatal("expected captureEver preserved for resume without replay")
	}
	if len(steps) != 0 {
		t.Fatalf("expected segment buffer cleared, got %d", len(steps))
	}
}

func TestEndCaptureResetsCaptureEver(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{{Action: "click", Selector: "#one"}}
	s.Bind(nil, &steps)
	s.InitBrowseMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	s.EndCapture()
	if s.CaptureEverEnabled() {
		t.Fatal("expected captureEver reset after EndCapture")
	}
	if len(steps) != 0 {
		t.Fatalf("expected steps cleared, got %d", len(steps))
	}
}

func TestShouldSyncAfterStopPreserveBuffer(t *testing.T) {
	s := NewLiveSession()
	s.InitBrowseMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	s.StopCapturePreserveBuffer()
	s.ResetCaptureSegment()
	if ShouldSyncRecordedStepsOnCaptureStart(s) {
		t.Fatal("expected no replay after preserve-buffer stop")
	}
}

func TestCompleteCaptureStopWithoutPage(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{{Action: "click", Selector: "#one"}}
	s.Bind(nil, &steps)
	s.InitBrowseMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	var events []RecordStepEvent
	stopped := CompleteCaptureStop(s, nil, func(event RecordStepEvent) {
		events = append(events, event)
	}, "manual", func(reason string) {
		if reason != "manual" {
			t.Fatalf("unexpected reason %q", reason)
		}
	})
	if !stopped {
		t.Fatal("expected stop")
	}
	if len(events) != 1 || events[0].Op != RecordStepSnapshot {
		t.Fatalf("expected snapshot, got %+v", events)
	}
	if s.CaptureEnabled() {
		t.Fatal("expected capture disabled")
	}
	if !s.CaptureEverEnabled() {
		t.Fatal("expected captureEver preserved")
	}
}

func TestCompleteCaptureStopWithoutNotifierKeepsBufferForPersist(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{{Action: "goto", Value: "https://example.com"}}
	s.Bind(nil, &steps)
	s.InitRecordMode()

	stopped := CompleteCaptureStop(s, nil, nil, "idle", nil)
	if !stopped {
		t.Fatal("expected stop")
	}
	if s.CaptureEnabled() {
		t.Fatal("expected capture disabled")
	}
	if len(steps) != 1 {
		t.Fatalf("expected buffer preserved for persistence, got %d", len(steps))
	}
}

func TestProcessToolbarStopWhenRecording(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{{Action: "click", Selector: "#one"}}
	s.Bind(nil, &steps)
	s.InitRecordMode()

	var stoppedReason string
	var events []RecordStepEvent
	ProcessToolbarAction("stop", s, nil, LiveOptions{
		Callbacks: LiveCallbacks{
			OnCaptureStop: func(reason string) { stoppedReason = reason },
		},
	}, func(event RecordStepEvent) {
		events = append(events, event)
	}, &steps)

	if stoppedReason != "manual" {
		t.Fatalf("expected manual stop, got %q", stoppedReason)
	}
	if s.CaptureEnabled() {
		t.Fatal("expected capture disabled after toolbar stop")
	}
	if len(events) != 1 || events[0].Op != RecordStepSnapshot {
		t.Fatalf("expected snapshot before stop, got %+v", events)
	}
}

func TestProcessToolbarStopDoesNotCloseBrowser(t *testing.T) {
	s := NewLiveSession()
	s.InitRecordMode()
	result := ProcessToolbarAction("stop", s, nil, LiveOptions{}, nil, nil)
	if result.CloseBrowser {
		t.Fatal("toolbar stop must not close browser when idle")
	}
}

func TestProcessToolbarActionNilSessionIsSafe(t *testing.T) {
	result := ProcessToolbarAction("stop", nil, nil, LiveOptions{}, nil, nil)
	if result.CloseBrowser {
		t.Fatal("nil session action must not close browser")
	}
}

func TestProcessToolbarRecordStartsCaptureAndEmitsStart(t *testing.T) {
	s := NewLiveSession()
	s.InitBrowseMode()
	started := false
	ProcessToolbarAction("record", s, nil, LiveOptions{
		Callbacks: LiveCallbacks{
			OnCaptureStart: func(bool) {
				started = true
			},
		},
	}, nil, nil)
	if !started {
		t.Fatal("expected capture start callback")
	}
	if !s.CaptureEnabled() {
		t.Fatal("expected capture enabled")
	}
}

func TestProcessToolbarStopWhilePausedEmitsSnapshotAndStop(t *testing.T) {
	s := NewLiveSession()
	steps := []RecordedStep{{Action: "click", Selector: "#one"}}
	s.Bind(nil, &steps)
	s.InitBrowseMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	s.Pause()
	var stoppedReason string
	var events []RecordStepEvent
	ProcessToolbarAction("stop", s, nil, LiveOptions{
		Callbacks: LiveCallbacks{
			OnCaptureStop: func(reason string) { stoppedReason = reason },
		},
	}, func(event RecordStepEvent) {
		events = append(events, event)
	}, &steps)
	if stoppedReason != "manual" {
		t.Fatalf("expected manual stop, got %q", stoppedReason)
	}
	if s.CaptureEnabled() || s.IsPaused() {
		t.Fatal("expected paused capture to stop fully")
	}
	if len(events) != 1 || events[0].Op != RecordStepSnapshot {
		t.Fatalf("expected snapshot before paused stop, got %+v", events)
	}
}

func TestProcessToolbarPauseResumePickerWhilePaused(t *testing.T) {
	s := NewLiveSession()
	s.InitRecordMode()
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	s.Pause()
	if !s.IsPaused() {
		t.Fatal("expected paused")
	}
	ProcessToolbarAction("resume", s, nil, LiveOptions{}, nil, nil)
	if s.IsPaused() {
		t.Fatal("expected resumed from toolbar action")
	}
	s.Pause()
	pickerCalled := false
	ProcessToolbarAction("picker", s, nil, LiveOptions{
		Callbacks: LiveCallbacks{
			OnPickerRequest: func() { pickerCalled = true },
		},
	}, nil, nil)
	if !pickerCalled {
		t.Fatal("expected picker request while paused")
	}
}
