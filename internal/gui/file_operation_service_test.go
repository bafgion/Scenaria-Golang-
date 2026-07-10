package gui

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOperationServiceReadAndSaveFeature(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(path, []byte("Функционал: Demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewFileOperationService(
		func(_ string) (string, error) { return path, nil },
		func() string { return root },
		func(run func() error) error { return run() },
		func(run func() error) error { return run() },
	)

	got, err := svc.ReadFeature(path)
	if err != nil {
		t.Fatalf("ReadFeature: %v", err)
	}
	if got != "Функционал: Demo" {
		t.Fatalf("unexpected content: %q", got)
	}
	if err := svc.SaveFeature(path, "Функционал: Updated"); err != nil {
		t.Fatalf("SaveFeature: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "Функционал: Updated" {
		t.Fatalf("unexpected saved content: %q", string(raw))
	}
	matches, err := filepath.Glob(filepath.Join(root, "demo.feature.tmp-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected atomic save cleanup, found temp files: %v", matches)
	}
}

func TestWriteFileAtomicNewFileSave(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "new.feature")
	if err := writeFileAtomic(target, []byte("Feature: New"), 0o644); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "Feature: New" {
		t.Fatalf("unexpected content: %q", string(raw))
	}
	matches, err := filepath.Glob(filepath.Join(root, "new.feature.tmp-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected temp cleanup, found: %v", matches)
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup for new file, stat err=%v", err)
	}
}

func TestWriteFileAtomicReplacesExistingFile(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(target, []byte("new"), 0o600); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "new" {
		t.Fatalf("unexpected content: %q", string(raw))
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup after successful replace, stat err=%v", err)
	}
}

func TestWriteFileAtomicRestoresOriginalOnFirstRenameFailure(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := atomicRename
	defer func() { atomicRename = origRename }()
	calls := 0
	atomicRename = func(oldpath, newpath string) error {
		calls++
		if calls == 1 {
			return errors.New("first rename failed")
		}
		return os.Rename(oldpath, newpath)
	}
	err := writeFileAtomic(target, []byte("new"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}
	raw, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(raw) != "old" {
		t.Fatalf("expected original content to remain, got %q", string(raw))
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup after first rename failure, stat err=%v", err)
	}
}

func TestWriteFileAtomicRestoresOriginalOnReplacementRenameFailure(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := atomicRename
	defer func() { atomicRename = origRename }()
	calls := 0
	atomicRename = func(oldpath, newpath string) error {
		calls++
		if calls == 2 {
			return errors.New("replace failed")
		}
		return os.Rename(oldpath, newpath)
	}
	err := writeFileAtomic(target, []byte("new"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}
	raw, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(raw) != "old" {
		t.Fatalf("expected original content restored, got %q", string(raw))
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup after restore, stat err=%v", err)
	}
}

func TestWriteFileAtomicRetainsBackupWhenRestoreCannotReplaceTarget(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := atomicRename
	defer func() { atomicRename = origRename }()
	calls := 0
	atomicRename = func(oldpath, newpath string) error {
		calls++
		switch calls {
		case 2, 3, 4:
			return errors.New("rename unavailable")
		default:
			return os.Rename(oldpath, newpath)
		}
	}
	err := writeFileAtomic(target, []byte("new"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}
	if raw, readErr := os.ReadFile(target); readErr == nil && string(raw) == "new" {
		t.Fatalf("replacement content must not be installed after failed restore")
	}
	backup, readErr := os.ReadFile(target + ".bak")
	if readErr != nil {
		t.Fatalf("expected retained backup: %v", readErr)
	}
	if string(backup) != "old" {
		t.Fatalf("expected backup to retain original content, got %q", string(backup))
	}
}

func TestFileOperationServiceDeleteFeature(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(path, []byte("Feature:"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewFileOperationService(
		func(p string) (string, error) { return p, nil },
		func() string { return root },
		func(run func() error) error { return run() },
		func(run func() error) error { return run() },
	)
	if err := svc.DeleteFeature(path); err != nil {
		t.Fatalf("DeleteFeature: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err=%v", err)
	}
}
