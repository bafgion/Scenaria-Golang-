package gui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

func TestPollVanessaRun_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	snap := NewService().PollVanessaRun(tmp, 3)
	if snap.TotalPlanned != 3 {
		t.Fatalf("got total %d", snap.TotalPlanned)
	}
	if len(snap.Cases) != 0 {
		t.Fatalf("expected no cases")
	}
}

func TestPollVanessaRun_WithJUnit(t *testing.T) {
	tmp := t.TempDir()
	junitDir := filepath.Join(tmp, "junit")
	if err := os.MkdirAll(junitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="suite">
  <testcase classname="Login" name="Успех"/>
</testsuite>`
	if err := os.WriteFile(filepath.Join(junitDir, "a.xml"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	snap := NewService().PollVanessaRun(tmp, 2)
	if len(snap.Cases) != 1 {
		t.Fatalf("expected 1 case, got %d", len(snap.Cases))
	}
	if snap.CompletedCases != 1 {
		t.Fatalf("completed %d", snap.CompletedCases)
	}
}

func TestPluginToVanessaRun_ScenarioAndFiles(t *testing.T) {
	req := PluginRunRequest{
		Tag:          "@smoke",
		Scenario:     "Login",
		VaDir:        "features/smoke",
		VaFiles:      "a.feature,b.feature",
		ReportAllure: true,
	}
	vReq := pluginToVanessaRun("/proj", req)
	if vReq.Tag != "@smoke" || len(vReq.ScenarioNames) != 1 || vReq.ScenarioNames[0] != "Login" {
		t.Fatalf("unexpected scenario request: %+v", vReq)
	}
	if !vReq.ReportAllure {
		t.Fatal("expected report allure")
	}
	if len(vReq.Paths) < 2 {
		t.Fatalf("expected paths from vaDir and vaFiles, got %v", vReq.Paths)
	}
}

func TestRunVanessaPluginCanceledOnProjectSwitch(t *testing.T) {
	root := t.TempDir()
	nextProject := t.TempDir()
	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	original := runVanessaBatch
	started := make(chan struct{})
	runVanessaBatch = func(ctx context.Context, req vanessa.RunRequest) (vanessa.BatchResult, error) {
		if req.ProjectRoot != root {
			t.Errorf("ProjectRoot = %q, want %q", req.ProjectRoot, root)
		}
		close(started)
		<-ctx.Done()
		return vanessa.BatchResult{Error: ctx.Err().Error()}, ctx.Err()
	}
	defer func() { runVanessaBatch = original }()

	done := make(chan VanessaRunResultDTO, 1)
	go func() {
		done <- svc.RunVanessaPlugin(PluginRunRequest{Name: "vanessa"})
	}()
	select {
	case <-started:
	case <-timeAfterTest():
		t.Fatal("Vanessa batch did not start")
	}

	if _, err := svc.OpenProject(nextProject); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}

	select {
	case result := <-done:
		if !strings.Contains(result.Error, context.Canceled.Error()) {
			t.Fatalf("expected context canceled result, got %+v", result)
		}
	case <-timeAfterTest():
		t.Fatal("Vanessa batch did not stop after project switch")
	}
}
