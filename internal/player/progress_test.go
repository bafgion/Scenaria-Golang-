package player

import (
	"context"
	"testing"
)

func TestRunProgressFromContext(t *testing.T) {
	var seen int
	ctx := WithRunProgress(context.Background(), func(ev RunProgressEvent) {
		if ev.Phase == ProgressScenarioStart {
			seen++
		}
	})
	emitRunProgress(ctx, RunProgressEvent{Phase: ProgressScenarioStart, Index: 1, Total: 2})
	if seen != 1 {
		t.Fatalf("expected progress callback, got %d", seen)
	}
}
