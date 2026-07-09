package gui

import (
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

func (c *CLIOps) ExportFeature(req ExportRequest) RunResult {
	input := strings.TrimSpace(req.InputPath)
	if input == "" {
		return RunResult{Error: "feature path is required"}
	}
	args := []string{input, "--output", req.Output, "--format", req.Format}
	if req.BaseURL != "" {
		args = append(args, "--base-url", req.BaseURL)
	}
	if req.Force {
		args = append(args, "--force")
	}
	out, err := c.Export(args)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (c *CLIOps) ImportFeatureJSON(req ImportRequest) RunResult {
	args := []string{req.JSONPath, "--output", req.OutputPath}
	if req.Force {
		args = append(args, "--force")
	}
	out, err := c.ImportJSON(args)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (c *CLIOps) RecordBaselineFeature(projectPath string, req BaselineRecordRequest) RunResult {
	path := strings.TrimSpace(projectPath)
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	output := strings.TrimSpace(req.Output)
	if output == "" {
		output = filepath.Join(path, "recorded.feature")
	} else {
		confined, err := paths.ConfineToProjectRoot(path, output)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		output = confined
	}
	featureName := strings.TrimSpace(req.FeatureName)
	if featureName == "" {
		featureName = "Записанный сценарий"
	}
	scenarioName := strings.TrimSpace(req.ScenarioName)
	if scenarioName == "" {
		scenarioName = "Базовый сценарий"
	}
	args := []string{
		"--output", output,
		"--feature", featureName,
		"--scenario", scenarioName,
	}
	for _, step := range req.Steps {
		step = strings.TrimSpace(step)
		if step != "" {
			args = append(args, "--step", step)
		}
	}
	out, err := c.Record(args)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}
