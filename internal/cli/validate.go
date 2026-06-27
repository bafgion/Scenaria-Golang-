package cli

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func RunValidate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: scenaria validate <path>")
	}
	target := args[0]

	files, err := collectFeatureFiles(target)
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

func collectFeatureFiles(target string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(target, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".feature") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan %q: %w", target, err)
	}
	return files, nil
}
