package player

import "testing"

func TestRunContextFailedLeafStep(t *testing.T) {
	ctx := NewRunContext(nil, 1, t.TempDir())
	if ctx.FailedLeafStep() != noFailedLeafStep {
		t.Fatalf("expected unset failed step, got %d", ctx.FailedLeafStep())
	}
	ctx.markFailedLeafStep()
	if ctx.FailedLeafStep() != 0 {
		t.Fatalf("expected failed step 0, got %d", ctx.FailedLeafStep())
	}
}
