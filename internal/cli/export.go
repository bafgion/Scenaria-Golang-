package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/exporter"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type exportOptions struct {
	input  string
	output string
	format string
}

func RunExport(args []string) error {
	opts, err := parseExportOptions(args)
	if err != nil {
		return err
	}

	store := scenario.NewFeatureStore()
	feature, err := store.Load(opts.input)
	if err != nil {
		return err
	}

	switch opts.format {
	case "json":
		document := exporter.NewFeatureExportDocument(opts.input, feature)
		if err := exporter.WriteFeatureJSON(opts.output, document); err != nil {
			return err
		}
	case "feature":
		if err := store.Save(opts.output, feature); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported export format %q", opts.format)
	}

	fmt.Printf("Exported %s -> %s (%s)\n", opts.input, opts.output, opts.format)
	return nil
}

func parseExportOptions(args []string) (exportOptions, error) {
	if len(args) == 0 {
		return exportOptions{}, fmt.Errorf("usage: scenaria export <feature-file> --output <file> [--format json|feature]")
	}

	opts := exportOptions{
		input: args[0],
	}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--output", "-o":
			if i+1 >= len(args) {
				return exportOptions{}, fmt.Errorf("--output requires a file path")
			}
			i++
			opts.output = args[i]
		case "--format":
			if i+1 >= len(args) {
				return exportOptions{}, fmt.Errorf("--format requires a value (json|feature)")
			}
			i++
			opts.format = strings.ToLower(args[i])
		default:
			return exportOptions{}, fmt.Errorf("unknown flag for export: %s", args[i])
		}
	}

	if opts.output == "" {
		return exportOptions{}, fmt.Errorf("--output is required")
	}
	if opts.format == "" {
		opts.format = inferExportFormat(opts.output)
	}
	if opts.format != "json" && opts.format != "feature" {
		return exportOptions{}, fmt.Errorf("unsupported export format %q (supported: json, feature)", opts.format)
	}
	return opts, nil
}

func inferExportFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return "json"
	case ".feature":
		return "feature"
	default:
		return "json"
	}
}
