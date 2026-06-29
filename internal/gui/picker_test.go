package gui

import "testing"

func TestPickerStepChoicesIncludeClickAndRaw(t *testing.T) {
	choices := PickerStepChoices("button.buy", "Допустим")
	if len(choices) < 2 {
		t.Fatalf("expected choices, got %d", len(choices))
	}
	if choices[0].Label != "Клик" || choices[0].StepBody != `нажимаю "button.buy"` {
		t.Fatalf("first choice: %+v", choices[0])
	}
	last := choices[len(choices)-1]
	if last.Label != "Только селектор" || last.StepBody != `"button.buy"` {
		t.Fatalf("last choice: %+v", last)
	}
}

func TestPickerStepChoicesEscapeQuotes(t *testing.T) {
	choices := PickerStepChoices(`input[name="email"]`, "Допустим")
	if choices[0].StepBody != `нажимаю 'input[name="email"]'` {
		t.Fatalf("got %q", choices[0].StepBody)
	}
}

func TestPickerStepChoicesRespectKeyword(t *testing.T) {
	first := PickerStepChoices("button.buy", "Допустим")
	next := PickerStepChoices("button.buy", "И")
	if first[0].Preview != `  Допустим нажимаю "button.buy"` {
		t.Fatalf("preview: %q", first[0].Preview)
	}
	if next[0].Preview != `  И нажимаю "button.buy"` {
		t.Fatalf("preview: %q", next[0].Preview)
	}
}
