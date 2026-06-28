//go:build desktop

package desktop

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func validateFeatureSteps(text string) string {
	lines := strings.Split(text, "\n")
	issues := make([]string, 0)
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
			issues = append(issues, fmt.Sprintf("строка %d: %v", i+1, err))
			if len(issues) >= 3 {
				break
			}
		}
	}
	if len(issues) == 0 {
		return "Шаги: OK"
	}
	return strings.Join(issues, "; ")
}

func isScenarioStructureLine(line string) bool {
	lower := strings.ToLower(line)
	prefixes := []string{
		"функция:", "функциональность:", "feature:", "сценарий:", "scenario:",
		"структура сценария:", "scenario outline:", "примеры:", "examples:",
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
