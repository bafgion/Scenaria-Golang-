package update

// Progress describes a single update step for the UI.
type Progress struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
	Percent int    `json:"percent"`
}

// Reporter receives update progress events.
type Reporter func(Progress)

func (r Reporter) report(stage, message string, percent int) {
	if r == nil {
		return
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	r(Progress{Stage: stage, Message: message, Percent: percent})
}
