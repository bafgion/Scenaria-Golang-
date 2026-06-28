package gherkin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var quoted = `((?:\\.|[^"])*)`

var (
	ifHeaderRe      = regexp.MustCompile(`(?i)^если\s+(.+)$`)
	repeatHeaderRe  = regexp.MustCompile(`(?i)^повторяю\s+(\d+)\s+раз(?:а)?$`)
	whileHeaderRe   = regexp.MustCompile(`(?i)^пока\s+(.+)$`)
	forEachHeaderRe = regexp.MustCompile(`(?i)^для\s+каждого\s+"` + quoted + `"\s+как\s+"` + quoted + `"$`)
	visibleCondRe   = regexp.MustCompile(`(?i)^вижу\s+"` + quoted + `"$`)
	hiddenCondRe    = regexp.MustCompile(`(?i)^не\s+вижу\s+"` + quoted + `"$`)
	urlContainsRe   = regexp.MustCompile(`(?i)^url\s+содержит\s+"` + quoted + `"$`)
	pageTextRe      = regexp.MustCompile(`(?i)^текст\s+на\s+странице\s+"` + quoted + `"$`)
	testClientRe    = regexp.MustCompile(`(?i)^я\s+подключаю\s+TestClient\s+"` + quoted + `"$`)
)

type blockHeader struct {
	Kind     string
	Condition *Condition
	Count    int
	Selector string
	Variable string
}

func parseCondition(expr string) (*Condition, error) {
	expr = strings.TrimSpace(expr)
	switch {
	case visibleCondRe.MatchString(expr):
		return &Condition{Type: "visible", Selector: unquote(visibleCondRe.FindStringSubmatch(expr)[1])}, nil
	case hiddenCondRe.MatchString(expr):
		return &Condition{Type: "hidden", Selector: unquote(hiddenCondRe.FindStringSubmatch(expr)[1])}, nil
	case urlContainsRe.MatchString(expr):
		return &Condition{Type: "url_contains", Value: unquote(urlContainsRe.FindStringSubmatch(expr)[1])}, nil
	case pageTextRe.MatchString(expr):
		return &Condition{Type: "page_text", Value: unquote(pageTextRe.FindStringSubmatch(expr)[1])}, nil
	default:
		return nil, fmt.Errorf("unknown condition %q", expr)
	}
}

func detectBlockHeader(step Step) (*blockHeader, error) {
	body := strings.TrimSpace(step.Text)
	if groups := ifHeaderRe.FindStringSubmatch(body); groups != nil {
		cond, err := parseCondition(groups[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", step.Line, err)
		}
		return &blockHeader{Kind: BlockIf, Condition: cond}, nil
	}
	if groups := whileHeaderRe.FindStringSubmatch(body); groups != nil {
		cond, err := parseCondition(groups[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", step.Line, err)
		}
		return &blockHeader{Kind: BlockWhile, Condition: cond}, nil
	}
	if groups := repeatHeaderRe.FindStringSubmatch(body); groups != nil {
		count, _ := strconv.Atoi(groups[1])
		if count < 1 {
			count = 1
		}
		return &blockHeader{Kind: BlockRepeat, Count: count}, nil
	}
	if groups := forEachHeaderRe.FindStringSubmatch(body); groups != nil {
		return &blockHeader{
			Kind:     BlockForEach,
			Selector: unquote(groups[1]),
			Variable: unquote(groups[2]),
		}, nil
	}
	return nil, nil
}

func ParseTestClientName(steps []Step) (string, error) {
	if len(steps) == 0 {
		return "", nil
	}
	if len(steps) != 1 {
		return "", fmt.Errorf("line %d: context block must contain exactly one TestClient step", steps[0].Line)
	}
	groups := testClientRe.FindStringSubmatch(strings.TrimSpace(steps[0].Text))
	if groups == nil {
		return "", fmt.Errorf("line %d: invalid TestClient step %q", steps[0].Line, steps[0].Text)
	}
	return unquote(groups[1]), nil
}

func IsTestClientStep(step Step) bool {
	return testClientRe.MatchString(strings.TrimSpace(step.Text))
}

func unquote(s string) string {
	return strings.ReplaceAll(s, `\"`, `"`)
}
