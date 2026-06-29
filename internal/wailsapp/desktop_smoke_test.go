//go:build desktop

package wailsapp_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDesktopSmokeScenariaGUI(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("desktop smoke runs on Windows with WebView2")
	}
	root, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(root, "scripts", "desktop-smoke.ps1")
	exe := os.Getenv("SCENARIA_GUI_EXE")
	if exe == "" {
		exe = filepath.Join(root, "build", "bin", "scenaria-gui.exe")
	}
	if _, err := os.Stat(exe); err != nil {
		t.Skipf("scenaria-gui.exe not built (%s); run wails build first", exe)
	}
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script, "-ExePath", exe)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("desktop smoke failed: %v\n%s", err, out)
	}
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "wails.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
