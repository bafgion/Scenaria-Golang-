package gui

import (
	"context"
	"errors"
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
	BrowseOnly       bool   `json:"browseOnly"`
}

// OpenBrowserRequest opens a Playwright browser without recording until capture is started explicitly.
type OpenBrowserRequest struct {
	URL              string `json:"url"`
	Headless         bool   `json:"headless"`
	TestClient       string `json:"testClient"`
	Output           string `json:"output"`
	IdleSeconds      int    `json:"idleSeconds"`
	FilterRecording  bool   `json:"filterRecording"`
	NavOnlyRecording bool   `json:"navOnlyRecording"`
	HoverRecord      bool   `json:"hoverRecord"`
	AppendTo         string `json:"appendTo"`
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

func (s *Service) OpenBrowser(req OpenBrowserRequest, emit func(string, any)) RunResult {
	return s.RecordLive(RecordRequest{
		URL:              req.URL,
		Output:           req.Output,
		IdleSeconds:      req.IdleSeconds,
		Headless:         req.Headless,
		FilterRecording:  req.FilterRecording,
		NavOnlyRecording: req.NavOnlyRecording,
		HoverRecord:      req.HoverRecord,
		AppendTo:         req.AppendTo,
		TestClient:       req.TestClient,
		FeatureName:      req.FeatureName,
		ScenarioName:     req.ScenarioName,
		BrowseOnly:       true,
	}, emit)
}

func (s *Service) liveRecordCallbacks(emit func(string, any)) recorder.LiveCallbacks {
	return recorder.LiveCallbacks{
		OnCaptureStart: func(resume bool) {
			if emit != nil {
				emit("record-started", map[string]any{"resume": resume})
			}
		},
		OnCaptureStop: func() {
			if emit != nil {
				emit("record-stopped", nil)
			}
		},
		OnPickerRequest: func() {
			if emit != nil {
				emit("toolbar-picker", nil)
			}
		},
		OnBrowserLost: func() {
			if emit != nil {
				emit("browser-lost", nil)
			}
		},
		OnStepRecorded: func(index int, line string) {
			if emit != nil {
				emit("record-step", map[string]any{"index": index, "line": line})
			}
		},
	}
}

func (s *Service) HasLiveBrowser() bool {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	return session != nil && session.BrowserAlive()
}

func (s *Service) emitLiveRecordedSteps(session *recorder.LiveSession, emit func(string, any)) {
	if emit == nil || session == nil {
		return
	}
	session.EachRecordedLine(func(index int, line string) {
		emit("record-step", map[string]any{"index": index, "line": line})
	})
}

func (s *Service) startCaptureOnExistingSession(req RecordRequest, emit func(string, any)) RunResult {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil || !session.BrowserAlive() {
		return RunResult{Error: "браузер не открыт"}
	}
	appCfg, err := s.loadAppSettings()
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	_ = session.ApplyRecorderOptions(
		req.FilterRecording,
		req.NavOnlyRecording,
		req.HoverRecord,
		appCfg.ScrollBeforeClick,
		appCfg.HoverRecordMinMs,
	)
	if session.CaptureEnabled() {
		if emit != nil {
			emit("record-started", map[string]any{"append": true, "sync": true})
		}
		return RunResult{}
	}
	replay := recorder.ShouldSyncRecordedStepsOnCaptureStart(session)
	if err := session.BeginCapture(); err != nil {
		return RunResult{Error: err.Error()}
	}
	if emit != nil {
		emit("record-started", map[string]any{"append": true, "resume": !replay})
		if replay {
			s.emitLiveRecordedSteps(session, emit)
		}
	}
	return RunResult{}
}

func (s *Service) RecordLive(req RecordRequest, emit func(string, any)) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	if strings.TrimSpace(req.URL) != "" && !strings.HasPrefix(strings.TrimSpace(req.URL), "http") {
		return RunResult{Error: "укажите корректный URL (https://…) или оставьте поле пустым"}
	}

	s.mu.RLock()
	existing := s.liveSession
	s.mu.RUnlock()
	if existing != nil && existing.BrowserAlive() && !req.BrowseOnly {
		return s.startCaptureOnExistingSession(req, emit)
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
	cleanURL := strings.TrimSpace(req.URL)
	if cleanURL != "" {
		cleanURL = httpauth.ApplyURLCredentials(cleanURL, appCfg)
		if err := s.saveAppSettings(appCfg); err != nil {
			return RunResult{Error: err.Error()}
		}
	}
	httpCreds := httpauth.PlaywrightHTTPCredentials(cleanURL, appCfg)

	s.mu.Lock()
	if s.recordCancel != nil {
		s.recordCancel()
	}
	session := recorder.NewLiveSession()
	s.liveSession = session
	var ctx context.Context
	var cancel context.CancelFunc
	if req.BrowseOnly {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(idle+30)*time.Second)
	}
	s.recordCancel = cancel
	s.recordEmit = emit
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.recordEmit = nil
		s.mu.Unlock()
	}()

	idleTimeout := time.Duration(idle) * time.Second
	if req.BrowseOnly {
		idleTimeout = 0
	}

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
		StartURL:          cleanURL,
		FeatureName:       featureName,
		ScenarioName:      scenarioName,
		OutputPath:        output,
		Headless:          req.Headless,
		IdleTimeout:       idleTimeout,
		Session:           session,
		AppendTo:          strings.TrimSpace(req.AppendTo),
		FilterImportant:   req.FilterRecording,
		NavOnly:           req.NavOnlyRecording,
		HoverRecord:       req.HoverRecord,
		ScrollBeforeClick: appCfg.ScrollBeforeClick,
		HoverRecordMinMs:  appCfg.HoverRecordMinMs,
		TestClient:        testClient,
		HTTPCredentials: httpCreds,
		BrowseOnly:      req.BrowseOnly,
		Callbacks: s.liveRecordCallbacks(emit),
	})

	s.mu.Lock()
	s.recordCancel = nil
	s.liveSession = nil
	s.mu.Unlock()

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return RunResult{Output: "Браузер закрыт."}
		}
		return RunResult{Error: fmt.Errorf("record: %w", err).Error()}
	}
	return RunResult{Output: fmt.Sprintf("Запись сохранена: %s\n", output)}
}

func (s *Service) BeginRecordingCapture() (bool, error) {
	s.mu.RLock()
	session := s.liveSession
	emit := s.recordEmit
	s.mu.RUnlock()
	if session == nil {
		return false, fmt.Errorf("браузер не открыт")
	}
	if session.CaptureEnabled() {
		return false, nil
	}
	replay := recorder.ShouldSyncRecordedStepsOnCaptureStart(session)
	if err := session.BeginCapture(); err != nil {
		return false, err
	}
	if replay && emit != nil {
		s.emitLiveRecordedSteps(session, emit)
	}
	return true, nil
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

type BrowserSessionDTO struct {
	BrowserOpen bool `json:"browserOpen"`
	Recording   bool `json:"recording"`
	Paused      bool `json:"paused"`
	StepCount   int  `json:"stepCount"`
}

func (s *Service) PollBrowserSession() BrowserSessionDTO {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil || !session.BrowserAlive() {
		return BrowserSessionDTO{}
	}
	return BrowserSessionDTO{
		BrowserOpen: true,
		Recording:   session.CaptureEnabled(),
		Paused:      session.IsPaused(),
		StepCount:   session.RecordedStepCount(),
	}
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

func (s *Service) StopRecordingCapture() error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return fmt.Errorf("браузер не открыт")
	}
	if !session.CaptureEnabled() {
		return nil
	}
	session.EndCapture()
	return nil
}

func (s *Service) CloseBrowser() {
	s.mu.Lock()
	cancel := s.recordCancel
	s.recordCancel = nil
	session := s.liveSession
	s.liveSession = nil
	s.recordEmit = nil
	s.mu.Unlock()
	if session != nil {
		session.Clear()
	}
	if cancel != nil {
		cancel()
	}
}

func (s *Service) CancelRecording() {
	s.CloseBrowser()
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

func (s *Service) UpdateRecordingOptions(filter, navOnly, hover, headless, scrollBefore bool, hoverMinMs int) error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil {
		return fmt.Errorf("запись не активна")
	}
	_ = session.ApplyRecorderOptions(filter, navOnly, hover, scrollBefore, hoverMinMs)
	session.RequestHeadless(headless)
	return nil
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
