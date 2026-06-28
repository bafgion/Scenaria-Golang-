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
	content = NormalizeFeatureText(content)
	feature := &Feature{}
	lines := strings.Split(content, "\n")

	var currentScenario *Scenario
	var currentStep *Step
	var currentExample *Example
	inBackground := false
	inExamples := false
	inDocString := false
	var pendingTags []string

	for idx, raw := range lines {
		lineNo := idx + 1
		line := strings.TrimSpace(raw)
		if inDocString {
			if line == `"""` {
				inDocString = false
				continue
			}
			if currentStep == nil {
				return nil, fmt.Errorf("line %d: docstring without step", lineNo)
			}
			if currentStep.DocString == "" {
				currentStep.DocString = raw
			} else {
				currentStep.DocString += "\n" + raw
			}
			continue
		}

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "@") {
			pendingTags = append(pendingTags, strings.Fields(line)...)
			continue
		}

		if strings.HasPrefix(line, "Функционал:") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "Функционал:"))
			if title == "" {
				return nil, fmt.Errorf("line %d: empty feature title", lineNo)
			}
			feature.Title = title
			feature.Line = lineNo
			if len(pendingTags) > 0 {
				feature.Tags = append(feature.Tags, pendingTags...)
				pendingTags = nil
			}
			currentScenario = nil
			currentStep = nil
			currentExample = nil
			inBackground = false
			inExamples = false
			continue
		}

		if strings.HasPrefix(line, "Контекст:") {
			feature.HasContextBlock = true
			inBackground = true
			currentScenario = nil
			currentStep = nil
			currentExample = nil
			inExamples = false
			continue
		}

		if strings.HasPrefix(line, "Сценарий:") || strings.HasPrefix(line, "Структура сценария:") {
			isOutline := strings.HasPrefix(line, "Структура сценария:")
			prefix := "Сценарий:"
			if isOutline {
				prefix = "Структура сценария:"
			}
			title := strings.TrimSpace(strings.TrimPrefix(line, prefix))
			if title == "" {
				return nil, fmt.Errorf("line %d: empty scenario title", lineNo)
			}
			scenario := Scenario{
				Title:     title,
				Line:      lineNo,
				IsOutline: isOutline,
			}
			if len(pendingTags) > 0 {
				scenario.Tags = append(scenario.Tags, pendingTags...)
				pendingTags = nil
			}
			feature.Scenarios = append(feature.Scenarios, scenario)
			currentScenario = &feature.Scenarios[len(feature.Scenarios)-1]
			currentStep = nil
			currentExample = nil
			inBackground = false
			inExamples = false
			continue
		}

		if strings.HasPrefix(line, "Примеры:") {
			if currentScenario == nil {
				return nil, fmt.Errorf("line %d: examples outside scenario", lineNo)
			}
			if !currentScenario.IsOutline {
				return nil, fmt.Errorf("line %d: examples are only valid for scenario outline", lineNo)
			}
			currentScenario.Examples = append(currentScenario.Examples, Example{Line: lineNo})
			currentScenario.ExamplesLine = lineNo
			currentExample = &currentScenario.Examples[len(currentScenario.Examples)-1]
			currentStep = nil
			inExamples = true
			inBackground = false
			continue
		}

		if line == `"""` {
			if currentStep == nil {
				return nil, fmt.Errorf("line %d: docstring without step", lineNo)
			}
			inDocString = true
			continue
		}

		if strings.HasPrefix(line, "|") {
			row, err := parseTableRow(line)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			if inExamples {
				if currentExample == nil {
					return nil, fmt.Errorf("line %d: examples table without section", lineNo)
				}
				currentExample.Rows = append(currentExample.Rows, row)
				continue
			}
			if currentStep == nil {
				return nil, fmt.Errorf("line %d: step table without step", lineNo)
			}
			currentStep.Table = append(currentStep.Table, row)
			continue
		}

		inExamples = false
		currentExample = nil

		step, ok := parseStep(line, lineNo)
		if !ok {
			return nil, fmt.Errorf("line %d: unsupported statement %q", lineNo, line)
		}
		step.Indent = lineIndentLevel(raw)

		switch {
		case inBackground:
			feature.Background = append(feature.Background, step)
			currentStep = &feature.Background[len(feature.Background)-1]
		case currentScenario != nil:
			currentScenario.Steps = append(currentScenario.Steps, step)
			currentStep = &currentScenario.Steps[len(currentScenario.Steps)-1]
		default:
			return nil, fmt.Errorf("line %d: step outside scenario/context", lineNo)
		}
	}

	if inDocString {
		return nil, fmt.Errorf("line %d: unterminated docstring", len(lines))
	}
	if len(pendingTags) > 0 {
		return nil, fmt.Errorf("line %d: dangling tags without section", len(lines))
	}

	if err := FinalizeStepTrees(feature); err != nil {
		return nil, err
	}

	return feature, nil
}

func parseStep(line string, lineNo int) (Step, bool) {
	for _, keyword := range stepKeywords {
		if strings.HasPrefix(line, keyword+" ") || line == keyword {
			text := strings.TrimSpace(strings.TrimPrefix(line, keyword))
			text = NormalizeLegacyHasTextEscapes(text)
			return Step{
				Keyword: keyword,
				Text:    text,
				Line:    lineNo,
			}, true
		}
	}
	text := strings.TrimSpace(line)
	if testClientRe.MatchString(text) {
		return Step{Keyword: "Допустим", Text: text, Line: lineNo}, true
	}
	if header, err := detectBlockHeader(Step{Text: text, Line: lineNo}); err == nil && header != nil {
		return Step{Keyword: "*", Text: text, Line: lineNo}, true
	}
	return Step{}, false
}

func parseTableRow(line string) ([]string, error) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return nil, fmt.Errorf("table row must start and end with |")
	}
	parts := strings.Split(trimmed, "|")
	if len(parts) < 3 {
		return nil, fmt.Errorf("empty table row")
	}

	row := make([]string, 0, len(parts)-2)
	for i := 1; i < len(parts)-1; i++ {
		row = append(row, strings.TrimSpace(parts[i]))
	}
	return row, nil
}
