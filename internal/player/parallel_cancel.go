package player

import "context"

// failFastParallelCancel stops sibling scenarios when the suite should not continue after a failure.
func failFastParallelCancel(runCtx context.Context, cancel context.CancelFunc, pool *browserPool) {
	if cancel == nil || ContinueOnFail(runCtx) {
		return
	}
	cancel()
	if pool != nil {
		pool.abortActiveSessions()
	}
}
