package gui

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/paths"
)

const startupTempMaxAge = 24 * time.Hour

func (s *Service) CleanupStartupTempArtifacts(projectRoot string) {
	if strings.TrimSpace(projectRoot) == "" {
		return
	}
	if err := cleanupStartupTempArtifactsForProject(projectRoot, startupTempMaxAge, time.Now()); err != nil {
		logx.Warn("startup temp cleanup failed", "project_root", projectRoot, "error", err)
	}
}

func (s *Service) CleanupGlobalStartupTemps(tempRoot string) {
	if strings.TrimSpace(tempRoot) == "" {
		return
	}
	if err := cleanupGlobalStartupTemps(tempRoot, startupTempMaxAge, time.Now()); err != nil {
		logx.Warn("startup global temp cleanup failed", "temp_root", tempRoot, "error", err)
	}
}

func cleanupStartupTempArtifactsForProject(projectRoot string, maxAge time.Duration, now time.Time) error {
	scenariaDir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return err
	}
	cutoff := now.Add(-maxAge)
	if err := cleanupProjectTempDir(filepath.Join(scenariaDir, "temp"), cutoff); err != nil {
		return err
	}
	if err := cleanupProjectRunTemps(filepath.Join(scenariaDir, "runs"), cutoff); err != nil {
		return err
	}
	return nil
}

func cleanupGlobalStartupTemps(tempRoot string, maxAge time.Duration, now time.Time) error {
	entries, err := os.ReadDir(tempRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	cutoff := now.Add(-maxAge)
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "scenaria-run-") {
			continue
		}
		full := filepath.Join(tempRoot, name)
		if stale, err := isOlderThan(full, cutoff); err != nil {
			return err
		} else if !stale {
			continue
		}
		if err := os.RemoveAll(full); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func cleanupProjectTempDir(tempDir string, cutoff time.Time) error {
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "run-") && !strings.HasPrefix(name, "tmp-") && !strings.HasPrefix(name, ".tmp-") {
			continue
		}
		full := filepath.Join(tempDir, name)
		if stale, err := isOlderThan(full, cutoff); err != nil {
			return err
		} else if !stale {
			continue
		}
		if err := os.RemoveAll(full); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func cleanupProjectRunTemps(runsDir string, cutoff time.Time) error {
	if _, err := os.Stat(runsDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return filepath.WalkDir(runsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == runsDir {
			return nil
		}
		name := d.Name()
		if !isScenariaTempRunArtifact(name) {
			return nil
		}
		stale, err := isOlderThan(path, cutoff)
		if err != nil {
			return err
		}
		if !stale {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
}

func isScenariaTempRunArtifact(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	switch {
	case strings.HasPrefix(name, ".scenaria-write-"):
		return true
	case strings.HasPrefix(name, ".allure-write-"):
		return true
	case strings.HasPrefix(name, ".screenshot-"):
		return true
	case strings.Contains(name, ".tmp-"):
		return true
	default:
		return false
	}
}

func isOlderThan(path string, cutoff time.Time) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.ModTime().Before(cutoff), nil
}
