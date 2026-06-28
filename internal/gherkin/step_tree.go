package gherkin

import "fmt"

func FinalizeStepTrees(feature *Feature) error {
	if feature == nil {
		return nil
	}
	var err error
	feature.Background, err = buildStepTree(feature.Background)
	if err != nil {
		return err
	}
	if feature.HasContextBlock {
		testClient, err := ParseTestClientName(feature.Background)
		if err != nil {
			return err
		}
		feature.TestClient = testClient
	}

	for i := range feature.Scenarios {
		feature.Scenarios[i].Steps, err = buildStepTree(feature.Scenarios[i].Steps)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildStepTree(flat []Step) ([]Step, error) {
	if len(flat) == 0 {
		return nil, nil
	}
	level := normalizeBaseIndent(flat)
	tree, _, err := parseStepsAtLevel(flat, 0, level)
	return tree, err
}

func normalizeBaseIndent(flat []Step) int {
	min := -1
	for _, step := range flat {
		if step.Indent < 0 {
			continue
		}
		if min < 0 || step.Indent < min {
			min = step.Indent
		}
	}
	if min < 0 {
		return 0
	}
	return min
}

func parseStepsAtLevel(flat []Step, start int, level int) ([]Step, int, error) {
	out := make([]Step, 0)
	i := start
	for i < len(flat) {
		step := flat[i]
		if step.Indent < level {
			break
		}
		if step.Indent > level {
			return nil, i, fmt.Errorf("line %d: unexpected step indent", step.Line)
		}

		header, err := detectBlockHeader(step)
		if err != nil {
			return nil, i, err
		}
		if header != nil {
			step.Block = header.Kind
			step.Condition = header.Condition
			step.RepeatCount = header.Count
			step.ForEachSelector = header.Selector
			step.ForEachVariable = header.Variable
			i++
			children, next, err := parseStepsAtLevel(flat, i, level+1)
			if err != nil {
				return nil, i, err
			}
			step.Children = children
			i = next
			out = append(out, step)
			continue
		}

		out = append(out, step)
		i++
	}
	return out, i, nil
}

func FlattenSteps(steps []Step) []Step {
	out := make([]Step, 0)
	for _, step := range steps {
		if step.Block != "" {
			out = append(out, step)
			out = append(out, FlattenSteps(step.Children)...)
			continue
		}
		out = append(out, step)
	}
	return out
}

func CountLeafSteps(steps []Step) int {
	total := 0
	for _, step := range steps {
		if step.Block != "" {
			total += CountLeafSteps(step.Children)
			continue
		}
		total++
	}
	return total
}

func CloneSteps(steps []Step) []Step {
	out := make([]Step, len(steps))
	for i, step := range steps {
		out[i] = step
		if len(step.Children) > 0 {
			out[i].Children = CloneSteps(step.Children)
		}
	}
	return out
}
