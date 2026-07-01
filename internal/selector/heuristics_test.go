package selector

import (
	"strings"
	"testing"
)

func TestRecorderScriptGolden(t *testing.T) {
	heuristics := HeuristicsJS
	for _, needle := range []string{
		"function buildInputSelector",
		"function buildClickSelector",
		"__scenariaHeuristics",
	} {
		if !strings.Contains(heuristics, needle) {
			t.Fatalf("heuristics script missing %q", needle)
		}
	}
	recorder := RecorderScript
	for _, needle := range []string{
		"MutationObserver",
		"elementFromPoint",
		"__scenariaRecorder",
		"shadowRoot",
	} {
		if !strings.Contains(recorder, needle) {
			t.Fatalf("recorder script missing %q", needle)
		}
	}
}

func TestBuildFromElementPriority(t *testing.T) {
	if got := BuildFromElement(ElementInfo{Tag: "button", ID: "submit"}); got != `#submit` {
		t.Fatalf("expected id selector, got %q", got)
	}
	if got := BuildFromElement(ElementInfo{Tag: "input", TestID: "login-email"}); got != `[data-testid="login-email"]` {
		t.Fatalf("unexpected testid selector: %q", got)
	}
	if got := BuildFromElement(ElementInfo{Tag: "button", Text: "Войти"}); got != `text="Войти"` {
		t.Fatalf("unexpected text selector: %q", got)
	}
	if got := BuildFromElement(ElementInfo{Tag: "button", Title: "Остановить запись"}); got != `[title="Остановить запись"]` {
		t.Fatalf("unexpected title selector: %q", got)
	}
	if got := BuildFromElement(ElementInfo{Tag: "canvas", TestID: "signature-pad"}); got != `[data-testid="signature-pad"]` {
		t.Fatalf("unexpected canvas testid selector: %q", got)
	}
}
