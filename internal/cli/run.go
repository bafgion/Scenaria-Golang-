package cli

import (
	"context"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type runOptions struct {
	target      string
	dryRun      bool
	summaryJSON string
}

func RunRun(args []string) error {
	opts, err := parseRunOptions(args)
	if err != nil {
		return err
	}

	store := scenario.NewFeatureStore()
	files, err := store.Discover(opts.target)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .feature files found in %q", opts.target)
	}

	var errorsCount int
	plan := player.ExecutionPlan{
		Features: make([]player.FeatureInput, 0, len(files)),
	}

	for _, path := range files {
		feature, loadErr := store.Load(path)
		if loadErr != nil {
			errorsCount++
			fmt.Printf("✗ %s: %v\n", path, loadErr)
			continue
		}

		issues := gherkin.ValidateFeature(feature)
		if len(issues) > 0 {
			errorsCount += len(issues)
			for _, issue := range issues {
				fmt.Printf("✗ %s:%d %s\n", path, issue.Line, issue.Message)
			}
			continue
		}

		plan.Features = append(plan.Features, player.FeatureInput{
			Path:    path,
			Feature: feature,
		})
		fmt.Printf("✓ %s\n", path)
	}

	if errorsCount > 0 {
		return fmt.Errorf("run preflight failed with %d issue(s)", errorsCount)
	}

	runner := player.NewRunner(opts.dryRun)
	result, err := runner.Execute(context.Background(), plan)
	if err != nil {
		return err
	}

	if opts.summaryJSON != "" {
		if writeErr := report.WriteRunSummary(opts.summaryJSON, report.FromExecutionResult(result)); writeErr != nil {
			return writeErr
		}
		fmt.Printf("Wrote summary report: %s\n", opts.summaryJSON)
	}

	fmt.Printf("Discovered %d file(s), %d scenario(s), %d step(s)\n", result.Files, result.Scenarios, result.Steps)
	if opts.dryRun {
		fmt.Println("Dry-run mode enabled: browser execution skipped")
		return nil
	}
	return nil
}

func parseRunOptions(args []string) (runOptions, error) {
	if len(args) == 0 {
		return runOptions{}, fmt.Errorf("usage: scenaria run <path> [--dry-run] [--summary-json <file>]")
	}
	opts := runOptions{target: args[0]}
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--dry-run":
			opts.dryRun = true
		case "--summary-json":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--summary-json requires a file path")
			}
			i++
			opts.summaryJSON = args[i]
		default:
			return runOptions{}, fmt.Errorf("unknown flag for run: %s", arg)
		}
	}
	return opts, nil
}
