package recorder

import (
	"regexp"
	"strings"
)

type RecordedStep struct {
	Action        string
	Selector      string
	Value         string
	Text          string
	Context       string
	InputType     string
	Generator     string
	HoverSelector string
	HoverText     string
}

func NormalizeSteps(steps []RecordedStep) []RecordedStep {
	if len(steps) == 0 {
		return nil
	}
	prepared := make([]RecordedStep, len(steps))
	for i, step := range steps {
		if step.Action == "fill" {
			prepared[i] = upgradeFillSelector(step)
		} else {
			prepared[i] = step
		}
	}
	merged := collapseFills(prepared)
	merged = collapseClickBeforeFill(merged)
	merged = collapseCheckboxNoise(merged)
	for i, step := range merged {
		merged[i] = upgradeCanvasStep(step)
		merged[i] = upgradeClickSelector(merged[i])
	}
	merged = collapseDuplicateScrolls(merged)
	merged = collapseDuplicateHovers(merged)
	merged = collapseRedundantHoverBeforeClick(merged)
	merged = collapseDuplicateClicks(merged)
	merged = dropDuplicateGotos(merged)
	return merged
}

func collapseFills(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	fillBySel := map[string]int{}
	selectBySel := map[string]int{}
	for _, step := range steps {
		switch step.Action {
		case "fill":
			sel := step.Selector
			if sel == "" {
				result = append(result, step)
				continue
			}
			if idx, ok := fillBySel[sel]; ok {
				result[idx] = step
			} else {
				fillBySel[sel] = len(result)
				result = append(result, step)
			}
		case "select":
			sel := step.Selector
			if sel == "" {
				result = append(result, step)
				continue
			}
			if idx, ok := selectBySel[sel]; ok {
				result[idx] = step
			} else {
				selectBySel[sel] = len(result)
				result = append(result, step)
			}
		case "click", "goto":
			fillBySel = map[string]int{}
			selectBySel = map[string]int{}
			result = append(result, step)
		default:
			result = append(result, step)
		}
	}
	return result
}

func collapseClickBeforeFill(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if len(result) > 0 && step.Action == "fill" && result[len(result)-1].Action == "click" {
			clickSel := result[len(result)-1].Selector
			fillSel := step.Selector
			if clickSel == fillSel || selectorIsFragile(clickSel) {
				result = result[:len(result)-1]
			}
		}
		if step.Action == "fill" {
			step = upgradeFillSelector(step)
		}
		result = append(result, step)
	}
	return result
}

func collapseCheckboxNoise(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if step.Action == "fill" && (step.InputType == "checkbox" || isCheckboxValue(step.Value)) {
			checked := isCheckedValue(step.Value)
			for len(result) > 0 && result[len(result)-1].Action == "click" {
				result = result[:len(result)-1]
			}
			normalized := upgradeCheckboxSelector(RecordedStep{
				Action:   checkboxAction(checked),
				Selector: step.Selector,
				Text:     step.Text,
			})
			if len(result) > 0 &&
				result[len(result)-1].Action == normalized.Action &&
				result[len(result)-1].Selector == normalized.Selector {
				continue
			}
			result = append(result, normalized)
			continue
		}
		if step.Action == "check" || step.Action == "uncheck" {
			for len(result) > 0 && result[len(result)-1].Action == "click" {
				result = result[:len(result)-1]
			}
			upgraded := upgradeCheckboxSelector(step)
			if len(result) > 0 &&
				result[len(result)-1].Action == upgraded.Action &&
				result[len(result)-1].Selector == upgraded.Selector {
				continue
			}
			result = append(result, upgraded)
			continue
		}
		result = append(result, step)
	}
	return result
}

func collapseDuplicateScrolls(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if step.Action == "scroll-to" && len(result) > 0 {
			last := result[len(result)-1]
			if last.Action == "scroll-to" && last.Selector == step.Selector {
				continue
			}
		}
		result = append(result, step)
	}
	return result
}

func collapseDuplicateClicks(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if step.Action == "click" && len(result) > 0 {
			last := result[len(result)-1]
			if last.Action == "click" && last.Selector == step.Selector {
				continue
			}
			if (last.Action == "check" || last.Action == "uncheck") && last.Selector == step.Selector {
				continue
			}
		}
		result = append(result, step)
	}
	return result
}

func collapseDuplicateHovers(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	seenSinceGoto := map[string]struct{}{}
	for _, step := range steps {
		if step.Action == "goto" {
			seenSinceGoto = map[string]struct{}{}
			result = append(result, step)
			continue
		}
		if step.Action == "hover" {
			if _, ok := seenSinceGoto[step.Selector]; ok {
				continue
			}
			seenSinceGoto[step.Selector] = struct{}{}
		}
		result = append(result, step)
	}
	return result
}

func collapseRedundantHoverBeforeClick(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if step.Action == "click" && len(result) > 0 {
			last := result[len(result)-1]
			if last.Action == "hover" && last.Selector == step.Selector {
				result = result[:len(result)-1]
			}
		}
		result = append(result, step)
	}
	return result
}

func dropDuplicateGotos(steps []RecordedStep) []RecordedStep {
	result := make([]RecordedStep, 0, len(steps))
	for _, step := range steps {
		if step.Action == "goto" && len(result) > 0 {
			last := result[len(result)-1]
			if last.Action == "goto" && last.Value == step.Value {
				continue
			}
		}
		result = append(result, step)
	}
	return result
}

func upgradeClickSelector(step RecordedStep) RecordedStep {
	if step.Action != "click" {
		return step
	}
	if strings.Contains(strings.ToLower(step.Selector), "canvas") {
		return step
	}
	ctx := strings.TrimSpace(step.Context)
	label := strings.TrimSpace(step.Text)
	if len(label) >= 2 && len(ctx) >= 6 {
		step.Selector = contextualButtonSelector(ctx, label)
		return step
	}
	if !selectorIsFragile(step.Selector) && step.Selector != "" {
		return step
	}
	if len(label) < 2 {
		return step
	}
	step.Selector = `text="` + escapeSelectorText(label) + `"`
	return step
}

func upgradeFillSelector(step RecordedStep) RecordedStep {
	if step.Action != "fill" || step.Selector == "" {
		return step
	}
	if !selectorIsFragile(step.Selector) && !isGenericPlaceholder(step.Selector) {
		return step
	}
	label := strings.TrimSpace(strings.TrimRight(step.Text, "*"))
	if len(label) < 2 {
		return step
	}
	step.Selector = `label:has-text("` + escapeSelectorText(label) + `")`
	return step
}

func upgradeCheckboxSelector(step RecordedStep) RecordedStep {
	if step.Action != "check" && step.Action != "uncheck" {
		return step
	}
	if !selectorIsFragile(step.Selector) {
		return step
	}
	text := strings.TrimSpace(step.Text)
	if len(text) < 4 {
		return step
	}
	snippet := text
	if len(snippet) > 60 {
		snippet = snippet[:40]
	}
	if len(snippet) < 4 {
		return step
	}
	step.Selector = `label:has-text("` + escapeSelectorText(snippet) + `")`
	return step
}

func upgradeCanvasStep(step RecordedStep) RecordedStep {
	selector := step.Selector
	if !strings.Contains(strings.ToLower(selector), "canvas") {
		return step
	}
	if step.Action == "click" && selectorIsFragile(selector) {
		return RecordedStep{
			Action:   "draw-signature",
			Selector: bestCanvasSelector(step),
			Text:     step.Text,
		}
	}
	if step.Action == "draw-signature" && selectorIsFragile(selector) {
		step.Selector = bestCanvasSelector(step)
	}
	return step
}

func bestCanvasSelector(step RecordedStep) string {
	text := strings.TrimSpace(step.Text)
	if len(text) >= 4 {
		snippet := text
		if len(snippet) > 40 {
			snippet = snippet[:30]
		}
		return `div:has-text("` + escapeSelectorText(snippet) + `") canvas`
	}
	return "canvas"
}

func contextualButtonSelector(context, label string) string {
	if len(context) > 80 {
		context = context[:60]
	}
	if len(label) > 60 {
		label = label[:40]
	}
	return `text="` + escapeSelectorText(context) + `" >> text="` + escapeSelectorText(label) + `"`
}

func selectorIsFragile(selector string) bool {
	lower := strings.ToLower(selector)
	if strings.Contains(lower, "nth-of-type") {
		return true
	}
	return regexp.MustCompile(`>\s*input\s*$`).MatchString(lower)
}

func isGenericPlaceholder(selector string) bool {
	match := regexp.MustCompile(`(?i)placeholder="([^"]+)"`).FindStringSubmatch(selector)
	if match == nil {
		return false
	}
	ph := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(match[1])), " ", "")
	generic := map[string]struct{}{
		"дд.мм.гггг": {}, "дд/мм/гггг": {}, "dd.mm.yyyy": {}, "mm/dd/yyyy": {},
		"__.__.____": {}, "--.--.----": {},
	}
	if _, ok := generic[ph]; ok {
		return true
	}
	return regexp.MustCompile(`(?i)^[дd_.\-/]{4,}$`).MatchString(ph)
}

func isCheckboxValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "on", "off":
		return true
	default:
		return false
	}
}

func isCheckedValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "on", "true", "1", "yes":
		return true
	default:
		return false
	}
}

func checkboxAction(checked bool) string {
	if checked {
		return "check"
	}
	return "uncheck"
}

func escapeSelectorText(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}
