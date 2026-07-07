package player

import (
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

// StepRecord captures one executed leaf step for HTML trace viewer and reports.
type StepRecord struct {
	Index         int    `json:"index"`
	Line          int    `json:"line,omitempty"`
	Keyword       string `json:"keyword,omitempty"`
	Text          string `json:"text"`
	Selector      string `json:"selector,omitempty"`
	Status        string `json:"status"`
	DurationMS    int64  `json:"duration_ms,omitempty"`
	Error         string `json:"error,omitempty"`
	Network       string `json:"network,omitempty"`
	PageContext   string `json:"page_context,omitempty"`
	DOMSnapshot   string `json:"dom_snapshot,omitempty"`
	A11ySnapshot  string `json:"a11y_snapshot,omitempty"`
	ScreenshotPNG []byte `json:"-"`
}

func (c *RunContext) beginLeafStep(step gherkin.Step) int {
	if c == nil {
		return -1
	}
	idx := len(c.stepRecords)
	c.stepRecords = append(c.stepRecords, StepRecord{
		Index:   idx,
		Line:    step.Line,
		Keyword: step.Keyword,
		Text:    step.Text,
		Status:  "running",
	})
	return idx
}

func (c *RunContext) completeLeafStep(idx int, selector string, started time.Time, session *browserSession, err error) {
	if c == nil || idx < 0 || idx >= len(c.stepRecords) {
		return
	}
	rec := &c.stepRecords[idx]
	rec.Selector = selector
	rec.DurationMS = time.Since(started).Milliseconds()
	if err != nil {
		rec.Status = "failed"
		rec.Error = err.Error()
		if session != nil {
			rec.Network = session.lastNetworkFailure()
			rec.PageContext = capturePageContext(session)
			rec.DOMSnapshot = captureDOMSnapshot(session)
			rec.A11ySnapshot = captureA11ySnapshot(session)
			rec.ScreenshotPNG = captureViewportScreenshot(session)
		}
		return
	}
	rec.Status = "passed"
	if session != nil {
		if net := session.lastNetworkFailure(); net != "" {
			rec.Network = net
		}
		if c != nil && c.stepScreenshots {
			rec.ScreenshotPNG = captureViewportScreenshot(session)
		}
	}
}

func (c *RunContext) StepRecords() []StepRecord {
	if c == nil || len(c.stepRecords) == 0 {
		return nil
	}
	out := make([]StepRecord, len(c.stepRecords))
	copy(out, c.stepRecords)
	return out
}

func actionSelector(action stepdsl.Action) string {
	switch action.Kind {
	case "click", "double-click", "hover", "clear", "check", "uncheck",
		"download-click", "assert-visible", "assert-hidden", "wait-visible",
		"wait-hidden", "remember-field", "drag-drop", "scroll-to":
		return action.Value1
	case "fill", "fill-generated", "select", "press-in", "upload", "assert-text":
		return action.Value2
	case "goto", "remember-url", "assert-url", "wait-url":
		return action.Value1
	default:
		return ""
	}
}
