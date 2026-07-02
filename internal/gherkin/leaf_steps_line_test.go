package gherkin

import "testing"

func TestLeafStepLineAtIndex(t *testing.T) {
	steps := []Step{
		{Line: 3, Text: "открыт"},
		{Block: "Если", Line: 4, Children: []Step{{Line: 5, Text: "клик"}}},
	}
	line, ok := LeafStepLineAtIndex(steps, 1)
	if !ok || line != 5 {
		t.Fatalf("expected line 5, got %d ok=%v", line, ok)
	}
}
