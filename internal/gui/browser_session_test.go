package gui

import "testing"

func TestPollBrowserSessionEmpty(t *testing.T) {
	svc := &Service{}
	dto := svc.PollBrowserSession()
	if dto.BrowserOpen || dto.Recording || dto.Paused || dto.StepCount != 0 {
		t.Fatalf("expected empty session, got %+v", dto)
	}
}
