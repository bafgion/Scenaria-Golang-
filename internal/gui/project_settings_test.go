package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestSaveProjectConfigHTMLReportOpenMode(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveProjectConfig(ProjectConfigDTO{HTMLReportOpenMode: "light"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := settings.LoadProjectConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTMLReportOpenMode != "light" {
		t.Fatalf("got %q", cfg.HTMLReportOpenMode)
	}
}

func TestOpenHTMLReportRespectsProjectLightMode(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	if err := os.MkdirAll(scenaria, 0o755); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(scenaria, "report.html")
	light := filepath.Join(scenaria, "report.light.html")
	if err := os.WriteFile(full, []byte("full"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(light, []byte("light"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := settings.SaveProjectConfig(root, settings.ProjectConfig{HTMLReportOpenMode: "light"}); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	result := svc.OpenHTMLReport("")
	if result.Error != "" {
		t.Fatalf("OpenHTMLReport: %s", result.Error)
	}
	if filepath.Base(result.Output) != "report.light.html" {
		t.Fatalf("expected light report, got %q", result.Output)
	}
}
