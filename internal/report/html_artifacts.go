package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const htmlArtifactMarker = ".scenaria-html-artifacts"

func prepareHTMLArtifactDirs(reportPath string) error {
	reportDir := filepath.Dir(reportPath)
	for _, name := range []string{"screenshots", "traces"} {
		if err := cleanHTMLArtifactDir(filepath.Join(reportDir, name)); err != nil {
			return err
		}
	}
	return nil
}

func cleanHTMLArtifactDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
				return fmt.Errorf("create html artifact dir %q: %w", dir, mkErr)
			}
			return writeHTMLArtifactMarker(dir)
		}
		return fmt.Errorf("read html artifact dir %q: %w", dir, err)
	}
	marked := false
	for _, entry := range entries {
		if entry.Name() == htmlArtifactMarker {
			marked = true
			break
		}
	}
	if !marked && len(entries) > 0 && !canAdoptLegacyHTMLArtifacts(dir, entries) {
		return fmt.Errorf("refusing to clean unmarked html artifact directory %q", dir)
	}
	for _, entry := range entries {
		if entry.Name() == htmlArtifactMarker {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			return fmt.Errorf("clean html artifact %q: %w", filepath.Join(dir, entry.Name()), err)
		}
	}
	return writeHTMLArtifactMarker(dir)
}

func canAdoptLegacyHTMLArtifacts(dir string, entries []os.DirEntry) bool {
	kind := filepath.Base(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			return false
		}
		name := entry.Name()
		switch kind {
		case "screenshots":
			if filepath.Ext(name) != ".png" || (!strings.Contains(name, "__scenario") && !strings.Contains(name, "__step_")) {
				return false
			}
		case "traces":
			if filepath.Ext(name) != ".zip" || !strings.Contains(name, "__") {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func writeHTMLArtifactMarker(dir string) error {
	path := filepath.Join(dir, htmlArtifactMarker)
	if err := os.WriteFile(path, []byte("Scenaria HTML report artifacts\n"), 0o644); err != nil {
		return fmt.Errorf("write html artifact marker %q: %w", path, err)
	}
	return nil
}
