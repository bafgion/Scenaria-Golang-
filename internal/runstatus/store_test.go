package runstatus

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreRecordAndLatest(t *testing.T) {
	tmp := t.TempDir()
	store, err := Open(tmp)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if err := store.Record(Entry{Path: "demo.feature", Success: true, Runner: "playwright"}); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	latest, err := store.Latest("demo.feature")
	if err != nil || latest == nil || !latest.Success {
		t.Fatalf("unexpected latest: %+v err=%v", latest, err)
	}
}

func TestWriteAtomicLockedRestoresOriginalOnReplacementFailure(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "run_status.json")
	if err := os.WriteFile(path, []byte(`[{"path":"old.feature","success":true}]`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := runStatusRename
	defer func() { runStatusRename = origRename }()
	calls := 0
	runStatusRename = func(oldpath, newpath string) error {
		calls++
		if calls == 2 {
			return errors.New("replace failed")
		}
		return os.Rename(oldpath, newpath)
	}
	err := writeAtomicLocked(path, []byte(`[{"path":"new.feature","success":true}]`+"\n"))
	if err == nil {
		t.Fatal("expected error")
	}
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(raw), "old.feature") {
		t.Fatalf("expected original run status restored, got %q", string(raw))
	}
	if _, err := os.Stat(path + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup after restore, stat err=%v", err)
	}
}

func TestWriteAtomicLockedRetainsBackupWhenRestoreCannotReplaceTarget(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "run_status.json")
	if err := os.WriteFile(path, []byte(`[{"path":"old.feature","success":true}]`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := runStatusRename
	defer func() { runStatusRename = origRename }()
	calls := 0
	runStatusRename = func(oldpath, newpath string) error {
		calls++
		switch calls {
		case 2, 3, 4:
			return errors.New("rename unavailable")
		default:
			return os.Rename(oldpath, newpath)
		}
	}
	err := writeAtomicLocked(path, []byte(`[{"path":"new.feature","success":true}]`+"\n"))
	if err == nil {
		t.Fatal("expected error")
	}
	if raw, readErr := os.ReadFile(path); readErr == nil && strings.Contains(string(raw), "new.feature") {
		t.Fatalf("replacement content must not be installed after failed restore: %q", string(raw))
	}
	backup, readErr := os.ReadFile(path + ".bak")
	if readErr != nil {
		t.Fatalf("expected retained backup: %v", readErr)
	}
	if !strings.Contains(string(backup), "old.feature") {
		t.Fatalf("expected backup to retain original run status, got %q", string(backup))
	}
}
