package gherkin

import (
	"fmt"
	"os"
	"strings"
)

var stepKeywords = []string{
	"Допустим",
	"Когда",
	"Тогда",
	"И",
	"Но",
	"*",
}

func ParseFeatureFile(path string) (*Feature, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read feature file: %w", err)
	}
	return ParseFeature(string(content))
}

func ParseFeature(content string) (*Feature, error) {
	feature := &Feature{}
	lines := strings.Split(content, "\n")

	var currentScenario *Scenario
	inBackground := false

	for idx, raw := range lines {
		lineNo := idx + 1
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "Функционал:") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "Функционал:"))
			if title == "" {
				return nil, fmt.Errorf("line %d: empty feature title", lineNo)
			}
			feature.Title = title
			currentScenario = nil
			inBackground = false
			continue
		}

		if strings.HasPrefix(line, "Контекст:") {
			inBackground = true
			currentScenario = nil
			continue
		}

		if strings.HasPrefix(line, "Сценарий:") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "Сценарий:"))
			if title == "" {
				return nil, fmt.Errorf("line %d: empty scenario title", lineNo)
			}
			feature.Scenarios = append(feature.Scenarios, Scenario{Title: title})
			currentScenario = &feature.Scenarios[len(feature.Scenarios)-1]
			inBackground = false
			continue
		}

		step, ok := parseStep(line, lineNo)
		if !ok {
			return nil, fmt.Errorf("line %d: unsupported statement %q", lineNo, line)
		}

		switch {
		case inBackground:
			feature.Background = append(feature.Background, step)
		case currentScenario != nil:
			currentScenario.Steps = append(currentScenario.Steps, step)
		default:
			return nil, fmt.Errorf("line %d: step outside scenario/context", lineNo)
		}
	}

	return feature, nil
}

func parseStep(line string, lineNo int) (Step, bool) {
	for _, keyword := range stepKeywords {
		if strings.HasPrefix(line, keyword+" ") || line == keyword {
			text := strings.TrimSpace(strings.TrimPrefix(line, keyword))
			return Step{
				Keyword: keyword,
				Text:    text,
				Line:    lineNo,
			}, true
		}
	}
	return Step{}, false
}
