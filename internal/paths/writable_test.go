package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWritableScenariaDirUsesProjectWhenWritable(t *testing.T) {
	tmp := t.TempDir()
	got, err := WritableScenariaDir(tmp)
	if err != nil {
		t.Fatalf("WritableScenariaDir: %v", err)
	}
	want := filepath.Join(tmp, ".scenaria")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRemapScenariaArtifactKeepsWritableProjectPath(t *testing.T) {
	tmp := t.TempDir()
	artifact := filepath.Join(tmp, ".scenaria", "report.html")
	got := RemapScenariaArtifact(tmp, artifact)
	if got != artifact {
		t.Fatalf("got %q want %q", got, artifact)
	}
}

func TestRemapScenariaArtifactFallsBackForReadOnlyProject(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("read-only directory chmod is unreliable on Windows")
	}
	tmp := t.TempDir()
	local := filepath.Join(tmp, ".scenaria")
	if err := os.MkdirAll(local, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(local, 0o755) })

	artifact := filepath.Join(local, "report.html")
	got := RemapScenariaArtifact(tmp, artifact)
	if got == artifact {
		t.Fatalf("expected remapped path, got original %q", got)
	}
	if filepath.Base(got) != "report.html" {
		t.Fatalf("expected report.html basename, got %q", got)
	}
	if !dirWritable(filepath.Dir(got)) {
		t.Fatalf("remapped dir is not writable: %s", filepath.Dir(got))
	}
}

func TestScenariaArtifactPath(t *testing.T) {
	tmp := t.TempDir()
	got, err := ScenariaArtifactPath(tmp, "junit.xml")
	if err != nil {
		t.Fatalf("ScenariaArtifactPath: %v", err)
	}
	want := filepath.Join(tmp, ".scenaria", "junit.xml")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
