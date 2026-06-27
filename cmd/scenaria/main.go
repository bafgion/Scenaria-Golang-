package main

import (
	"fmt"
	"os"

	"github.com/bafgion/scenaria-golang/internal/cli"
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
		"record": {
			name:        "record",
			description: "Create baseline recorded scenario",
			run:         cli.RunRecord,
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
	fmt.Println("  run       Execute scenarios (`--dry-run`, `--summary-json`, `--junit`, `--engine`)")
	fmt.Println("  validate  Validate scenario files and project consistency")
	fmt.Println("  export    Export scenarios (`--output`, `--format json|feature`)")
	fmt.Println("  record    Record baseline scenario (`--output`, `--step`, ...)")
	fmt.Println("  version   Print version")
	fmt.Println("  help      Show this help")
	fmt.Println()
	fmt.Println("Install globally:")
	fmt.Println("  go install ./cmd/scenaria")
	return nil
}

func runVersion(_ []string) error {
	fmt.Println("scenaria dev")
	return nil
}

func runNotImplemented(name string) func(args []string) error {
	return func(_ []string) error {
		return fmt.Errorf("command %q is not implemented yet (see docs/MIGRATION_PLAN.md)", name)
	}
}
