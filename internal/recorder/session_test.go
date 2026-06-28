package recorder

import "testing"

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
