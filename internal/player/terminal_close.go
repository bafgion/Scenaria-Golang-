package player

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type terminalCloseError struct {
	action    string
	remaining int
}

func (e *terminalCloseError) Error() string {
	return fmt.Sprintf("scenario ended by %s with %d remaining step(s)", e.action, e.remaining)
}

func remainingLeafStepCount(steps []gherkin.Step, from int) int {
	if from >= len(steps) {
		return 0
	}
	return gherkin.CountLeafSteps(steps[from:])
}

func (e *StepExecutor) finishIfBrowserClosed(
	runCtx *RunContext,
	steps []gherkin.Step,
	from int,
	terminalAction string,
) error {
	if from >= len(steps) {
		return nil
	}
	remaining := remainingLeafStepCount(steps, from)
	if remaining == 0 {
		return nil
	}
	if runCtx != nil {
		runCtx.appendSkippedLeafSteps(steps[from:])
	}
	return &terminalCloseError{action: terminalAction, remaining: remaining}
}

func terminalActionFromStep(step gherkin.Step) string {
	if step.Block != "" {
		return ""
	}
	action, err := stepdsl.Parse(step)
	if err != nil {
		return ""
	}
	switch action.Kind {
	case "close-browser", "close-tab":
		return action.Kind
	default:
		return ""
	}
}
