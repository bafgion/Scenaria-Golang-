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
	target            string
	dryRun            bool
	summaryJSON       string
	junitPath         string
	engine            string
	browser           string
	headed            bool
	baseURL           string
	installPlaywright bool
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

	runner, err := buildRunner(opts)
	if err != nil {
		return err
	}
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
	if opts.junitPath != "" {
		if writeErr := report.WriteJUnit(opts.junitPath, result); writeErr != nil {
			return writeErr
		}
		fmt.Printf("Wrote JUnit report: %s\n", opts.junitPath)
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
		return runOptions{}, fmt.Errorf("usage: scenaria run <path> [--dry-run] [--summary-json <file>] [--junit <file>] [--engine stub|playwright] [--browser chromium|firefox|webkit] [--headed] [--base-url <url>] [--install-playwright]")
	}
	opts := runOptions{
		target:  args[0],
		engine:  "stub",
		browser: "chromium",
	}
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
		case "--junit":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--junit requires a file path")
			}
			i++
			opts.junitPath = args[i]
		case "--engine":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--engine requires a value (stub|playwright)")
			}
			i++
			opts.engine = args[i]
		case "--browser":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--browser requires a value (chromium|firefox|webkit)")
			}
			i++
			opts.browser = args[i]
		case "--headed":
			opts.headed = true
		case "--base-url":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--base-url requires a URL value")
			}
			i++
			opts.baseURL = args[i]
		case "--install-playwright":
			opts.installPlaywright = true
		default:
			return runOptions{}, fmt.Errorf("unknown flag for run: %s", arg)
		}
	}
	return opts, nil
}

func buildRunner(opts runOptions) (player.Runner, error) {
	if opts.dryRun {
		return player.DryRunner{}, nil
	}
	switch opts.engine {
	case "stub":
		return player.BrowserRunner{
			Executor: player.StubBrowserExecutor{},
		}, nil
	case "playwright":
		return player.BrowserRunner{
			Executor: player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
				BrowserName: opts.browser,
				Headless:    !opts.headed,
				BaseURL:     opts.baseURL,
				AutoInstall: opts.installPlaywright,
			}),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported run engine %q (supported: stub, playwright)", opts.engine)
	}
}
