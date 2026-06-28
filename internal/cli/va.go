package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

type vaOptions struct {
	project  string
	paths    []string
	tag      string
	dryRun   bool
	scenario string
}

func RunVA(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		return printVAHelp()
	}
	if args[0] != "run" {
		return fmt.Errorf("unknown va subcommand %q (supported: run)", args[0])
	}
	opts, err := parseVAOptions(args[1:])
	if err != nil {
		return err
	}
	projectRoot := opts.project
	if projectRoot == "" {
		projectRoot = paths.InferProjectRoot(opts.paths)
	}
	result, err := vanessa.Run(vanessa.RunRequest{
		ProjectRoot:   projectRoot,
		Paths:         opts.paths,
		Tag:           opts.tag,
		ScenarioNames: splitCSV(opts.scenario),
		DryRun:        opts.dryRun,
	})
	if err != nil && result.Error == "" {
		return err
	}
	for _, c := range result.Cases {
		mark := "✓"
		if !c.Success {
			mark = "✗"
		}
		fmt.Printf("%s %s: %s\n", mark, c.Name, c.Message)
	}
	if result.RunDir != "" {
		fmt.Printf("Run directory: %s\n", result.RunDir)
	}
	if !result.Success {
		if result.Error != "" {
			return fmt.Errorf("vanessa run failed: %s", result.Error)
		}
		return fmt.Errorf("vanessa run failed with exit code %d", result.ExitCode)
	}
	return nil
}

func parseVAOptions(args []string) (vaOptions, error) {
	opts := vaOptions{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--project requires a path")
			}
			opts.project = args[i]
		case "--dir":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--dir requires a path")
			}
			opts.paths = append(opts.paths, args[i])
		case "--files":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--files requires a path")
			}
			opts.paths = append(opts.paths, splitCSV(args[i])...)
		case "--tag":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--tag requires a value")
			}
			opts.tag = args[i]
		case "--scenario":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--scenario requires a value")
			}
			opts.scenario = args[i]
		case "--dry-run":
			opts.dryRun = true
		default:
			if !strings.HasPrefix(args[i], "-") && len(opts.paths) == 0 {
				opts.paths = append(opts.paths, args[i])
				continue
			}
			return vaOptions{}, fmt.Errorf("unknown flag for va run: %s", args[i])
		}
	}
	if len(opts.paths) == 0 && opts.project != "" {
		opts.paths = []string{opts.project}
	}
	if len(opts.paths) == 0 && opts.project == "" {
		if wd, err := os.Getwd(); err == nil {
			opts.project = wd
			opts.paths = []string{wd}
		}
	}
	if len(opts.paths) == 0 {
		return vaOptions{}, fmt.Errorf("usage: scenaria va run [--project <dir>] [--dir <dir>] [--files <paths>] [--tag <tag>] [--dry-run]")
	}
	return opts, nil
}

func printVAHelp() error {
	fmt.Println("Vanessa Automation runner (1C)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  scenaria va run [--project <dir>] [--dir <features>] [--files a.feature,b.feature]")
	fmt.Println("                  [--tag smoke] [--scenario \"Name\"] [--dry-run]")
	fmt.Println()
	fmt.Println("Configure platform in .scenaria/vanessa.json:")
	fmt.Println(`  {"platform_executable":"C:\\Program Files\\1cv8\\bin\\1cv8.exe","epf_path":"C:\\vanessa\\vanessa-automation.epf"}`)
	return nil
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, piece := range parts {
		if piece = strings.TrimSpace(piece); piece != "" {
			out = append(out, piece)
		}
	}
	return out
}
