package cli

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func RunValidate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: scenaria validate <path>")
	}
	target := args[0]

	store := scenario.NewFeatureStore()
	files, err := store.Discover(target)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .feature files found in %q", target)
	}

	var errorsCount int
	for _, path := range files {
		feature, parseErr := gherkin.ParseFeatureFile(path)
		if parseErr != nil {
			errorsCount++
			fmt.Printf("✗ %s: %v\n", path, parseErr)
			continue
		}
		issues := gherkin.ValidateFeature(feature)
		if len(issues) == 0 {
			fmt.Printf("✓ %s\n", path)
			continue
		}

		errorsCount += len(issues)
		for _, issue := range issues {
			fmt.Printf("✗ %s:%d %s\n", path, issue.Line, issue.Message)
		}
	}

	if errorsCount > 0 {
		return fmt.Errorf("validation failed with %d issue(s)", errorsCount)
	}
	fmt.Printf("Validated %d file(s): no issues found\n", len(files))
	return nil
}
