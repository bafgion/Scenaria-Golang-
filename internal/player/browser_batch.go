package player

import "context"

func notStartedScenarioResult(runCase RunCase) ScenarioResult {
	return scenarioResultFromCase(runCase, "not-started", MsgScenarioNotStartedBrowserClosed)
}

func appendNotStartedScenario(
	ctx context.Context,
	result *ExecutionResult,
	runCase RunCase,
	index, total int,
) {
	runResult := notStartedScenarioResult(runCase)
	result.ScenarioResults = append(result.ScenarioResults, runResult)
	recordScenarioRunStatus(ctx, runResult)
	emitRunProgress(ctx, RunProgressEvent{
		Phase:       ProgressScenarioDone,
		Index:       index,
		Total:       total,
		FeaturePath: runCase.FeaturePath,
		Scenario:    runCase.Name,
		Success:     false,
		Message:     runResult.Message,
	})
}

func appendRemainingNotStarted(
	ctx context.Context,
	result *ExecutionResult,
	plan ExecutionPlan,
	fromIndex, total int,
) {
	for j := fromIndex; j < len(plan.Cases); j++ {
		runCase := plan.Cases[j]
		emitRunProgress(ctx, RunProgressEvent{
			Phase:       ProgressScenarioStart,
			Index:       j + 1,
			Total:       total,
			FeaturePath: runCase.FeaturePath,
			Scenario:    runCase.Name,
		})
		appendNotStartedScenario(ctx, result, runCase, j+1, total)
	}
}
