package recorder

import (
	"fmt"
	"strings"
)

func EventToRecordedStep(eventType string, detail map[string]string) (RecordedStep, bool) {
	eventType = strings.ToLower(strings.TrimSpace(eventType))
	switch eventType {
	case "click":
		sel := strings.TrimSpace(detail["selector"])
		if sel == "" {
			sel = BuildSelectorFromDetail(detail)
		}
		if sel == "" {
			return RecordedStep{}, false
		}
		return RecordedStep{
			Action:   "click",
			Selector: sel,
			Text:     detail["text"],
			Context:  detail["contexttext"],
		}, true
	case "draw-signature":
		sel := strings.TrimSpace(detail["selector"])
		if sel == "" {
			sel = BuildSelectorFromDetail(detail)
		}
		if sel == "" {
			return RecordedStep{}, false
		}
		return RecordedStep{Action: "draw-signature", Selector: sel}, true
	case "input", "fill", "change":
		sel := strings.TrimSpace(detail["selector"])
		if sel == "" {
			sel = BuildSelectorFromDetail(detail)
		}
		value := detail["value"]
		if sel == "" || value == "" {
			return RecordedStep{}, false
		}
		return RecordedStep{
			Action:    "fill",
			Selector:  sel,
			Value:     value,
			Text:      detail["captiontext"],
			InputType: detail["inputtype"],
		}, true
	case "hover":
		sel := strings.TrimSpace(detail["selector"])
		if sel == "" {
			sel = BuildSelectorFromDetail(detail)
		}
		if sel == "" {
			return RecordedStep{}, false
		}
		return RecordedStep{Action: "hover", Selector: sel}, true
	case "goto":
		url := strings.TrimSpace(detail["url"])
		if url == "" {
			return RecordedStep{}, false
		}
		return RecordedStep{Action: "goto", Value: url}, true
	default:
		return RecordedStep{}, false
	}
}

func RecordedStepsToLines(steps []RecordedStep) []string {
	out := make([]string, 0, len(steps))
	for _, step := range steps {
		if line, ok := RecordedStepToLine(step); ok {
			out = append(out, line)
		}
	}
	return out
}

func RecordedStepToLine(step RecordedStep) (string, bool) {
	switch step.Action {
	case "goto":
		if step.Value == "" {
			return "", false
		}
		return fmt.Sprintf(`открыт "%s"`, escapeStepText(step.Value)), true
	case "click":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`нажимаю "%s"`, escapeStepText(step.Selector)), true
	case "fill":
		if step.Selector == "" || step.Value == "" {
			return "", false
		}
		return fmt.Sprintf(`ввожу "%s" в "%s"`, escapeStepText(step.Value), escapeStepText(step.Selector)), true
	case "draw-signature":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`рисую подпись в "%s"`, escapeStepText(step.Selector)), true
	case "hover":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`навожу "%s"`, escapeStepText(step.Selector)), true
	default:
		return "", false
	}
}

func escapeStepText(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
