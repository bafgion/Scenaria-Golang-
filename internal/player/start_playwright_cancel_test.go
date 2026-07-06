package player

import (
	"context"
	"testing"
)

func TestStartPlaywrightCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, stop, err := startPlaywright(ctx)
	if err == nil {
		if stop != nil {
			stop()
		}
		t.Fatal("expected context error")
	}
}
