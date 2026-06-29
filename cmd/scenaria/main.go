package main

import (
	"fmt"
	"os"

	"github.com/bafgion/scenaria-golang/internal/cli"
	"github.com/bafgion/scenaria-golang/internal/version"
)

type command struct {
	name        string
	description string
	run         func(args []string) error
}

func main() {
	commands := map[string]command{
		"help": {
			name:        "help",
			description: "Show available commands",
			run:         runHelp,
		},
		"run": {
			name:        "run",
			description: "Execute scenarios",
			run:         cli.RunRun,
		},
		"validate": {
			name:        "validate",
			description: "Validate scenarios",
			run:         cli.RunValidate,
		},
		"export": {
			name:        "export",
			description: "Export scenarios",
			run:         cli.RunExport,
		},
		"import-json": {
			name:        "import-json",
			description: "Import scenario from JSON export",
			run:         cli.RunImportJSON,
		},
		"record": {
			name:        "record",
			description: "Create baseline recorded scenario",
			run:         cli.RunRecord,
		},
		"init": {
			name:        "init",
			description: "Initialize .scenaria project scaffold",
			run:         cli.RunInit,
		},
		"update": {
			name:        "update",
			description: "Check for application updates",
			run:         cli.RunUpdate,
		},
		"plugins": {
			name:        "plugins",
			description: "Manage plugins",
			run:         cli.RunPlugins,
		},
		"va": {
			name:        "va",
			description: "Vanessa Automation (1C) runner",
			run:         runVA,
		},
		"version": {
			name:        "version",
			description: "Print CLI version",
			run:         runVersion,
		},
	}

	args := os.Args[1:]
	if len(args) == 0 {
		_ = runHelp(nil)
		return
	}
	if args[0] == "--help" || args[0] == "-h" {
		_ = runHelp(nil)
		return
	}
	if args[0] == "--version" {
		_ = runVersion(nil)
		return
	}

	cmdName := args[0]
	cmd, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmdName)
		_ = runHelp(nil)
		os.Exit(1)
	}

	if err := cmd.run(args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runHelp(_ []string) error {
	fmt.Println("scenaria - Go migration CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  scenaria <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run       Execute scenarios (`--dry-run`, `--summary-json`, `--junit`, `--allure`, `--trace`, `--video`, `--engine`)")
	fmt.Println("  validate  Validate scenario files and project consistency")
	fmt.Println("  export    Export scenarios (`--output`, `--format json|feature|ts|python`)")
	fmt.Println("  import-json  Import JSON export to .feature (`--output`, `--force`)")
	fmt.Println("  record    Record baseline scenario (`--output`, `--step`, `--live`)")
	fmt.Println("  init      Initialize project scaffold (`.scenaria/`)")
	fmt.Println("  update    Check for updates (`--check`)")
	fmt.Println("  plugins   Manage plugins (`list`, `install`, `uninstall`)")
	fmt.Println("  va        Vanessa Automation runner (`va run`)")
	fmt.Println("  version   Print version")
	fmt.Println("  help      Show this help")
	fmt.Println()
	fmt.Println("Install globally:")
	fmt.Println("  go install ./cmd/scenaria")
	return nil
}

func runVersion(_ []string) error {
	fmt.Println(version.String())
	return nil
}

func runVA(args []string) error {
	return cli.RunVA(args)
}
