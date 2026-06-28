package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInferProjectRoot(t *testing.T) {
	tmp := t.TempDir()
	featureDir := filepath.Join(tmp, "features")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	login := filepath.Join(featureDir, "login.feature")
	logout := filepath.Join(featureDir, "logout.feature")
	for _, path := range []string{login, logout} {
		if err := os.WriteFile(path, []byte("Функционал: Demo\nСценарий: S1\n"), 0o644); err != nil {
			t.Fatalf("write feature: %v", err)
		}
	}

	root := InferProjectRoot([]string{login, logout})
	if root != featureDir {
		t.Fatalf("unexpected root: %q", root)
	}
}
