package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectConfig struct {
	DefaultRunner      string `json:"default_runner"`
	FeaturesRoot       string `json:"features_root"`
	BaseURL            string `json:"base_url"`
	VAParamsBase       string `json:"va_params_base"`
	NavWaitUntil       string `json:"nav_wait_until,omitempty"`
	HTMLReportOpenMode string `json:"html_report_open_mode,omitempty"` // "full" or "light"
}

func DefaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		DefaultRunner: "playwright",
		FeaturesRoot:  "features",
		VAParamsBase:  ".scenaria/va-params.base.json",
	}
}

func ProjectConfigPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".scenaria", "project.json")
}

func LoadProjectConfig(projectRoot string) (ProjectConfig, error) {
	cfg := DefaultProjectConfig()
	if projectRoot == "" {
		return cfg, nil
	}
	path := ProjectConfigPath(projectRoot)
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read project config: %w", err)
	}
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return cfg, fmt.Errorf("decode project config: %w", err)
	}
	if cfg.DefaultRunner == "" || cfg.DefaultRunner == "ask" {
		cfg.DefaultRunner = "playwright"
	}
	if cfg.FeaturesRoot == "" {
		cfg.FeaturesRoot = "features"
	}
	if cfg.VAParamsBase == "" {
		cfg.VAParamsBase = ".scenaria/va-params.base.json"
	}
	cfg.HTMLReportOpenMode = NormalizeHTMLReportOpenMode(cfg.HTMLReportOpenMode)
	return cfg, nil
}

// NormalizeHTMLReportOpenMode returns "full" or "light".
func NormalizeHTMLReportOpenMode(mode string) string {
	if strings.EqualFold(strings.TrimSpace(mode), "light") {
		return "light"
	}
	return "full"
}

func SaveProjectConfig(projectRoot string, cfg ProjectConfig) error {
	dir := filepath.Join(projectRoot, ".scenaria")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return writeJSON(ProjectConfigPath(projectRoot), cfg)
}
