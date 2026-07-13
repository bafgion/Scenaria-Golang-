package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

type vaOptions struct {
	project            string
	paths              []string
	tag                string
	excludeTags        []string
	dryRun             bool
	scenario           string
	platformExecutable string
	epfPath            string
	ibConnection       string
	allure             bool
	rerunFailedRunDir  string
	installEPF         bool
	epfDownloadURL     string
	epfDestination     string
}

func RunVA(args []string) error {
	return RunVAContextWithOutput(context.Background(), args, nil)
}

func RunVAWithOutput(args []string, out io.Writer) error {
	return RunVAContextWithOutput(context.Background(), args, out)
}

func RunVAContext(ctx context.Context, args []string) error {
	return RunVAContextWithOutput(ctx, args, nil)
}

func RunVAContextWithOutput(ctx context.Context, args []string, out io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		return printVAHelp(out)
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
	result, err := vanessa.RunContext(ctx, vanessa.RunRequest{
		ProjectRoot:        projectRoot,
		Paths:              opts.paths,
		Tag:                opts.tag,
		ExcludeTags:        opts.excludeTags,
		ScenarioNames:      splitCSV(opts.scenario),
		DryRun:             opts.dryRun,
		PlatformExecutable: opts.platformExecutable,
		EPFPath:            opts.epfPath,
		IBConnection:       opts.ibConnection,
		ReportAllure:       opts.allure,
		RerunFailedRunDir:  opts.rerunFailedRunDir,
		InstallEPF:         opts.installEPF,
		EPFDownloadURL:     opts.epfDownloadURL,
		EPFDestination:     opts.epfDestination,
	})
	if err != nil && result.Error == "" {
		return err
	}
	for _, c := range result.Cases {
		mark := "✓"
		if !c.Success {
			mark = "✗"
		}
		cliPrintf(out, "%s %s: %s\n", mark, c.Name, c.Message)
	}
	if result.RunDir != "" {
		cliPrintf(out, "Run directory: %s\n", result.RunDir)
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
		case "--exclude-tag":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--exclude-tag requires a value")
			}
			opts.excludeTags = append(opts.excludeTags, splitCSV(args[i])...)
		case "--scenario":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--scenario requires a value")
			}
			opts.scenario = args[i]
		case "--platform-exe":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--platform-exe requires a path")
			}
			opts.platformExecutable = args[i]
		case "--epf":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--epf requires a path")
			}
			opts.epfPath = args[i]
		case "--ib":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--ib requires a connection string")
			}
			opts.ibConnection = args[i]
		case "--allure":
			opts.allure = true
		case "--rerun-failed":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--rerun-failed requires a previous run directory")
			}
			opts.rerunFailedRunDir = args[i]
		case "--epf-install":
			opts.installEPF = true
		case "--epf-url":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--epf-url requires a URL")
			}
			opts.epfDownloadURL = args[i]
		case "--epf-dest":
			i++
			if i >= len(args) {
				return vaOptions{}, fmt.Errorf("--epf-dest requires a path")
			}
			opts.epfDestination = args[i]
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
		return vaOptions{}, fmt.Errorf("usage: scenaria va run [--project <dir>] [--dir <dir>] [--files <paths>] [--tag <tag>] [--exclude-tag <tag>] [--dry-run]")
	}
	return opts, nil
}

func printVAHelp(out io.Writer) error {
	cliPrintln(out, "Vanessa Automation runner (1C)")
	cliPrintln(out)
	cliPrintln(out, "Usage:")
	cliPrintln(out, "  scenaria va run [--project <dir>] [--dir <features>] [--files a.feature,b.feature]")
	cliPrintln(out, "                  [--tag smoke] [--exclude-tag wip] [--scenario \"Name\"]")
	cliPrintln(out, "                  [--platform-exe <path>] [--epf <path>] [--ib <conn>] [--allure]")
	cliPrintln(out, "                  [--rerun-failed <run-dir>] [--epf-install] [--epf-url <url>] [--dry-run]")
	cliPrintln(out)
	cliPrintln(out, "Configure platform in .scenaria/vanessa.json:")
	cliPrintln(out, `  {"platform_executable":"C:\\Program Files\\1cv8\\bin\\1cv8.exe","epf_path":"C:\\vanessa\\vanessa-automation.epf"}`)
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
