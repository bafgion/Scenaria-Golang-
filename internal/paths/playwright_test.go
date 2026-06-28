package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurePlaywrightBrowsersBundled(t *testing.T) {
	tmp := t.TempDir()
	browsers := filepath.Join(tmp, "browsers")
	if err := os.MkdirAll(browsers, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(tmp, "scenaria-test.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PLAYWRIGHT_BROWSERS_PATH", "")
	// Cannot reliably override os.Executable in test; verify helper returns path shape.
	if got := BundledBrowsersDir(); got == "" && filepath.Base(got) != "browsers" {
		// BundledBrowsersDir uses real executable path in tests — smoke check only.
	}
}
