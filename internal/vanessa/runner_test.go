package vanessa

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateSettingsMissingPaths(t *testing.T) {
	issues := ValidateSettings(DefaultSettings())
	if len(issues) == 0 {
		t.Fatal("expected validation issues for empty settings")
	}
}

func TestMergeVAParamsDry(t *testing.T) {
	tmp := t.TempDir()
	cfg := DefaultSettings()
	cfg.ReportJUnit = true
	runDir := filepath.Join(tmp, "run-1")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	merged, path, err := MergeVAParams(cfg, RunRequest{ProjectRoot: tmp, Paths: []string{tmp}}, runDir)
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}
	if merged == nil || path == "" {
		t.Fatalf("unexpected merge result: %+v %q", merged, path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("VAParams not written: %v", err)
	}
}

func TestParseJUnitFileEmpty(t *testing.T) {
	if cases := parseJUnitFile("missing.xml"); cases != nil {
		t.Fatalf("expected nil for missing file")
	}
}

func TestRunContextCancelsPlatformProcess(t *testing.T) {
	tmp := t.TempDir()
	epfPath := filepath.Join(tmp, "vanessa.epf")
	if err := os.WriteFile(epfPath, []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}
	markerPath := filepath.Join(tmp, "helper.started")
	cfg := DefaultSettings()
	cfg.PlatformExecutable = os.Args[0]
	cfg.PlatformMode = "-test.run=TestVanessaContextHelperProcess"
	cfg.PlatformExtraArgs = []string{"--", "--scenaria-helper-marker", markerPath}
	cfg.EPFPath = epfPath
	cfg.RunsDir = filepath.Join(tmp, "runs")
	cfg.ProcessTimeoutSec = 60
	if err := SaveProjectSettings(tmp, cfg); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct {
		result BatchResult
		err    error
	}, 1)
	go func() {
		result, err := RunContext(ctx, RunRequest{ProjectRoot: tmp, Paths: []string{tmp}})
		done <- struct {
			result BatchResult
			err    error
		}{result: result, err: err}
	}()

	waitForFile(t, markerPath)
	cancel()

	select {
	case got := <-done:
		if !errors.Is(got.err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got result=%+v err=%v", got.result, got.err)
		}
		if got.result.ExitCode != -1 {
			t.Fatalf("ExitCode = %d, want -1", got.result.ExitCode)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("RunContext did not return after cancellation")
	}
}

func TestVanessaContextHelperProcess(t *testing.T) {
	markerPath := ""
	for i, arg := range os.Args {
		if arg == "--scenaria-helper-marker" && i+1 < len(os.Args) {
			markerPath = os.Args[i+1]
			break
		}
	}
	if markerPath == "" {
		return
	}
	if err := os.WriteFile(markerPath, []byte("started"), 0o644); err != nil {
		os.Exit(2)
	}
	for {
		time.Sleep(time.Second)
	}
}

func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.After(3 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %s", path)
		case <-ticker.C:
			if _, err := os.Stat(path); err == nil {
				return
			}
		}
	}
}
