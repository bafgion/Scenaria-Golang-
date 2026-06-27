package stepdsl

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type Action struct {
	Kind   string
	Value1 string
	Value2 string
}

func Parse(step gherkin.Step) (Action, error) {
	text := strings.TrimSpace(step.Text)
	lower := strings.ToLower(text)
	values := extractQuotedValues(text)

	switch {
	case strings.HasPrefix(lower, "нажимаю клавишу "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: key press step must contain key in quotes", step.Line)
		}
		return Action{Kind: "press", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "жду ") || strings.HasPrefix(lower, "ожидаю "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: wait step must contain duration in quotes", step.Line)
		}
		return Action{Kind: "wait", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "url содержит ") || strings.HasPrefix(lower, "адрес содержит "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: url assertion step must contain expected fragment in quotes", step.Line)
		}
		return Action{Kind: "assert-url-contains", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "открываю ") || strings.HasPrefix(lower, "перехожу на ") || strings.HasPrefix(lower, "перейти на "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: navigation step must contain URL in quotes", step.Line)
		}
		return Action{Kind: "goto", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "нажимаю ") || strings.HasPrefix(lower, "кликаю "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: click step must contain selector in quotes", step.Line)
		}
		return Action{Kind: "click", Value1: values[0]}, nil
	case strings.HasPrefix(lower, "ввожу "):
		if len(values) < 2 {
			return Action{}, fmt.Errorf("line %d: fill step must contain value and selector in quotes", step.Line)
		}
		return Action{Kind: "fill", Value1: values[0], Value2: values[1]}, nil
	case strings.HasPrefix(lower, "вижу ") || strings.HasPrefix(lower, "должен видеть "):
		if len(values) < 1 {
			return Action{}, fmt.Errorf("line %d: assertion step must contain expected text in quotes", step.Line)
		}
		return Action{Kind: "assert-text", Value1: values[0]}, nil
	default:
		return Action{}, fmt.Errorf("line %d: unsupported step text %q", step.Line, step.Text)
	}
}

func ResolveURL(raw string, baseURL string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	if strings.TrimSpace(baseURL) == "" {
		return trimmed
	}
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasPrefix(trimmed, "/") {
		return base + trimmed
	}
	return base + "/" + strings.TrimLeft(trimmed, "/")
}

func extractQuotedValues(s string) []string {
	values := make([]string, 0)
	for {
		start := strings.Index(s, `"`)
		if start < 0 {
			break
		}
		s = s[start+1:]
		end := strings.Index(s, `"`)
		if end < 0 {
			break
		}
		values = append(values, s[:end])
		s = s[end+1:]
	}
	return values
}
