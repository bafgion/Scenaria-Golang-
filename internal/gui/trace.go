package gui

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

func (s *Service) resolveTracePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		root := s.ProjectPath()
		if root == "" {
			return "", fmt.Errorf("open a project folder first")
		}
		resolved, err := paths.ScenariaArtifactPath(root, "traces")
		if err != nil {
			return "", err
		}
		path = resolved
	} else {
		root := s.ProjectPath()
		if root == "" {
			return "", fmt.Errorf("open a project folder first")
		}
		confined, err := paths.ConfineToProjectRoot(root, path)
		if err != nil {
			return "", fmt.Errorf("trace path outside project: %w", err)
		}
		path = confined
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("trace not found: %s", path)
	}
	if info.IsDir() {
		latest, err := newestTraceZip(path)
		if err != nil {
			return "", err
		}
		return latest, nil
	}
	if !strings.HasSuffix(strings.ToLower(path), ".zip") {
		return "", fmt.Errorf("trace must be a .zip file or traces directory")
	}
	return path, nil
}

func newestTraceZip(dir string) (string, error) {
	var (
		latestPath string
		latestTime time.Time
		found      bool
	)
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(p), ".zip") {
			return nil
		}
		info, statErr := d.Info()
		if statErr != nil {
			return statErr
		}
		if !found || info.ModTime().After(latestTime) {
			found = true
			latestTime = info.ModTime()
			latestPath = p
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("no trace .zip files in %s", dir)
	}
	return latestPath, nil
}

// OpenTrace launches Playwright trace viewer for the latest or given trace archive.
func (s *Service) OpenTrace(path string) RunResult {
	tracePath, err := s.resolveTracePath(path)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	abs, err := filepath.Abs(tracePath)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	if err := startPlaywrightShowTrace(abs); err != nil {
		return RunResult{Error: err.Error()}
	}
	return RunResult{Output: abs}
}

func startPlaywrightShowTrace(traceZip string) error {
	if _, err := exec.LookPath("playwright"); err == nil {
		cmd := exec.Command("playwright", "show-trace", traceZip)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start playwright show-trace: %w", err)
		}
		return nil
	}
	if _, err := exec.LookPath("npx"); err == nil {
		cmd := exec.Command("npx", "playwright", "show-trace", traceZip)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start npx playwright show-trace: %w", err)
		}
		return nil
	}
	return fmt.Errorf("playwright CLI not found in PATH — install Node.js and run: npx playwright install")
}
