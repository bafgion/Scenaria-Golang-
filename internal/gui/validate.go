package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type ValidationIssue struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
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
	lines := strings.Split(text, "\n")
	issues := make([]ValidationIssue, 0)
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if isTagLine(line) || isTableLine(line) {
			continue
		}
		if isScenarioStructureLine(line) {
			continue
		}
		stepText := stripGherkinKeyword(line)
		if stepText == "" {
			continue
		}
		if _, err := stepdsl.Parse(gherkin.Step{Line: i + 1, Text: stepText}); err != nil {
			issues = append(issues, ValidationIssue{
				Line:    i + 1,
				Message: fmt.Sprintf("%v", err),
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
		"если ", "повторяю ", "пока ", "для каждого ", "конец",
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

func stripGherkinKeyword(line string) string {
	for _, keyword := range []string{"Допустим ", "Когда ", "Тогда ", "И ", "Но ", "Given ", "When ", "Then ", "And ", "But "} {
		if strings.HasPrefix(line, keyword) {
			return strings.TrimSpace(strings.TrimPrefix(line, keyword))
		}
	}
	return line
}
