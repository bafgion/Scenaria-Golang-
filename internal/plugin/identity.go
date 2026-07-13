package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

// ValidatePluginID accepts only stable path-segment identifiers for project addons.
func ValidatePluginID(id string) error {
	if id == "" {
		return fmt.Errorf("plugin id is required")
	}
	if id != strings.TrimSpace(id) {
		return fmt.Errorf("invalid plugin id %q: leading or trailing whitespace is not allowed", id)
	}
	if id == "." || id == ".." {
		return fmt.Errorf("invalid plugin id %q", id)
	}
	if filepath.IsAbs(id) || filepath.VolumeName(id) != "" {
		return fmt.Errorf("invalid plugin id %q: paths are not allowed", id)
	}
	if strings.HasSuffix(id, ".") {
		return fmt.Errorf("invalid plugin id %q: trailing dot is not allowed", id)
	}
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '.', r == '-', r == '_':
		default:
			return fmt.Errorf("invalid plugin id %q: only letters, digits, dot, dash and underscore are allowed", id)
		}
	}
	return nil
}

func addonsRoot(projectRoot string) (string, error) {
	project, err := resolveProjectRoot(projectRoot)
	if err != nil {
		return "", err
	}
	guard := paths.PathGuard{Root: project}
	root := filepath.Join(project, "addons")
	if _, err := os.Lstat(root); err == nil {
		return guard.ResolveExisting("addons")
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat addons root: %w", err)
	}
	return root, nil
}

func resolveProjectRoot(projectRoot string) (string, error) {
	projectRoot = strings.TrimSpace(projectRoot)
	if projectRoot == "" {
		return "", fmt.Errorf("project root is required")
	}
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return filepath.Clean(resolved), nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("resolve project root %q: %w", abs, err)
	}
	parent := filepath.Dir(abs)
	name := filepath.Base(abs)
	return paths.PathGuard{Root: parent}.ResolveNewFile(name)
}

func addonPath(projectRoot, pluginID string) (string, error) {
	if err := ValidatePluginID(pluginID); err != nil {
		return "", err
	}
	root, err := addonsRoot(projectRoot)
	if err != nil {
		return "", err
	}
	dest := filepath.Clean(filepath.Join(root, pluginID))
	if _, err := os.Lstat(dest); err == nil {
		return paths.PathGuard{Root: root}.ResolveExisting(pluginID)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat plugin path: %w", err)
	}
	if _, err := os.Lstat(root); err == nil {
		return paths.PathGuard{Root: root}.ResolveNewFile(pluginID)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat addons root: %w", err)
	}
	return dest, nil
}
