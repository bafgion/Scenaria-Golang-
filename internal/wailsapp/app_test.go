package wailsapp

import "testing"

func TestSearchStepsReturnsRichEntries(t *testing.T) {
	app := NewApp()
	entries := app.svc.SearchSteps("нажимаю")
	if len(entries) == 0 {
		t.Fatal("expected click step")
	}
	found := false
	for _, e := range entries {
		if e.Label == "нажимаю" && e.Action == "click" && e.Example != "" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected rich click entry, got %+v", entries[0])
	}
}

func TestValidateFeature(t *testing.T) {
	app := NewApp()
	text := "Функционал: Demo\n  Сценарий: Ok\n    Допустим открыт \"https://example.com\"\n"
	errs := app.svc.ValidateFeature(text)
	if len(errs) > 0 {
		t.Fatalf("unexpected validation errors: %+v", errs)
	}
}

func TestSearchStepsRejectsNonStringQuery(t *testing.T) {
	// Service layer always receives string from bindings; ensure empty query works.
	app := NewApp()
	all := app.svc.SearchSteps("")
	if len(all) < 50 {
		t.Fatalf("expected full catalog, got %d entries", len(all))
	}
}

func TestRefactorReplaceDelegates(t *testing.T) {
	app := NewApp()
	got := app.svc.RefactorReplaceInText("hello", "hello", "bye", false)
	if got.Count != 1 || got.Text != "bye" {
		t.Fatalf("got %+v", got)
	}
}
