package stepcatalog_test

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

func TestEntriesFromGolden(t *testing.T) {
	entries := stepcatalog.Entries()
	if len(entries) < 50 {
		t.Fatalf("expected catalog from snippets+golden, got %d entries", len(entries))
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

func TestEntryHasExampleAndParameters(t *testing.T) {
	entries := stepcatalog.Entries()
	var click *stepcatalog.Entry
	for i := range entries {
		if entries[i].Action == "click" && entries[i].Label == "нажимаю" {
			click = &entries[i]
			break
		}
	}
	if click == nil {
		t.Fatal("expected click snippet in catalog")
	}
	if click.Example == "" || click.Example != click.Template {
		t.Fatalf("expected example text, got %q", click.Example)
	}
	if len(click.Parameters) == 0 {
		t.Fatal("expected parameters for click step")
	}
}
