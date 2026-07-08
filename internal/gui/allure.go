package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type allureServeState struct {
	mu  sync.Mutex
	cmd *exec.Cmd
	dir string
}

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

func AllureCLIAvailable() bool {
	_, err := exec.LookPath("allure")
	return err == nil
}

type AllureStatusDTO struct {
	Installed  bool   `json:"installed"`
	Running    bool   `json:"running"`
	ResultsDir string `json:"resultsDir"`
}

func (s *Service) AllureStatus(dir string) AllureStatusDTO {
	status := AllureStatusDTO{Installed: AllureCLIAvailable()}
	resultsDir, err := s.resolveAllureDir(dir)
	if err == nil {
		status.ResultsDir = resultsDir
	}
	s.allureServe.mu.Lock()
	status.Running = s.allureServe.cmd != nil && s.allureServe.cmd.Process != nil
	if status.Running && s.allureServe.dir != "" {
		status.ResultsDir = s.allureServe.dir
	}
	s.allureServe.mu.Unlock()
	return status
}

func (s *Service) ServeAllure(dir string) RunResult {
	resultsDir, err := s.resolveAllureDir(dir)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	if !AllureCLIAvailable() {
		return RunResult{
			Error: "allure CLI not found in PATH — install from https://docs.qameta.io/allure/#_installing_a_commandline",
		}
	}
	s.allureServe.mu.Lock()
	if s.allureServe.cmd != nil && s.allureServe.cmd.Process != nil {
		activeDir := s.allureServe.dir
		s.allureServe.mu.Unlock()
		if activeDir == "" {
			activeDir = resultsDir
		}
		return RunResult{Output: fmt.Sprintf("Allure serve уже запущен (%s). Если окно закрыто — нажмите «Allure снова».\n", activeDir)}
	}
	cmd := exec.Command("allure", "serve", resultsDir)
	if err := cmd.Start(); err != nil {
		s.allureServe.mu.Unlock()
		return RunResult{Error: fmt.Sprintf("start allure serve: %v", err)}
	}
	s.allureServe.cmd = cmd
	s.allureServe.dir = resultsDir
	s.allureServe.mu.Unlock()
	go func() {
		_ = cmd.Wait()
		s.allureServe.mu.Lock()
		if s.allureServe.cmd == cmd {
			s.allureServe.cmd = nil
			s.allureServe.dir = ""
		}
		s.allureServe.mu.Unlock()
	}()
	return RunResult{Output: fmt.Sprintf("Allure serve: %s\n", resultsDir)}
}

// stopAllureServe terminates a background `allure serve` started from this GUI service.
func (s *Service) stopAllureServe() {
	if s == nil {
		return
	}
	s.allureServe.mu.Lock()
	defer s.allureServe.mu.Unlock()
	if s.allureServe.cmd != nil && s.allureServe.cmd.Process != nil {
		_ = s.allureServe.cmd.Process.Kill()
	}
	s.allureServe.cmd = nil
	s.allureServe.dir = ""
}

func (s *Service) OpenHTMLReport(path string) RunResult {
	path = strings.TrimSpace(path)
	root := s.ProjectPath()
	if root == "" {
		return RunResult{Error: "open a project folder first"}
	}
	if path == "" {
		resolved, err := paths.ScenariaArtifactPath(root, "report.html")
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		path = resolved
	} else if !filepath.IsAbs(path) {
		confined, err := paths.ConfineToProjectRoot(root, path)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		path = paths.RemapScenariaArtifact(root, confined)
	}
	if _, err := os.Stat(path); err != nil {
		return RunResult{Error: fmt.Sprintf("report not found: %s", path)}
	}
	openMode := "full"
	if cfg, err := settings.LoadProjectConfig(root); err == nil {
		openMode = cfg.HTMLReportOpenMode
	}
	path = report.PreferredHTMLReportPath(path, openMode)
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
