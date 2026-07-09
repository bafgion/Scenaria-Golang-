package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/report"
)

func TestFileURL(t *testing.T) {
	got := fileURL(`C:\proj\.scenaria\report.html`)
	want := "file:///C:/proj/.scenaria/report.html"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestServeAllure_MissingDir(t *testing.T) {
	dir := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	result := svc.ServeAllure("")
	if result.Error == "" {
		t.Fatal("expected error for missing allure results")
	}
}

func TestOpenHTMLReport_Missing(t *testing.T) {
	dir := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	result := svc.OpenHTMLReport("")
	if result.Error == "" {
		t.Fatal("expected error for missing report")
	}
}

func TestOpenHTMLReportPrefersLatestRunPointer(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	legacy := filepath.Join(scenaria, "report.html")
	latestHTML := filepath.Join(scenaria, "runs", "run-1", "report.html")
	if err := os.MkdirAll(filepath.Dir(latestHTML), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacy, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(latestHTML, []byte("fresh"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := report.WriteLatestRunPointer(root, report.RunArtifactLayout{
		RunID:    "run-1",
		Dir:      filepath.Dir(latestHTML),
		HTMLPath: latestHTML,
	}); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	result := svc.OpenHTMLReport(legacy)
	if result.Error != "" {
		t.Fatalf("OpenHTMLReport: %s", result.Error)
	}
	if result.Output != latestHTML {
		t.Fatalf("OpenHTMLReport = %q, want %q", result.Output, latestHTML)
	}
}

func TestOpenHTMLReport_Found(t *testing.T) {
	dir := t.TempDir()
	scenaria := filepath.Join(dir, ".scenaria")
	if err := os.MkdirAll(scenaria, 0o755); err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(scenaria, "report.html")
	if err := os.WriteFile(report, []byte("<html></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	result := svc.OpenHTMLReport("")
	if result.Error != "" {
		t.Fatalf("unexpected error: %s", result.Error)
	}
	if result.Output == "" {
		t.Fatal("expected report path")
	}
	if filepath.Base(result.Output) != "report.html" {
		t.Fatalf("expected report path, got %q", result.Output)
	}
}
