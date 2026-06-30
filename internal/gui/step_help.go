package gui

import (
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func stepCatalogEntryFrom(entry stepcatalog.Entry) StepCatalogEntry {
	return StepCatalogEntry{
		Label:       entry.Label,
		Action:      entry.Action,
		Category:    entry.Category,
		Description: entry.Description,
		Template:    entry.Template,
		Example:     entry.Example,
		Parameters:  entry.Parameters,
		Help:        entry.Help,
	}
}

// DescribeEditorLine returns catalog help for a single editor line (with optional Gherkin keyword).
func DescribeEditorLine(line string) (StepCatalogEntry, bool) {
	line = trimLine(line)
	if line == "" || isCommentLine(line) || isScenarioStructureLine(line) || isTagLine(line) || isTableLine(line) {
		return StepCatalogEntry{}, false
	}
	_, stepText := splitStepKeyword(line)
	if stepText == "" {
		return StepCatalogEntry{}, false
	}
	if action, err := stepdsl.Parse(gherkin.Step{Text: stepText}); err == nil {
		if entry, ok := stepcatalog.LookupByAction(action.Kind); ok {
			return stepCatalogEntryFrom(entry), true
		}
	}
	if entry, ok := stepcatalog.LookupByStepText(stepText); ok {
		return stepCatalogEntryFrom(entry), true
	}
	return StepCatalogEntry{}, false
}

func trimLine(line string) string {
	for len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
		line = line[1:]
	}
	for len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
		line = line[:len(line)-1]
	}
	return line
}

func isCommentLine(line string) bool {
	return len(line) > 0 && line[0] == '#'
}
