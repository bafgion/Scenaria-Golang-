package runstatus

import (
	"testing"
)

func TestStoreRecordBatch(t *testing.T) {
	tmp := t.TempDir()
	store, err := Open(tmp)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	batch := []Entry{
		{Path: "a.feature::One", Success: true, Runner: "playwright"},
		{Path: "b.feature::Two", Success: false, Runner: "playwright"},
		{Path: "c.feature::Three", Success: true, Runner: "playwright"},
	}
	if err := store.RecordBatch(batch); err != nil {
		t.Fatalf("RecordBatch failed: %v", err)
	}
	got, err := store.List(0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].Path != "c.feature::Three" || got[1].Path != "b.feature::Two" || got[2].Path != "a.feature::One" {
		t.Fatalf("unexpected order: %+v", got)
	}
}
