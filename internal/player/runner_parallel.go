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

	if pwExec, ok := r.Executor.(*PlaywrightExecutor); ok && poolEligibleForPlan(pwExec.options, plan) {
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
					finalizeParallelScenario(runCtx, &mu, results, i, len(plan.Cases), rc, failed)
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
					runResult := terminalScenarioResult(rc, ScenarioResult{}, err)
					finalizeParallelScenario(runCtx, &mu, results, i, len(plan.Cases), rc, runResult)
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						failFastParallelCancel(runCtx, cancel, pool)
					}
					mu.Unlock()
					continue
				}
				runResult, err := exec.ExecuteScenarioOnSession(runCtx, slot.session, scenarioInputFromCase(rc))
				pool.release(runCtx, slot, exec.options, runResult, rc, runID)
				runResult = terminalScenarioResult(rc, runResult, err)

				mu.Lock()
				runErr := terminalScenarioError(runResult, err)
				if runErr != nil {
					if firstErr == nil {
						firstErr = runErr
						failFastParallelCancel(runCtx, cancel, pool)
					}
				}
				mu.Unlock()
				finalizeParallelScenario(runCtx, &mu, results, i, len(plan.Cases), rc, runResult)
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
					finalizeParallelScenario(runCtx, &mu, results, i, len(plan.Cases), rc, failed)
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
				runResult = terminalScenarioResult(rc, runResult, err)
				mu.Lock()
				runErr := terminalScenarioError(runResult, err)
				if runErr != nil {
					if firstErr == nil {
						firstErr = runErr
						failFastParallelCancel(runCtx, cancel, nil)
					}
				}
				mu.Unlock()
				finalizeParallelScenario(runCtx, &mu, results, i, len(plan.Cases), rc, runResult)
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
	var priorClosedBrowser bool
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

		if i == 0 {
			if session == nil || !session.alive() {
				if attached != nil {
					runResult := ScenarioResult{
						FeaturePath: runCase.FeaturePath,
						Scenario:    runCase.Name,
						Status:      "failed",
						Message:     MsgBrowserClosed,
						FailedStep:  failedStepIndex(0),
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
		} else if priorClosedBrowser {
			if attached != nil {
				appendNotStartedScenario(ctx, &result, runCase, i+1, total)
				continue
			}
			if err := openSession(); err != nil {
				return result, err
			}
		} else if session != nil && session.alive() {
			if err := session.resetForScenario(); err != nil {
				if attached != nil {
					appendNotStartedScenario(ctx, &result, runCase, i+1, total)
					continue
				}
				if IsBrowserSessionClosed(err) {
					if err := openSession(); err != nil {
						return result, err
					}
				} else {
					runResult := ScenarioResult{
						FeaturePath: runCase.FeaturePath,
						Scenario:    runCase.Name,
						Status:      "failed",
						Message:     UserFacingBrowserError(err),
						FailedStep:  failedStepIndex(0),
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
		} else if attached != nil {
			appendNotStartedScenario(ctx, &result, runCase, i+1, total)
			continue
		} else if err := openSession(); err != nil {
			return result, err
		}

		runResult, err := exec.ExecuteScenarioOnSession(ctx, session, scenarioInputFromCase(runCase))
		if err != nil {
			if runResult.Scenario == "" {
				runResult = ScenarioResult{
					FeaturePath: runCase.FeaturePath,
					Scenario:    runCase.Name,
					Status:      "failed",
					Message:     UserFacingBrowserError(err),
					FailedStep:  failedStepIndex(0),
				}
			} else if runResult.Status == "" {
				runResult.Status = "failed"
				if runResult.Message == "" {
					runResult.Message = UserFacingBrowserError(err)
				}
			}
			runResult.Message = NormalizeScenarioMessage(runResult.Message)
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
			priorClosedBrowser = ScenarioBlocksSessionReuse(runResult, runCase, session)
			continue
		}
		if runResult.Status == "failed" {
			runResult.Message = NormalizeScenarioMessage(runResult.Message)
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
			priorClosedBrowser = ScenarioBlocksSessionReuse(runResult, runCase, session)
			continue
		}
		result.ScenarioResults = append(result.ScenarioResults, runResult)
		recordScenarioRunStatus(ctx, runResult)
		emitRunProgress(ctx, RunProgressEvent{
			Phase: ProgressScenarioDone, Index: i + 1, Total: total,
			FeaturePath: runCase.FeaturePath, Scenario: runCase.Name, Success: true,
		})
		priorClosedBrowser = ScenarioBlocksSessionReuse(runResult, runCase, session)
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

func terminalScenarioResult(runCase RunCase, runResult ScenarioResult, err error) ScenarioResult {
	if runResult.FeaturePath == "" {
		runResult.FeaturePath = runCase.FeaturePath
	}
	if runResult.Scenario == "" {
		runResult.Scenario = runCase.Name
	}
	if runResult.CaseID == "" {
		runResult.CaseID = runCase.CaseID
	}
	if runResult.ExampleIndex == 0 {
		runResult.ExampleIndex = runCase.ExampleIndex
	}
	if err != nil {
		if runResult.Status == "" {
			runResult.Status = "failed"
		}
		if runResult.Message == "" {
			runResult.Message = UserFacingBrowserError(err)
		}
		if IsBrowserSessionClosed(err) && runResult.FailedStep == nil {
			runResult.FailedStep = failedStepIndex(0)
		}
	} else if runResult.Status == "" {
		runResult.Status = "passed"
	}
	runResult.Message = NormalizeScenarioMessage(runResult.Message)
	return runResult
}

func terminalScenarioError(runResult ScenarioResult, err error) error {
	if err != nil {
		return err
	}
	if runResult.Status == "failed" {
		return fmt.Errorf("scenario %q failed: %s", runResult.Scenario, runResult.Message)
	}
	return nil
}

func finalizeParallelScenario(
	ctx context.Context,
	mu *sync.Mutex,
	results []ScenarioResult,
	index, total int,
	runCase RunCase,
	runResult ScenarioResult,
) {
	runResult = terminalScenarioResult(runCase, runResult, nil)
	if index >= 0 && index < len(results) {
		if mu != nil {
			mu.Lock()
		}
		results[index] = runResult
		if mu != nil {
			mu.Unlock()
		}
	}
	recordScenarioRunStatus(ctx, runResult)
	emitRunProgress(ctx, RunProgressEvent{
		Phase:       ProgressScenarioDone,
		Index:       index + 1,
		Total:       total,
		CaseID:      runResult.CaseID,
		FeaturePath: runResult.FeaturePath,
		Scenario:    runResult.Scenario,
		Success:     runResult.Status == "passed",
		Message:     runResult.Message,
	})
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
	var mu sync.Mutex
	for i := range results {
		if results[i].Scenario != "" || i >= len(cases) {
			continue
		}
		finalizeParallelScenario(ctx, &mu, results, i, len(cases), cases[i], canceledScenarioResult(cases[i], err))
	}
}
