package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CanonicalPath returns the absolute resolved form of path.
// For not-yet-created files it resolves the nearest existing ancestor via symlinks.
func CanonicalPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return filepath.Clean(resolved), nil
	} else if !os.IsNotExist(err) {
		return "", err
	}
	dir := filepath.Dir(abs)
	for {
		if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
			rel, err := filepath.Rel(dir, abs)
			if err != nil {
				return "", err
			}
			return filepath.Clean(filepath.Join(resolvedDir, rel)), nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Clean(abs), nil
		}
		dir = parent
	}
}

// SamePath reports whether two paths refer to the same location.
func SamePath(a, b string) bool {
	ca, ea := CanonicalPath(a)
	cb, eb := CanonicalPath(b)
	if ea == nil && eb == nil {
		return ca == cb || strings.EqualFold(ca, cb)
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
