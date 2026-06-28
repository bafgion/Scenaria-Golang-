package stepdsl

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestParseWaitLoneSecondsSuffix(t *testing.T) {
	step := gherkin.Step{Line: 1, Text: "жду 2 с"}
	action, err := Parse(step)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if action.Kind != "wait" || action.Value1 != "2000ms" {
		t.Fatalf("unexpected action: %+v", action)
	}
}

func TestParseSwitchTabOneBasedStorage(t *testing.T) {
	step := gherkin.Step{Line: 1, Text: "переключаюсь на вкладку 3"}
	action, err := Parse(step)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if action.IntVal != 2 {
		t.Fatalf("expected 0-based index 2, got %d", action.IntVal)
	}
}
