package cli

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func RunRun(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: scenaria run <path> [--dry-run]")
	}
	target := args[0]
	dryRun := hasFlag(args[1:], "--dry-run")

	store := scenario.NewFeatureStore()
	files, err := store.Discover(target)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .feature files found in %q", target)
	}

	var (
		errorsCount int
		scenarios   int
		steps       int
	)

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

		scenarios += len(feature.Scenarios)
		steps += len(feature.Background)
		for _, s := range feature.Scenarios {
			steps += len(s.Steps)
		}
		fmt.Printf("✓ %s\n", path)
	}

	if errorsCount > 0 {
		return fmt.Errorf("run preflight failed with %d issue(s)", errorsCount)
	}

	fmt.Printf("Discovered %d file(s), %d scenario(s), %d step(s)\n", len(files), scenarios, steps)
	if dryRun {
		fmt.Println("Dry-run mode enabled: browser execution skipped")
		return nil
	}
	return fmt.Errorf("browser execution engine is not implemented yet; retry with --dry-run")
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
