package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
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

	svc := NewReportService(func() string { return root }, nil, nil)
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
	svc := NewReportService(func() string { return root }, nil, nil)
	if !svc.ArtifactExists(file) {
		t.Fatal("expected existing file")
	}
	if svc.ArtifactExists(filepath.Join(root, "missing")) {
		t.Fatal("expected missing path to be false")
	}
}

func TestReportServiceWriteRunReportsSkipsLatestPointerOnWriteError(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	if err := os.MkdirAll(filepath.Join(scenaria, "runs", "run-1", "report.html"), 0o755); err != nil {
		t.Fatal(err)
	}
	svc := NewReportService(func() string { return root }, func() string { return "run-1" }, nil)
	plan := player.ExecutionPlan{Cases: []player.RunCase{{FeaturePath: "a.feature", Name: "Fail"}}}
	result := player.ExecutionResult{Mode: "browser", RunID: "run-1"}

	layout, err := svc.WriteRunReports(root, RunRequest{HTMLPath: filepath.Join(scenaria, "report.html")}, plan, result)
	if err == nil {
		t.Fatal("expected report write error")
	}
	if layout.RunID != "run-1" {
		t.Fatalf("unexpected layout: %+v", layout)
	}
	if got, readErr := report.ReadLatestRunPointer(root); readErr != nil {
		t.Fatal(readErr)
	} else if got != nil {
		t.Fatalf("latest pointer should not be written on report failure: %+v", got)
	}
}
