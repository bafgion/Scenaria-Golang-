package player

import (
	"context"

	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

type runStatusHook struct {
	store  *runstatus.Store
	runner string
}

type runStatusCtxKey struct{}

// WithRunStatusHook enables incremental run_status.json writes after each scenario.
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
	entry := runstatus.Entry{
		Path:    result.FeaturePath + "::" + result.Scenario,
		Success: result.Status == "passed",
		Message: result.Message,
		Runner:  h.runner,
	}
	if result.FailedStep != nil {
		entry.FailedStep = result.FailedStep
	}
	_ = h.store.Record(entry)
}
