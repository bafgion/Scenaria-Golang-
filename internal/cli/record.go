package cli

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type recordOptions struct {
	output       string
	featureTitle string
	scenario     string
	steps        []string
}

func RunRecord(args []string) error {
	opts, err := parseRecordOptions(args)
	if err != nil {
		return err
	}

	store := scenario.NewFeatureStore()
	feature := &gherkin.Feature{
		Title: opts.featureTitle,
		Scenarios: []gherkin.Scenario{
			{
				Title: opts.scenario,
				Steps: make([]gherkin.Step, 0, len(opts.steps)),
			},
		},
	}

	for _, text := range opts.steps {
		feature.Scenarios[0].Steps = append(feature.Scenarios[0].Steps, gherkin.Step{
			Keyword: "Когда",
			Text:    text,
		})
	}
	if len(feature.Scenarios[0].Steps) == 0 {
		feature.Scenarios[0].Steps = append(feature.Scenarios[0].Steps, gherkin.Step{
			Keyword: "Когда",
			Text:    "выполняю действие",
		})
	}

	if err := store.Save(opts.output, feature); err != nil {
		return err
	}

	fmt.Printf("Recorded baseline scenario: %s\n", opts.output)
	return nil
}

func parseRecordOptions(args []string) (recordOptions, error) {
	opts := recordOptions{
		featureTitle: "Записанный сценарий",
		scenario:     "Базовый сценарий",
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--output", "-o":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--output requires a file path")
			}
			i++
			opts.output = args[i]
		case "--feature":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--feature requires a title")
			}
			i++
			opts.featureTitle = args[i]
		case "--scenario":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--scenario requires a title")
			}
			i++
			opts.scenario = args[i]
		case "--step":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--step requires text")
			}
			i++
			opts.steps = append(opts.steps, args[i])
		default:
			return recordOptions{}, fmt.Errorf("unknown flag for record: %s", args[i])
		}
	}
	if opts.output == "" {
		return recordOptions{}, fmt.Errorf("usage: scenaria record --output <file.feature> [--feature <title>] [--scenario <title>] [--step <text> ...]")
	}
	return opts, nil
}
