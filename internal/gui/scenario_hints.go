package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

const playwrightChainSep = " >> "

type ScenarioHintDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StepIndex   int    `json:"stepIndex"`
	Severity    string `json:"severity"`
	AutoFixable bool   `json:"autoFixable"`
}

type ScenarioHintFixRequest struct {
	Text      string `json:"text"`
	HintID    string `json:"hintId"`
	StepIndex int    `json:"stepIndex"`
}

type scenarioStepLine struct {
	lineNo  int
	raw     string
	indent  string
	keyword string
	body    string
}

func AnalyzeScenarioHints(text string) []ScenarioHintDTO {
	lines := parseScenarioStepLines(text)
	hints := make([]ScenarioHintDTO, 0)
	for i, line := range lines {
		action, err := stepdsl.Parse(gherkin.Step{Text: line.body})
		if err != nil {
			continue
		}
		switch action.Kind {
		case "click", "double-click":
			if hint := menuHoverHint(lines, i, action.Value1); hint != nil {
				hints = append(hints, *hint)
			}
		case "goto":
			if i > 0 {
				prev, err := stepdsl.Parse(gherkin.Step{Text: lines[i-1].body})
				if err == nil && prev.Kind == "goto" {
					hints = append(hints, ScenarioHintDTO{
						ID:          "duplicate_goto",
						Title:       "Два шага «открыт» подряд",
						StepIndex:   i,
						Severity:    "warning",
						AutoFixable: true,
					})
				}
			}
		}
	}
	return hints
}

func menuHoverHint(lines []scenarioStepLine, index int, selector string) *ScenarioHintDTO {
	selector = strings.TrimSpace(selector)
	if !strings.Contains(selector, playwrightChainSep) {
		return nil
	}
	if index > 0 {
		prev, err := stepdsl.Parse(gherkin.Step{Text: lines[index-1].body})
		if err == nil && prev.Kind == "hover" {
			return nil
		}
	}
	return &ScenarioHintDTO{
		ID:          "menu_hover",
		Title:       "Клик по меню без предшествующего «навожу»",
		StepIndex:   index,
		Severity:    "warning",
		AutoFixable: true,
	}
}

func ApplyScenarioHintFix(req ScenarioHintFixRequest) RefactorResult {
	switch req.HintID {
	case "menu_hover":
		return applyMenuHoverFix(req.Text, req.StepIndex)
	case "duplicate_goto":
		return removeScenarioStepLine(req.Text, req.StepIndex)
	default:
		return RefactorResult{Text: req.Text}
	}
}

func applyMenuHoverFix(text string, stepIndex int) RefactorResult {
	lines := parseScenarioStepLines(text)
	if stepIndex < 0 || stepIndex >= len(lines) {
		return RefactorResult{Text: text}
	}
	target := lines[stepIndex]
	action, err := stepdsl.Parse(gherkin.Step{Text: target.body})
	if err != nil || (action.Kind != "click" && action.Kind != "double-click") {
		return RefactorResult{Text: text}
	}
	parts := strings.Split(action.Value1, playwrightChainSep)
	if len(parts) < 2 {
		return RefactorResult{Text: text}
	}
	hoverSel := strings.TrimSpace(parts[0])
	clickSel := strings.TrimSpace(parts[len(parts)-1])
	if hoverSel == "" || clickSel == "" {
		return RefactorResult{Text: text}
	}

	verb := "нажимаю"
	if action.Kind == "double-click" {
		verb = "дважды нажимаю"
	}
	hoverBody := fmt.Sprintf(`навожу "%s"`, escapeGherkinString(hoverSel))
	clickBody := fmt.Sprintf(`%s "%s"`, verb, escapeGherkinString(clickSel))
	hoverLine := target.indent + target.keyword + hoverBody
	clickLine := target.indent + target.keyword + clickBody

	all := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	lineNo := target.lineNo - 1
	if lineNo < 0 || lineNo >= len(all) {
		return RefactorResult{Text: text}
	}
	updated := append([]string{}, all[:lineNo]...)
	updated = append(updated, hoverLine, clickLine)
	updated = append(updated, all[lineNo+1:]...)
	return RefactorResult{Text: joinLines(updated, text), Count: 1}
}

func removeScenarioStepLine(text string, stepIndex int) RefactorResult {
	lines := parseScenarioStepLines(text)
	if stepIndex < 0 || stepIndex >= len(lines) {
		return RefactorResult{Text: text}
	}
	all := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	lineNo := lines[stepIndex].lineNo - 1
	if lineNo < 0 || lineNo >= len(all) {
		return RefactorResult{Text: text}
	}
	updated := append(all[:lineNo], all[lineNo+1:]...)
	return RefactorResult{Text: joinLines(updated, text), Count: 1}
}

func parseScenarioStepLines(text string) []scenarioStepLine {
	rawLines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	out := make([]scenarioStepLine, 0)
	for i, raw := range rawLines {
		stripped := strings.TrimSpace(raw)
		if stripped == "" || strings.HasPrefix(stripped, "#") {
			continue
		}
		if isScenarioStructureLine(stripped) {
			continue
		}
		keyword, body := splitStepKeyword(stripped)
		if body == "" && keyword == "" {
			continue
		}
		prefix := keyword
		if prefix != "" {
			prefix += " "
		}
		out = append(out, scenarioStepLine{
			lineNo:  i + 1,
			raw:     raw,
			indent:  leadingIndent(raw),
			keyword: prefix,
			body:    body,
		})
	}
	return out
}

func escapeGherkinString(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
