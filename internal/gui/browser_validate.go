package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

func (s *Service) ValidateBrowser(req ValidateRequest) ([]ValidationIssue, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	if s.validateCancel != nil {
		s.validateCancel()
	}
	s.validateGen++
	myGen := s.validateGen
	s.validateCancel = cancel
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		if s.validateGen == myGen {
			s.validateCancel = nil
		}
		s.mu.Unlock()
		cancel()
	}()
	return s.ValidateBrowserContext(ctx, req)
}

func (s *Service) ValidateBrowserContext(ctx context.Context, req ValidateRequest) ([]ValidationIssue, error) {
	path := s.ProjectPath()
	if path == "" {
		return nil, fmt.Errorf("open a project folder first")
	}
	targets := req.Targets
	if len(targets) == 0 {
		targets = []string{path}
	}
	store := scenario.NewFeatureStore()
	validator := selector.Validator{Headless: true}
	appCfg, _ := s.loadAppSettings()
	baseURL := ""
	if root := paths.InferProjectRoot(targets); root != "" {
		if cfg, err := settings.LoadProjectConfig(root); err == nil {
			baseURL = cfg.BaseURL
		}
	}
	browserName := strings.TrimSpace(req.Browser)
	if browserName == "" {
		browserName = "chromium"
	}
	if appCfg != nil && strings.TrimSpace(appCfg.Browser) != "" {
		browserName = appCfg.Browser
	}
	headless := true
	if appCfg != nil {
		headless = appCfg.Headless
	}

	out := make([]ValidationIssue, 0)
	type featureSnapshot struct {
		path    string
		feature *gherkin.Feature
		loadErr error
	}
	snapshots := make([]featureSnapshot, 0, len(targets))
	if err := s.withProjectFSReadLock(func() error {
		for _, featurePath := range targets {
			feature, err := store.Load(featurePath)
			snapshots = append(snapshots, featureSnapshot{path: featurePath, feature: feature, loadErr: err})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	for _, snap := range snapshots {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		if snap.loadErr != nil {
			out = append(out, ValidationIssue{Line: 1, Message: snap.loadErr.Error(), Status: "missing"})
			continue
		}
		fileIssues := make([]ValidationIssue, 0)
		for _, issue := range gherkin.ValidateFeature(snap.feature) {
			fileIssues = append(fileIssues, ValidationIssue{
				Line:    issue.Line,
				Message: issue.Message,
				Status:  "missing",
			})
		}
		syntaxIssues, err := validator.ValidateFeature(snap.path, snap.feature)
		if err != nil {
			fileIssues = append(fileIssues, ValidationIssue{Line: 1, Message: err.Error(), Status: "missing"})
		} else {
			for _, issue := range syntaxIssues {
				fileIssues = append(fileIssues, ValidationIssue{
					Line:     issue.Line,
					Message:  issue.Message,
					Selector: issue.Selector,
					Status:   "missing",
				})
			}
		}
		if len(fileIssues) > 0 {
			out = append(out, fileIssues...)
			continue
		}
		stepResults, err := validator.ValidateFeatureInBrowserDetailed(ctx, snap.path, snap.feature, selector.BrowserValidateOptions{
			BrowserName: browserName,
			Headless:    headless,
			BaseURL:     baseURL,
		})
		if err != nil {
			out = append(out, ValidationIssue{Line: 1, Message: err.Error(), Status: "missing"})
			continue
		}
		for _, step := range stepResults {
			if step.Status == "skipped" {
				continue
			}
			out = append(out, ValidationIssue{
				Line:     step.Line,
				Message:  step.Message,
				Selector: step.Selector,
				Status:   step.Status,
				StepText: step.StepText,
			})
		}
	}
	return out, nil
}
