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
	lines := strings.Split(text, "\n")
	issues := make([]ValidationIssue, 0)
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
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

func stripGherkinKeyword(line string) string {
	for _, keyword := range []string{"Допустим ", "Когда ", "Тогда ", "И ", "Но ", "Given ", "When ", "Then ", "And ", "But "} {
		if strings.HasPrefix(line, keyword) {
			return strings.TrimSpace(strings.TrimPrefix(line, keyword))
		}
	}
	return line
}
