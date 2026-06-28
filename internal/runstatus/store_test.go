package runstatus

import (
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
