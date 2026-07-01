package recorder

import "strings"

// ApplyCoalescedStep merges step into steps using live-recording coalescing rules.
// The second return value is nil when the incoming step is skipped (duplicate hover/click/tab).
func ApplyCoalescedStep(steps []RecordedStep, step RecordedStep) ([]RecordedStep, *RecordedStep) {
	step, ok := polishIncomingStep(step)
	if !ok {
		return steps, nil
	}

	if step.Action == "fill" && (step.InputType == "checkbox" || isCheckboxValue(step.Value)) {
		checked := isCheckedValue(step.Value)
		return ApplyCoalescedStep(steps, RecordedStep{
			Action:   checkboxAction(checked),
			Selector: step.Selector,
			Text:     step.Text,
		})
	}

	if len(steps) == 0 {
		if step.Action == "fill" {
			step = upgradeFillSelector(step)
		}
		out := append(steps, step)
		return out, &step
	}

	last := steps[len(steps)-1]
	action := step.Action

	if action == "fill" && last.Action == "fill" && last.Selector == step.Selector {
		upgraded := upgradeFillSelector(step)
		out := append(steps[:len(steps)-1], upgraded)
		return out, &upgraded
	}

	if action == "fill" && last.Action == "click" &&
		(last.Selector == step.Selector || selectorIsFragile(last.Selector)) {
		upgraded := upgradeFillSelector(step)
		out := append(steps[:len(steps)-1], upgraded)
		return out, &upgraded
	}

	if action == "press" && last.Action == "fill" && last.Selector == step.Selector &&
		strings.EqualFold(strings.TrimSpace(step.Value), "tab") {
		return steps, nil
	}

	if action == "check" || action == "uncheck" {
		step = upgradeCheckboxSelector(step)
		if last.Action == "click" {
			out := append(steps[:len(steps)-1], step)
			return out, &step
		}
		if last.Action == action && last.Selector == step.Selector {
			return steps, nil
		}
		if (last.Action == "check" || last.Action == "uncheck") && last.Selector == step.Selector {
			out := append(steps[:len(steps)-1], step)
			return out, &step
		}
	}

	if action == "hover" {
		if last.Action == "hover" && last.Selector == step.Selector {
			return steps, nil
		}
		if hoverSeenSinceGoto(steps, step.Selector) {
			return steps, nil
		}
	}

	if action == "click" {
		if last.Action == "hover" && last.Selector == step.Selector {
			out := append(steps[:len(steps)-1], step)
			return out, &step
		}
		if last.Action == "click" && last.Selector == step.Selector {
			return steps, nil
		}
	}

	if action == "goto" && last.Action == "goto" && last.Value == step.Value {
		return steps, nil
	}

	out := append(steps, step)
	if step.Action == "fill" {
		step = upgradeFillSelector(step)
		out[len(out)-1] = step
	}
	lastStep := out[len(out)-1]
	return out, &lastStep
}

type StepNotifier func(index int, line string)

func appendRecordedStep(steps *[]RecordedStep, step RecordedStep, notify StepNotifier) {
	step, ok := polishIncomingStep(step)
	if !ok {
		return
	}
	if step.Action == "click" {
		hoverSel := strings.TrimSpace(step.HoverSelector)
		clickSel := strings.TrimSpace(step.Selector)
		if hoverSel != "" && hoverSel != clickSel {
			needsHover := len(*steps) == 0 ||
				((*steps)[len(*steps)-1].Action != "hover" || (*steps)[len(*steps)-1].Selector != hoverSel) &&
					!hoverSeenSinceGoto(*steps, hoverSel)
			if needsHover {
				hover := RecordedStep{Action: "hover", Selector: hoverSel, Text: step.HoverText}
				updated, emitted := ApplyCoalescedStep(*steps, hover)
				*steps = updated
				emitRecordedStep(notify, *steps, emitted)
			}
		}
	}
	updated, emitted := ApplyCoalescedStep(*steps, step)
	*steps = updated
	emitRecordedStep(notify, *steps, emitted)
}

func emitRecordedStep(notify StepNotifier, steps []RecordedStep, emitted *RecordedStep) {
	if notify == nil || emitted == nil {
		return
	}
	idx := len(steps) - 1
	for i := range steps {
		if &steps[i] == emitted {
			idx = i
			break
		}
	}
	if line, ok := RecordedStepToLine(*emitted); ok {
		notify(idx, line)
	}
}

func appendCoalescedStep(steps *[]RecordedStep, step RecordedStep, notify StepNotifier) {
	appendRecordedStep(steps, step, notify)
}

func appendGotoStep(steps *[]RecordedStep, url string, notify StepNotifier) {
	url = strings.TrimSpace(url)
	if url == "" {
		return
	}
	step := RecordedStep{Action: "goto", Value: url}
	updated, emitted := ApplyCoalescedStep(*steps, step)
	*steps = updated
	emitRecordedStep(notify, *steps, emitted)
}

func hoverSeenSinceGoto(steps []RecordedStep, selector string) bool {
	for i := len(steps) - 1; i >= 0; i-- {
		if steps[i].Action == "goto" {
			break
		}
		if steps[i].Action == "hover" && steps[i].Selector == selector {
			return true
		}
	}
	return false
}
