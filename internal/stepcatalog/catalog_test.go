package stepcatalog_test

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

func TestEntriesFromGolden(t *testing.T) {
	entries := stepcatalog.Entries()
	if len(entries) < 40 {
		t.Fatalf("expected catalog from golden fixture, got %d entries", len(entries))
	}
	found := false
	for _, entry := range entries {
		if entry.Template == `рисую подпись в "#sign"` {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected draw-signature template in catalog")
	}
}

func TestSearchFilters(t *testing.T) {
	results := stepcatalog.Search("телефон")
	if len(results) == 0 {
		t.Fatal("expected phone generator in search results")
	}
}
