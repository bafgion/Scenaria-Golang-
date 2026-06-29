package recorder

import (
	"fmt"
	"strings"
)

func recordedStepLine(step RecordedStep) (string, bool) {
	switch step.Action {
	case "upload":
		if step.Selector == "" || step.Value == "" {
			return "", false
		}
		return fmt.Sprintf(`загружаю файл "%s" в "%s"`, escapeStepText(step.Value), escapeStepText(step.Selector)), true
	case "check":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`отмечаю "%s"`, escapeStepText(step.Selector)), true
	case "uncheck":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`снимаю отметку с "%s"`, escapeStepText(step.Selector)), true
	case "press":
		key := strings.TrimSpace(step.Value)
		if key == "" {
			return "", false
		}
		return fmt.Sprintf(`нажимаю клавишу "%s"`, escapeStepText(key)), true
	case "press-in":
		key := strings.TrimSpace(step.Value)
		if key == "" || step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`нажимаю клавишу "%s" в "%s"`, escapeStepText(key), escapeStepText(step.Selector)), true
	case "scroll-to":
		if step.Selector == "" {
			return "", false
		}
		return fmt.Sprintf(`скроллю к "%s"`, escapeStepText(step.Selector)), true
	case "drag-drop":
		if step.Selector == "" || step.Value == "" {
			return "", false
		}
		return fmt.Sprintf(`перетаскиваю "%s" в "%s"`, escapeStepText(step.Selector), escapeStepText(step.Value)), true
	default:
		return "", false
	}
}
