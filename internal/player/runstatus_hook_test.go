package player

import (
	"context"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

func TestRunStatusHookRecordsScenario(t *testing.T) {
	dir := t.TempDir()
	store, err := runstatus.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := WithRunStatusHook(context.Background(), store, "playwright")
	recordScenarioRunStatus(ctx, ScenarioResult{
		FeaturePath: "features/login.feature",
		Scenario:    "Успешный вход",
		Status:      "passed",
	})
	entries, err := store.List(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(entries))
	}
	if entries[0].Path != "features/login.feature::Успешный вход" {
		t.Fatalf("path = %q", entries[0].Path)
	}
	if !entries[0].Success || entries[0].Runner != "playwright" {
		t.Fatalf("entry = %+v", entries[0])
	}
}
