package recorder

import "testing"

func TestLiveSessionHoldForTestRun(t *testing.T) {
	s := NewLiveSession()
	if s.TestRunHeld() {
		t.Fatal("expected no test hold initially")
	}
	release := s.HoldForTestRun()
	if !s.TestRunHeld() {
		t.Fatal("expected hold active")
	}
	release()
	if s.TestRunHeld() {
		t.Fatal("expected hold released")
	}
}
