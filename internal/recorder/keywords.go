package recorder

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

var recordedStepKeywords = []string{
	"Допустим", "Дано", "Когда", "Тогда", "И", "Но",
	"Given", "When", "Then", "And", "But",
}

// AssignRecordedStepKeywords maps bare recorder lines to Gherkin steps using Russian
// recording style: first step is Допустим, continuations are И.
func AssignRecordedStepKeywords(lines []string, existingStepCount int) []gherkin.Step {
	steps := make([]gherkin.Step, 0, len(lines))
	for i, line := range lines {
		body, keyword := splitRecordedStepKeyword(line)
		if body == "" {
			continue
		}
		if keyword == "" {
			if existingStepCount+i == 0 {
				keyword = "Допустим"
			} else {
				keyword = "И"
			}
		}
		steps = append(steps, gherkin.Step{Keyword: keyword, Text: body})
	}
	return steps
}

func splitRecordedStepKeyword(line string) (body, keyword string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", ""
	}
	for _, kw := range recordedStepKeywords {
		prefix := kw + " "
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix)), kw
		}
	}
	return trimmed, ""
}
