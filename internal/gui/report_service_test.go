package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReportServiceProjectArtifacts(t *testing.T) {
	root := t.TempDir()
	sc := filepath.Join(root, ".scenaria")
	if err := os.MkdirAll(filepath.Join(sc, "traces"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sc, "traces", "trace.zip"), []byte("zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sc, "report.html"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewReportService(func() string { return root })
	art := svc.ProjectArtifacts()
	if art.HTMLReport == "" {
		t.Fatal("expected html report path")
	}
	if art.TracesDir == "" {
		t.Fatal("expected traces dir path")
	}
}

func TestReportServiceArtifactExists(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "x.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewReportService(func() string { return root })
	if !svc.ArtifactExists(file) {
		t.Fatal("expected existing file")
	}
	if svc.ArtifactExists(filepath.Join(root, "missing")) {
		t.Fatal("expected missing path to be false")
	}
}
