package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type recordOptions struct {
	output       string
	featureTitle string
	scenario     string
	steps        []string
	live         bool
	startURL     string
	headless     bool
	idleSeconds  int
}

func RunRecord(args []string) error {
	opts, err := parseRecordOptions(args)
	if err != nil {
		return err
	}

	if opts.live {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.idleSeconds+15)*time.Second)
		defer cancel()
		return recorder.RecordLive(ctx, recorder.LiveOptions{
			StartURL:     opts.startURL,
			FeatureName:  opts.featureTitle,
			ScenarioName: opts.scenario,
			OutputPath:   opts.output,
			Headless:     opts.headless,
			IdleTimeout:  time.Duration(opts.idleSeconds) * time.Second,
		})
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
		headless:     false,
		idleSeconds:  30,
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
		case "--live":
			opts.live = true
		case "--url":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--url requires a URL")
			}
			i++
			opts.startURL = args[i]
		case "--headless":
			opts.headless = true
		case "--idle":
			if i+1 >= len(args) {
				return recordOptions{}, fmt.Errorf("--idle requires seconds")
			}
			i++
			if _, err := fmt.Sscanf(args[i], "%d", &opts.idleSeconds); err != nil {
				return recordOptions{}, fmt.Errorf("--idle expects integer seconds")
			}
		default:
			return recordOptions{}, fmt.Errorf("unknown flag for record: %s", args[i])
		}
	}
	if opts.output == "" {
		return recordOptions{}, fmt.Errorf("usage: scenaria record --output <file.feature> [--feature <title>] [--scenario <title>] [--step <text> ...] [--live --url <url> [--idle 30] [--headless]]")
	}
	if opts.live && opts.startURL == "" {
		return recordOptions{}, fmt.Errorf("--live requires --url")
	}
	return opts, nil
}
