package cli

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/plugin"
)

type pluginsOptions struct {
	project string
	name    string
	source  string
}

func RunPlugins(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: scenaria plugins <list|install|uninstall> [flags]")
	}
	subcommand := args[0]
	opts, err := parsePluginsOptions(args[1:])
	if err != nil {
		return err
	}
	if opts.project == "" {
		opts.project = "."
	}

	switch subcommand {
	case "list":
		plugins, err := plugin.List(opts.project)
		if err != nil {
			return err
		}
		if len(plugins) == 0 {
			fmt.Println("No plugins installed")
			return nil
		}
		for _, entry := range plugins {
			fmt.Printf("- %s (%s)\n", entry.Name, entry.Source)
		}
		return nil
	case "install":
		if opts.name == "" || opts.source == "" {
			return fmt.Errorf("plugins install requires --name and --source")
		}
		if err := plugin.FetchAndInstall(opts.project, opts.name, opts.source); err != nil {
			return err
		}
		fmt.Printf("Installed plugin %q from %s\n", opts.name, opts.source)
		return nil
	case "uninstall":
		if opts.name == "" {
			return fmt.Errorf("plugins uninstall requires --name")
		}
		removed, err := plugin.Uninstall(opts.project, opts.name)
		if err != nil {
			return err
		}
		if !removed {
			fmt.Printf("Plugin %q was not installed\n", opts.name)
			return nil
		}
		fmt.Printf("Uninstalled plugin %q\n", opts.name)
		return nil
	default:
		return fmt.Errorf("unknown plugins subcommand: %s", subcommand)
	}
}

func parsePluginsOptions(args []string) (pluginsOptions, error) {
	opts := pluginsOptions{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 >= len(args) {
				return pluginsOptions{}, fmt.Errorf("--project requires a path")
			}
			i++
			opts.project = args[i]
		case "--name":
			if i+1 >= len(args) {
				return pluginsOptions{}, fmt.Errorf("--name requires a value")
			}
			i++
			opts.name = args[i]
		case "--source":
			if i+1 >= len(args) {
				return pluginsOptions{}, fmt.Errorf("--source requires a value")
			}
			i++
			opts.source = args[i]
		default:
			return pluginsOptions{}, fmt.Errorf("unknown plugins flag: %s", args[i])
		}
	}
	return opts, nil
}
