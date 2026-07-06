package gui

import (
	"testing"
	"time"
)

func TestDefaultRunTimeoutIsReasonable(t *testing.T) {
	if DefaultRunTimeout < 5*time.Minute {
		t.Fatalf("run timeout too short: %v", DefaultRunTimeout)
	}
	if DefaultRunTimeout > 2*time.Hour {
		t.Fatalf("run timeout too long: %v", DefaultRunTimeout)
	}
}
