package player

import (
	"context"
	"fmt"
	"sync"
)

// BrowserRunner executes scenarios sequentially or with a worker pool.
type BrowserRunner struct {
	Executor        BrowserExecutor
	ParallelWorkers int
}

func (r BrowserRunner) Execute(ctx context.Context, plan ExecutionPlan) (ExecutionResult, error) {
	if r.Executor == nil {
		return ExecutionResult{}, fmt.Errorf("browser runner: executor is nil")
	}
	files, scenarios, steps, _ := SummarizePlan(plan)
	result := ExecutionResult{
		Mode:      "browser",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}
	workers := r.ParallelWorkers
	if workers < 1 {
		workers = 1
	}
	if workers == 1 || len(plan.Cases) <= 1 {
		for _, runCase := range plan.Cases {
			runResult, err := r.Executor.ExecuteScenario(ctx, scenarioInputFromCase(runCase))
			if err != nil {
				return ExecutionResult{}, err
			}
			result.ScenarioResults = append(result.ScenarioResults, runResult)
		}
		return result, nil
	}

	results := make([]ScenarioResult, len(plan.Cases))
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	for index, runCase := range plan.Cases {
		wg.Add(1)
		go func(i int, rc RunCase) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			runResult, err := r.Executor.ExecuteScenario(ctx, scenarioInputFromCase(rc))
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				if runResult.Scenario == "" {
					runResult = ScenarioResult{
						FeaturePath: rc.FeaturePath,
						Scenario:    rc.Name,
						Status:      "failed",
						Message:     err.Error(),
					}
				} else if runResult.Status == "" {
					runResult.Status = "failed"
					if runResult.Message == "" {
						runResult.Message = err.Error()
					}
				}
			}
			results[i] = runResult
		}(index, runCase)
	}
	wg.Wait()
	if firstErr != nil {
		return ExecutionResult{}, firstErr
	}
	result.ScenarioResults = results
	return result, nil
}

func scenarioInputFromCase(runCase RunCase) ScenarioInput {
	return ScenarioInput{
		FeaturePath:  runCase.FeaturePath,
		ScenarioName: runCase.Name,
		Steps:        runCase.Steps,
		TestClient:   runCase.TestClient,
		Variables:    runCase.Variables,
		ProjectRoot:  runCase.ProjectRoot,
		StartStep:    runCase.StartStep,
		EndStep:      runCase.EndStep,
	}
}
