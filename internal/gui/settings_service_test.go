package gui

import (
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestSettingsServiceLoadSettingsNilStoreDefaults(t *testing.T) {
	svc := NewSettingsService(nil)

	cfg, err := svc.LoadSettings()
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}
	if cfg.Browser != "chromium" {
		t.Fatalf("expected chromium default, got %q", cfg.Browser)
	}
}

func TestSettingsServiceSaveSettingsRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	svc := NewSettingsService(settings.NewStore(path))
	if err := settings.SaveAppSettings(path, &settings.AppSettings{Browser: "chromium"}); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	if err := svc.SaveSettings(AppSettingsDTO{
		Browser:            "firefox",
		ParallelWorkers:    4,
		StepsPanelHeight:   170,
		RecentProjects:     []string{"a", "b"},
		ChecklistDismissed: true,
	}); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	cfg, err := svc.LoadSettings()
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}
	if cfg.Browser != "firefox" {
		t.Fatalf("expected firefox, got %q", cfg.Browser)
	}
	if cfg.ParallelWorkers != 4 {
		t.Fatalf("expected workers=4, got %d", cfg.ParallelWorkers)
	}
	if !cfg.ChecklistDismissed {
		t.Fatal("expected checklist dismissed persisted")
	}
}
