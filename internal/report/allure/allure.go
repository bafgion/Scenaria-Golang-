// Package allure writes Allure 2 result files for scenaria runs.
// See docs/ROADMAP.md phase 3.
package allure

import (
	"fmt"
	"os"
	"path/filepath"
)

type Options struct {
	OutputDir string
}

// WritePlaceholder creates the output directory; full writer is planned in v0.14.
func WritePlaceholder(opts Options) error {
	if opts.OutputDir == "" {
		return fmt.Errorf("allure output directory is required")
	}
	return os.MkdirAll(filepath.Clean(opts.OutputDir), 0o755)
}
