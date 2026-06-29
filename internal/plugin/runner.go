package plugin

import (
	"fmt"
	"strings"
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
	if IsVanessa(desc) {
		return vaRunTarget(projectRoot, dryRun), nil
	}
	if len(desc.Commands) == 0 {
		return RunTarget{}, fmt.Errorf("plugin %q defines no commands", desc.ID)
	}
	return resolveCommand(projectRoot, strings.TrimSpace(desc.Commands[0]), dryRun)
}

func resolveCommand(projectRoot, command string, dryRun bool) (RunTarget, error) {
	parts := strings.Fields(command)
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

func vaRunTarget(projectRoot string, dryRun bool) RunTarget {
	args := []string{"run", "--project", projectRoot}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return RunTarget{Runner: "va", Args: args}
}
