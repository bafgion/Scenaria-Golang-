package recorder

import "testing"

func TestLiveSessionBrowseThenCapture(t *testing.T) {
	s := NewLiveSession()
	s.InitBrowseMode()
	if s.CaptureEnabled() {
		t.Fatal("browse mode should not capture")
	}
	if s.CaptureEverEnabled() {
		t.Fatal("browse mode should not set captureEver")
	}
	if err := s.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	if !s.CaptureEnabled() {
		t.Fatal("expected capture enabled after BeginCapture")
	}
}

func TestLiveSessionPauseResume(t *testing.T) {
	s := NewLiveSession()
	if s.IsPaused() {
		t.Fatal("expected running by default")
	}
	s.Pause()
	if !s.IsPaused() {
		t.Fatal("expected paused")
	}
	s.Resume()
	if s.IsPaused() {
		t.Fatal("expected resumed")
	}
}

func TestLiveSessionRecordMode(t *testing.T) {
	s := NewLiveSession()
	s.InitRecordMode()
	if !s.CaptureEnabled() {
		t.Fatal("record mode should capture immediately")
	}
}
