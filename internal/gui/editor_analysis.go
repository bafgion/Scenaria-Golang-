package gui

import "strings"

// editorAnalysisStep is a shared editor-time snapshot used by validation and steps panel.
type editorAnalysisStep struct {
	Line       int
	Raw        string
	Indent     string
	Keyword    string
	StepText   string
	Action     stepMatch
	ParseErr   error
	TestClient bool
	Recovered  bool
}

func analyzeEditorSteps(text string) []editorAnalysisStep {
	lines := strings.Split(text, "\n")
	out := make([]editorAnalysisStep, 0, len(lines))
	for i, raw := range lines {
		line := normalizeEditorLine(raw)
		if line == "" || isCommentLine(line) {
			continue
		}
		if isTagLine(line) || isTableLine(line) || isScenarioStructureLine(line) {
			continue
		}
		keyword, stepText := splitStepKeyword(line)
		if stepText == "" && keyword == "" {
			continue
		}
		parsed := matchStepText(stepText, i+1)
		step := editorAnalysisStep{
			Line:       i + 1,
			Raw:        raw,
			Indent:     leadingIndent(raw),
			Keyword:    keyword,
			StepText:   stepText,
			Action:     parsed,
			ParseErr:   parsed.ParseErr,
			TestClient: parsed.TestClient,
			Recovered:  parsed.Recovered,
		}
		out = append(out, step)
	}
	return out
}

func normalizeEditorLine(line string) string {
	return strings.TrimSpace(line)
}

func isCommentLine(line string) bool {
	return strings.HasPrefix(line, "#")
}
