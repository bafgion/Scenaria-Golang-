package gui

import "testing"

func TestUpdateRecordingOptionsRequiresSession(t *testing.T) {
	svc := NewService()
	err := svc.UpdateRecordingOptions(true, false, false)
	if err == nil {
		t.Fatal("expected error when no live session")
	}
}
