package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeBrowserEngine(t *testing.T) {
	if NormalizeBrowserEngine("Firefox") != "firefox" {
		t.Fatal("firefox")
	}
	if NormalizeBrowserEngine("") != "chromium" {
		t.Fatal("default chromium")
	}
}

func TestFindBrowserExecutableChromium(t *testing.T) {
	tmp := t.TempDir()
	folder := filepath.Join(tmp, "chromium-1200")
	exe := filepath.Join(folder, "chrome-win64", "chrome.exe")
	if err := os.MkdirAll(filepath.Dir(exe), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exe, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, ok := FindBrowserExecutable(tmp, "chromium")
	if !ok || got != exe {
		t.Fatalf("got %q ok=%v", got, ok)
	}
}

func TestBrowserInstallStatusMissing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("PLAYWRIGHT_BROWSERS_PATH", tmp)
	installed, detail := BrowserInstallStatus("firefox")
	if installed {
		t.Fatalf("expected missing firefox, got %q", detail)
	}
	if detail == "" {
		t.Fatal("expected detail message")
	}
}

func TestEnsurePlaywrightEngineSkipsWhenPresent(t *testing.T) {
	tmp := t.TempDir()
	folder := filepath.Join(tmp, "chromium-1200")
	exe := filepath.Join(folder, "chrome-win64", "chrome.exe")
	if err := os.MkdirAll(filepath.Dir(exe), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exe, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PLAYWRIGHT_BROWSERS_PATH", tmp)
	if err := EnsurePlaywrightEngine("chromium"); err != nil {
		t.Fatalf("EnsurePlaywrightEngine: %v", err)
	}
}
