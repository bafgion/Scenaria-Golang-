package gui

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestListPluginsRequiresProject(t *testing.T) {
	svc := NewService()
	if _, err := svc.ListPlugins(); err == nil {
		t.Fatal("expected error without project")
	}
}

func TestPluginLifecycleViaService(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	zipPath := filepath.Join(t.TempDir(), "demo.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	entry, err := w.Create("plugin.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte(`{"name":"demo"}`)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := svc.InstallPlugin("demo", zipPath); err != nil {
		t.Fatalf("install: %v", err)
	}
	plugins, err := svc.ListPlugins()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" {
		t.Fatalf("got %+v", plugins)
	}
	if err := svc.UninstallPlugin("demo"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
}

func TestPluginServiceRejectsMalformedPluginID(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	err := svc.InstallPlugin("../outside", filepath.Join(t.TempDir(), "missing.zip"))
	if err == nil || !strings.Contains(err.Error(), "invalid plugin id") {
		t.Fatalf("expected invalid plugin id error, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "addons")); !os.IsNotExist(statErr) {
		t.Fatalf("addons directory should not be created for malformed id, stat err=%v", statErr)
	}

	err = svc.UninstallPlugin("../outside")
	if err == nil || !strings.Contains(err.Error(), "invalid plugin id") {
		t.Fatalf("expected invalid plugin id uninstall error, got %v", err)
	}
}

func TestPluginServiceCancelActiveCancelsRunInvocation(t *testing.T) {
	root := t.TempDir()
	writePluginDescriptorForTest(t, root, "demo", `{"id":"demo","structuredRuns":[{"runner":"run","args":["--scenario","Scenario With Spaces"]}]}`)
	runner := newBlockingPluginRunner()
	svc := NewPluginService(func() string { return root }, func() pluginCLIRunner { return runner })

	done := make(chan RunResult, 1)
	go func() {
		done <- svc.Run(PluginRunRequest{Name: "demo"})
	}()
	args := runner.waitStarted(t)
	want := []string{root, "--scenario", "Scenario With Spaces"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("run args = %#v, want %#v", args, want)
	}

	svc.CancelActive()

	select {
	case result := <-done:
		if !strings.Contains(result.Error, context.Canceled.Error()) {
			t.Fatalf("expected context canceled result, got %+v", result)
		}
	case <-timeAfterTest():
		t.Fatal("service Run did not return after cancellation")
	}
}

func TestPluginServiceCancelPluginOnlyCancelsMatchingInvocation(t *testing.T) {
	root := t.TempDir()
	writePluginDescriptorForTest(t, root, "demo", `{"id":"demo","commands":["run"]}`)
	writePluginDescriptorForTest(t, root, "other", `{"id":"other","commands":["run"]}`)
	demoRunner := newBlockingPluginRunner()
	otherRunner := newBlockingPluginRunner()
	svc := NewPluginService(
		func() string { return root },
		func() pluginCLIRunner {
			if demoRunner.startedCount() == 0 {
				return demoRunner
			}
			return otherRunner
		},
	)

	demoDone := make(chan RunResult, 1)
	otherDone := make(chan RunResult, 1)
	go func() { demoDone <- svc.Run(PluginRunRequest{Name: "demo"}) }()
	_ = demoRunner.waitStarted(t)
	go func() { otherDone <- svc.Run(PluginRunRequest{Name: "other"}) }()
	_ = otherRunner.waitStarted(t)

	svc.CancelPlugin("demo")

	select {
	case <-demoDone:
	case <-timeAfterTest():
		t.Fatal("matching plugin invocation was not canceled")
	}
	select {
	case <-otherDone:
		t.Fatal("non-matching plugin invocation was canceled")
	default:
	}
	svc.CancelActive()
	select {
	case <-otherDone:
	case <-timeAfterTest():
		t.Fatal("remaining plugin invocation was not canceled")
	}
}

func writePluginDescriptorForTest(t *testing.T, root, name, payload string) {
	t.Helper()
	dir := filepath.Join(root, "addons", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
}

type blockingPluginRunner struct {
	started  chan []string
	canceled chan struct{}
	once     sync.Once
	mu       sync.Mutex
	count    int
}

func newBlockingPluginRunner() *blockingPluginRunner {
	return &blockingPluginRunner{
		started:  make(chan []string, 1),
		canceled: make(chan struct{}),
	}
}

func (r *blockingPluginRunner) VA(args []string) (string, error) {
	return "", nil
}

func (r *blockingPluginRunner) Run(ctx context.Context, args []string) (string, error) {
	copied := append([]string(nil), args...)
	r.mu.Lock()
	r.count++
	r.mu.Unlock()
	r.started <- copied
	<-ctx.Done()
	r.once.Do(func() { close(r.canceled) })
	return "canceled", ctx.Err()
}

func (r *blockingPluginRunner) waitStarted(t *testing.T) []string {
	t.Helper()
	select {
	case args := <-r.started:
		return args
	case <-timeAfterTest():
		t.Fatal("plugin runner did not start")
	}
	return nil
}

func (r *blockingPluginRunner) startedCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

func timeAfterTest() <-chan time.Time {
	return time.After(time.Second)
}
