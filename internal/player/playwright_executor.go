package player

import (
	"context"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func (e *PlaywrightExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	exec := NewStepExecutor(ExecutorOptions{
		BaseURL:           e.options.BaseURL,
		MaxLoopIterations: e.options.MaxLoopIterations,
	})
	prompt := e.options.PromptEmailCode
	if prompt == nil {
		prompt = emailCodePrompter()
	}
	steps := input.Steps
	if input.StartStep != -1 || input.EndStep != -1 {
		start := input.StartStep
		if start < 0 {
			start = 0
		}
		steps = gherkin.ApplyStepRange(steps, start, input.EndStep)
	}
	return e.executeWithSession(ctx, input, func(ctx context.Context, session *browserSession) error {
		runCtx := NewRunContext(input.Variables, time.Now().UnixNano(), input.ProjectRoot, WithPromptEmailCode(prompt))
		return exec.ExecuteSteps(ctx, session, steps, runCtx)
	})
}
