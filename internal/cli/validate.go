package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/selector"
)

type validateOptions struct {
	target      string
	json        string
	browser     bool
	browserName string
	headless    bool
	baseURL     string
}

func RunValidate(args []string) error {
	opts, err := parseValidateOptions(args)
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

	validator := selector.Validator{}
	type caseResult struct {
		Path    string `json:"path"`
		Success bool   `json:"success"`
		Issues  []string `json:"issues"`
	}
	results := make([]caseResult, 0, len(files))
	var errorsCount int

	for _, path := range files {
		feature, parseErr := gherkin.ParseFeatureFile(path)
		if parseErr != nil {
			errorsCount++
			fmt.Printf("✗ %s: %v\n", path, parseErr)
			results = append(results, caseResult{Path: path, Success: false, Issues: []string{parseErr.Error()}})
			continue
		}

		issues := make([]string, 0)
		for _, issue := range gherkin.ValidateFeature(feature) {
			issues = append(issues, fmt.Sprintf("line %d: %s", issue.Line, issue.Message))
		}
		selectorIssues, selErr := validator.ValidateFeature(path, feature)
		if selErr != nil {
			issues = append(issues, selErr.Error())
		}
		for _, issue := range selectorIssues {
			issues = append(issues, fmt.Sprintf("line %d: selector %q: %s", issue.Line, issue.Selector, issue.Message))
		}
		if opts.browser && len(issues) == 0 {
			browserIssues, browserErr := validator.ValidateFeatureInBrowser(context.Background(), path, feature, selector.BrowserValidateOptions{
				BrowserName: opts.browserName,
				Headless:    opts.headless,
				BaseURL:     opts.baseURL,
			})
			if browserErr != nil {
				issues = append(issues, browserErr.Error())
			}
			for _, issue := range browserIssues {
				if issue.Selector != "" {
					issues = append(issues, fmt.Sprintf("line %d: selector %q: %s", issue.Line, issue.Selector, issue.Message))
				} else {
					issues = append(issues, fmt.Sprintf("line %d: %s", issue.Line, issue.Message))
				}
			}
		}

		if len(issues) == 0 {
			fmt.Printf("✓ %s\n", path)
			results = append(results, caseResult{Path: path, Success: true})
			continue
		}

		errorsCount += len(issues)
		for _, issue := range issues {
			fmt.Printf("✗ %s %s\n", path, issue)
		}
		results = append(results, caseResult{Path: path, Success: false, Issues: issues})
	}

	if opts.json != "" {
		payload, err := json.MarshalIndent(map[string]any{
			"success": errorsCount == 0,
			"cases":   results,
		}, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(opts.json, append(payload, '\n'), 0o644); err != nil {
			return err
		}
		fmt.Printf("JSON: %s\n", opts.json)
	}

	if errorsCount > 0 {
		return fmt.Errorf("validation failed with %d issue(s)", errorsCount)
	}
	fmt.Printf("Validated %d file(s): no issues found\n", len(files))
	return nil
}

func parseValidateOptions(args []string) (validateOptions, error) {
	if len(args) == 0 {
		return validateOptions{}, fmt.Errorf("usage: scenaria validate <path> [--json <file>] [--browser [chromium|firefox|webkit]] [--headed] [--base-url <url>]")
	}
	opts := validateOptions{target: args[0], headless: true, browserName: "chromium"}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--json":
			if i+1 >= len(args) {
				return validateOptions{}, fmt.Errorf("--json requires a file path")
			}
			i++
			opts.json = args[i]
		case "--browser":
			opts.browser = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				switch strings.ToLower(args[i+1]) {
				case "chromium", "firefox", "webkit":
					i++
					opts.browserName = args[i]
				}
			}
		case "--headed":
			opts.headless = false
		case "--base-url":
			if i+1 >= len(args) {
				return validateOptions{}, fmt.Errorf("--base-url requires a URL value")
			}
			i++
			opts.baseURL = args[i]
		default:
			return validateOptions{}, fmt.Errorf("unknown flag for validate: %s", args[i])
		}
	}
	return opts, nil
}
