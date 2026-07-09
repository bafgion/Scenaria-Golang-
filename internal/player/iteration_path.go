package player

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

// IterationFrame identifies one loop iteration in runtime execution.
type IterationFrame struct {
	Kind  string `json:"kind"`
	Index int    `json:"index"`
}

// DefaultMaxActionAttempts caps leaf-step executions per scenario.
const DefaultMaxActionAttempts = 10_000

func (c *RunContext) copyIterationPath() []IterationFrame {
	if c == nil || len(c.iterationPath) == 0 {
		return nil
	}
	out := make([]IterationFrame, len(c.iterationPath))
	copy(out, c.iterationPath)
	return out
}

func (c *RunContext) withIteration(kind string, index int, fn func() error) error {
	if c == nil {
		return fn()
	}
	c.iterationPath = append(c.iterationPath, IterationFrame{Kind: kind, Index: index})
	defer func() {
		c.iterationPath = c.iterationPath[:len(c.iterationPath)-1]
	}()
	return fn()
}

func (c *RunContext) consumeActionAttempt(limit int) error {
	if c == nil {
		return nil
	}
	c.actionAttempts++
	if limit > 0 && c.actionAttempts > limit {
		return fmt.Errorf("scenario exceeded max action attempts (%d)", limit)
	}
	return nil
}

func (c *RunContext) appendSkippedLeafSteps(steps []gherkin.Step) {
	if c == nil {
		return
	}
	for _, step := range gherkin.LeafSteps(steps) {
		if gherkin.IsTestClientStep(step) {
			continue
		}
		idx := len(c.stepRecords)
		c.stepRecords = append(c.stepRecords, StepRecord{
			Index:         idx,
			Line:          step.Line,
			Keyword:       step.Keyword,
			Text:          step.Text,
			Status:        "skipped",
			IterationPath: c.copyIterationPath(),
		})
	}
}

func (c *RunContext) setStepTerminalAction(idx int, action string) {
	if c == nil || idx < 0 || idx >= len(c.stepRecords) {
		return
	}
	c.stepRecords[idx].TerminalAction = action
}

// FormatIterationPath renders repeat[3].for_each[2] for reports.
func FormatIterationPath(path []IterationFrame) string {
	if len(path) == 0 {
		return ""
	}
	parts := make([]string, len(path))
	for i, frame := range path {
		parts[i] = fmt.Sprintf("%s[%d]", frame.Kind, frame.Index)
	}
	return strings.Join(parts, ".")
}

func screenshotPrefix(iterationPath []IterationFrame, base string) string {
	label := FormatIterationPath(iterationPath)
	if label == "" {
		return base
	}
	return base + "-" + strings.ReplaceAll(label, ".", "_")
}
