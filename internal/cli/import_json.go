package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/exporter"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type importJSONOptions struct {
	input  string
	output string
	force  bool
}

func RunImportJSON(args []string) error {
	opts, err := parseImportJSONOptions(args)
	if err != nil {
		return err
	}
	doc, err := exporter.ReadFeatureJSON(opts.input)
	if err != nil {
		return err
	}
	if !opts.force {
		if _, err := os.Stat(opts.output); err == nil {
			return fmt.Errorf("output file %q already exists (use --force)", opts.output)
		}
	}
	store := scenario.NewFeatureStore()
	if err := store.Save(opts.output, doc.Feature); err != nil {
		return err
	}
	fmt.Printf("Imported %s -> %s\n", opts.input, opts.output)
	return nil
}

func parseImportJSONOptions(args []string) (importJSONOptions, error) {
	if len(args) == 0 {
		return importJSONOptions{}, fmt.Errorf("usage: scenaria import-json <file.json> --output <file.feature> [--force]")
	}
	opts := importJSONOptions{input: args[0]}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--output", "-o":
			i++
			if i >= len(args) {
				return importJSONOptions{}, fmt.Errorf("--output requires a file path")
			}
			opts.output = args[i]
		case "--force":
			opts.force = true
		default:
			return importJSONOptions{}, fmt.Errorf("unknown flag for import-json: %s", args[i])
		}
	}
	if opts.output == "" {
		base := strings.TrimSuffix(filepath.Base(opts.input), filepath.Ext(opts.input))
		opts.output = base + ".feature"
	}
	if !strings.HasSuffix(strings.ToLower(opts.output), ".feature") {
		opts.output += ".feature"
	}
	return opts, nil
}
