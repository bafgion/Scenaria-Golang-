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
	appCfg, _ := settings.LoadDefaultAppSettings()
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
	for _, featurePath := range targets {
		feature, err := store.Load(featurePath)
		if err != nil {
			out = append(out, ValidationIssue{Line: 1, Message: err.Error(), Status: "missing"})
			continue
		}
		fileIssues := make([]ValidationIssue, 0)
		for _, issue := range gherkin.ValidateFeature(feature) {
			fileIssues = append(fileIssues, ValidationIssue{
				Line:    issue.Line,
				Message: issue.Message,
				Status:  "missing",
			})
		}
		syntaxIssues, err := validator.ValidateFeature(featurePath, feature)
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
		stepResults, err := validator.ValidateFeatureInBrowserDetailed(context.Background(), featurePath, feature, selector.BrowserValidateOptions{
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
