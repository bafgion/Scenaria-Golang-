package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type ValidationIssue struct {
	Line       int    `json:"line"`
	Message    string `json:"message"`
	Selector   string `json:"selector,omitempty"`
	Status     string `json:"status,omitempty"`
	StepText   string `json:"stepText,omitempty"`
	Mode       string `json:"mode,omitempty"`
	ActionKind string `json:"actionKind,omitempty"`
	MatchCount int    `json:"matchCount,omitempty"`
	Limitation string `json:"limitation,omitempty"`
}

func ValidateFeatureContent(text string) []ValidationIssue {
	feature, err := gherkin.ParseFeature(text)
	if err == nil {
		return validateParsedFeature(feature)
	}
	return validateFeatureLines(text)
}

func validateParsedFeature(feature *gherkin.Feature) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	for _, issue := range gherkin.ValidateFeature(feature) {
		issues = append(issues, ValidationIssue{
			Line:    issue.Line,
			Message: issue.Message,
		})
	}
	validateStepList(feature.Background, &issues)
	for _, scenario := range feature.Scenarios {
		validateStepList(scenario.Steps, &issues)
	}
	return issues
}

func validateFeatureLines(text string) []ValidationIssue {
	analysis := analyzeEditorSteps(text)
	issues := make([]ValidationIssue, 0, len(analysis))
	for _, step := range analysis {
		if step.TestClient {
			continue
		}
		if step.ParseErr != nil {
			issues = append(issues, ValidationIssue{
				Line:    step.Line,
				Message: fmt.Sprintf("%v", step.ParseErr),
			})
		}
	}
	return issues
}

func validateStepList(steps []gherkin.Step, issues *[]ValidationIssue) {
	for _, step := range steps {
		if step.Block != "" || step.Condition != nil {
			validateStepList(step.Children, issues)
			continue
		}
		if strings.TrimSpace(step.Text) == "" {
			validateStepList(step.Children, issues)
			continue
		}
		if gherkin.IsTestClientStep(step) {
			validateStepList(step.Children, issues)
			continue
		}
		if _, err := stepdsl.Parse(step); err != nil {
			*issues = append(*issues, ValidationIssue{
				Line:    step.Line,
				Message: fmt.Sprintf("%v", err),
			})
		}
		validateStepList(step.Children, issues)
	}
}

func isScenarioStructureLine(line string) bool {
	lower := strings.ToLower(line)
	prefixes := []string{
		"функция:", "функциональность:", "функционал:", "feature:", "сценарий:", "scenario:",
		"структура сценария:", "scenario outline:", "примеры:", "examples:",
		"контекст:", "background:",
		"если ", "if ", "повторяю ", "repeat ", "пока ", "while ", "для каждого ", "for each ", "конец", "end if",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func isTagLine(line string) bool {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if !strings.HasPrefix(part, "@") {
			return false
		}
	}
	return true
}

func isTableLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "|")
}
