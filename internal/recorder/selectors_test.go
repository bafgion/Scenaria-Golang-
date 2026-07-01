package recorder

import "testing"

func TestCanonicalizeRecordedSelector(t *testing.T) {
	cases := map[string]string{
		`button:has-text("Войти")`:                          `text="Войти"`,
		`a:has-text("Каталог")`:                             `text="Каталог"`,
		`div:has-text("Меню") >> button:has-text("Товары")`: `text="Меню" >> text="Товары"`,
		`[data-testid="submit"]`:                            `[data-testid="submit"]`,
	}
	for in, want := range cases {
		if got := canonicalizeRecordedSelector(in); got != want {
			t.Fatalf("canonicalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEventToRecordedStepCanonicalizesHasText(t *testing.T) {
	step, ok := EventToRecordedStep("click", map[string]string{
		"selector": `button:has-text("OK")`,
		"text":     "OK",
	})
	if !ok || step.Selector != `text="OK"` {
		t.Fatalf("unexpected step: %+v ok=%v", step, ok)
	}
}
