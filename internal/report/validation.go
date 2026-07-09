package report

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type htmlStepValidation struct {
	Status     string `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	Mode       string `json:"mode,omitempty"`
	ActionKind string `json:"action_kind,omitempty"`
	MatchCount int    `json:"match_count,omitempty"`
	Limitation string `json:"limitation,omitempty"`
}

type ValidationCollectOptions struct {
	ProjectRoot  string
	FeaturePaths []string
	BrowserName  string
	Headless     bool
	BaseURL      string
	Mode         string
}

func stepValidationFromSelector(step selector.StepValidation) htmlStepValidation {
	return htmlStepValidation{
		Status:     step.Status,
		Message:    step.Message,
		Mode:       step.Mode,
		ActionKind: step.ActionKind,
		MatchCount: step.MatchCount,
		Limitation: step.Limitation,
	}
}

func uniqueFeaturePaths(result player.ExecutionResult) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, sr := range result.ScenarioResults {
		p := normalizeFeaturePath(sr.FeaturePath)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func normalizeFeaturePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	clean := filepath.Clean(path)
	if clean == "." {
		return path
	}
	return clean
}

func pathsMatchFeature(stored, featurePath string) bool {
	a := normalizeFeaturePath(stored)
	b := normalizeFeaturePath(featurePath)
	if a == "" || b == "" {
		return false
	}
	if a == b {
		return true
	}
	return filepath.Base(a) == filepath.Base(b) && filepath.Base(a) != "" && filepath.Base(a) != "."
}

func CollectStepValidations(ctx context.Context, opts ValidationCollectOptions) (map[string]map[int]htmlStepValidation, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	out := map[string]map[int]htmlStepValidation{}
	if len(opts.FeaturePaths) == 0 {
		return out, nil
	}
	baseURL := strings.TrimSpace(opts.BaseURL)
	if baseURL == "" && strings.TrimSpace(opts.ProjectRoot) != "" {
		if cfg, err := settings.LoadProjectConfig(opts.ProjectRoot); err == nil {
			baseURL = strings.TrimSpace(cfg.BaseURL)
		}
	}
	browserName := strings.TrimSpace(opts.BrowserName)
	if browserName == "" {
		browserName = "chromium"
	}
	mode := selector.NormalizeValidationMode(opts.Mode)
	validator := selector.Validator{Headless: opts.Headless}
	store := scenario.NewFeatureStore()
	browserOpts := selector.BrowserValidateOptions{
		BrowserName: browserName,
		Headless:    opts.Headless,
		BaseURL:     baseURL,
		Mode:        mode,
	}

	for _, featurePath := range opts.FeaturePaths {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		featurePath = normalizeFeaturePath(featurePath)
		if featurePath == "" {
			continue
		}
		feature, err := store.Load(featurePath)
		if err != nil {
			continue
		}
		if len(gherkin.ValidateFeature(feature)) > 0 {
			continue
		}
		if syntaxIssues, err := validator.ValidateFeature(featurePath, feature); err != nil || len(syntaxIssues) > 0 {
			continue
		}
		steps, err := validator.ValidateFeatureInBrowserDetailed(ctx, featurePath, feature, browserOpts)
		if err != nil {
			continue
		}
		byLine := map[int]htmlStepValidation{}
		for _, step := range steps {
			if step.Status == "skipped" {
				continue
			}
			byLine[step.Line] = stepValidationFromSelector(step)
		}
		if len(byLine) > 0 {
			out[featurePath] = byLine
		}
	}
	return out, nil
}

func applyStepValidations(sc *htmlScenario, byLine map[int]htmlStepValidation) {
	if sc == nil || len(byLine) == 0 {
		return
	}
	for i := range sc.Steps {
		line := sc.Steps[i].Line
		if line <= 0 {
			continue
		}
		if v, ok := byLine[line]; ok {
			copied := v
			sc.Steps[i].Validation = &copied
			if sc.ValidationLimitation == "" && copied.Limitation != "" {
				sc.ValidationLimitation = copied.Limitation
				sc.ValidationMode = copied.Mode
			}
		}
	}
}

func resolveStepValidations(result player.ExecutionResult, opts HTMLOptions) map[string]map[int]htmlStepValidation {
	if len(opts.StepValidations) > 0 {
		return opts.StepValidations
	}
	if opts.SkipStepValidation || result.Mode != "browser" {
		return nil
	}
	paths := uniqueFeaturePaths(result)
	if len(paths) == 0 {
		return nil
	}
	baseURL := strings.TrimSpace(opts.ValidationBaseURL)
	if baseURL == "" && strings.TrimSpace(opts.ProjectRoot) != "" {
		if cfg, err := settings.LoadProjectConfig(opts.ProjectRoot); err == nil {
			baseURL = strings.TrimSpace(cfg.BaseURL)
		}
	}
	collected, err := CollectStepValidations(context.Background(), ValidationCollectOptions{
		ProjectRoot:  opts.ProjectRoot,
		FeaturePaths: paths,
		BrowserName:  opts.ValidationBrowser,
		Headless:     opts.ValidationHeadless,
		BaseURL:      baseURL,
		Mode:         opts.ValidationMode,
	})
	if err != nil {
		return nil
	}
	return collected
}
