package gui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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

	result, _, err := svc.runInProcess(context.Background(), RunRequest{
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

func TestRun_DryRunReturnsRunScopedHTMLPath(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "smoke.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: smoke\nСценарий: тест\nКогда открыт \"https://example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(tmp, ".scenaria", "report.html")
	if err := os.MkdirAll(filepath.Dir(legacy), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacy, []byte("stale legacy report"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(tmp); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}

	result := svc.Run(RunRequest{
		Targets:  []string{featurePath},
		DryRun:   true,
		HTMLPath: legacy,
	}, nil)
	if result.Error != "" {
		t.Fatalf("Run dry-run: %s", result.Error)
	}
	got := firstNonEmpty(result.ReportPath, result.HTMLPath)
	if got == "" {
		t.Fatal("expected html path in run result")
	}
	if got == legacy {
		t.Fatalf("expected run-scoped report path, got legacy %q", got)
	}
	if !strings.Contains(filepath.ToSlash(got), "/runs/run-") {
		t.Fatalf("expected path under .scenaria/runs/, got %q", got)
	}
	if _, err := os.Stat(got); err != nil {
		t.Fatalf("report not written at %q: %v", got, err)
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
	_, _, err := svc.runInProcess(context.Background(), RunRequest{
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
