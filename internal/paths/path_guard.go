package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathGuard resolves project-scoped paths and rejects physical escapes through symlinks/junctions.
type PathGuard struct {
	Root string
}

func (g PathGuard) ResolveExisting(path string) (string, error) {
	raw, err := g.resolveInput(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(raw)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", path, err)
	}
	if err := g.EnsureContained(resolved); err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func (g PathGuard) ResolveNewFile(path string) (string, error) {
	raw, err := g.resolveInput(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(raw); err == nil {
		return g.ResolveExisting(raw)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat path %q: %w", path, err)
	}
	parent := filepath.Dir(raw)
	for {
		resolvedParent, err := filepath.EvalSymlinks(parent)
		if err == nil {
			if err := g.EnsureContained(resolvedParent); err != nil {
				return "", err
			}
			suffix, err := filepath.Rel(parent, raw)
			if err != nil {
				return "", fmt.Errorf("resolve relative path: %w", err)
			}
			return filepath.Join(filepath.Clean(resolvedParent), suffix), nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("resolve parent %q: %w", parent, err)
		}
		next := filepath.Dir(parent)
		if next == parent {
			return "", fmt.Errorf("resolve parent %q: %w", parent, err)
		}
		parent = next
	}
}

func (g PathGuard) ResolveExistingOrNew(path string) (string, error) {
	raw, err := g.resolveInput(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(raw); err == nil {
		return g.ResolveExisting(raw)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat path %q: %w", path, err)
	}
	return g.ResolveNewFile(raw)
}

func (g PathGuard) EnsureContained(path string) error {
	root, err := g.resolvedRoot()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	rel, err := filepath.Rel(root, filepath.Clean(abs))
	if err != nil || !isContainedRel(rel) {
		return fmt.Errorf("path must be inside project: %s", path)
	}
	return nil
}

func (g PathGuard) resolveInput(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, err := g.rootAbs()
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Clean(filepath.Join(root, path)), nil
}

func (g PathGuard) resolvedRoot() (string, error) {
	root, err := g.rootAbs()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve root %q: %w", root, err)
	}
	return filepath.Clean(resolved), nil
}

func (g PathGuard) rootAbs() (string, error) {
	root := strings.TrimSpace(g.Root)
	if root == "" {
		return "", fmt.Errorf("project root is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}
	return filepath.Clean(abs), nil
}

func isContainedRel(rel string) bool {
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)))
}
