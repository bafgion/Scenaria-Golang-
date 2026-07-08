package report

import (
	"path/filepath"
	"testing"
)

func TestLayoutRunArtifactsUsesPerRunDir(t *testing.T) {
	root := t.TempDir()
	layout, err := LayoutRunArtifacts(root, "run-42", RunArtifactInputs(
		".scenaria/report.html",
		".scenaria/junit.xml",
		".scenaria/summary.json",
		".scenaria/allure-results",
		".scenaria/traces",
		".scenaria/videos",
	))
	if err != nil {
		t.Fatal(err)
	}
	if layout.RunID != "run-42" {
		t.Fatalf("runID = %q", layout.RunID)
	}
	if !filepath.IsAbs(layout.Dir) {
		t.Fatalf("expected absolute dir, got %q", layout.Dir)
	}
	if filepath.Base(layout.HTMLPath) != "report.html" {
		t.Fatalf("html = %q", layout.HTMLPath)
	}
	if filepath.Base(filepath.Dir(layout.HTMLPath)) != "run-42" {
		t.Fatalf("html dir = %q", layout.HTMLPath)
	}
}

func TestWriteLatestRunPointer(t *testing.T) {
	root := t.TempDir()
	layout := RunArtifactLayout{RunID: "run-1", Dir: filepath.Join(root, ".scenaria", "runs", "run-1"), HTMLPath: "report.html"}
	if err := WriteLatestRunPointer(root, layout); err != nil {
		t.Fatal(err)
	}
	got, err := ReadLatestRunPointer(root)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.RunID != "run-1" {
		t.Fatalf("unexpected latest run pointer: %+v", got)
	}
}
