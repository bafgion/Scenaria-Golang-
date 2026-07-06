package gui

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunInProcess_DryRunWritesHTMLReport(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "smoke.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: smoke\nСценарий: тест\nКогда открыт \"https://example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	scenaria := filepath.Join(tmp, ".scenaria")
	if err := os.MkdirAll(scenaria, 0o755); err != nil {
		t.Fatal(err)
	}
	htmlPath := filepath.Join(scenaria, "report.html")

	svc := NewService()
	if _, err := svc.OpenProject(tmp); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}

	result, err := svc.runInProcess(context.Background(), RunRequest{
		Targets:  []string{featurePath},
		DryRun:   true,
		HTMLPath: htmlPath,
	}, nil)
	if err != nil {
		t.Fatalf("runInProcess dry-run: %v", err)
	}
	if result.Mode != "dry-run" {
		t.Fatalf("expected dry-run mode, got %q", result.Mode)
	}
	if _, err := os.Stat(htmlPath); err != nil {
		t.Fatalf("HTML report not written: %v", err)
	}
}

func TestRunInProcess_DryRunSkipsReportWhenPathEmpty(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "smoke.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: smoke\nСценарий: тест\nКогда открыт \"https://example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(tmp); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}
	_, err := svc.runInProcess(context.Background(), RunRequest{
		Targets: []string{featurePath},
		DryRun:  true,
	}, nil)
	if err != nil {
		t.Fatalf("runInProcess dry-run: %v", err)
	}
	report := filepath.Join(tmp, ".scenaria", "report.html")
	if _, statErr := os.Stat(report); statErr == nil {
		t.Fatal("expected no default HTML report without htmlPath")
	}
}
