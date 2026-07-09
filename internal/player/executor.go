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
	RetryPolicy       RetryPolicy
	MaxActionAttempts int // 0 = DefaultMaxActionAttempts
}

func (e *StepExecutor) maxActionAttempts() int {
	if e != nil && e.options.MaxActionAttempts > 0 {
		return e.options.MaxActionAttempts
	}
	return DefaultMaxActionAttempts
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
		if page, err := session.currentPage(); err == nil {
			runCtx.SetPage(page)
		}
	}
	for i, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := e.executeStep(ctx, session, step, runCtx); err != nil {
			return fmt.Errorf("line %d: %w", step.Line, err)
		}
		if session.isClosed() {
			terminal := terminalActionFromStep(step)
			if terminal == "" {
				terminal = "close-browser"
			}
			if err := e.finishIfBrowserClosed(runCtx, steps, i+1, terminal); err != nil {
				return err
			}
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
		ok, err := evaluateRunCondition(runCtx, step.Condition)
		if err != nil {
			return err
		}
		if ok {
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
			ok, err := evaluateRunCondition(runCtx, step.Condition)
			if err != nil {
				return err
			}
			if !ok {
				break
			}
			iterations++
			iterErr := runCtx.withIteration("while", iterations, func() error {
				return e.ExecuteSteps(ctx, session, step.Children, runCtx)
			})
			if iterErr != nil {
				return iterErr
			}
			if session.isClosed() {
				return nil
			}
		}
		if iterations >= limit {
			ok, err := evaluateRunCondition(runCtx, step.Condition)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			return fmt.Errorf("while loop exceeded max loop iterations")
		}
		return nil
	case gherkin.BlockRepeat:
		count := step.RepeatCount
		if count < 1 {
			return fmt.Errorf("repeat count must be at least 1")
		}
		limit := e.maxLoopIterations()
		if count > limit {
			return fmt.Errorf("repeat count %d exceeds max loop iterations %d", count, limit)
		}
		for i := 0; i < count; i++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			iterErr := runCtx.withIteration("repeat", i+1, func() error {
				return e.ExecuteSteps(ctx, session, step.Children, runCtx)
			})
			if iterErr != nil {
				return iterErr
			}
			if session.isClosed() {
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

	if runCtx != nil {
		if err := runCtx.consumeActionAttempt(e.maxActionAttempts()); err != nil {
			return err
		}
	}

	idx := -1
	started := time.Now()
	if runCtx != nil {
		idx = runCtx.beginLeafStep(step)
	}

	action, err := stepdsl.Parse(step)
	if err != nil {
		if runCtx != nil && idx >= 0 {
			runCtx.completeLeafStep(idx, "", started, session, err)
		}
		return e.failLeafStep(runCtx, err)
	}
	if runCtx != nil {
		if action.Value1, err = runCtx.ResolveText(action.Value1); err != nil {
			runCtx.completeLeafStep(idx, actionSelector(action), started, session, err)
			return e.failLeafStep(runCtx, err)
		}
		if action.Value2, err = runCtx.ResolveText(action.Value2); err != nil {
			runCtx.completeLeafStep(idx, actionSelector(action), started, session, err)
			return e.failLeafStep(runCtx, err)
		}
	}
	selector := actionSelector(action)
	if session != nil {
		session.clearNetworkFailure()
	}
	if err := e.runAction(ctx, session, action, runCtx, idx); err != nil {
		if runCtx != nil && idx >= 0 {
			runCtx.completeLeafStep(idx, selector, started, session, err)
		}
		return e.failLeafStep(runCtx, err)
	}
	if runCtx != nil {
		if action.Kind == "close-browser" || action.Kind == "close-tab" {
			runCtx.setStepTerminalAction(idx, action.Kind)
		}
		if idx >= 0 {
			runCtx.completeLeafStep(idx, selector, started, session, nil)
		}
		runCtx.RecordStep(step)
	}
	return nil
}

func evaluateRunCondition(runCtx *RunContext, cond *gherkin.Condition) (bool, error) {
	if runCtx == nil {
		return false, nil
	}
	return runCtx.EvaluateConditionResult(cond)
}

func (e *StepExecutor) executeForEach(ctx context.Context, session *browserSession, step gherkin.Step, runCtx *RunContext) error {
	if runCtx == nil {
		return fmt.Errorf("for_each requires run context")
	}
	selector, err := runCtx.ResolveText(step.ForEachSelector)
	if err != nil {
		return err
	}
	page, err := session.currentPage()
	if err != nil {
		return err
	}
	// for_each uses live DOM: count is re-checked before each iteration.
	locators, err := page.Locator(selector).Count()
	if err != nil {
		return fmt.Errorf("for_each locator %q: %w", selector, err)
	}
	for index := 0; index < locators; index++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if index >= e.maxLoopIterations() {
			return fmt.Errorf("превышен лимит итераций for_each (%d)", e.maxLoopIterations())
		}
		currentCount, err := page.Locator(selector).Count()
		if err != nil {
			return fmt.Errorf("for_each index %d selector %q: count failed: %w", index, selector, err)
		}
		if index >= currentCount {
			return fmt.Errorf("for_each index %d selector %q: element removed during iteration (count %d)", index, selector, currentCount)
		}
		locator := page.Locator(selector).Nth(index)
		text, err := locator.InnerText()
		if err != nil {
			return fmt.Errorf("for_each index %d selector %q: %w", index, selector, err)
		}
		text = strings.TrimSpace(text)
		if text == "" {
			return fmt.Errorf("for_each index %d selector %q: empty inner text", index, selector)
		}
		runCtx.Remember(step.ForEachVariable, text)
		iterErr := runCtx.withIteration("for_each", index+1, func() error {
			return e.ExecuteSteps(ctx, session, step.Children, runCtx)
		})
		if iterErr != nil {
			return iterErr
		}
		if session.isClosed() {
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
