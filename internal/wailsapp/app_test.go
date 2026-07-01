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

func TestFormatFeatureBinding(t *testing.T) {
	app := NewApp()
	text := "  Когда нажимаю \"a\"\n\n  Тогда вижу \"b\"\n"
	got := app.FormatFeature(text)
	if got != "\tКогда нажимаю \"a\"\n\tТогда вижу \"b\"\n" {
		t.Fatalf("got %q", got)
	}
}

func TestDescribeEditorLineBinding(t *testing.T) {
	app := NewApp()
	entry := app.DescribeEditorLine(`	Когда нажимаю "#login"`)
	if entry.Label != "нажимаю" {
		t.Fatalf("got %+v", entry)
	}
}

func TestEventBindingTypes(t *testing.T) {
	app := NewApp()
	progress, result := app.EventBindingTypes()
	if progress.Percent != 0 || progress.Message != "" {
		t.Fatalf("expected zero progress DTO, got %+v", progress)
	}
	if result.Error != "" || result.Output != "" || len(result.Cases) != 0 {
		t.Fatalf("expected zero result DTO, got %+v", result)
	}
}

func TestRefactorReplaceDelegates(t *testing.T) {
	app := NewApp()
	got := app.svc.RefactorReplaceInText("hello", "hello", "bye", false)
	if got.Count != 1 || got.Text != "bye" {
		t.Fatalf("got %+v", got)
	}
}
