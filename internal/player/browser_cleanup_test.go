package player

import (
	"context"
	"testing"
)

func TestStartPlaywrightStopsOnCancel(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	ctx, cancel := context.WithCancel(context.Background())
	_, stop, err := startPlaywright(ctx)
	if err != nil {
		t.Fatalf("startPlaywright: %v", err)
	}
	cancel()
	stop()
}
