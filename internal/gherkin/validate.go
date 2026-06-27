package gherkin

func ValidateFeature(feature *Feature) []Issue {
	if feature == nil {
		return []Issue{{Line: 0, Message: "feature is nil"}}
	}

	var issues []Issue
	if feature.Title == "" {
		issues = append(issues, Issue{Line: 1, Message: "missing feature title"})
	}
	if len(feature.Scenarios) == 0 {
		line := feature.Line
		if line == 0 {
			line = 1
		}
		issues = append(issues, Issue{Line: line, Message: "feature has no scenarios"})
	}

	for _, scenario := range feature.Scenarios {
		if scenario.Title == "" {
			line := scenario.Line
			if line == 0 {
				line = feature.Line
			}
			issues = append(issues, Issue{Line: line, Message: "scenario has empty title"})
		}
		if len(scenario.Steps) == 0 {
			line := scenario.Line
			if line == 0 {
				line = feature.Line
			}
			issues = append(issues, Issue{Line: line, Message: "scenario has no steps: " + scenario.Title})
		}
		if scenario.IsOutline {
			if len(scenario.Examples) == 0 {
				line := scenario.Line
				if line == 0 {
					line = feature.Line
				}
				issues = append(issues, Issue{Line: line, Message: "outline scenario has no examples: " + scenario.Title})
				continue
			}
			for _, example := range scenario.Examples {
				if len(example.Rows) < 2 {
					line := example.Line
					if line == 0 {
						line = scenario.ExamplesLine
					}
					issues = append(issues, Issue{Line: line, Message: "examples must include header and at least one data row"})
				}
			}
		}
	}
	return issues
}
