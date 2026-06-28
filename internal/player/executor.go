package player

import (
	"context"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type ExecutorOptions struct {
	BaseURL string
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
	switch step.Block {
	case gherkin.BlockIf:
		if runCtx != nil && runCtx.EvaluateCondition(step.Condition) {
			return e.ExecuteSteps(ctx, session, step.Children, runCtx)
		}
		return nil
	case gherkin.BlockWhile:
		iterations := 0
		for iterations < MaxLoopIterations {
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
		if iterations >= MaxLoopIterations && runCtx != nil && runCtx.EvaluateCondition(step.Condition) {
			return fmt.Errorf("превышен лимит итераций цикла «пока»")
		}
		return nil
	case gherkin.BlockRepeat:
		count := step.RepeatCount
		if count < 1 {
			count = 1
		}
		if count > MaxLoopIterations {
			count = MaxLoopIterations
		}
		for i := 0; i < count; i++ {
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

	action, err := stepdsl.Parse(step)
	if err != nil {
		return err
	}
	if runCtx != nil {
		if action.Value1, err = runCtx.ResolveText(action.Value1); err != nil {
			return err
		}
		if action.Value2, err = runCtx.ResolveText(action.Value2); err != nil {
			return err
		}
	}
	if err := executeAction(ctx, session, action, e.options.BaseURL, runCtx); err != nil {
		return err
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
	locators, err := session.page.Locator(selector).All()
	if err != nil {
		return fmt.Errorf("for_each locator failed: %w", err)
	}
	for index, locator := range locators {
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
