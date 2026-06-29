package gui

import (
	"strings"
	"testing"
)

const textReplaceSample = `Функционал: UI
@smoke
Сценарий: Demo
	# old url
	Допустим открыт "https://staging.example.com"
	И нажимаю "buy"
`

func TestLineIsReplaceableSkipsHeadersTagsComments(t *testing.T) {
	if LineIsReplaceable("Функционал: UI", true) {
		t.Fatal("feature header should not be replaceable")
	}
	if LineIsReplaceable("@smoke", true) {
		t.Fatal("tag should not be replaceable")
	}
	if LineIsReplaceable("\t# note", true) {
		t.Fatal("comment should not be replaceable")
	}
	if !LineIsReplaceable(`	Допустим открыт "x"`, true) {
		t.Fatal("step line should be replaceable")
	}
}

func TestReplaceAllInTextStepsOnly(t *testing.T) {
	got := ReplaceAllInText(textReplaceSample, "staging.example.com", "prod.example.com", true, true)
	if got.Count != 1 {
		t.Fatalf("count=%d want 1", got.Count)
	}
	if strings.Contains(got.Text, "staging") {
		t.Fatalf("staging should be replaced: %q", got.Text)
	}
	if !strings.Contains(got.Text, "@smoke") || !strings.Contains(got.Text, "Функционал: UI") {
		t.Fatalf("headers/tags preserved: %q", got.Text)
	}
}

func TestReplaceAllInTextPreservesIndents(t *testing.T) {
	got := ReplaceAllInText(textReplaceSample, "buy", "checkout", true, true)
	if got.Count != 1 {
		t.Fatalf("count=%d", got.Count)
	}
	if !strings.Contains(got.Text, `	И нажимаю "checkout"`) {
		t.Fatalf("indent preserved: %q", got.Text)
	}
}

func TestReplaceAllInTextStepsOnlySkipsScenarioTitle(t *testing.T) {
	text := "Функционал: old title\nСценарий: old title\n\tДопустим открыт \"https://x.com\"\n"
	got := ReplaceAllInText(text, "old title", "new title", true, true)
	if got.Count != 0 {
		t.Fatalf("expected no replacements in headers, got %d: %q", got.Count, got.Text)
	}
}
