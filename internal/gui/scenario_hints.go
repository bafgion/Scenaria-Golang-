package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

const playwrightChainSep = " >> "

var (
	fragileNthRE     = regexp.MustCompile(`(?i)nth-of-type`)
	fragileInputTail = regexp.MustCompile(`>\s*input\s*$`)
)

var assertKinds = map[string]struct{}{
	"assert-text": {}, "assert-visible": {}, "assert-hidden": {}, "assert-url": {},
	"assert-url-contains": {}, "assert-tab-count": {}, "assert-download-contains": {},
}

var interactiveAfterGoto = map[string]struct{}{
	"click": {}, "double-click": {}, "fill": {}, "fill-generated": {}, "type": {},
	"select": {}, "check": {}, "uncheck": {}, "download-click": {},
}

var waitKinds = map[string]struct{}{
	"wait": {}, "wait-visible": {}, "wait-hidden": {},
}

type ScenarioHintDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StepIndex   int    `json:"stepIndex"`
	Line        int    `json:"line"`
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

type parsedScenarioStep struct {
	line   scenarioStepLine
	action stepdsl.Action
	err    error
}

func AnalyzeScenarioHints(text string) []ScenarioHintDTO {
	steps := parseScenarioSteps(text)
	hints := make([]ScenarioHintDTO, 0)
	for i, step := range steps {
		if step.err != nil {
			continue
		}
		action := step.action
		switch action.Kind {
		case "click", "double-click":
			if hint := menuHoverHint(steps, i, action.Value1); hint != nil {
				hints = append(hints, *hint)
			}
			if hint := divClickHint(steps, i, action.Value1); hint != nil {
				hints = append(hints, *hint)
			}
		case "goto":
			if i > 0 {
				prev := steps[i-1]
				if prev.err == nil && prev.action.Kind == "goto" {
					hints = append(hints, scenarioHint(
						"duplicate_goto", "Два шага «открыт» подряд", i, steps, "warning", true,
					))
				}
			}
			if hint := gotoNoWaitHint(steps, i); hint != nil {
				hints = append(hints, *hint)
			}
		case "fill", "fill-generated", "type":
			if hint := fillNoAssertHint(steps, i); hint != nil {
				hints = append(hints, *hint)
			}
		}
		if sel := selectorFromAction(action); sel != "" {
			if isFragileSelector(sel) {
				hints = append(hints, scenarioHint(
					"fragile_selector", "Хрупкий CSS-селектор — добавьте data-testid или id",
					i, steps, "warning", false,
				))
			}
			if strings.Count(sel, playwrightChainSep) >= 2 {
				hints = append(hints, scenarioHint(
					"long_chain", "Длинная цепочка селекторов (>>) — возможно, стоит разбить",
					i, steps, "info", false,
				))
			}
		}
	}
	return hints
}

func scenarioHint(id, title string, index int, steps []parsedScenarioStep, severity string, autoFixable bool) ScenarioHintDTO {
	line := 1
	if index >= 0 && index < len(steps) {
		line = steps[index].line.lineNo
	}
	return ScenarioHintDTO{
		ID: id, Title: title, StepIndex: index, Line: line,
		Severity: severity, AutoFixable: autoFixable,
	}
}

func parseScenarioSteps(text string) []parsedScenarioStep {
	lines := parseScenarioStepLines(text)
	out := make([]parsedScenarioStep, 0, len(lines))
	for _, line := range lines {
		action, err := stepdsl.Parse(gherkin.Step{Text: line.body})
		out = append(out, parsedScenarioStep{line: line, action: action, err: err})
	}
	return out
}

func menuHoverHint(steps []parsedScenarioStep, index int, selector string) *ScenarioHintDTO {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return nil
	}
	hasChain := strings.Contains(selector, playwrightChainSep)
	if !hasChain {
		return nil
	}
	if index > 0 {
		prev := steps[index-1]
		if prev.err == nil && prev.action.Kind == "hover" {
			return nil
		}
	}
	return &ScenarioHintDTO{
		ID:          "menu_hover",
		Title:       "Клик по меню без предшествующего «навожу»",
		StepIndex:   index,
		Line:        steps[index].line.lineNo,
		Severity:    "warning",
		AutoFixable: true,
	}
}

// SplitPlaywrightChainSelector splits a chained Playwright selector into hover and click parts.
func SplitPlaywrightChainSelector(selector string) (hover string, click string, ok bool) {
	text := strings.TrimSpace(selector)
	if !strings.Contains(text, playwrightChainSep) {
		return "", "", false
	}
	parts := make([]string, 0)
	for _, part := range strings.Split(text, playwrightChainSep) {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], parts[len(parts)-1], true
}

// ProposeMenuHoverSelectors returns hover/click selectors for a menu click step.
func ProposeMenuHoverSelectors(clickSelector, hoverSelector string) (hover string, click string, ok bool) {
	clickSelector = strings.TrimSpace(clickSelector)
	hoverSelector = strings.TrimSpace(hoverSelector)
	if clickSelector == "" {
		return "", "", false
	}
	if hover, click, ok := SplitPlaywrightChainSelector(clickSelector); ok {
		if hoverSelector != "" {
			hover = hoverSelector
		}
		return hover, click, true
	}
	if hoverSelector != "" {
		return hoverSelector, clickSelector, true
	}
	return "", "", false
}

func divClickHint(steps []parsedScenarioStep, index int, selector string) *ScenarioHintDTO {
	lower := strings.ToLower(strings.TrimSpace(selector))
	if !strings.Contains(lower, "div:has-text") {
		return nil
	}
	if strings.Contains(lower, "button") || strings.Contains(lower, " a:") || strings.HasPrefix(lower, "a:") {
		return nil
	}
	hint := scenarioHint(
		"div_click", "Клик по div:has-text — уточните до button или a",
		index, steps, "info", false,
	)
	return &hint
}

func fillNoAssertHint(steps []parsedScenarioStep, index int) *ScenarioHintDTO {
	for j := index + 1; j < len(steps) && j <= index+3; j++ {
		if steps[j].err != nil {
			continue
		}
		if _, ok := assertKinds[steps[j].action.Kind]; ok {
			return nil
		}
	}
	hint := scenarioHint(
		"fill_no_assert", "Ввод без проверки в следующих шагах",
		index, steps, "info", true,
	)
	return &hint
}

func gotoNoWaitHint(steps []parsedScenarioStep, index int) *ScenarioHintDTO {
	if index+1 >= len(steps) {
		return nil
	}
	next := steps[index+1]
	if next.err != nil {
		return nil
	}
	if _, ok := interactiveAfterGoto[next.action.Kind]; !ok {
		return nil
	}
	if _, ok := waitKinds[next.action.Kind]; ok {
		return nil
	}
	hint := scenarioHint(
		"goto_no_wait", "После перехода сразу действие — добавьте «жду» или «жду появления»",
		index, steps, "info", true,
	)
	return &hint
}

func selectorFromAction(action stepdsl.Action) string {
	switch action.Kind {
	case "click", "hover", "double-click", "download-click", "press":
		return action.Value1
	case "fill", "type":
		if action.Value2 != "" {
			return action.Value2
		}
		return action.Value1
	case "fill-generated":
		return action.Value2
	case "press-in":
		return action.Value2
	default:
		return ""
	}
}

func isFragileSelector(selector string) bool {
	lower := strings.ToLower(strings.TrimSpace(selector))
	if lower == "" {
		return false
	}
	if fragileNthRE.MatchString(lower) {
		return true
	}
	if fragileInputTail.MatchString(lower) {
		return true
	}
	if strings.Contains(lower, "[data-testid=") || strings.Contains(lower, "#") {
		return false
	}
	if strings.Contains(lower, "nth-child") {
		return true
	}
	return false
}

func ApplyScenarioHintFix(req ScenarioHintFixRequest) RefactorResult {
	switch req.HintID {
	case "menu_hover":
		return applyMenuHoverFix(req.Text, req.StepIndex)
	case "duplicate_goto":
		return removeScenarioStepLine(req.Text, req.StepIndex)
	case "goto_no_wait":
		return applyGotoNoWaitFix(req.Text, req.StepIndex)
	case "fill_no_assert":
		return applyFillNoAssertFix(req.Text, req.StepIndex)
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
	var hoverSel, clickSel string
	if len(parts) < 2 {
		if hover, click, ok := ProposeMenuHoverSelectors(action.Value1, ""); ok {
			hoverSel, clickSel = hover, click
		} else {
			return RefactorResult{Text: text}
		}
	} else {
		hoverSel = strings.TrimSpace(parts[0])
		clickSel = strings.TrimSpace(parts[len(parts)-1])
	}
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

func applyGotoNoWaitFix(text string, stepIndex int) RefactorResult {
	steps := parseScenarioSteps(text)
	if stepIndex < 0 || stepIndex >= len(steps) || stepIndex+1 >= len(steps) {
		return RefactorResult{Text: text}
	}
	gotoStep := steps[stepIndex]
	next := steps[stepIndex+1]
	if gotoStep.err != nil || gotoStep.action.Kind != "goto" || next.err != nil {
		return RefactorResult{Text: text}
	}

	waitBody := `жду 1 с`
	if sel := selectorFromAction(next.action); sel != "" {
		waitBody = fmt.Sprintf(`жду появления "%s"`, escapeGherkinString(sel))
	}
	waitLine := gotoStep.line.indent + gotoStep.line.keyword + waitBody

	all := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	gotoLineNo := gotoStep.line.lineNo - 1
	if gotoLineNo < 0 || gotoLineNo >= len(all) {
		return RefactorResult{Text: text}
	}
	updated := append([]string{}, all[:gotoLineNo+1]...)
	updated = append(updated, waitLine)
	updated = append(updated, all[gotoLineNo+1:]...)
	return RefactorResult{Text: joinLines(updated, text), Count: 1}
}

func applyFillNoAssertFix(text string, stepIndex int) RefactorResult {
	steps := parseScenarioSteps(text)
	if stepIndex < 0 || stepIndex >= len(steps) {
		return RefactorResult{Text: text}
	}
	step := steps[stepIndex]
	if step.err != nil {
		return RefactorResult{Text: text}
	}

	var assertBody string
	switch step.action.Kind {
	case "fill", "type":
		value := step.action.Value1
		selector := step.action.Value2
		if value == "" || selector == "" {
			return RefactorResult{Text: text}
		}
		assertBody = fmt.Sprintf(`проверяю текст "%s" в "%s"`, escapeGherkinString(value), escapeGherkinString(selector))
	case "fill-generated":
		selector := step.action.Value2
		if selector == "" {
			return RefactorResult{Text: text}
		}
		assertBody = fmt.Sprintf(`вижу "%s"`, escapeGherkinString(selector))
	default:
		return RefactorResult{Text: text}
	}

	assertLine := step.line.indent + step.line.keyword + assertBody
	all := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	lineNo := step.line.lineNo - 1
	if lineNo < 0 || lineNo >= len(all) {
		return RefactorResult{Text: text}
	}
	updated := append([]string{}, all[:lineNo+1]...)
	updated = append(updated, assertLine)
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
