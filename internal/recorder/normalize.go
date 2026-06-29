package recorder

import (
	"regexp"
	"strings"
)

type RecordedStep struct {
	Action     string
	Selector   string
	Value      string
	Text       string
	Context    string
	InputType  string
	Generator  string
}

func NormalizeSteps(steps []RecordedStep) []RecordedStep {
	if len(steps) == 0 {
		return nil
	}
	out := make([]RecordedStep, 0, len(steps))
	fillBySel := map[string]int{}
	selectBySel := map[string]int{}
	for _, step := range steps {
		step = upgradeFillSelector(step)
		switch step.Action {
		case "fill":
			sel := step.Selector
			if sel == "" {
				out = append(out, step)
				continue
			}
			if idx, ok := fillBySel[sel]; ok {
				out[idx] = step
			} else {
				fillBySel[sel] = len(out)
				out = append(out, step)
			}
		case "select":
			sel := step.Selector
			if sel == "" {
				out = append(out, step)
				continue
			}
			if idx, ok := selectBySel[sel]; ok {
				out[idx] = step
			} else {
				selectBySel[sel] = len(out)
				out = append(out, step)
			}
		case "click":
			if len(out) > 0 && out[len(out)-1].Action == "click" && out[len(out)-1].Selector == step.Selector {
				continue
			}
			if len(out) > 0 && (out[len(out)-1].Action == "check" || out[len(out)-1].Action == "uncheck") && out[len(out)-1].Selector == step.Selector {
				continue
			}
			out = append(out, upgradeClickSelector(step))
			fillBySel = map[string]int{}
			selectBySel = map[string]int{}
		case "scroll-to":
			if len(out) > 0 && out[len(out)-1].Action == "scroll-to" && out[len(out)-1].Selector == step.Selector {
				continue
			}
			out = append(out, step)
		case "goto":
			if len(out) > 0 && out[len(out)-1].Action == "goto" && out[len(out)-1].Value == step.Value {
				continue
			}
			out = append(out, step)
			fillBySel = map[string]int{}
			selectBySel = map[string]int{}
		default:
			out = append(out, step)
		}
	}
	return out
}

func upgradeClickSelector(step RecordedStep) RecordedStep {
	if step.Action != "click" {
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
	step.Selector = `button:has-text("` + escapeSelectorText(label) + `")`
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

func contextualButtonSelector(context, label string) string {
	if len(context) > 80 {
		context = context[:60]
	}
	if len(label) > 60 {
		label = label[:40]
	}
	return `div:has-text("` + escapeSelectorText(context) + `") >> button:has-text("` + escapeSelectorText(label) + `")`
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

func escapeSelectorText(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}
