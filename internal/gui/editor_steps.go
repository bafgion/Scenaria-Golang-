package gui

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

// EditorStepRow is one Gherkin step line for the steps strip (Python StepsStrip).
type EditorStepRow struct {
	Line    int    `json:"line"`
	Keyword string `json:"keyword"`
	Text    string `json:"text"`
	Kind    string `json:"kind"`
	Action  string `json:"action"`
	Element string `json:"element"`
	Value   string `json:"value"`
	Error   string `json:"error"`
}

func ParseEditorSteps(text string) []EditorStepRow {
	return editorStepsFromAnalysis(analyzeEditorSteps(text))
}

func splitStepKeyword(line string) (keyword, rest string) {
	for _, kw := range []string{
		"Допустим ", "Дано ", "Когда ", "Тогда ", "И ", "Но ",
		"Given ", "When ", "Then ", "And ", "But ",
	} {
		if strings.HasPrefix(line, kw) {
			return strings.TrimSpace(kw), strings.TrimSpace(strings.TrimPrefix(line, kw))
		}
	}
	return "", line
}

func actionDisplayName(kind string) string {
	switch kind {
	case "goto":
		return "Открыть"
	case "click", "double-click":
		return "Нажать"
	case "hover":
		return "Навести"
	case "fill", "type":
		return "Ввести"
	case "assert-text", "assert-text-regex", "assert-visible", "assert-enabled", "assert-disabled", "assert-selected":
		return "Проверить"
	case "press", "press-in":
		return "Клавиша"
	default:
		if kind == "" {
			return "Шаг"
		}
		return kind
	}
}

func actionFields(a stepdsl.Action) (element, value string) {
	switch a.Kind {
	case "goto", "assert-visible", "assert-enabled", "assert-disabled", "assert-selected", "remember-text", "remember-url":
		value = a.Value1
	case "assert-text", "assert-text-regex":
		value = a.Value1
		element = a.Value2
	case "click", "hover", "double-click", "fill", "type", "press", "download-click":
		element = a.Value1
		if a.Value2 != "" {
			value = a.Value2
		}
	case "press-in":
		element = a.Value2
		value = a.Value1
	case "remember-field":
		element = a.Value1
		value = a.Value2
	default:
		if a.Value1 != "" && a.Value2 != "" {
			element = a.Value1
			value = a.Value2
		} else if a.Value1 != "" {
			value = a.Value1
		}
	}
	return element, value
}
