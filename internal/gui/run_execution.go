package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

// RunExecutionHost supplies project-scoped dependencies for in-process run orchestration.
type RunExecutionHost struct {
	ProjectPath       func() string
	CurrentSession    func() *RunSession
	PrepareSnapshot   func(targets []string) ([]player.FeatureInput, error)
	LoadAppSettings   func() (*settings.AppSettings, error)
	HasLiveBrowser    func() bool
	LiveSession       func() *recorder.LiveSession
	Reports           *ReportService
}

func (s *RunService) CanReuseLiveBrowser(req RunRequest, hasLiveBrowser bool, plan player.ExecutionPlan) bool {
	if !req.ReuseLiveBrowser {
		return false
	}
	if req.DryRun || !hasLiveBrowser {
		return false
	}
	if len(plan.Cases) > 1 {
		return false
	}
	if player.PlanContainsCloseBrowser(plan) {
		return false
	}
	workers := req.Workers
	if workers < 1 {
		workers = 1
	}
	if workers > 1 {
		return false
	}
	if strings.TrimSpace(req.TraceDir) != "" || strings.TrimSpace(req.VideoDir) != "" {
		return false
	}
	browser := strings.ToLower(strings.TrimSpace(req.Browser))
	if browser != "" && browser != "chromium" {
		return false
	}
	return true
}

func (s *RunService) RunOnLiveBrowser(
	ctx context.Context,
	exec *player.PlaywrightExecutor,
	plan player.ExecutionPlan,
	navWait string,
	session *recorder.LiveSession,
) (player.ExecutionResult, error) {
	if session == nil || !session.BrowserAlive() {
		return player.ExecutionResult{}, fmt.Errorf("браузер не открыт")
	}
	page, ok := session.ActivePage()
	if !ok {
		return player.ExecutionResult{}, fmt.Errorf("браузер не открыт")
	}

	wasRecording := session.CaptureEnabled()
	wasPaused := session.IsPaused()
	if wasRecording && !wasPaused {
		session.Pause()
	}
	releaseTestHold := session.HoldForTestRun()
	defer releaseTestHold()
	defer func() {
		if wasRecording && !wasPaused {
			session.Resume()
		}
	}()

	browserSession, err := player.AttachToPage(page, navWait)
	if err != nil {
		return player.ExecutionResult{}, err
	}

	runner := player.BrowserRunner{Executor: exec}
	return runner.ExecuteSequentialAttached(ctx, exec, plan, browserSession)
}

func (s *RunService) ExecuteInProcess(
	ctx context.Context,
	req RunRequest,
	emit EventEmitter,
	host RunExecutionHost,
) (player.ExecutionResult, report.RunArtifactLayout, error) {
	var artifactLayout report.RunArtifactLayout
	targets := req.Targets
	if len(targets) == 0 {
		if path := host.ProjectPath(); path != "" {
			targets = []string{path}
		}
	}
	if len(targets) == 0 {
		return player.ExecutionResult{}, artifactLayout, fmt.Errorf("no files to run")
	}

	featureInputs, err := host.PrepareSnapshot(targets)
	if err != nil {
		return player.ExecutionResult{}, artifactLayout, err
	}

	plan := player.BuildExecutionPlanWithTestClient(featureInputs, req.Tag, req.Scenario, req.Vars, req.TestClient)
	if req.StartStep >= 0 || req.EndStep >= 0 {
		for i := range plan.Cases {
			if req.StartStep >= 0 {
				plan.Cases[i].StartStep = req.StartStep
			}
			if req.EndStep >= 0 {
				plan.Cases[i].EndStep = req.EndStep
			}
		}
	}
	if len(plan.Cases) == 0 {
		if req.Scenario != "" {
			return player.ExecutionResult{}, artifactLayout, fmt.Errorf("no scenarios found with name %q", req.Scenario)
		}
		if req.Tag != "" {
			return player.ExecutionResult{}, artifactLayout, fmt.Errorf("no scenarios found with tag %q", req.Tag)
		}
		return player.ExecutionResult{}, artifactLayout, fmt.Errorf("no runnable scenarios found")
	}

	ctx = player.WithRunProgress(ctx, func(ev player.RunProgressEvent) {
		if emit != nil {
			emit("run-progress", ev)
			if ev.Phase == player.ProgressScenarioDone {
				emit("run-results-changed", nil)
			}
		}
	})
	ctx = player.WithContinueOnFail(ctx, req.ContinueOnFail)
	if host.CurrentSession != nil {
		if session := host.CurrentSession(); session != nil {
			ctx = player.WithRunID(ctx, session.RunID)
		}
	}

	root := paths.InferProjectRoot(targets)
	reports := host.Reports
	if root != "" && reports != nil {
		req, artifactLayout, err = reports.LayoutRunRequestArtifacts(root, req)
		if err != nil {
			return player.ExecutionResult{}, artifactLayout, err
		}
	}
	statusIncremental := false
	if root != "" && !req.DryRun {
		if st, err := runstatus.Open(root); err == nil {
			engine := resolveGUIEngine(req, targets)
			ctx = player.WithRunStatusHook(ctx, st, engine)
			statusIncremental = true
		}
	}

	if req.DryRun {
		runner := player.DryRunner{}
		result, err := runner.Execute(ctx, plan)
		if reports == nil {
			return result, artifactLayout, err
		}
		return reports.FinalizeRunReports(root, req, plan, result, err, false)
	}

	workers := req.Workers
	if workers < 1 {
		workers = 1
	}

	appCfg, _ := host.LoadAppSettings()
	httpCreds := player.ResolveRunHTTPCredentials(req.BaseURL, plan, appCfg)
	navWait, err := resolveRunNavWait(paths.InferProjectRoot(targets), appCfg)
	if err != nil {
		return player.ExecutionResult{}, artifactLayout, err
	}

	hasLiveBrowser := host.HasLiveBrowser != nil && host.HasLiveBrowser()
	reuseLive := s.CanReuseLiveBrowser(req, hasLiveBrowser, plan)

	exec := player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
		BrowserName:       firstNonEmpty(req.Browser, "chromium"),
		Headless:          !req.Headed,
		BaseURL:           req.BaseURL,
		AutoInstall:       req.InstallPW,
		SlowMo:            float64(req.SlowMo),
		TraceDir:          req.TraceDir,
		VideoDir:          req.VideoDir,
		HTTPCredentials:   httpCreds,
		MaxLoopIterations: appCfgMaxLoops(appCfg),
		NavWaitUntil:      navWait,
		CloseAfterRun:     !reuseLive,
		StepScreenshots:   req.HTMLPath != "" && !req.HTMLLightMode,
	})

	var result player.ExecutionResult
	var runErr error
	if reuseLive {
		var session *recorder.LiveSession
		if host.LiveSession != nil {
			session = host.LiveSession()
		}
		result, runErr = s.RunOnLiveBrowser(ctx, exec, plan, navWait, session)
	} else {
		runner := player.BrowserRunner{Executor: exec, ParallelWorkers: workers}
		result, runErr = runner.Execute(ctx, plan)
	}
	if reports == nil {
		return result, artifactLayout, runErr
	}
	return reports.FinalizeRunReports(root, req, plan, result, runErr, statusIncremental)
}
