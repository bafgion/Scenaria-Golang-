package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/version"
)

func (s *Service) runExecutionHost() RunExecutionHost {
	return RunExecutionHost{
		ProjectPath:     s.ProjectPath,
		CurrentSession:  s.CurrentRunSession,
		PrepareSnapshot: s.prepareRunFeatureSnapshot,
		LoadAppSettings: s.loadAppSettings,
		HasLiveBrowser:  s.HasLiveBrowser,
		LiveSession:     s.recorderOps().LiveSession,
		Reports:         s.reporter(),
	}
}

func (s *Service) runInProcess(ctx context.Context, req RunRequest, emit EventEmitter) (player.ExecutionResult, report.RunArtifactLayout, error) {
	return s.runner().ExecuteInProcess(ctx, req, emit, s.runExecutionHost())
}

func (s *Service) prepareRunFeatureSnapshot(targets []string) ([]player.FeatureInput, error) {
	store := scenario.NewFeatureStore()
	files := make([]string, 0)
	featureInputs := make([]player.FeatureInput, 0, len(files))
	if err := s.withProjectFSReadLock(func() error {
		for _, target := range targets {
			discovered, err := store.Discover(target)
			if err != nil {
				return err
			}
			files = append(files, discovered...)
		}
		files = dedupePaths(files)
		if len(files) == 0 {
			return fmt.Errorf("no .feature files found in %v", targets)
		}
		featureInputs = make([]player.FeatureInput, 0, len(files))
		for _, path := range files {
			feature, err := store.Load(path)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			if issues := gherkin.ValidateFeature(feature); len(issues) > 0 {
				return fmt.Errorf("%s: %s", path, issues[0].Message)
			}
			featureInputs = append(featureInputs, player.FeatureInput{Path: path, Feature: feature})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return featureInputs, nil
}

func (s *Service) finalizeGUIReports(
	root string,
	req RunRequest,
	plan player.ExecutionPlan,
	result player.ExecutionResult,
	runErr error,
	statusIncremental bool,
) (player.ExecutionResult, report.RunArtifactLayout, error) {
	return s.reporter().FinalizeRunReports(root, req, plan, result, runErr, statusIncremental)
}

func (s *Service) canReuseLiveBrowser(req RunRequest) bool {
	return s.runner().CanReuseLiveBrowser(req, s.HasLiveBrowser())
}

func resolveRunNavWait(projectRoot string, appCfg *settings.AppSettings) (string, error) {
	navWait := settings.ResolveNavWaitUntil(projectRoot, appCfg)
	if _, err := player.ParseNavWaitUntil(navWait); err != nil {
		return "", err
	}
	return navWait, nil
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

func scenarioResultsToEntries(results []player.ScenarioResult, runner string) []RunResultEntry {
	if len(results) == 0 {
		return nil
	}
	at := time.Now().UTC().Format(time.RFC3339)
	out := make([]RunResultEntry, 0, len(results))
	for _, r := range results {
		success := player.ResultStatusIsSuccessful(r.Status)
		out = append(out, RunResultEntry{
			Path:       r.FeaturePath + "::" + r.Scenario,
			Success:    success,
			Status:     r.Status,
			Message:    r.Message,
			Runner:     runner,
			At:         at,
			FailedStep: r.FailedStep,
		})
	}
	return out
}
