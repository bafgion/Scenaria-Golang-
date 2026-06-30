package gui

import "testing"

func TestUpdateStartURLs(t *testing.T) {
	text := "  Допустим открыт \"https://old.com\"\n"
	got := UpdateStartURLs(text, "https://new.com")
	if got.Count != 1 {
		t.Fatalf("count=%d", got.Count)
	}
	if got.Text != "  Допустим открыт \"https://new.com\"\n" {
		t.Fatalf("text=%q", got.Text)
	}
}

func TestCollapseBlankLinesBetweenSteps(t *testing.T) {
	text := "  Когда нажимаю \"a\"\n\n  Тогда вижу \"b\"\n"
	got := CollapseBlankLinesBetweenSteps(text)
	if got != "  Когда нажимаю \"a\"\n  Тогда вижу \"b\"\n" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatFeature(t *testing.T) {
	text := "  Когда нажимаю \"a\"\n\n  Тогда вижу \"b\"\n"
	got := FormatFeature(text)
	want := "\tКогда нажимаю \"a\"\n\tТогда вижу \"b\"\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestReplaceInText(t *testing.T) {
	got := ReplaceInText("hello Hello", "hello", "bye", false)
	if got.Count != 2 {
		t.Fatalf("count=%d want 2", got.Count)
	}
	if got.Text != "bye bye" {
		t.Fatalf("text=%q", got.Text)
	}
}

func TestReplaceInTextCaseSensitive(t *testing.T) {
	got := ReplaceInText("hello Hello", "hello", "bye", true)
	if got.Count != 1 || got.Text != "bye Hello" {
		t.Fatalf("got %+v", got)
	}
}
