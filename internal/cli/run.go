package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/report/allure"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/version"
)

type runOptions struct {
	targets           []string
	dryRun            bool
	summaryJSON       string
	junitPath         string
	htmlPath          string
	allureDir         string
	traceDir          string
	videoDir          string
	engine            string
	browser           string
	headed            bool
	baseURL           string
	installPlaywright bool
	tag               string
	variables         map[string]string
	testClient        string
	slowMo            float64
	workers           int
}

func RunRun(args []string) error {
	opts, err := parseRunOptions(args)
	if err != nil {
		return err
	}

	store := scenario.NewFeatureStore()
	files := make([]string, 0)
	for _, target := range opts.targets {
		discovered, err := store.Discover(target)
		if err != nil {
			return err
		}
		files = append(files, discovered...)
	}
	if len(files) == 0 {
		return fmt.Errorf("no .feature files found in %v", opts.targets)
	}

	var errorsCount int
	featureInputs := make([]player.FeatureInput, 0, len(files))

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

		featureInputs = append(featureInputs, player.FeatureInput{
			Path:    path,
			Feature: feature,
		})
		fmt.Printf("✓ %s\n", path)
	}

	if errorsCount > 0 {
		return fmt.Errorf("run preflight failed with %d issue(s)", errorsCount)
	}

	plan := player.BuildExecutionPlanWithTestClient(featureInputs, opts.tag, opts.variables, opts.testClient)
	if len(plan.Cases) == 0 {
		if opts.tag != "" {
			return fmt.Errorf("no scenarios found with tag %q in %v", opts.tag, opts.targets)
		}
		return fmt.Errorf("no runnable scenarios found in %v", opts.targets)
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
	if opts.htmlPath != "" {
		if writeErr := report.WriteHTML(opts.htmlPath, result); writeErr != nil {
			return writeErr
		}
		fmt.Printf("Wrote HTML report: %s\n", opts.htmlPath)
	}
	if opts.allureDir != "" {
		if writeErr := allure.WriteResults(opts.allureDir, result); writeErr != nil {
			return writeErr
		}
		fmt.Printf("Wrote Allure results: %s\n", opts.allureDir)
	}

	if !opts.dryRun {
		if root := paths.InferProjectRoot(opts.targets); root != "" {
			if store, storeErr := runstatus.Open(root); storeErr == nil {
				for _, scenarioResult := range result.ScenarioResults {
					_ = store.Record(runstatus.Entry{
						Path:    scenarioResult.FeaturePath + "::" + scenarioResult.Scenario,
						Success: scenarioResult.Status == "passed",
						Message: scenarioResult.Message,
						Runner:  opts.engine,
					})
				}
			}
		}
	}

	fmt.Printf("Discovered %d file(s), %d scenario(s), %d step(s) [%s]\n", result.Files, result.Scenarios, result.Steps, version.String())
	if opts.dryRun {
		fmt.Println("Dry-run mode enabled: browser execution skipped")
		return nil
	}
	return nil
}

func parseRunOptions(args []string) (runOptions, error) {
	if len(args) == 0 {
		return runOptions{}, fmt.Errorf("usage: scenaria run <path> [more paths...] [--dry-run] [--summary-json <file>] [--junit <file>] [--allure <dir>] [--trace <dir>] [--video <dir>] [--engine stub|playwright] [--browser chromium|firefox|webkit] [--headed] [--base-url <url>] [--install-playwright] [--tag <tag>] [--var NAME=VALUE] [--slow-mo <ms>]")
	}
	opts := runOptions{
		engine:  "",
		browser: "chromium",
	}
	if cfg, err := settings.LoadDefaultAppSettings(); err == nil && cfg != nil {
		if cfg.Browser != "" {
			opts.browser = cfg.Browser
		}
		opts.headed = !cfg.Headless
		if cfg.ParallelWorkers > 0 {
			opts.workers = cfg.ParallelWorkers
		}
		if cfg.MaxLoopIterations > 0 {
			player.SetMaxLoopIterations(cfg.MaxLoopIterations)
		}
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			opts.targets = append(opts.targets, arg)
			continue
		}
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
		case "--html":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--html requires a file path")
			}
			i++
			opts.htmlPath = args[i]
		case "--allure":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--allure requires a directory path")
			}
			i++
			opts.allureDir = args[i]
		case "--trace":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--trace requires a directory path")
			}
			i++
			opts.traceDir = args[i]
		case "--video":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--video requires a directory path")
			}
			i++
			opts.videoDir = args[i]
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
		case "--tag":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--tag requires a tag name")
			}
			i++
			opts.tag = args[i]
		case "--test-client":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--test-client requires a client name")
			}
			i++
			opts.testClient = args[i]
		case "--var":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--var requires NAME=VALUE")
			}
			i++
			if opts.variables == nil {
				opts.variables = map[string]string{}
			}
			pair := args[i]
			eq := strings.Index(pair, "=")
			if eq <= 0 {
				return runOptions{}, fmt.Errorf("--var expects NAME=VALUE, got %q", pair)
			}
			opts.variables[pair[:eq]] = pair[eq+1:]
		case "--slow-mo":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--slow-mo requires milliseconds value")
			}
			i++
			var slowMo float64
			if _, err := fmt.Sscanf(args[i], "%f", &slowMo); err != nil || slowMo < 0 {
				return runOptions{}, fmt.Errorf("--slow-mo expects non-negative number, got %q", args[i])
			}
			opts.slowMo = slowMo
		case "--workers":
			if i+1 >= len(args) {
				return runOptions{}, fmt.Errorf("--workers requires a positive integer")
			}
			i++
			var workers int
			if _, err := fmt.Sscanf(args[i], "%d", &workers); err != nil || workers < 1 {
				return runOptions{}, fmt.Errorf("--workers expects integer >= 1, got %q", args[i])
			}
			opts.workers = workers
		default:
			return runOptions{}, fmt.Errorf("unknown flag for run: %s", arg)
		}
	}
	if len(opts.targets) == 0 {
		return runOptions{}, fmt.Errorf("at least one path is required")
	}
	if opts.baseURL == "" {
		if root := paths.InferProjectRoot(opts.targets); root != "" {
			if cfg, err := settings.LoadProjectConfig(root); err == nil && cfg.BaseURL != "" {
				opts.baseURL = cfg.BaseURL
			}
		}
	}
	return opts, nil
}

func resolveRunEngine(target, engine string) string {
	if engine != "" {
		return engine
	}
	root := paths.InferProjectRoot([]string{target})
	if root != "" {
		if cfg, err := settings.LoadProjectConfig(root); err == nil && cfg.DefaultRunner != "" {
			return cfg.DefaultRunner
		}
	}
	return "playwright"
}

func buildRunner(opts runOptions) (player.Runner, error) {
	if opts.dryRun {
		return player.DryRunner{}, nil
	}
	engine := ""
	if len(opts.targets) > 0 {
		engine = resolveRunEngine(opts.targets[0], opts.engine)
	} else {
		engine = resolveRunEngine("", opts.engine)
	}
	switch engine {
	case "stub":
		return player.BrowserRunner{
			Executor:        player.StubBrowserExecutor{},
			ParallelWorkers: opts.workers,
		}, nil
	case "playwright":
		return player.BrowserRunner{
			Executor: player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
				BrowserName: opts.browser,
				Headless:    !opts.headed,
				BaseURL:     opts.baseURL,
				AutoInstall: opts.installPlaywright,
				SlowMo:      opts.slowMo,
				TraceDir:    opts.traceDir,
				VideoDir:    opts.videoDir,
			}),
			ParallelWorkers: opts.workers,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported run engine %q (supported: stub, playwright)", opts.engine)
	}
}
