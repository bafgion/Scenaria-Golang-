package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

var (
	allureServeMu  sync.Mutex
	allureServeCmd *exec.Cmd
)

func (s *Service) defaultAllureDir() (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return paths.ScenariaArtifactPath(path, "allure-results")
}

func (s *Service) resolveAllureDir(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return s.defaultAllureDir()
	}
	if !filepath.IsAbs(dir) {
		path := s.ProjectPath()
		if path == "" {
			return "", fmt.Errorf("open a project folder first")
		}
		dir = filepath.Join(path, dir)
	}
	if !s.ArtifactExists(dir) {
		return "", fmt.Errorf("allure results not found: %s", dir)
	}
	return dir, nil
}

func (s *Service) ServeAllure(dir string) RunResult {
	resultsDir, err := s.resolveAllureDir(dir)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	if _, err := exec.LookPath("allure"); err != nil {
		return RunResult{
			Error: "allure CLI not found in PATH — install from https://docs.qameta.io/allure/",
		}
	}
	allureServeMu.Lock()
	if allureServeCmd != nil && allureServeCmd.Process != nil {
		allureServeMu.Unlock()
		return RunResult{Output: "Allure serve уже запущен\n"}
	}
	cmd := exec.Command("allure", "serve", resultsDir)
	if err := cmd.Start(); err != nil {
		allureServeMu.Unlock()
		return RunResult{Error: fmt.Sprintf("start allure serve: %v", err)}
	}
	allureServeCmd = cmd
	allureServeMu.Unlock()
	go func() {
		_ = cmd.Wait()
		allureServeMu.Lock()
		if allureServeCmd == cmd {
			allureServeCmd = nil
		}
		allureServeMu.Unlock()
	}()
	return RunResult{Output: fmt.Sprintf("Allure serve: %s\n", resultsDir)}
}

func (s *Service) OpenHTMLReport(path string) RunResult {
	path = strings.TrimSpace(path)
	if path == "" {
		root := s.ProjectPath()
		if root == "" {
			return RunResult{Error: "open a project folder first"}
		}
		resolved, err := paths.ScenariaArtifactPath(root, "report.html")
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		path = resolved
	} else if !filepath.IsAbs(path) {
		root := s.ProjectPath()
		if root == "" {
			return RunResult{Error: "open a project folder first"}
		}
		path = filepath.Join(root, path)
	}
	if _, err := os.Stat(path); err != nil {
		return RunResult{Error: fmt.Sprintf("report not found: %s", path)}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	return RunResult{Output: abs}
}

func fileURL(path string) string {
	slash := filepath.ToSlash(path)
	if strings.HasPrefix(slash, "/") {
		return "file://" + slash
	}
	return "file:///" + slash
}
