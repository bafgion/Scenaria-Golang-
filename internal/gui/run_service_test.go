package gui

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

func TestRunServiceListRunResults(t *testing.T) {
	root := t.TempDir()
	store, err := runstatus.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(runstatus.Entry{
		Path:    "demo.feature::A",
		Success: false,
		Message: "boom",
		Runner:  "playwright",
	}); err != nil {
		t.Fatal(err)
	}

	svc := NewRunService(func() string { return root })
	entries, err := svc.ListRunResults(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path == "" || entries[0].Runner != "playwright" {
		t.Fatalf("unexpected entry: %#v", entries[0])
	}
}

func TestRunServiceFlakyMetrics(t *testing.T) {
	root := t.TempDir()
	store, err := runstatus.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(runstatus.Entry{
		Path:    "demo.feature::A",
		Success: false,
		Runner:  "playwright",
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Record(runstatus.Entry{
		Path:    "demo.feature::A",
		Success: true,
		Runner:  "playwright",
	}); err != nil {
		t.Fatal(err)
	}

	svc := NewRunService(func() string { return root })
	metrics, err := svc.FlakyMetrics(200)
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics.Scenarios) == 0 {
		t.Fatal("expected scenario stats")
	}
}
