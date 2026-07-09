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

// ScenarioEndedWithCloseBrowser reports whether a finished scenario executed close-browser.
func ScenarioEndedWithCloseBrowser(result ScenarioResult) bool {
	for _, rec := range result.StepRecords {
		if rec.TerminalAction == "close-browser" {
			return true
		}
	}
	return false
}

// ScenarioBlocksSessionReuse reports whether the next scenario must not reuse the current browser session.
func ScenarioBlocksSessionReuse(result ScenarioResult, runCase RunCase, session *browserSession) bool {
	if ScenarioEndedWithCloseBrowser(result) {
		return true
	}
	if result.Status == "passed" && stepsContainCloseBrowser(runCase.Steps) {
		return true
	}
	return result.Status == "passed" && session != nil && session.isClosed()
}
