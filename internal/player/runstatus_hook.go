package player

import (
	"context"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

type runStatusHook struct {
	store   *runstatus.Store
	runner  string
	pending []runstatus.Entry
	mu      sync.Mutex
}

type runStatusCtxKey struct{}

// WithRunStatusHook buffers run_status rows during a run; call FlushRunStatus when the run ends.
func WithRunStatusHook(ctx context.Context, store *runstatus.Store, runner string) context.Context {
	if store == nil {
		return ctx
	}
	return context.WithValue(ctx, runStatusCtxKey{}, &runStatusHook{store: store, runner: runner})
}

func runStatusHookFrom(ctx context.Context) (*runStatusHook, bool) {
	h, ok := ctx.Value(runStatusCtxKey{}).(*runStatusHook)
	return h, ok && h != nil && h.store != nil
}

func recordScenarioRunStatus(ctx context.Context, result ScenarioResult) {
	h, ok := runStatusHookFrom(ctx)
	if !ok {
		return
	}
	entry := RunstatusEntry(result, h.runner)
	h.mu.Lock()
	h.pending = append(h.pending, entry)
	h.mu.Unlock()
}

// FlushRunStatus writes buffered run_status rows in one batch.
func FlushRunStatus(ctx context.Context) {
	h, ok := runStatusHookFrom(ctx)
	if !ok {
		return
	}
	h.mu.Lock()
	pending := h.pending
	h.pending = nil
	h.mu.Unlock()
	if len(pending) == 0 {
		return
	}
	_ = h.store.RecordBatch(pending)
}
