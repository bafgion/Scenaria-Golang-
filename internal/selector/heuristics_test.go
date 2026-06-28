package selector

import "testing"

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
	if got := BuildFromElement(ElementInfo{Tag: "canvas", TestID: "signature-pad"}); got != `[data-testid="signature-pad"]` {
		t.Fatalf("unexpected canvas testid selector: %q", got)
	}
}
