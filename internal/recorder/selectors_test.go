package recorder

import "testing"

func TestCanonicalizeRecordedSelectorKeepsHasText(t *testing.T) {
	in := `button:has-text("Войти")`
	if got := canonicalizeRecordedSelector(in); got != in {
		t.Fatalf("got %q want %q", got, in)
	}
}

func TestClickHasTextSelector(t *testing.T) {
	got := ClickHasTextSelector("Далее", "button")
	if got != `button:has-text("Далее")` {
		t.Fatalf("got %q", got)
	}
	got = ClickHasTextSelector("Каталог", "a")
	if got != `a:has-text("Каталог")` {
		t.Fatalf("got %q", got)
	}
}

func TestContextualClickSelector(t *testing.T) {
	got := ContextualClickSelector("Без договора", "Выбрать")
	want := `div:has-text("Без договора") >> button:has-text("Выбрать")`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestIsUnstableSelectorValue(t *testing.T) {
	if !isUnstableSelectorValue("radix-:r3:") {
		t.Fatal("expected unstable radix id")
	}
	if isUnstableSelectorValue("login-form") {
		t.Fatal("expected stable id")
	}
	if isUnstableSelectorValue("dzhinsy_relaxed_rl200_iz_liotsella_svetlo_zheltogo_tsveta") {
		t.Fatal("expected product slug id to be stable")
	}
	if !isUnstableSelectorValue("2f1a9c0b7e8d6a5c4b3a2910fedcba98") {
		t.Fatal("expected long hex id to be unstable")
	}
}

func TestEventToRecordedStepKeepsHasText(t *testing.T) {
	step, ok := EventToRecordedStep("click", map[string]string{
		"selector": `button:has-text("OK")`,
		"text":     "OK",
	})
	if !ok || step.Selector != `button:has-text("OK")` {
		t.Fatalf("unexpected step: %+v ok=%v", step, ok)
	}
}
