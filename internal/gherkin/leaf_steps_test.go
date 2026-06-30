package gherkin

import "testing"

func TestLeafStepsAndApplyStepRange(t *testing.T) {
	steps := []Step{
		{Line: 3, Text: `открыт "https://example.com"`},
		{Line: 4, Text: `нажимаю "#ok"`},
		{Line: 5, Block: BlockIf, Children: []Step{
			{Line: 6, Text: `жду 1 сек`},
		}},
	}
	leaves := LeafSteps(steps)
	if len(leaves) != 3 {
		t.Fatalf("expected 3 leaf steps, got %d", len(leaves))
	}
	sliced := ApplyStepRange(steps, 1, -1)
	if len(sliced) != 2 || sliced[0].Line != 4 {
		t.Fatalf("unexpected slice: %+v", sliced)
	}
	idx, ok := LeafStepIndexAtLine(steps, 4)
	if !ok || idx != 1 {
		t.Fatalf("index at line 4: %d ok=%v", idx, ok)
	}
}

func TestScenarioContainingLine(t *testing.T) {
	feature := &Feature{
		Scenarios: []Scenario{
			{Title: "A", Line: 2},
			{Title: "B", Line: 5},
		},
	}
	if sc, ok := ScenarioContainingLine(feature, 3); !ok || sc.Title != "A" {
		t.Fatalf("line 3: %+v ok=%v", sc, ok)
	}
	if sc, ok := ScenarioContainingLine(feature, 5); !ok || sc.Title != "B" {
		t.Fatalf("line 5: %+v ok=%v", sc, ok)
	}
}
