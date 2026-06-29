package gui

import (
	"strings"
	"testing"
)

func TestAnalyzeScenarioHintsMenuHover(t *testing.T) {
	text := "\tКогда нажимаю \"nav >> a.menu\"\n"
	hints := AnalyzeScenarioHints(text)
	if len(hints) != 1 || hints[0].ID != "menu_hover" {
		t.Fatalf("got %+v", hints)
	}
}

func TestApplyMenuHoverFix(t *testing.T) {
	text := "\tКогда нажимаю \"nav >> a.menu\"\n"
	result := ApplyScenarioHintFix(ScenarioHintFixRequest{
		Text: text, HintID: "menu_hover", StepIndex: 0,
	})
	if result.Count != 1 {
		t.Fatalf("count=%d", result.Count)
	}
	if !strings.Contains(result.Text, `навожу "nav"`) || !strings.Contains(result.Text, `нажимаю "a.menu"`) {
		t.Fatalf("got %q", result.Text)
	}
}

func TestAnalyzeDuplicateGoto(t *testing.T) {
	text := "\tДопустим открыт \"https://a.com\"\n\tИ открыт \"https://b.com\"\n"
	hints := AnalyzeScenarioHints(text)
	if len(hints) != 1 || hints[0].ID != "duplicate_goto" {
		t.Fatalf("got %+v", hints)
	}
}
