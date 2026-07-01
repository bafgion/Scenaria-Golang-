package player

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/paths"
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
	runID := fmt.Sprintf("%x", time.Now().UnixNano())
	logx.Info("run started", "run_id", runID, "scenarios", len(plan.Cases), "workers", workers)

	if workers == 1 || len(plan.Cases) <= 1 {
		for _, runCase := range plan.Cases {
			if err := ctx.Err(); err != nil {
				return result, err
			}
			runResult, err := r.Executor.ExecuteScenario(ctx, scenarioInputFromCase(runCase))
			if err != nil {
				if runResult.Scenario == "" {
					runResult = ScenarioResult{
						FeaturePath: runCase.FeaturePath,
						Scenario:    runCase.Name,
						Status:      "failed",
						Message:     err.Error(),
					}
				} else if runResult.Status == "" {
					runResult.Status = "failed"
					if runResult.Message == "" {
						runResult.Message = err.Error()
					}
				}
				result.ScenarioResults = append(result.ScenarioResults, runResult)
				return result, executionFailure(err, result)
			}
			if runResult.Status == "failed" {
				result.ScenarioResults = append(result.ScenarioResults, runResult)
				runErr := fmt.Errorf("scenario %q failed: %s", runCase.Name, runResult.Message)
				return result, executionFailure(runErr, result)
			}
			result.ScenarioResults = append(result.ScenarioResults, runResult)
		}
		return result, nil
	}

	if pwExec, ok := r.Executor.(*PlaywrightExecutor); ok && poolEligible(pwExec.options) {
		return r.executeParallelWithPool(ctx, result, pwExec, plan, workers, runID)
	}
	return r.executeParallel(ctx, result, plan, workers, runID)
}

func (r BrowserRunner) executeParallelWithPool(
	ctx context.Context,
	result ExecutionResult,
	exec *PlaywrightExecutor,
	plan ExecutionPlan,
	workers int,
	runID string,
) (ExecutionResult, error) {
	if exec.options.AutoInstall {
		if err := paths.EnsurePlaywrightEngine(exec.options.BrowserName); err != nil {
			return result, fmt.Errorf("playwright install failed: %w", err)
		}
	} else {
		paths.ConfigurePlaywrightBrowsersForEngine(exec.options.BrowserName)
	}

	pool, err := newBrowserPool(ctx, exec.options, workers)
	if err != nil {
		return result, err
	}
	defer pool.Close()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

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

			if err := runCtx.Err(); err != nil {
				mu.Lock()
				results[i] = ScenarioResult{
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
					Status:      "failed",
					Message:     err.Error(),
				}
				mu.Unlock()
				return
			}

			slot, err := pool.acquire(runCtx)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel()
				}
				results[i] = ScenarioResult{
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
					Status:      "failed",
					Message:     err.Error(),
				}
				mu.Unlock()
				return
			}
			runResult, err := exec.ExecuteScenarioOnSession(runCtx, slot.session, scenarioInputFromCase(rc))
			pool.release(slot)

			mu.Lock()
			defer mu.Unlock()
			scenarioFailed := err != nil || runResult.Status == "failed"
			if scenarioFailed {
				if firstErr == nil {
					if err != nil {
						firstErr = err
					} else {
						firstErr = fmt.Errorf("scenario %q failed: %s", rc.Name, runResult.Message)
					}
					cancel()
				}
				if err != nil && runResult.Scenario == "" {
					runResult = ScenarioResult{
						FeaturePath: rc.FeaturePath,
						Scenario:    rc.Name,
						Status:      "failed",
						Message:     err.Error(),
					}
				} else if runResult.Status == "" {
					runResult.Status = "failed"
					if runResult.Message == "" && err != nil {
						runResult.Message = err.Error()
					}
				}
			}
			results[i] = runResult
		}(index, runCase)
	}
	wg.Wait()

	result.ScenarioResults = results
	if firstErr != nil {
		logx.Warn("run failed", "run_id", runID, "error", firstErr)
		return result, executionFailure(firstErr, result)
	}
	if err := ctx.Err(); err != nil {
		return result, executionFailure(err, result)
	}
	return result, nil
}

func (r BrowserRunner) executeParallel(
	ctx context.Context,
	result ExecutionResult,
	plan ExecutionPlan,
	workers int,
	runID string,
) (ExecutionResult, error) {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return ExecutionResult{}, err
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

			if err := runCtx.Err(); err != nil {
				mu.Lock()
				results[i] = ScenarioResult{
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
					Status:      "failed",
					Message:     err.Error(),
				}
				mu.Unlock()
				return
			}

			runResult, err := r.Executor.ExecuteScenario(runCtx, scenarioInputFromCase(rc))
			mu.Lock()
			defer mu.Unlock()
			scenarioFailed := err != nil || runResult.Status == "failed"
			if scenarioFailed {
				if firstErr == nil {
					if err != nil {
						firstErr = err
					} else {
						firstErr = fmt.Errorf("scenario %q failed: %s", rc.Name, runResult.Message)
					}
					cancel()
				}
				if err != nil {
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
			}
			results[i] = runResult
		}(index, runCase)
	}
	wg.Wait()
	result.ScenarioResults = results
	if firstErr != nil {
		logx.Warn("run failed", "run_id", runID, "error", firstErr)
		return result, executionFailure(firstErr, result)
	}
	if err := ctx.Err(); err != nil {
		return result, executionFailure(err, result)
	}
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
