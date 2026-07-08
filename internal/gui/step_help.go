package gui

import "github.com/bafgion/scenaria-golang/internal/stepcatalog"

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
	line = normalizeEditorLine(line)
	if line == "" || isCommentLine(line) || isScenarioStructureLine(line) || isTagLine(line) || isTableLine(line) {
		return StepCatalogEntry{}, false
	}
	_, stepText := splitStepKeyword(line)
	if stepText == "" {
		return StepCatalogEntry{}, false
	}
	match := matchStepText(stepText, 1)
	if match.ParseErr == nil && match.CatalogOK {
		return stepCatalogEntryFrom(match.Catalog), true
	}
	return StepCatalogEntry{}, false
}
