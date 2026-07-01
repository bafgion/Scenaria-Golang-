package player

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupDownloads(t *testing.T) {
	root := t.TempDir()
	ctx := NewRunContext(nil, 99, root)
	dir := ctx.DownloadDir()
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("download dir missing: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file.bin"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx.CleanupDownloads()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected download dir removed, stat err=%v", err)
	}
}
