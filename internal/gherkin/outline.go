package gherkin

import (
	"regexp"
	"strings"
)

var placeholderKeyRE = regexp.MustCompile(`<([^>]+)>`)

type RunnableScenario struct {
	Title        string
	Tags         []string
	Steps        []Step
	ExampleIndex int
}

func ExpandFeature(feature *Feature) []RunnableScenario {
	return ExpandFeatureAtPath(feature, "")
}

func ExpandFeatureAtPath(feature *Feature, featurePath string) []RunnableScenario {
	if feature == nil {
		return nil
	}
	out := make([]RunnableScenario, 0)
	for _, scenario := range feature.Scenarios {
		out = append(out, ExpandScenarioAtPath(feature, scenario, featurePath)...)
	}
	return out
}

func ExpandScenario(feature *Feature, scenario Scenario) []RunnableScenario {
	return ExpandScenarioAtPath(feature, scenario, "")
}

func ExpandScenarioAtPath(feature *Feature, scenario Scenario, featurePath string) []RunnableScenario {
	tags := MergeTags(feature.Tags, scenario.Tags)
	if !scenario.IsOutline {
		return []RunnableScenario{{
			Title: scenario.Title,
			Tags:  tags,
			Steps: mergeBackgroundSteps(feature.Background, scenario.Steps),
		}}
	}

	out := make([]RunnableScenario, 0)
	for _, example := range scenario.Examples {
		if len(example.Rows) < 2 {
			continue
		}
		header := example.Rows[0]
		for rowIndex, row := range example.Rows[1:] {
			values := exampleRowValues(header, row)
			steps := mergeBackgroundSteps(feature.Background, expandSteps(scenario.Steps, values))
			title := scenario.Title
			if sample := firstRowSample(header, row); sample != "" {
				title = scenario.Title + " — " + sample
			}
			out = append(out, RunnableScenario{
				Title:        title,
				Tags:         tags,
				Steps:        steps,
				ExampleIndex: rowIndex + 1,
			})
		}
	}
	if len(out) == 0 && scenario.IsOutline && strings.TrimSpace(featurePath) != "" {
		if rows, err := LoadScenarioParams(featurePath, scenario.Title); err == nil && len(rows) > 0 {
			out = expandScenarioFromParamRows(feature, scenario, tags, rows)
		}
	}
	return out
}

func mergeBackgroundSteps(background []Step, scenarioSteps []Step) []Step {
	out := make([]Step, 0, len(background)+len(scenarioSteps))
	out = append(out, cloneSteps(background)...)
	out = append(out, cloneSteps(scenarioSteps)...)
	return out
}

func cloneSteps(steps []Step) []Step {
	return CloneSteps(steps)
}

func expandSteps(steps []Step, values map[string]string) []Step {
	out := make([]Step, len(steps))
	for i, step := range steps {
		out[i] = expandStep(step, values)
	}
	return out
}

func expandStep(step Step, values map[string]string) Step {
	expanded := step
	expanded.Text = substitutePlaceholders(step.Text, values)
	if len(step.Children) > 0 {
		expanded.Children = expandSteps(step.Children, values)
	}
	return expanded
}

func substitutePlaceholders(text string, values map[string]string) string {
	for key, value := range values {
		text = strings.ReplaceAll(text, "<"+key+">", value)
	}
	return text
}

func exampleRowValues(header []string, row []string) map[string]string {
	values := make(map[string]string, len(header))
	for i, key := range header {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if i < len(row) {
			values[key] = row[i]
		}
	}
	return values
}

func firstRowSample(header []string, row []string) string {
	return firstValuesSampleOrdered(exampleRowValues(header, row), header)
}

func placeholderKeysFromSteps(steps []Step) []string {
	seen := make(map[string]bool)
	var keys []string
	var walk func([]Step)
	walk = func(ss []Step) {
		for _, step := range ss {
			for _, match := range placeholderKeyRE.FindAllStringSubmatch(step.Text, -1) {
				key := match[1]
				if !seen[key] {
					seen[key] = true
					keys = append(keys, key)
				}
			}
			if len(step.Children) > 0 {
				walk(step.Children)
			}
		}
	}
	walk(steps)
	return keys
}

func firstValuesSampleOrdered(values map[string]string, keys []string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values[key]); value != "" {
			return value
		}
	}
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
