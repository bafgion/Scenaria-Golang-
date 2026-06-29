package gui

import "strings"

type ExportPreview struct {
	StepCount     int               `json:"stepCount"`
	ScenarioTitle string            `json:"scenarioTitle"`
	Issues        []ValidationIssue `json:"issues"`
	Hints         []ScenarioHintDTO `json:"hints"`
}

func (s *Service) PreviewExport(text string) ExportPreview {
	return ExportPreview{
		StepCount:     len(ParseEditorSteps(text)),
		ScenarioTitle: featureTitleFromText(text),
		Issues:        ValidateFeatureContent(text),
		Hints:         AnalyzeScenarioHints(text),
	}
}

func featureTitleFromText(text string) string {
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		lower := strings.ToLower(line)
		for _, prefix := range []string{"функциональность:", "функция:", "feature:"} {
			if strings.HasPrefix(lower, prefix) {
				return strings.TrimSpace(line[len(prefix):])
			}
		}
	}
	return ""
}
