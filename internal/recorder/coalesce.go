package recorder

import "strings"

// ApplyCoalescedStep merges step into steps using live-recording coalescing rules.
// The second return value is nil when the incoming step is skipped (duplicate hover/click/tab).
func ApplyCoalescedStep(steps []RecordedStep, step RecordedStep) ([]RecordedStep, *RecordedStep) {
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

	if action == "fill" && last.Action == "click" && last.Selector == step.Selector {
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

	if action == "hover" && last.Action == "hover" && last.Selector == step.Selector {
		return steps, nil
	}

	if action == "click" && last.Action == "click" && last.Selector == step.Selector {
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

func appendRecordedStep(steps *[]RecordedStep, step RecordedStep) {
	if step.Action == "click" {
		hoverSel := strings.TrimSpace(step.HoverSelector)
		if hoverSel != "" {
			needsHover := len(*steps) == 0 ||
				(*steps)[len(*steps)-1].Action != "hover" ||
				(*steps)[len(*steps)-1].Selector != hoverSel
			if needsHover {
				hover := RecordedStep{Action: "hover", Selector: hoverSel, Text: step.HoverText}
				updated, _ := ApplyCoalescedStep(*steps, hover)
				*steps = updated
			}
		}
	}
	updated, _ := ApplyCoalescedStep(*steps, step)
	*steps = updated
}

func appendCoalescedStep(steps *[]RecordedStep, step RecordedStep) {
	appendRecordedStep(steps, step)
}
