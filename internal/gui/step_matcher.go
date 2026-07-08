package gui

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type stepMatch struct {
	Action     stepdsl.Action
	ParseErr   error
	Catalog    stepcatalog.Entry
	CatalogOK  bool
	TestClient bool
	Recovered  bool
	Normalized string
}

func matchStepText(stepText string, line int) stepMatch {
	stepText = normalizeEditorStepText(stepText)
	if stepText == "" {
		return stepMatch{}
	}
	if gherkin.IsTestClientStep(gherkin.Step{Text: stepText}) {
		return stepMatch{TestClient: true}
	}

	action, err := stepdsl.Parse(gherkin.Step{Line: line, Text: stepText})
	recovered := false
	normalized := stepText
	if err != nil {
		normalized = recoverEditorStepText(stepText)
		if normalized != stepText {
			if recoveredAction, recoveredErr := stepdsl.Parse(gherkin.Step{Line: line, Text: normalized}); recoveredErr == nil {
				action = recoveredAction
				err = nil
				recovered = true
			}
		}
	}

	match := stepMatch{
		Action:     action,
		ParseErr:   err,
		Recovered:  recovered,
		Normalized: normalized,
	}
	if err == nil {
		if entry, ok := stepcatalog.LookupByAction(action.Kind); ok {
			match.Catalog = entry
			match.CatalogOK = true
		}
	}
	return match
}

func normalizeEditorStepText(stepText string) string {
	return strings.TrimSpace(stepText)
}

func recoverEditorStepText(stepText string) string {
	replacer := strings.NewReplacer(
		"“", `"`,
		"”", `"`,
		"„", `"`,
		"«", `"`,
		"»", `"`,
		"’", "'",
		"‘", "'",
		"Ё", "Е",
		"ё", "е",
	)
	return replacer.Replace(strings.TrimSpace(stepText))
}
