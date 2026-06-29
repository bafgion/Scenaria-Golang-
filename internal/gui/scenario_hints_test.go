package gui

import (
	"strings"
	"testing"
)

func TestAnalyzeScenarioHintsSetsLine(t *testing.T) {
	text := "Функция: demo\n\tСценарий: s\n\tКогда нажимаю \"nav >> a.menu\"\n"
	hints := AnalyzeScenarioHints(text)
	if len(hints) != 1 {
		t.Fatalf("got %+v", hints)
	}
	if hints[0].Line != 3 {
		t.Fatalf("expected line 3, got %d", hints[0].Line)
	}
}

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

func TestProposeMenuHoverSelectorsChain(t *testing.T) {
	hover, click, ok := ProposeMenuHoverSelectors(`div:has-text("Категория") >> a:has-text("Категория")`, "")
	if !ok || hover != `div:has-text("Категория")` || click != `a:has-text("Категория")` {
		t.Fatalf("got %q %q ok=%v", hover, click, ok)
	}
}

func TestProposeMenuHoverSelectorsExplicitHover(t *testing.T) {
	hover, click, ok := ProposeMenuHoverSelectors("a.sub", "nav.menu")
	if !ok || hover != "nav.menu" || click != "a.sub" {
		t.Fatalf("got %q %q ok=%v", hover, click, ok)
	}
}

func TestAnalyzeDuplicateGoto(t *testing.T) {
	text := "\tДопустим открыт \"https://a.com\"\n\tИ открыт \"https://b.com\"\n"
	hints := AnalyzeScenarioHints(text)
	if len(hints) != 1 || hints[0].ID != "duplicate_goto" {
		t.Fatalf("got %+v", hints)
	}
}

func TestAnalyzeFragileSelector(t *testing.T) {
	text := "\tКогда нажимаю \"div > input:nth-of-type(2)\"\n"
	hints := AnalyzeScenarioHints(text)
	found := false
	for _, h := range hints {
		if h.ID == "fragile_selector" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected fragile_selector, got %+v", hints)
	}
}

func TestAnalyzeFillNoAssert(t *testing.T) {
	text := "\tКогда ввожу \"alice\" в \"#name\"\n\tИ нажимаю \"#submit\"\n"
	hints := AnalyzeScenarioHints(text)
	found := false
	for _, h := range hints {
		if h.ID == "fill_no_assert" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected fill_no_assert, got %+v", hints)
	}
}

func TestAnalyzeGotoNoWait(t *testing.T) {
	text := "\tДопустим открыт \"https://example.com\"\n\tКогда нажимаю \"#login\"\n"
	hints := AnalyzeScenarioHints(text)
	found := false
	for _, h := range hints {
		if h.ID == "goto_no_wait" {
			found = true
			if !h.AutoFixable {
				t.Fatalf("expected auto-fixable goto_no_wait")
			}
		}
	}
	if !found {
		t.Fatalf("expected goto_no_wait, got %+v", hints)
	}
}

func TestApplyFillNoAssertFix(t *testing.T) {
	text := "\tКогда ввожу \"alice\" в \"#name\"\n\tИ нажимаю \"#submit\"\n"
	result := ApplyScenarioHintFix(ScenarioHintFixRequest{
		Text: text, HintID: "fill_no_assert", StepIndex: 0,
	})
	if result.Count != 1 {
		t.Fatalf("count=%d", result.Count)
	}
	if !strings.Contains(result.Text, `проверяю текст "alice" в "#name"`) {
		t.Fatalf("got %q", result.Text)
	}
}

func TestApplyGotoNoWaitFix(t *testing.T) {
	text := "\tДопустим открыт \"https://example.com\"\n\tКогда нажимаю \"#login\"\n"
	result := ApplyScenarioHintFix(ScenarioHintFixRequest{
		Text: text, HintID: "goto_no_wait", StepIndex: 0,
	})
	if result.Count != 1 {
		t.Fatalf("count=%d", result.Count)
	}
	if !strings.Contains(result.Text, `жду появления "#login"`) {
		t.Fatalf("got %q", result.Text)
	}
}

func TestApplyDuplicateGotoFix(t *testing.T) {
	text := "\tДопустим открыт \"https://a.com\"\n\tИ открыт \"https://b.com\"\n"
	result := ApplyScenarioHintFix(ScenarioHintFixRequest{
		Text: text, HintID: "duplicate_goto", StepIndex: 1,
	})
	if result.Count != 1 {
		t.Fatalf("count=%d", result.Count)
	}
	if strings.Contains(result.Text, "https://b.com") {
		t.Fatalf("second goto should be removed: %q", result.Text)
	}
	if !strings.Contains(result.Text, "https://a.com") {
		t.Fatalf("first goto should remain: %q", result.Text)
	}
}

func TestCollectAllHintsDetectsIssues(t *testing.T) {
	text := "" +
		"\tДопустим открыт \"https://a.com\"\n" +
		"\tИ открыт \"https://b.com\"\n" +
		"\tКогда нажимаю \"div:has-text(\\\"Menu\\\") >> a.item >> button.ok\"\n" +
		"\tИ ввожу \"x\" в \"#email\"\n"
	hints := AnalyzeScenarioHints(text)
	ids := map[string]bool{}
	for _, h := range hints {
		ids[h.ID] = true
	}
	for _, want := range []string{"duplicate_goto", "long_chain", "fill_no_assert"} {
		if !ids[want] {
			t.Fatalf("missing hint %q, got %+v", want, hints)
		}
	}
}

func TestAnalyzeLongChain(t *testing.T) {
	text := "\tКогда нажимаю \"a >> b >> c\"\n"
	hints := AnalyzeScenarioHints(text)
	found := false
	for _, h := range hints {
		if h.ID == "long_chain" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected long_chain, got %+v", hints)
	}
}

func TestAnalyzeDivClick(t *testing.T) {
	text := "\tКогда нажимаю \"div:has-text(\\\"Каталог\\\")\"\n"
	hints := AnalyzeScenarioHints(text)
	found := false
	for _, h := range hints {
		if h.ID == "div_click" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected div_click, got %+v", hints)
	}
}

func TestIsFragileSelectorIgnoresStableSelectors(t *testing.T) {
	if isFragileSelector(`button:has-text("Войти")`) {
		t.Fatal("plain has-text should not be fragile")
	}
	if isFragileSelector(`#login-form input`) {
		t.Fatal("id selector should not be fragile")
	}
}
