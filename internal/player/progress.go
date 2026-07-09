package player

import "context"

// RunProgressPhase marks run lifecycle events for GUI progress UI.
type RunProgressPhase string

const (
	ProgressScenarioStart RunProgressPhase = "scenario_start"
	ProgressScenarioDone  RunProgressPhase = "scenario_done"
)

// RunProgressEvent is emitted during Execute for live progress in the IDE.
type RunProgressEvent struct {
	Phase       RunProgressPhase `json:"phase"`
	Index       int              `json:"index"`
	Total       int              `json:"total"`
	CaseID      string           `json:"caseId,omitempty"`
	FeaturePath string           `json:"featurePath,omitempty"`
	Scenario    string           `json:"scenario,omitempty"`
	Success     bool             `json:"success,omitempty"`
	Message     string           `json:"message,omitempty"`
}

type progressCtxKey struct{}

// WithRunProgress attaches a progress reporter to ctx (GUI runs only).
func WithRunProgress(ctx context.Context, fn func(RunProgressEvent)) context.Context {
	if fn == nil {
		return ctx
	}
	return context.WithValue(ctx, progressCtxKey{}, fn)
}

func emitRunProgress(ctx context.Context, ev RunProgressEvent) {
	if ctx == nil {
		return
	}
	fn, _ := ctx.Value(progressCtxKey{}).(func(RunProgressEvent))
	if fn != nil {
		fn(ev)
	}
}
