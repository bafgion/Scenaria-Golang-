package gui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type RecordRequest struct {
	URL              string `json:"url"`
	Output           string `json:"output"`
	IdleSeconds      int    `json:"idleSeconds"`
	Headless         bool   `json:"headless"`
	FilterRecording  bool   `json:"filterRecording"`
	NavOnlyRecording bool   `json:"navOnlyRecording"`
	HoverRecord      bool   `json:"hoverRecord"`
	AppendTo         string `json:"appendTo"`
	TestClient       string `json:"testClient"`
	FeatureName      string `json:"featureName"`
	ScenarioName     string `json:"scenarioName"`
}

type BaselineRecordRequest struct {
	Output       string   `json:"output"`
	FeatureName  string   `json:"featureName"`
	ScenarioName string   `json:"scenarioName"`
	Steps        []string `json:"steps"`
}

type ExportRequest struct {
	InputPath string `json:"inputPath"`
	Output    string `json:"output"`
	Format    string `json:"format"`
	BaseURL   string `json:"baseURL"`
	Force     bool   `json:"force"`
}

type ImportRequest struct {
	JSONPath   string `json:"jsonPath"`
	OutputPath string `json:"outputPath"`
	Force      bool   `json:"force"`
}

func (s *Service) ValidateFeature(text string) []ValidationIssue {
	return ValidateFeatureContent(text)
}

func (s *Service) Export(req ExportRequest) RunResult {
	input := strings.TrimSpace(req.InputPath)
	if input == "" {
		return RunResult{Error: "feature path is required"}
	}
	args := []string{input, "--output", req.Output, "--format", req.Format}
	if req.BaseURL != "" {
		args = append(args, "--base-url", req.BaseURL)
	}
	if req.Force {
		args = append(args, "--force")
	}
	out, err := captureCLI(func() error { return cliRunExport(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) ImportJSON(req ImportRequest) RunResult {
	args := []string{req.JSONPath, "--output", req.OutputPath}
	if req.Force {
		args = append(args, "--force")
	}
	out, err := captureCLI(func() error {
		return cliRunImportJSON(args)
	})
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) RecordLive(req RecordRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return RunResult{Error: "start URL is required"}
	}
	output := strings.TrimSpace(req.Output)
	if output == "" {
		output = filepath.Join(path, "recorded.feature")
	} else if !filepath.IsAbs(output) {
		output = filepath.Join(path, output)
	}
	idle := req.IdleSeconds
	if idle <= 0 {
		idle = 30
	}

	appCfg, err := s.loadAppSettings()
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	cleanURL := httpauth.ApplyURLCredentials(url, appCfg)
	if err := s.saveAppSettings(appCfg); err != nil {
		return RunResult{Error: err.Error()}
	}
	httpCreds := httpauth.PlaywrightHTTPCredentials(cleanURL, appCfg)

	s.mu.Lock()
	if s.recordCancel != nil {
		s.recordCancel()
	}
	session := recorder.NewLiveSession()
	s.liveSession = session
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(idle+30)*time.Second)
	s.recordCancel = cancel
	s.mu.Unlock()

	var testClient *settings.TestClient
	if name := strings.TrimSpace(req.TestClient); name != "" {
		client, err := settings.LoadTestClientByName(path, name)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		testClient = client
	}

	featureName := strings.TrimSpace(req.FeatureName)
	if featureName == "" {
		featureName = "Записанный сценарий"
	}
	scenarioName := strings.TrimSpace(req.ScenarioName)
	if scenarioName == "" {
		scenarioName = "Запись"
	}

	err = recorder.RecordLive(ctx, recorder.LiveOptions{
		StartURL:        cleanURL,
		FeatureName:     featureName,
		ScenarioName:    scenarioName,
		OutputPath:      output,
		Headless:        req.Headless,
		IdleTimeout:     time.Duration(idle) * time.Second,
		Session:         session,
		AppendTo:        strings.TrimSpace(req.AppendTo),
		FilterImportant: req.FilterRecording,
		NavOnly:         req.NavOnlyRecording,
		HoverRecord:     req.HoverRecord,
		TestClient:      testClient,
		HTTPCredentials: httpCreds,
	})

	s.mu.Lock()
	s.recordCancel = nil
	s.liveSession = nil
	s.mu.Unlock()

	if err != nil {
		return RunResult{Error: fmt.Errorf("record: %w", err).Error()}
	}
	return RunResult{Output: fmt.Sprintf("Запись сохранена: %s\n", output)}
}

func (s *Service) RecordBaseline(req BaselineRecordRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	output := strings.TrimSpace(req.Output)
	if output == "" {
		output = filepath.Join(path, "recorded.feature")
	} else if !filepath.IsAbs(output) {
		output = filepath.Join(path, output)
	}
	featureName := strings.TrimSpace(req.FeatureName)
	if featureName == "" {
		featureName = "Записанный сценарий"
	}
	scenarioName := strings.TrimSpace(req.ScenarioName)
	if scenarioName == "" {
		scenarioName = "Базовый сценарий"
	}
	args := []string{
		"--output", output,
		"--feature", featureName,
		"--scenario", scenarioName,
	}
	for _, step := range req.Steps {
		step = strings.TrimSpace(step)
		if step != "" {
			args = append(args, "--step", step)
		}
	}
	out, err := captureCLI(func() error { return cliRunRecord(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) PauseRecording() {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session != nil {
		session.Pause()
	}
}

func (s *Service) ResumeRecording() {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session != nil {
		session.Resume()
	}
}

func (s *Service) IsRecordingPaused() bool {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return false
	}
	return session.IsPaused()
}

func (s *Service) CancelRecording() {
	s.mu.Lock()
	cancel := s.recordCancel
	s.recordCancel = nil
	session := s.liveSession
	s.liveSession = nil
	s.mu.Unlock()
	if session != nil {
		session.Clear()
	}
	if cancel != nil {
		cancel()
	}
}

func (s *Service) FocusBrowser() error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return fmt.Errorf("браузер не открыт")
	}
	return session.FocusBrowser()
}

func (s *Service) UpdateRecordingOptions(filter, navOnly, hover bool) error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return fmt.Errorf("запись не активна")
	}
	return session.ApplyRecorderConfig(filter, navOnly, hover)
}

func (s *Service) UndoRecordedStep() bool {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return false
	}
	return session.UndoLastStep()
}

// cli hooks for testing
var (
	cliRunExport     = func(args []string) error { return runExport(args) }
	cliRunImportJSON = func(args []string) error { return runImportJSON(args) }
	cliRunRecord     = func(args []string) error { return runRecord(args) }
	cliRunVA         = func(args []string) error { return runVA(args) }
)
