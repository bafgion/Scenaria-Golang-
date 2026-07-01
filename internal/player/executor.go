package player

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type ExecutorOptions struct {
	BaseURL           string
	MaxLoopIterations int
	MaxActionRetries  int // 0 = default; <0 = disable retries
	RetryBackoff      time.Duration
}

func (e *StepExecutor) maxLoopIterations() int {
	if e != nil && e.options.MaxLoopIterations > 0 {
		return e.options.MaxLoopIterations
	}
	return DefaultMaxLoopIterations
}

type StepExecutor struct {
	options ExecutorOptions
}

func NewStepExecutor(options ExecutorOptions) *StepExecutor {
	return &StepExecutor{options: options}
}

func (e *StepExecutor) ExecuteSteps(ctx context.Context, session *browserSession, steps []gherkin.Step, runCtx *RunContext) error {
	if runCtx != nil {
		runCtx.SetPage(session.page)
	}
	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := e.executeStep(ctx, session, step, runCtx); err != nil {
			return fmt.Errorf("line %d: %w", step.Line, err)
		}
		if session.closed {
			return nil
		}
	}
	return nil
}

func (e *StepExecutor) executeStep(ctx context.Context, session *browserSession, step gherkin.Step, runCtx *RunContext) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	switch step.Block {
	case gherkin.BlockIf:
		if runCtx != nil && runCtx.EvaluateCondition(step.Condition) {
			return e.ExecuteSteps(ctx, session, step.Children, runCtx)
		}
		return nil
	case gherkin.BlockWhile:
		iterations := 0
		limit := e.maxLoopIterations()
		for iterations < limit {
			if err := ctx.Err(); err != nil {
				return err
			}
			if runCtx == nil || !runCtx.EvaluateCondition(step.Condition) {
				break
			}
			iterations++
			if err := e.ExecuteSteps(ctx, session, step.Children, runCtx); err != nil {
				return err
			}
			if session.closed {
				return nil
			}
		}
		if iterations >= limit && runCtx != nil && runCtx.EvaluateCondition(step.Condition) {
			return fmt.Errorf("превышен лимит итераций цикла «пока»")
		}
		return nil
	case gherkin.BlockRepeat:
		count := step.RepeatCount
		if count < 1 {
			count = 1
		}
		limit := e.maxLoopIterations()
		if count > limit {
			count = limit
		}
		for i := 0; i < count; i++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := e.ExecuteSteps(ctx, session, step.Children, runCtx); err != nil {
				return err
			}
			if session.closed {
				return nil
			}
		}
		return nil
	case gherkin.BlockForEach:
		return e.executeForEach(ctx, session, step, runCtx)
	}

	if gherkin.IsTestClientStep(step) {
		return nil
	}

	action, err := stepdsl.Parse(step)
	if err != nil {
		return e.failLeafStep(runCtx, err)
	}
	if runCtx != nil {
		if action.Value1, err = runCtx.ResolveText(action.Value1); err != nil {
			return e.failLeafStep(runCtx, err)
		}
		if action.Value2, err = runCtx.ResolveText(action.Value2); err != nil {
			return e.failLeafStep(runCtx, err)
		}
	}
	if err := e.runAction(ctx, session, action, runCtx); err != nil {
		return e.failLeafStep(runCtx, err)
	}
	if runCtx != nil {
		runCtx.RecordStep(step)
	}
	return nil
}

func (e *StepExecutor) executeForEach(ctx context.Context, session *browserSession, step gherkin.Step, runCtx *RunContext) error {
	if runCtx == nil {
		return fmt.Errorf("for_each requires run context")
	}
	selector, err := runCtx.ResolveText(step.ForEachSelector)
	if err != nil {
		return err
	}
	locators, err := session.page.Locator(selector).Count()
	if err != nil {
		return fmt.Errorf("for_each locator failed: %w", err)
	}
	for index := 0; index < locators; index++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		locator := session.page.Locator(selector).Nth(index)
		text, _ := locator.InnerText()
		text = strings.TrimSpace(text)
		if text == "" {
			text = fmt.Sprintf("%d", index+1)
		}
		runCtx.Remember(step.ForEachVariable, text)
		if err := e.ExecuteSteps(ctx, session, step.Children, runCtx); err != nil {
			return err
		}
		if session.closed {
			return nil
		}
	}
	return nil
}

func (e *StepExecutor) failLeafStep(runCtx *RunContext, err error) error {
	if runCtx != nil {
		runCtx.markFailedLeafStep()
	}
	return err
}
