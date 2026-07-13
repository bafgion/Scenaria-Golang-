package gui

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDuplicateFeatureRejectsUnsafeDestinationNames(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(src, []byte("Feature: Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(filepath.Dir(root), "outside.feature")
	if err := os.WriteFile(outside, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := newFileOpsTestService(root)

	invalid := []string{
		"../outside",
		`..\outside`,
		"sub/name",
		`sub\name`,
		`C:\outside`,
		`C:outside`,
		`\\server\share`,
		"CON",
		"COM1.txt",
		"LPT1.feature",
		".feature",
		"demo.",
		"demo ",
		"bad:name",
		"bad*name",
		"bad\nname",
	}
	for _, name := range invalid {
		if _, err := svc.DuplicateFeature(src, name); err == nil {
			t.Fatalf("DuplicateFeature(%q) expected error", name)
		}
	}
	assertGUIFileContent(t, outside, "keep")
}

func TestRenameFeatureRejectsUnsafeDestinationNames(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(src, []byte("Feature: Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := newFileOpsTestService(root)

	for _, name := range []string{"../outside", `C:outside`, "CON", "bad:name", "demo "} {
		if _, err := svc.RenameFeature(src, name); err == nil {
			t.Fatalf("RenameFeature(%q) expected error", name)
		}
	}
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("source should remain after invalid rename: %v", err)
	}
}

func TestDuplicateFeatureCollisionPastHundredDoesNotOverwrite(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(src, []byte("Feature: Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	createCollisionSeries(t, root, "demo-copy", ".feature", 105)
	svc := newFileOpsTestService(root)

	got, err := svc.DuplicateFeature(src, "")
	if err != nil {
		t.Fatalf("DuplicateFeature: %v", err)
	}
	if filepath.Base(got) != "demo-copy-106.feature" {
		t.Fatalf("got %q", got)
	}
	assertGUIFileContent(t, filepath.Join(root, "demo-copy.feature"), "existing-1")
	assertGUIFileContent(t, filepath.Join(root, "demo-copy-99.feature"), "existing-99")
	assertGUIFileContent(t, got, "Feature: Demo\n")
}

func TestImportFeaturesCollisionPastHundredDoesNotOverwrite(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	src := filepath.Join(external, "imported.feature")
	if err := os.WriteFile(src, []byte("Feature: Imported\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	createCollisionSeries(t, root, "imported", ".feature", 105)
	svc := newFileOpsTestService(root)

	paths, err := svc.ImportFeatures(root, []string{src})
	if err != nil {
		t.Fatalf("ImportFeatures: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected one import, got %d", len(paths))
	}
	if filepath.Base(paths[0]) != "imported-106.feature" {
		t.Fatalf("got %q", paths[0])
	}
	assertGUIFileContent(t, filepath.Join(root, "imported.feature"), "existing-1")
	assertGUIFileContent(t, filepath.Join(root, "imported-99.feature"), "existing-99")
	assertGUIFileContent(t, paths[0], "Feature: Imported\n")
}

func TestCopyFeatureFileExclusivePreservesExistingDestination(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.feature")
	dest := filepath.Join(root, "dest.feature")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFeatureFileExclusive(src, dest); err == nil {
		t.Fatal("expected existing destination error")
	}
	assertGUIFileContent(t, dest, "old")
}

func TestCopyFeatureFileExclusiveCleansPartialDestination(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src-dir")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(root, "dest.feature")
	err := copyFeatureFileExclusive(srcDir, dest)
	if err == nil {
		t.Fatal("expected copy from directory to fail")
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("partial destination should be removed, stat err=%v copy err=%v", statErr, err)
	}
}

func newFileOpsTestService(root string) *FileOperationService {
	return NewFileOperationService(
		func(path string) (string, error) {
			if filepath.IsAbs(path) {
				return filepath.Clean(path), nil
			}
			return filepath.Join(root, path), nil
		},
		func() string { return root },
		func(run func() error) error { return run() },
		func(run func() error) error { return run() },
	)
}

func createCollisionSeries(t *testing.T, dir, base, ext string, count int) {
	t.Helper()
	for i := 1; i <= count; i++ {
		name := base + ext
		if i > 1 {
			name = base + "-" + strconv.Itoa(i) + ext
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte("existing-"+strconv.Itoa(i)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertGUIFileContent(t *testing.T, path, want string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(raw) != want {
		t.Fatalf("%s = %q, want %q", path, string(raw), want)
	}
}
