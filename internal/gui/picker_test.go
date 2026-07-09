package gui

import "testing"

func TestPickerSuggestedChoiceIndex(t *testing.T) {
	if got := PickerSuggestedChoiceIndex("fill"); got != 1 {
		t.Fatalf("fill index: got %d want 1", got)
	}
	if got := PickerSuggestedChoiceIndex("click"); got != 0 {
		t.Fatalf("click index: got %d want 0", got)
	}
	if got := PickerSuggestedChoiceIndex("unknown"); got != 0 {
		t.Fatalf("unknown index: got %d want 0", got)
	}
}

func TestPickerStepChoicesIncludesFillAndSelect(t *testing.T) {
	choices := PickerStepChoices(`#email`, "Допустим")
	labels := make([]string, len(choices))
	for i, c := range choices {
		labels[i] = c.Label
	}
	for _, want := range []string{"Клик", "Заполнить", "Выбрать"} {
		found := false
		for _, label := range labels {
			if label == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing %q in %v", want, labels)
		}
	}
}
