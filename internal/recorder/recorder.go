package recorder

import (
	"fmt"
	"os"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/selector"
)

type Options struct {
	FeatureName  string
	ScenarioName string
	Steps        []string
}

func WriteFeature(path string, opts Options) error {
	if strings.TrimSpace(opts.FeatureName) == "" {
		return fmt.Errorf("feature name is required")
	}
	if strings.TrimSpace(opts.ScenarioName) == "" {
		opts.ScenarioName = "Recorded scenario"
	}

	feature := &gherkin.Feature{
		Title: opts.FeatureName,
		Scenarios: []gherkin.Scenario{
			{
				Title: opts.ScenarioName,
				Steps: make([]gherkin.Step, 0, len(opts.Steps)),
			},
		},
	}
	for _, line := range opts.Steps {
		step, ok := parseRecordedStep(line)
		if !ok {
			return fmt.Errorf("unsupported recorded step %q", line)
		}
		feature.Scenarios[0].Steps = append(feature.Scenarios[0].Steps, step)
	}

	content, err := gherkin.SerializeFeature(feature)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func parseRecordedStep(line string) (gherkin.Step, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return gherkin.Step{}, false
	}
	keywords := []string{"Допустим", "Когда", "Тогда", "И", "Но"}
	for _, keyword := range keywords {
		if strings.HasPrefix(trimmed, keyword+" ") {
			return gherkin.Step{
				Keyword: keyword,
				Text:    strings.TrimSpace(strings.TrimPrefix(trimmed, keyword)),
			}, true
		}
	}
	return gherkin.Step{Keyword: "Когда", Text: trimmed}, true
}

// EventsToStep converts a recorder DOM event into a Gherkin step line.
func EventsToStep(eventType string, detail map[string]string) (string, bool) {
	sel := strings.TrimSpace(detail["selector"])
	if sel == "" {
		sel = BuildSelectorFromDetail(detail)
	}
	switch eventType {
	case "click":
		if sel == "" {
			return "", false
		}
		return `нажимаю "` + sel + `"`, true
	case "input":
		value := detail["value"]
		if sel == "" || value == "" {
			return "", false
		}
		return `ввожу "` + value + `" в "` + sel + `"`, true
	default:
		return "", false
	}
}

func BuildSelectorFromDetail(detail map[string]string) string {
	return selector.BuildFromElement(selector.ElementInfo{
		Tag:         detail["tag"],
		ID:          detail["id"],
		Name:        detail["name"],
		Text:        detail["text"],
		TestID:      detail["testid"],
		Placeholder: detail["placeholder"],
		Role:        detail["role"],
		Label:       detail["captiontext"],
		AriaLabel:   detail["arialabel"],
		Type:        detail["inputtype"],
	})
}

// RecorderScript is injected into the browser during interactive recording.
var RecorderScript = selector.RecorderHeuristicsJS
