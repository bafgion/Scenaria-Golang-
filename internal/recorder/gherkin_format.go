package recorder

import (
	"regexp"
	"strings"
)

var recordedStepKeywordRE = regexp.MustCompile(`(?i)^(Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+`)

// StripRecordedStepKeyword removes a leading Gherkin keyword from a step body.
func StripRecordedStepKeyword(body string) string {
	return strings.TrimSpace(recordedStepKeywordRE.ReplaceAllString(strings.TrimSpace(body), ""))
}

func pickRecordedStepKeyword(sessionIndex int) string {
	if sessionIndex <= 0 {
		return "Допустим"
	}
	return "И"
}

// FormatRecordedGherkinLine renders a tab-indented Gherkin step line.
func FormatRecordedGherkinLine(body, keyword string) string {
	bare := StripRecordedStepKeyword(body)
	if bare == "" {
		return "\t" + keyword + " выполняю действие"
	}
	return "\t" + keyword + " " + bare
}

// FormatRecordedStepsAsGherkin converts recorded steps to formatted Gherkin lines.
func FormatRecordedStepsAsGherkin(steps []RecordedStep) []string {
	bare := RecordedStepsToLines(steps)
	out := make([]string, 0, len(bare))
	for i, line := range bare {
		out = append(out, FormatRecordedGherkinLine(line, pickRecordedStepKeyword(i)))
	}
	return out
}
