package plugin

import (
	"fmt"
	"strings"
	"unicode"
)

type RunTarget struct {
	Runner string
	Args   []string
}

func ResolveRun(projectRoot string, desc Descriptor, dryRun bool) (RunTarget, error) {
	projectRoot = strings.TrimSpace(projectRoot)
	if projectRoot == "" {
		return RunTarget{}, fmt.Errorf("project path is required")
	}
	if specs := structuredCommands(desc); len(specs) > 0 {
		return resolveCommandSpec(projectRoot, specs[0], dryRun)
	}
	if IsVanessa(desc) {
		return vaRunTarget(projectRoot, dryRun), nil
	}
	if len(desc.Commands) == 0 {
		return RunTarget{}, fmt.Errorf("plugin %q defines no commands", desc.ID)
	}
	return resolveCommand(projectRoot, strings.TrimSpace(desc.Commands[0]), dryRun)
}

func resolveCommand(projectRoot, command string, dryRun bool) (RunTarget, error) {
	parts, err := splitLegacyCommand(command)
	if err != nil {
		return RunTarget{}, err
	}
	if len(parts) == 0 {
		return RunTarget{}, fmt.Errorf("empty plugin command")
	}
	switch parts[0] {
	case "va":
		if len(parts) < 2 || parts[1] != "run" {
			return RunTarget{}, fmt.Errorf("unsupported plugin command %q (expected va run)", command)
		}
		return vaRunTarget(projectRoot, dryRun), nil
	case "run":
		args := append([]string{projectRoot}, parts[1:]...)
		if dryRun {
			args = append(args, "--dry-run")
		}
		return RunTarget{Runner: "run", Args: args}, nil
	default:
		return RunTarget{}, fmt.Errorf("unsupported plugin command %q (supported: va run, run)", command)
	}
}

func resolveCommandSpec(projectRoot string, spec CommandSpec, dryRun bool) (RunTarget, error) {
	runner := strings.TrimSpace(spec.Runner)
	if runner == "" {
		return RunTarget{}, fmt.Errorf("empty plugin command runner")
	}
	args := append([]string(nil), spec.Args...)
	switch runner {
	case "va":
		if len(args) == 0 {
			return vaRunTarget(projectRoot, dryRun), nil
		}
		if args[0] != "run" {
			return RunTarget{}, fmt.Errorf("unsupported plugin command %q (expected va run)", runner)
		}
		if !containsFlag(args, "--project") {
			args = append([]string{"run", "--project", projectRoot}, args[1:]...)
		}
		if dryRun && !containsFlag(args, "--dry-run") {
			args = append(args, "--dry-run")
		}
		return RunTarget{Runner: "va", Args: args}, nil
	case "run":
		args = append([]string{projectRoot}, args...)
		if dryRun && !containsFlag(args, "--dry-run") {
			args = append(args, "--dry-run")
		}
		return RunTarget{Runner: "run", Args: args}, nil
	default:
		return RunTarget{}, fmt.Errorf("unsupported plugin command runner %q (supported: va, run)", runner)
	}
}

func vaRunTarget(projectRoot string, dryRun bool) RunTarget {
	args := []string{"run", "--project", projectRoot}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return RunTarget{Runner: "va", Args: args}
}

func containsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func splitLegacyCommand(command string) ([]string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, nil
	}
	parts := make([]string, 0)
	var current strings.Builder
	inQuote := rune(0)
	escaped := false
	hadToken := false
	for _, r := range command {
		switch {
		case escaped:
			if r != '"' && r != '\'' && r != '\\' {
				current.WriteRune('\\')
			}
			current.WriteRune(r)
			hadToken = true
			escaped = false
		case r == '\\' && inQuote != '\'':
			escaped = true
			hadToken = true
		case inQuote != 0:
			if r == inQuote {
				inQuote = 0
			} else {
				current.WriteRune(r)
			}
			hadToken = true
		case r == '"' || r == '\'':
			inQuote = r
			hadToken = true
		case unicode.IsSpace(r):
			if hadToken {
				parts = append(parts, current.String())
				current.Reset()
				hadToken = false
			}
		default:
			current.WriteRune(r)
			hadToken = true
		}
	}
	if escaped {
		current.WriteRune('\\')
	}
	if inQuote != 0 {
		return nil, fmt.Errorf("unterminated quote in plugin command %q", command)
	}
	if hadToken {
		parts = append(parts, current.String())
	}
	return parts, nil
}
