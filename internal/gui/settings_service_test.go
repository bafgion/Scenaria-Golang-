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

func TestSettingsServiceSaveSettingsTrimsPersistedStrings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	svc := NewSettingsService(settings.NewStore(path))
	if err := settings.SaveAppSettings(path, &settings.AppSettings{Browser: "chromium"}); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	if err := svc.SaveSettings(AppSettingsDTO{
		Browser:        "chromium",
		NavWaitUntil:   "  domcontentloaded  ",
		SessionProject: "  /projects/demo  ",
		ActiveTab:      "  demo.feature  ",
		StartURL:       "  https://example.com  ",
	}); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	cfg, err := svc.LoadSettings()
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}
	if cfg.NavWaitUntil != "domcontentloaded" {
		t.Fatalf("expected trimmed nav wait, got %q", cfg.NavWaitUntil)
	}
	if cfg.SessionProject != "/projects/demo" {
		t.Fatalf("expected trimmed session project, got %q", cfg.SessionProject)
	}
	if cfg.ActiveTab != "demo.feature" {
		t.Fatalf("expected trimmed active tab, got %q", cfg.ActiveTab)
	}
	if cfg.StartURL != "https://example.com" {
		t.Fatalf("expected trimmed start url, got %q", cfg.StartURL)
	}
}

func TestSettingsServiceRecentsRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	svc := NewSettingsService(settings.NewStore(path))
	if err := settings.SaveAppSettings(path, &settings.AppSettings{Browser: "chromium"}); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	if err := svc.RememberRecentProject("/projects/a"); err != nil {
		t.Fatalf("remember project: %v", err)
	}
	if err := svc.RememberRecentFeature("/features/login.feature"); err != nil {
		t.Fatalf("remember feature: %v", err)
	}

	recents := svc.LoadRecents()
	if len(recents.Projects) != 1 || recents.Projects[0] != "/projects/a" {
		t.Fatalf("unexpected projects: %#v", recents.Projects)
	}
	if len(recents.Features) != 1 || recents.Features[0] != "/features/login.feature" {
		t.Fatalf("unexpected features: %#v", recents.Features)
	}
}

func TestSettingsServiceHTTPAuthRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	svc := NewSettingsService(settings.NewStore(path))
	if err := settings.SaveAppSettings(path, &settings.AppSettings{Browser: "chromium"}); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	if err := svc.SaveHTTPAuth(HTTPAuthRequest{
		Host:     "example.com",
		Username: "user",
		Password: "secret",
	}); err != nil {
		t.Fatalf("save http auth: %v", err)
	}

	hosts, err := svc.ListHTTPAuthHosts()
	if err != nil {
		t.Fatalf("list hosts: %v", err)
	}
	if len(hosts) != 1 || hosts[0] != "example.com" {
		t.Fatalf("unexpected hosts: %#v", hosts)
	}

	creds, err := svc.HTTPAuthForHost("example.com")
	if err != nil {
		t.Fatalf("auth for host: %v", err)
	}
	if creds.Username != "user" || !creds.HasPassword {
		t.Fatalf("unexpected creds: %#v", creds)
	}

	if err := svc.RemoveHTTPAuth("example.com"); err != nil {
		t.Fatalf("remove http auth: %v", err)
	}
	hosts, err = svc.ListHTTPAuthHosts()
	if err != nil {
		t.Fatalf("list hosts after remove: %v", err)
	}
	if len(hosts) != 0 {
		t.Fatalf("expected empty hosts, got %#v", hosts)
	}
}
