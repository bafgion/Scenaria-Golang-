package recorder

import "strings"

var genericTagSelectors = map[string]struct{}{
	"img": {}, "div": {}, "span": {}, "p": {}, "section": {}, "main": {},
	"body": {}, "html": {}, "a": {}, "button": {}, "input": {}, "canvas": {},
	"ul": {}, "li": {}, "form": {}, "label": {},
}

func isGenericTagSelector(selector string) bool {
	_, ok := genericTagSelectors[strings.ToLower(strings.TrimSpace(selector))]
	return ok
}

// polishIncomingStep normalizes selectors and drops low-value recorder noise before coalescing.
func polishIncomingStep(step RecordedStep) (RecordedStep, bool) {
	switch step.Action {
	case "click":
		step = upgradeClickSelector(step)
		step.Selector = canonicalizeRecordedSelector(step.Selector)
		if step.HoverSelector != "" {
			step.HoverSelector = canonicalizeRecordedSelector(step.HoverSelector)
		}
	case "fill":
		step = upgradeFillSelector(step)
	case "check", "uncheck":
		step = upgradeCheckboxSelector(step)
	case "hover", "scroll-to":
		step.Selector = canonicalizeRecordedSelector(step.Selector)
		if isGenericTagSelector(step.Selector) {
			return step, false
		}
	case "goto":
		step.Value = strings.TrimSpace(step.Value)
		return step, step.Value != ""
	case "press":
		return step, strings.TrimSpace(step.Value) != ""
	}
	if step.Selector == "" && step.Action != "goto" && step.Action != "press" {
		return step, false
	}
	return step, true
}
