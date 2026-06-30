package gherkin

// LeafSteps returns executable leaf steps in document order (block headers omitted).
func LeafSteps(steps []Step) []Step {
	out := make([]Step, 0)
	for _, step := range steps {
		if step.Block != "" {
			out = append(out, LeafSteps(step.Children)...)
			continue
		}
		out = append(out, step)
	}
	return out
}

// LeafStepIndexAtLine returns the 0-based leaf index for an exact source line.
func LeafStepIndexAtLine(steps []Step, line int) (int, bool) {
	leaves := LeafSteps(steps)
	for i, step := range leaves {
		if step.Line == line {
			return i, true
		}
	}
	return 0, false
}

// ApplyStepRange slices leaf steps from start through end (inclusive).
// Negative start is treated as 0; negative end runs through the last step.
func ApplyStepRange(steps []Step, start, end int) []Step {
	leaves := LeafSteps(steps)
	if len(leaves) == 0 {
		return leaves
	}
	if start < 0 {
		start = 0
	}
	if start >= len(leaves) {
		start = len(leaves) - 1
	}
	if end < 0 || end >= len(leaves) {
		end = len(leaves) - 1
	}
	if end < start {
		end = start
	}
	return leaves[start : end+1]
}

// ScenarioContainingLine finds the scenario block that owns a source line.
func ScenarioContainingLine(feature *Feature, line int) (Scenario, bool) {
	if feature == nil || line < 1 {
		return Scenario{}, false
	}
	for i, sc := range feature.Scenarios {
		end := int(^uint(0) >> 1)
		if i+1 < len(feature.Scenarios) {
			end = feature.Scenarios[i+1].Line - 1
		}
		if line >= sc.Line && line <= end {
			return sc, true
		}
	}
	return Scenario{}, false
}
