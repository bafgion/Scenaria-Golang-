package gui

import (
	"os"
	"path/filepath"
	"testing"
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
	if result.Output == "" || result.Output[:7] != "file://" {
		t.Fatalf("expected file url, got %q", result.Output)
	}
}
