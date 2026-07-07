package gui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/report/allure"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/version"
)

func (s *Service) runInProcess(ctx context.Context, req RunRequest, emit EventEmitter) (player.ExecutionResult, error) {
	targets := req.Targets
	if len(targets) == 0 {
		if path := s.ProjectPath(); path != "" {
			targets = []string{path}
		}
	}
	if len(targets) == 0 {
		return player.ExecutionResult{}, fmt.Errorf("нет файлов для запуска — откройте сценарий или проект")
	}

	store := scenario.NewFeatureStore()
	files := make([]string, 0)
	for _, target := range targets {
		discovered, err := store.Discover(target)
		if err != nil {
			return player.ExecutionResult{}, err
		}
		files = append(files, discovered...)
	}
	files = dedupePaths(files)
	if len(files) == 0 {
		return player.ExecutionResult{}, fmt.Errorf("no .feature files found in %v", targets)
	}

	featureInputs := make([]player.FeatureInput, 0, len(files))
	for _, path := range files {
		feature, err := store.Load(path)
		if err != nil {
			return player.ExecutionResult{}, fmt.Errorf("%s: %w", path, err)
		}
		if issues := gherkin.ValidateFeature(feature); len(issues) > 0 {
			return player.ExecutionResult{}, fmt.Errorf("%s: %s", path, issues[0].Message)
		}
		featureInputs = append(featureInputs, player.FeatureInput{Path: path, Feature: feature})
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
			return player.ExecutionResult{}, fmt.Errorf("no scenarios found with name %q", req.Scenario)
		}
		if req.Tag != "" {
			return player.ExecutionResult{}, fmt.Errorf("no scenarios found with tag %q", req.Tag)
		}
		return player.ExecutionResult{}, fmt.Errorf("no runnable scenarios found")
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

	root := paths.InferProjectRoot(targets)
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
		return s.finalizeGUIReports(root, req, plan, result, err, false)
	}

	workers := req.Workers
	if workers < 1 {
		workers = 1
	}

	appCfg, _ := settings.LoadDefaultAppSettings()
	httpCreds := player.ResolveRunHTTPCredentials(req.BaseURL, plan, appCfg)
	navWait, err := resolveRunNavWait(paths.InferProjectRoot(targets), appCfg)
	if err != nil {
		return player.ExecutionResult{}, err
	}

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
		CloseAfterRun:     !s.canReuseLiveBrowser(req),
		StepScreenshots:   req.HTMLPath != "" && !req.HTMLLightMode,
	})

	var result player.ExecutionResult
	var runErr error
	if s.canReuseLiveBrowser(req) {
		result, runErr = s.runOnLiveBrowser(ctx, exec, plan, navWait)
	} else {
		runner := player.BrowserRunner{Executor: exec, ParallelWorkers: workers}
		result, runErr = runner.Execute(ctx, plan)
	}
	return s.finalizeGUIReports(root, req, plan, result, runErr, statusIncremental)
}

func (s *Service) finalizeGUIReports(
	root string,
	req RunRequest,
	plan player.ExecutionPlan,
	result player.ExecutionResult,
	runErr error,
	statusIncremental bool,
) (player.ExecutionResult, error) {
	if root != "" {
		req = remapRunArtifacts(root, req)
	}
	reportErr := s.writeGUIReports(root, req, plan, result)
	if reportErr != nil && runErr != nil {
		return result, errors.Join(runErr, reportErr)
	}
	if reportErr != nil {
		return result, reportErr
	}
	if root != "" && !statusIncremental {
		recordGUIStatus(root, req, result)
	}
	return result, runErr
}

func (s *Service) canReuseLiveBrowser(req RunRequest) bool {
	if !req.ReuseLiveBrowser {
		return false
	}
	if req.DryRun || !s.HasLiveBrowser() {
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

func resolveRunNavWait(projectRoot string, appCfg *settings.AppSettings) (string, error) {
	navWait := settings.ResolveNavWaitUntil(projectRoot, appCfg)
	if _, err := player.ParseNavWaitUntil(navWait); err != nil {
		return "", err
	}
	return navWait, nil
}

func (s *Service) runOnLiveBrowser(ctx context.Context, exec *player.PlaywrightExecutor, plan player.ExecutionPlan, navWait string) (player.ExecutionResult, error) {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
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
	result, err := runner.ExecuteSequentialAttached(ctx, exec, plan, browserSession)
	return result, err
}

func (s *Service) formatRunOutput(result player.ExecutionResult, runErr error) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Discovered %d file(s), %d scenario(s), %d step(s) [%s]\n",
		result.Files, result.Scenarios, result.Steps, version.String())
	if runErr != nil {
		fmt.Fprintf(&b, "Run failed: %v\n", runErr)
	}
	return b.String()
}

func dedupePaths(pathsIn []string) []string {
	seen := make(map[string]struct{}, len(pathsIn))
	out := make([]string, 0, len(pathsIn))
	for _, p := range pathsIn {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func appCfgMaxLoops(cfg *settings.AppSettings) int {
	if cfg == nil || cfg.MaxLoopIterations <= 0 {
		return 100
	}
	return cfg.MaxLoopIterations
}

func resolveGUIEngine(req RunRequest, targets []string) string {
	if strings.TrimSpace(req.Engine) != "" {
		return req.Engine
	}
	root := paths.InferProjectRoot(targets)
	if root != "" {
		if cfg, err := settings.LoadProjectConfig(root); err == nil && cfg.DefaultRunner != "" {
			return cfg.DefaultRunner
		}
	}
	return "playwright"
}

func remapRunArtifacts(root string, req RunRequest) RunRequest {
	req.HTMLPath = paths.RemapScenariaArtifact(root, req.HTMLPath)
	req.JUnitPath = paths.RemapScenariaArtifact(root, req.JUnitPath)
	req.SummaryJSON = paths.RemapScenariaArtifact(root, req.SummaryJSON)
	req.AllureDir = paths.RemapScenariaArtifact(root, req.AllureDir)
	req.TraceDir = paths.RemapScenariaArtifact(root, req.TraceDir)
	req.VideoDir = paths.RemapScenariaArtifact(root, req.VideoDir)
	return req
}

func (s *Service) writeGUIReports(projectRoot string, req RunRequest, plan player.ExecutionPlan, result player.ExecutionResult) error {
	var prevSummary *report.RunSummaryDetailed
	if req.SummaryJSON != "" {
		prevSummary = report.ReadPreviousSummary(req.SummaryJSON)
	}
	if req.SummaryJSON != "" {
		if err := report.WriteRunSummaryDetailed(req.SummaryJSON, report.FromExecutionResultDetailed(result)); err != nil {
			return err
		}
	}
	if req.JUnitPath != "" {
		if err := report.WriteJUnit(req.JUnitPath, result); err != nil {
			return err
		}
	}
	if req.HTMLPath != "" {
		bridgeURL, bridgeToken := s.ReportBridgeCredentials()
		if err := report.WriteHTML(req.HTMLPath, result, report.HTMLOptions{
			Plan:            plan,
			ProjectRoot:     projectRoot,
			LightMode:       req.HTMLLightMode,
			BridgeURL:       bridgeURL,
			BridgeToken:     bridgeToken,
			Locale:          req.ReportLocale,
			ReportDir:       filepath.Dir(req.HTMLPath),
			PreviousSummary: prevSummary,
		}); err != nil {
			return err
		}
	}
	if req.AllureDir != "" {
		if err := allure.WriteResults(req.AllureDir, result); err != nil {
			return err
		}
	}
	return nil
}

func recordGUIStatus(root string, req RunRequest, result player.ExecutionResult) {
	store, err := runstatus.Open(root)
	if err != nil {
		return
	}
	engine := resolveGUIEngine(req, nil)
	for _, scenarioResult := range result.ScenarioResults {
		_ = store.Record(player.RunstatusEntry(scenarioResult, engine))
	}
}

func scenarioResultsToEntries(results []player.ScenarioResult, runner string) []RunResultEntry {
	if len(results) == 0 {
		return nil
	}
	at := time.Now().UTC().Format(time.RFC3339)
	out := make([]RunResultEntry, 0, len(results))
	for _, r := range results {
		success := r.Status == "passed" || r.Status == "dry-run"
		out = append(out, RunResultEntry{
			Path:       r.FeaturePath + "::" + r.Scenario,
			Success:    success,
			Message:    r.Message,
			Runner:     runner,
			At:         at,
			FailedStep: r.FailedStep,
		})
	}
	return out
}
