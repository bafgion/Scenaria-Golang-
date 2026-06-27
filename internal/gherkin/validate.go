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
		issues = append(issues, Issue{Line: 1, Message: "feature has no scenarios"})
	}

	for _, scenario := range feature.Scenarios {
		if scenario.Title == "" {
			issues = append(issues, Issue{Line: 1, Message: "scenario has empty title"})
		}
		if len(scenario.Steps) == 0 {
			issues = append(issues, Issue{Line: 1, Message: "scenario has no steps: " + scenario.Title})
		}
	}
	return issues
}
