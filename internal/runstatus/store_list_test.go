package runstatus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreList(t *testing.T) {
	dir := t.TempDir()
	store := &Store{path: filepath.Join(dir, "run_status.json")}
	for i := 0; i < 3; i++ {
		if err := store.Record(Entry{Path: "a.feature::S", Success: i == 0}); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	got, err := store.List(2)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, err := os.Stat(store.path); err != nil {
		t.Fatalf("status file missing: %v", err)
	}
}
