package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AppSettings struct {
	Browser             string `json:"browser"`
	Headless            bool   `json:"headless"`
	RecordingHoverMode  bool   `json:"recording_hover_mode"`
	RecordingFilterMode bool   `json:"recording_filter_mode"`
	NavOnlyRecording    bool   `json:"nav_only_recording"`
	ParallelWorkers     int    `json:"parallel_workers"`
	SlowMo              int    `json:"slow_mo"`
	MaxLoopIterations   int    `json:"max_loop_iterations"`
	ToolbarCompact      bool   `json:"toolbar_compact"`
	StepsPanelVisible   bool     `json:"steps_panel_visible"`
	StepsPanelHeight    int      `json:"steps_panel_height"`
	SidebarWidth        int      `json:"sidebar_width"`
	RecentProjects      []string `json:"recent_projects"`
	RecentFeatures      []string `json:"recent_features"`
	SessionProject      string   `json:"session_project,omitempty"`
	OpenTabs            []string `json:"open_tabs,omitempty"`
	ActiveTab           string   `json:"active_tab,omitempty"`
	UntitledTabs        []UntitledTabSession `json:"untitled_tabs,omitempty"`
	ScrollBeforeClick   bool     `json:"scroll_before_click"`
	HoverRecordMinMs    int      `json:"hover_record_min_ms"`
	SelectorClickStrategies []string `json:"selector_click_strategies,omitempty"`
	SelectorInputStrategies []string `json:"selector_input_strategies,omitempty"`
	CheckUpdatesOnStartup *bool  `json:"check_updates_on_startup,omitempty"`
	HTTPAuth            map[string]HTTPAuthEntry `json:"http_auth,omitempty"`
	Editor              EditorSettings             `json:"editor,omitempty"`
}

// UntitledTabSession stores in-memory editor tab state across app restarts.
type UntitledTabSession struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type HTTPAuthEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TestClient struct {
	Name         string            `json:"name"`
	BaseURL      string            `json:"base_url"`
	Cookies      []Cookie          `json:"cookies"`
	LocalStorage map[string]string `json:"local_storage"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	HTTPOnly bool   `json:"http_only"`
	Secure   bool   `json:"secure"`
}

func DefaultAppSettingsPath() string {
	if base := os.Getenv("APPDATA"); base != "" {
		return filepath.Join(base, "Scenaria", "settings.json")
	}
	config, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(config, "Scenaria", "settings.json")
}

func LoadDefaultAppSettings() (*AppSettings, error) {
	path := DefaultAppSettingsPath()
	if path == "" {
		return &AppSettings{Browser: "chromium"}, nil
	}
	cfg, err := LoadAppSettings(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AppSettings{Browser: "chromium"}, nil
		}
		return nil, err
	}
	if strings.TrimSpace(cfg.Browser) == "" {
		cfg.Browser = "chromium"
	}
	return cfg, nil
}

func CheckUpdatesOnStartupEnabled(cfg *AppSettings) bool {
	if cfg == nil || cfg.CheckUpdatesOnStartup == nil {
		return true
	}
	return *cfg.CheckUpdatesOnStartup
}

func LoadAppSettings(path string) (*AppSettings, error) {
	var cfg AppSettings
	if err := readJSON(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveAppSettings(path string, cfg *AppSettings) error {
	return writeJSON(path, cfg)
}

func LoadTestClient(path string) (*TestClient, error) {
	var client TestClient
	if err := readJSON(path, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

func SaveTestClient(path string, client *TestClient) error {
	return writeJSON(path, client)
}

func TestClientPath(projectRoot, name string) (string, error) {
	if projectRoot == "" || name == "" {
		return "", fmt.Errorf("project root and test client name are required")
	}
	return filepath.Join(projectRoot, ".scenaria", "test_clients", name+".json"), nil
}

func readJSON(path string, dst any) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read json file %q: %w", path, err)
	}
	if err := json.Unmarshal(payload, dst); err != nil {
		return fmt.Errorf("decode json file %q: %w", path, err)
	}
	return nil
}

func writeJSON(path string, src any) error {
	payload, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return fmt.Errorf("encode json file %q: %w", path, err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write json file %q: %w", path, err)
	}
	return nil
}
