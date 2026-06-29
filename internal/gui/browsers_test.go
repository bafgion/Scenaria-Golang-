package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBrowserInstallStatusDTO(t *testing.T) {
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

	svc := NewService()
	status := svc.BrowserInstallStatus("chromium")
	if !status.Installed {
		t.Fatalf("expected chromium installed, got %+v", status)
	}
	if status.Label != "Chromium" {
		t.Fatalf("label: %q", status.Label)
	}
}
