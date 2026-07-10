package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

type AppSettings struct {
	Browser                 string                   `json:"browser"`
	Headless                bool                     `json:"headless"`
	RecordingHoverMode      bool                     `json:"recording_hover_mode"`
	RecordingFilterMode     bool                     `json:"recording_filter_mode"`
	NavOnlyRecording        bool                     `json:"nav_only_recording"`
	DisableRecordURLWait    bool                     `json:"disable_record_url_wait,omitempty"`
	ParallelWorkers         int                      `json:"parallel_workers"`
	SlowMo                  int                      `json:"slow_mo"`
	MaxLoopIterations       int                      `json:"max_loop_iterations"`
	NavWaitUntil            string                   `json:"nav_wait_until,omitempty"`
	ToolbarCompact          bool                     `json:"toolbar_compact"`
	StepsPanelVisible       bool                     `json:"steps_panel_visible"`
	StepsPanelHeight        int                      `json:"steps_panel_height"`
	SidebarWidth            int                      `json:"sidebar_width"`
	RecentProjects          []string                 `json:"recent_projects"`
	RecentFeatures          []string                 `json:"recent_features"`
	SessionProject          string                   `json:"session_project,omitempty"`
	OpenTabs                []string                 `json:"open_tabs,omitempty"`
	ActiveTab               string                   `json:"active_tab,omitempty"`
	UntitledTabs            []UntitledTabSession     `json:"untitled_tabs,omitempty"`
	ScrollBeforeClick       bool                     `json:"scroll_before_click"`
	HoverRecordMinMs        int                      `json:"hover_record_min_ms"`
	SelectorClickStrategies []string                 `json:"selector_click_strategies,omitempty"`
	SelectorInputStrategies []string                 `json:"selector_input_strategies,omitempty"`
	LibraryHeuristicsMUI    *bool                    `json:"library_heuristics_mui,omitempty"`
	LibraryHeuristicsAnt    *bool                    `json:"library_heuristics_ant,omitempty"`
	CheckUpdatesOnStartup   *bool                    `json:"check_updates_on_startup,omitempty"`
	HTTPAuth                map[string]HTTPAuthEntry `json:"http_auth,omitempty"`
	Editor                  EditorSettings           `json:"editor,omitempty"`
	ChecklistDismissed      bool                     `json:"checklist_dismissed,omitempty"`
	WelcomePlayedSuccess    bool                     `json:"welcome_played_success,omitempty"`
	OnboardingCompleted     bool                     `json:"onboarding_completed,omitempty"`
	OnboardingDismissed     bool                     `json:"onboarding_dismissed,omitempty"`
	OnboardingVersion       int                      `json:"onboarding_version,omitempty"`
	StartURL                string                   `json:"start_url,omitempty"`
	RunDialogConfirmed      bool                     `json:"run_dialog_confirmed,omitempty"`
	PickerDuringRecording   bool                     `json:"picker_during_recording,omitempty"`
	UILocale                string                   `json:"ui_locale,omitempty"`
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

var (
	atomicCreateTemp = os.CreateTemp
	atomicRename     = os.Rename
	atomicRemove     = os.Remove
)

func DefaultAppSettingsPath() string {
	primary := filepath.Join(paths.AppDataDir(), "settings.json")
	if paths.FileExists(primary) {
		return primary
	}
	for _, legacy := range paths.LegacySettingsPaths() {
		if paths.FileExists(legacy) {
			return legacy
		}
	}
	return primary
}

func LoadDefaultAppSettings() (*AppSettings, error) {
	path := DefaultAppSettingsPath()
	if path == "" {
		return &AppSettings{Browser: "chromium"}, nil
	}
	cfg, err := LoadAppSettings(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
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
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create json dir %q: %w", dir, err)
		}
	}
	data := append(payload, '\n')
	tmp, err := atomicCreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp json file %q: %w", path, err)
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = atomicRemove(tmpPath)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("write temp json file %q: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp json file %q: %w", path, err)
	}
	backupPath := ""
	backupCreated := false
	if _, err := os.Stat(path); err == nil {
		backupPath = path + ".bak"
		if err := atomicRemove(backupPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale backup %q: %w", backupPath, err)
		}
		if err := atomicRename(path, backupPath); err != nil {
			return fmt.Errorf("backup json file %q: %w", path, err)
		}
		backupCreated = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat json file %q: %w", path, err)
	}
	if err := atomicRename(tmpPath, path); err != nil {
		if backupCreated {
			if restoreErr := atomicRename(backupPath, path); restoreErr != nil {
				if copyErr := restoreJSONBackupByCopy(backupPath, path); copyErr != nil {
					return fmt.Errorf("replace json file %q: %w (restore backup %q failed: %v, atomic copy restore failed: %w)", path, err, backupPath, restoreErr, copyErr)
				}
			}
			_ = atomicRemove(backupPath)
		}
		return fmt.Errorf("replace json file %q: %w", path, err)
	}
	if backupCreated {
		_ = atomicRemove(backupPath)
	}
	cleanupTemp = false
	return nil
}

func restoreJSONBackupByCopy(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := atomicCreateTemp(filepath.Dir(dest), filepath.Base(dest)+".restore-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = atomicRemove(tmpPath)
		}
	}()
	if _, err := io.Copy(tmp, in); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := atomicRename(tmpPath, dest); err != nil {
		return err
	}
	cleanupTemp = false
	return nil
}
