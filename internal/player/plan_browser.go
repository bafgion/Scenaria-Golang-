package player

import "github.com/bafgion/scenaria-golang/internal/gherkin"

// PlanContainsCloseBrowser reports whether any scenario in the plan ends with a close-browser step.
func PlanContainsCloseBrowser(plan ExecutionPlan) bool {
	for _, runCase := range plan.Cases {
		if stepsContainCloseBrowser(runCase.Steps) {
			return true
		}
	}
	return false
}

func stepsContainCloseBrowser(steps []gherkin.Step) bool {
	for _, step := range steps {
		if step.Block != "" {
			if stepsContainCloseBrowser(step.Children) {
				return true
			}
			continue
		}
		if terminalActionFromStep(step) == "close-browser" {
			return true
		}
	}
	return false
}
