package report

import (
	"os"
	"path/filepath"
	"strings"
)

func writeScreenshotArtifact(dir, name string, data []byte) string {
	if len(data) == 0 || strings.TrimSpace(dir) == "" || strings.TrimSpace(name) == "" {
		return ""
	}
	name = filepath.Base(name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ""
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return ""
	}
	return filepath.ToSlash(filepath.Join(filepath.Base(dir), name))
}
