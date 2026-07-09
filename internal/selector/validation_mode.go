package selector

// ValidationMode controls how browser selector validation reaches each step.
type ValidationMode string

const (
	ValidationModeStatic ValidationMode = "static"
	ValidationModeFlow   ValidationMode = "flow"
)

const (
	staticValidationLimitation = "проверка только на текущей странице; шаги сценария не выполняются"
	flowValidationLimitation   = "выполняются только безопасные шаги (goto, click, hover, scroll) перед проверкой"
)

func NormalizeValidationMode(mode string) ValidationMode {
	switch ValidationMode(mode) {
	case ValidationModeFlow:
		return ValidationModeFlow
	default:
		return ValidationModeStatic
	}
}

func validationLimitation(mode ValidationMode) string {
	if mode == ValidationModeFlow {
		return flowValidationLimitation
	}
	return staticValidationLimitation
}
