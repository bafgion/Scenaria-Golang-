package vanessa

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	PlatformExecutable string   `json:"platform_executable"`
	PlatformMode       string   `json:"platform_mode"`
	PlatformExtraArgs  []string `json:"platform_extra_args"`
	EPFPath            string   `json:"epf_path"`
	IBConnection       string   `json:"ib_connection_string"`
	User               string   `json:"user"`
	Password           string   `json:"password"`
	RunsDir            string   `json:"runs_dir"`
	ProcessTimeoutSec  int      `json:"process_timeout_sec"`
	ReportJUnit        bool     `json:"report_junit"`
	ReportAllure       bool     `json:"report_allure"`
	ProjectBaseParams  string   `json:"project_base_params"`
	DryRunOnly         bool     `json:"dry_run_only"`
	LogEncoding        string   `json:"log_encoding"`
}

func DefaultSettings() Settings {
	return Settings{
		PlatformMode:      "TestManager",
		ProcessTimeoutSec: 3600,
		ReportJUnit:       true,
		ProjectBaseParams: ".scenaria/va-params.base.json",
		LogEncoding:       "auto",
	}
}

func LoadSettings(projectRoot string) (Settings, error) {
	cfg := DefaultSettings()
	if projectRoot != "" {
		path := filepath.Join(projectRoot, ".scenaria", "vanessa.json")
		if err := readJSON(path, &cfg); err == nil {
			return cfg, nil
		}
	}
	globalPath := filepath.Join(os.Getenv("APPDATA"), "Scenaria", "settings.json")
	payload, err := os.ReadFile(globalPath)
	if err != nil {
		return cfg, nil
	}
	var root map[string]any
	if err := json.Unmarshal(payload, &root); err != nil {
		return cfg, nil
	}
	plugins, _ := root["plugins"].(map[string]any)
	vanessa, _ := plugins["vanessa"].(map[string]any)
	if vanessa == nil {
		return cfg, nil
	}
	raw, err := json.Marshal(vanessa)
	if err != nil {
		return cfg, nil
	}
	_ = json.Unmarshal(raw, &cfg)
	return cfg, nil
}

func ValidateSettings(cfg Settings) []string {
	issues := make([]string, 0)
	if info, err := os.Stat(cfg.PlatformExecutable); err != nil || info.IsDir() {
		issues = append(issues, "platform executable not found")
	}
	if info, err := os.Stat(cfg.EPFPath); err != nil || info.IsDir() {
		issues = append(issues, "Vanessa EPF not found")
	}
	return issues
}

func readJSON(path string, dst any) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, dst)
}

func ResolveRunsDir(cfg Settings) string {
	if dir := filepath.Clean(cfg.RunsDir); dir != "" && dir != "." {
		return dir
	}
	base, _ := os.UserConfigDir()
	if base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, "Scenaria", "vanessa-runs")
}

func BaseParamsPath(projectRoot string, cfg Settings) string {
	if projectRoot == "" {
		return ""
	}
	rel := cfg.ProjectBaseParams
	if rel == "" {
		if pc, err := loadProjectVAParamsBase(projectRoot); err == nil && pc != "" {
			rel = pc
		}
	}
	if rel == "" {
		rel = ".scenaria/va-params.base.json"
	}
	return filepath.Join(projectRoot, filepath.FromSlash(rel))
}

func loadProjectVAParamsBase(projectRoot string) (string, error) {
	path := filepath.Join(projectRoot, ".scenaria", "project.json")
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var raw map[string]string
	if err := json.Unmarshal(payload, &raw); err != nil {
		return "", err
	}
	return raw["va_params_base"], nil
}

func SaveProjectSettings(projectRoot string, cfg Settings) error {
	dir := filepath.Join(projectRoot, ".scenaria")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create .scenaria: %w", err)
	}
	payload, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "vanessa.json"), append(payload, '\n'), 0o644)
}
