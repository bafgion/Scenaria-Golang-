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

type indexedRunCase struct {
	index   int
	runCase RunCase
}

func (r BrowserRunner) Execute(ctx context.Context, plan ExecutionPlan) (result ExecutionResult, err error) {
	defer FlushRunStatus(ctx)
	if r.Executor == nil {
		return ExecutionResult{}, fmt.Errorf("browser runner: executor is nil")
	}
	files, scenarios, steps, _ := SummarizePlan(plan)
	result = ExecutionResult{
		Mode:      "browser",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}
	workers := r.ParallelWorkers
	if workers < 1 {
		workers = 1
	}
	runID := RunIDFromContext(ctx)
	if runID == "" {
		runID = fmt.Sprintf("run-%d", time.Now().UnixNano())
	}
	defer func() { stampRunID(&result, runID) }()
	logx.Info("run started", "run_id", runID, "scenarios", len(plan.Cases), "workers", workers)
	logx.Debug("parallel run scheduled", "run_id", runID, "scheduled_cases", len(plan.Cases), "workers", workers)

	if workers == 1 || len(plan.Cases) <= 1 {
		if pwExec, ok := r.Executor.(*PlaywrightExecutor); ok {
			result, err = r.executeSequentialSession(ctx, result, pwExec, plan, nil)
			return result, err
		}
		total := len(plan.Cases)
		var firstErr error
		for i, runCase := range plan.Cases {
			if err := ctx.Err(); err != nil {
				return result, err
			}
			emitRunProgress(ctx, RunProgressEvent{
				Phase:       ProgressScenarioStart,
				Index:       i + 1,
				Total:       total,
				FeaturePath: runCase.FeaturePath,
				Scenario:    runCase.Name,
			})
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
				recordScenarioRunStatus(ctx, runResult)
				emitRunProgress(ctx, RunProgressEvent{
					Phase:       ProgressScenarioDone,
					Index:       i + 1,
					Total:       total,
					FeaturePath: runCase.FeaturePath,
					Scenario:    runCase.Name,
					Success:     false,
					Message:     runResult.Message,
				})
				if !ContinueOnFail(ctx) {
					return result, executionFailure(err, result)
				}
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			if runResult.Status == "failed" {
				result.ScenarioResults = append(result.ScenarioResults, runResult)
				recordScenarioRunStatus(ctx, runResult)
				emitRunProgress(ctx, RunProgressEvent{
					Phase:       ProgressScenarioDone,
					Index:       i + 1,
					Total:       total,
					FeaturePath: runCase.FeaturePath,
					Scenario:    runCase.Name,
					Success:     false,
					Message:     runResult.Message,
				})
				runErr := fmt.Errorf("scenario %q failed: %s", runCase.Name, runResult.Message)
				if !ContinueOnFail(ctx) {
					return result, executionFailure(runErr, result)
				}
				if firstErr == nil {
					firstErr = runErr
				}
				continue
			}
			result.ScenarioResults = append(result.ScenarioResults, runResult)
			recordScenarioRunStatus(ctx, runResult)
			emitRunProgress(ctx, RunProgressEvent{
				Phase:       ProgressScenarioDone,
				Index:       i + 1,
				Total:       total,
				FeaturePath: runCase.FeaturePath,
				Scenario:    runCase.Name,
				Success:     true,
			})
		}
		if firstErr != nil {
			return result, executionFailure(firstErr, result)
		}
		return result, nil
	}

	if pwExec, ok := r.Executor.(*PlaywrightExecutor); ok && poolEligible(pwExec.options) {
		result, err = r.executeParallelWithPool(ctx, result, pwExec, plan, workers, runID)
		return result, err
	}
	result, err = r.executeParallel(ctx, result, plan, workers, runID)
	return result, err
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

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool, err := newBrowserPool(runCtx, exec.options, workers)
	if err != nil {
		return result, err
	}
	defer pool.Close()

	results := make([]ScenarioResult, len(plan.Cases))
	jobs := make(chan indexedRunCase)
	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				i, rc := job.index, job.runCase
				if err := runCtx.Err(); err != nil {
					failed := canceledScenarioResult(rc, err)
					mu.Lock()
					results[i] = failed
					mu.Unlock()
					recordScenarioRunStatus(runCtx, failed)
					emitRunProgress(runCtx, RunProgressEvent{
						Phase:       ProgressScenarioDone,
						Index:       i + 1,
						Total:       len(plan.Cases),
						FeaturePath: rc.FeaturePath,
						Scenario:    rc.Name,
						Success:     false,
						Message:     err.Error(),
					})
					continue
				}

				emitRunProgress(runCtx, RunProgressEvent{
					Phase:       ProgressScenarioStart,
					Index:       i + 1,
					Total:       len(plan.Cases),
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
				})

				slot, err := pool.acquire(runCtx)
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						failFastParallelCancel(runCtx, cancel, pool)
					}
					results[i] = ScenarioResult{
						FeaturePath: rc.FeaturePath,
						Scenario:    rc.Name,
						Status:      "failed",
						Message:     err.Error(),
					}
					mu.Unlock()
					continue
				}
				runResult, err := exec.ExecuteScenarioOnSession(runCtx, slot.session, scenarioInputFromCase(rc))
				pool.release(slot)

				mu.Lock()
				scenarioFailed := err != nil || runResult.Status == "failed"
				if scenarioFailed {
					if firstErr == nil {
						if err != nil {
							firstErr = err
						} else {
							firstErr = fmt.Errorf("scenario %q failed: %s", rc.Name, runResult.Message)
						}
						failFastParallelCancel(runCtx, cancel, pool)
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
				mu.Unlock()

				recordScenarioRunStatus(runCtx, runResult)
				emitRunProgress(runCtx, RunProgressEvent{
					Phase:       ProgressScenarioDone,
					Index:       i + 1,
					Total:       len(plan.Cases),
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
					Success:     !scenarioFailed,
					Message:     runResult.Message,
				})
			}
		}()
	}
	enqueueJobs(runCtx, jobs, plan.Cases)
	close(jobs)
	wg.Wait()
	fillCanceledResults(runCtx, results, plan.Cases)

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
	jobs := make(chan indexedRunCase)
	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				i, rc := job.index, job.runCase
				if err := runCtx.Err(); err != nil {
					failed := canceledScenarioResult(rc, err)
					mu.Lock()
					results[i] = failed
					mu.Unlock()
					recordScenarioRunStatus(runCtx, failed)
					emitRunProgress(runCtx, RunProgressEvent{
						Phase:       ProgressScenarioDone,
						Index:       i + 1,
						Total:       len(plan.Cases),
						FeaturePath: rc.FeaturePath,
						Scenario:    rc.Name,
						Success:     false,
						Message:     err.Error(),
					})
					continue
				}

				emitRunProgress(runCtx, RunProgressEvent{
					Phase:       ProgressScenarioStart,
					Index:       i + 1,
					Total:       len(plan.Cases),
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
				})

				runResult, err := r.Executor.ExecuteScenario(runCtx, scenarioInputFromCase(rc))
				mu.Lock()
				scenarioFailed := err != nil || runResult.Status == "failed"
				if scenarioFailed {
					if firstErr == nil {
						if err != nil {
							firstErr = err
						} else {
							firstErr = fmt.Errorf("scenario %q failed: %s", rc.Name, runResult.Message)
						}
						failFastParallelCancel(runCtx, cancel, nil)
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
				mu.Unlock()

				recordScenarioRunStatus(runCtx, runResult)
				emitRunProgress(runCtx, RunProgressEvent{
					Phase:       ProgressScenarioDone,
					Index:       i + 1,
					Total:       len(plan.Cases),
					FeaturePath: rc.FeaturePath,
					Scenario:    rc.Name,
					Success:     !scenarioFailed,
					Message:     runResult.Message,
				})
			}
		}()
	}
	enqueueJobs(runCtx, jobs, plan.Cases)
	close(jobs)
	wg.Wait()
	fillCanceledResults(runCtx, results, plan.Cases)
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

func (r BrowserRunner) executeSequentialPlaywright(
	ctx context.Context,
	result ExecutionResult,
	exec *PlaywrightExecutor,
	plan ExecutionPlan,
) (ExecutionResult, error) {
	return r.executeSequentialSession(ctx, result, exec, plan, nil)
}

// ExecuteSequentialAttached runs scenarios on an already open browser page (IDE live session).
func (r BrowserRunner) ExecuteSequentialAttached(
	ctx context.Context,
	exec *PlaywrightExecutor,
	plan ExecutionPlan,
	session *browserSession,
) (ExecutionResult, error) {
	files, scenarios, steps, _ := SummarizePlan(plan)
	result := ExecutionResult{
		Mode:      "browser",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}
	return r.executeSequentialSession(ctx, result, exec, plan, session)
}

func (r BrowserRunner) executeSequentialSession(
	ctx context.Context,
	result ExecutionResult,
	exec *PlaywrightExecutor,
	plan ExecutionPlan,
	attached *browserSession,
) (ExecutionResult, error) {
	if attached == nil {
		if exec.options.AutoInstall {
			if err := paths.EnsurePlaywrightEngine(exec.options.BrowserName); err != nil {
				return result, fmt.Errorf("playwright install failed: %w", err)
			}
		} else {
			paths.ConfigurePlaywrightBrowsersForEngine(exec.options.BrowserName)
		}
	}

	total := len(plan.Cases)
	var firstErr error
	session := attached
	var stopPW func()
	var stopWatch func()

	openSession := func() error {
		if session != nil && session.alive() {
			return nil
		}
		if attached != nil {
			return fmt.Errorf("браузер не открыт")
		}
		if stopWatch != nil {
			stopWatch()
			stopWatch = nil
		}
		if stopPW != nil {
			if session != nil {
				session.close()
				time.Sleep(playwrightDrainDelay)
			}
			stopPW()
			stopPW = nil
		}
		pw, stop, err := startPlaywright(ctx)
		if err != nil {
			return fmt.Errorf("start playwright: %w", err)
		}
		stopPW = stop
		next, err := newBrowserSession(pw, exec.options)
		if err != nil {
			stopPW()
			stopPW = nil
			return err
		}
		session = next
		stopWatch = session.watchContext(ctx)
		return nil
	}
	defer func() {
		if stopWatch != nil {
			stopWatch()
		}
		if attached == nil && exec.options.CloseAfterRun && session != nil {
			session.close()
			time.Sleep(playwrightDrainDelay)
		}
		if stopPW != nil && exec.options.CloseAfterRun {
			stopPW()
		}
	}()

	if session == nil {
		if err := openSession(); err != nil {
			return result, err
		}
	} else if attached != nil {
		stopWatch = session.watchContext(ctx)
	}

	for i, runCase := range plan.Cases {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		emitRunProgress(ctx, RunProgressEvent{
			Phase:       ProgressScenarioStart,
			Index:       i + 1,
			Total:       total,
			FeaturePath: runCase.FeaturePath,
			Scenario:    runCase.Name,
		})

		if i > 0 && attached == nil {
			if err := session.resetForScenario(); err != nil {
				runResult := ScenarioResult{
					FeaturePath: runCase.FeaturePath,
					Scenario:    runCase.Name,
					Status:      "failed",
					Message:     err.Error(),
				}
				result.ScenarioResults = append(result.ScenarioResults, runResult)
				recordScenarioRunStatus(ctx, runResult)
				emitRunProgress(ctx, RunProgressEvent{
					Phase: ProgressScenarioDone, Index: i + 1, Total: total,
					FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: false, Message: runResult.Message,
				})
				if !ContinueOnFail(ctx) {
					return result, executionFailure(err, result)
				}
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
		}

		if session == nil || !session.alive() {
			if session != nil && session.external {
				runResult := ScenarioResult{
					FeaturePath: runCase.FeaturePath,
					Scenario:    runCase.Name,
					Status:      "failed",
					Message:     "браузер закрыт — откройте браузер или уберите шаг «закрываю браузер»",
				}
				result.ScenarioResults = append(result.ScenarioResults, runResult)
				recordScenarioRunStatus(ctx, runResult)
				emitRunProgress(ctx, RunProgressEvent{
					Phase: ProgressScenarioDone, Index: i + 1, Total: total,
					FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: false, Message: runResult.Message,
				})
				err := fmt.Errorf("scenario %q failed: %s", runCase.Name, runResult.Message)
				if !ContinueOnFail(ctx) {
					return result, executionFailure(err, result)
				}
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			if err := openSession(); err != nil {
				return result, err
			}
		}

		runResult, err := exec.ExecuteScenarioOnSession(ctx, session, scenarioInputFromCase(runCase))
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
			recordScenarioRunStatus(ctx, runResult)
			emitRunProgress(ctx, RunProgressEvent{
				Phase: ProgressScenarioDone, Index: i + 1, Total: total,
				FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: false, Message: runResult.Message,
			})
			if !ContinueOnFail(ctx) {
				return result, executionFailure(err, result)
			}
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if runResult.Status == "failed" {
			result.ScenarioResults = append(result.ScenarioResults, runResult)
			recordScenarioRunStatus(ctx, runResult)
			emitRunProgress(ctx, RunProgressEvent{
				Phase: ProgressScenarioDone, Index: i + 1, Total: total,
				FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: false, Message: runResult.Message,
			})
			runErr := fmt.Errorf("scenario %q failed: %s", runCase.Name, runResult.Message)
			if !ContinueOnFail(ctx) {
				return result, executionFailure(runErr, result)
			}
			if firstErr == nil {
				firstErr = runErr
			}
			continue
		}
		result.ScenarioResults = append(result.ScenarioResults, runResult)
		recordScenarioRunStatus(ctx, runResult)
		emitRunProgress(ctx, RunProgressEvent{
			Phase: ProgressScenarioDone, Index: i + 1, Total: total,
			FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: true,
		})
	}
	if firstErr != nil {
		return result, executionFailure(firstErr, result)
	}
	return result, nil
}

func scenarioInputFromCase(runCase RunCase) ScenarioInput {
	return ScenarioInput{
		CaseID:       runCase.CaseID,
		FeaturePath:  runCase.FeaturePath,
		ScenarioName: runCase.Name,
		ExampleIndex: runCase.ExampleIndex,
		Steps:        runCase.Steps,
		TestClient:   runCase.TestClient,
		Variables:    runCase.Variables,
		ProjectRoot:  runCase.ProjectRoot,
		StartStep:    runCase.StartStep,
		EndStep:      runCase.EndStep,
	}
}

func canceledScenarioResult(runCase RunCase, err error) ScenarioResult {
	result := scenarioResultFromCase(runCase, "canceled", err.Error())
	return result
}

func enqueueJobs(ctx context.Context, jobs chan<- indexedRunCase, cases []RunCase) {
	for index, runCase := range cases {
		select {
		case <-ctx.Done():
			return
		case jobs <- indexedRunCase{index: index, runCase: runCase}:
		}
	}
}

func fillCanceledResults(ctx context.Context, results []ScenarioResult, cases []RunCase) {
	err := ctx.Err()
	if err == nil {
		return
	}
	for i := range results {
		if results[i].Scenario != "" || i >= len(cases) {
			continue
		}
		results[i] = canceledScenarioResult(cases[i], err)
	}
}
