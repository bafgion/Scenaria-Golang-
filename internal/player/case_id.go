package player

import (
	"fmt"
	"path/filepath"
	"strings"
)

// BuildCaseID returns a stable identity for a runnable scenario case.
func BuildCaseID(featurePath, scenario string, exampleIndex int) string {
	fp := strings.ReplaceAll(filepath.ToSlash(strings.TrimSpace(featurePath)), "\\", "/")
	id := fp + "::" + strings.TrimSpace(scenario)
	if exampleIndex > 0 {
		id += fmt.Sprintf("#%d", exampleIndex)
	}
	return id
}

func scenarioResultFromCase(runCase RunCase, status, message string) ScenarioResult {
	return ScenarioResult{
		FeaturePath:  runCase.FeaturePath,
		Scenario:     runCase.Name,
		CaseID:       runCase.CaseID,
		ExampleIndex: runCase.ExampleIndex,
		Status:       status,
		Message:      message,
	}
}

func stampRunID(result *ExecutionResult, runID string) {
	if result == nil || runID == "" {
		return
	}
	result.RunID = runID
	for i := range result.ScenarioResults {
		result.ScenarioResults[i].RunID = runID
	}
}
