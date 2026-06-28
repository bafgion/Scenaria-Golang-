package player

import (
	"context"
	"time"
)

func (e *PlaywrightExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	exec := NewStepExecutor(ExecutorOptions{BaseURL: e.options.BaseURL})
	return e.executeWithSession(ctx, input, func(ctx context.Context, session *browserSession) error {
		runCtx := NewRunContext(input.Variables, time.Now().UnixNano(), input.ProjectRoot)
		return exec.ExecuteSteps(ctx, session, input.Steps, runCtx)
	})
}
