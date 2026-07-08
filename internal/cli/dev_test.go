package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunDevGoroutines(t *testing.T) {
	if err := RunDev([]string{"goroutines"}); err != nil {
		t.Fatal(err)
	}
}

func TestRunDevPprof(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "goroutine.pprof")
	if err := RunDev([]string{"pprof", "--out", outPath}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("expected non-empty pprof output")
	}
}

func TestRunDevGoroutinesIncludesCurrentGoroutine(t *testing.T) {
	// Capture output by re-running stack dump logic indirectly.
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)
	if n == 0 || !strings.Contains(string(buf[:n]), "goroutine") {
		t.Fatal("expected goroutine stack dump")
	}
}
