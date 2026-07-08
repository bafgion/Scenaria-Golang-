package gui

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

// EditorAnalysisDTO bundles validation, steps strip, and optional scenario hints from one editor pass.
type EditorAnalysisDTO struct {
	Issues []ValidationIssue `json:"issues"`
	Steps  []EditorStepRow   `json:"steps"`
	Hints  []ScenarioHintDTO `json:"hints"`
}

// AnalyzeEditorContent runs a single analyzeEditorSteps pass and derives editor UI payloads from it.
func AnalyzeEditorContent(text string, includeHints bool) EditorAnalysisDTO {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	analysis := analyzeEditorSteps(normalized)

	dto := EditorAnalysisDTO{
		Steps: editorStepsFromAnalysis(analysis),
	}
	if feature, err := gherkin.ParseFeature(normalized); err == nil {
		dto.Issues = validateParsedFeature(feature)
	} else {
		dto.Issues = validationIssuesFromAnalysis(analysis)
	}
	if includeHints {
		dto.Hints = AnalyzeScenarioHintsFromSteps(scenarioStepsFromAnalysis(analysis))
	}
	return dto
}

func editorStepsFromAnalysis(analysis []editorAnalysisStep) []EditorStepRow {
	out := make([]EditorStepRow, 0, len(analysis))
	for _, step := range analysis {
		row := EditorStepRow{
			Line:    step.Line,
			Keyword: step.Keyword,
			Action:  step.StepText,
			Text:    step.StepText,
		}
		if step.TestClient {
			row.Kind = "test-client"
			row.Action = "TestClient"
			out = append(out, row)
			continue
		}
		if step.ParseErr != nil {
			row.Error = step.ParseErr.Error()
			out = append(out, row)
			continue
		}
		row.Kind = step.Action.Action.Kind
		row.Action = actionDisplayName(step.Action.Action.Kind)
		row.Element, row.Value = actionFields(step.Action.Action)
		out = append(out, row)
	}
	return out
}

func validationIssuesFromAnalysis(analysis []editorAnalysisStep) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	for _, step := range analysis {
		if step.TestClient || step.ParseErr == nil {
			continue
		}
		issues = append(issues, ValidationIssue{
			Line:    step.Line,
			Message: step.ParseErr.Error(),
		})
	}
	return issues
}

func scenarioStepsFromAnalysis(analysis []editorAnalysisStep) []parsedScenarioStep {
	out := make([]parsedScenarioStep, 0, len(analysis))
	for _, step := range analysis {
		line := scenarioStepLine{
			lineNo: step.Line,
			raw:    step.Raw,
			indent: step.Indent,
			body:   step.StepText,
		}
		if step.Keyword != "" {
			line.keyword = step.Keyword + " "
		}
		out = append(out, parsedScenarioStep{
			line:   line,
			action: step.Action.Action,
			err:    step.ParseErr,
		})
	}
	return out
}
