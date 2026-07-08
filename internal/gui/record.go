package gui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/paths"
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
	return s.editorAnalyzer().ValidateFeature(text)
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
	out, err := s.cliRunner().Export(args)
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
	out, err := s.cliRunner().ImportJSON(args)
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

func resolveRecordTargetPath(projectPath string, req RecordRequest) string {
	if appendTo := strings.TrimSpace(req.AppendTo); appendTo != "" {
		if confined, err := paths.ConfineToProjectRoot(projectPath, appendTo); err == nil {
			return confined
		}
	}
	output := strings.TrimSpace(req.Output)
	if output == "" {
		return filepath.Join(projectPath, "recorded.feature")
	}
	if confined, err := paths.ConfineToProjectRoot(projectPath, output); err == nil {
		return confined
	}
	return filepath.Join(projectPath, "recorded.feature")
}

func (s *Service) guardedRecordEmit(gen uint64, targetPath string, emit func(string, any)) func(string, any) {
	if emit == nil {
		return nil
	}
	return func(name string, payload any) {
		s.mu.RLock()
		active := s.recordGen == gen
		s.mu.RUnlock()
		if !active {
			logx.Debug("stale record event skipped", "event", name, "expected_gen", gen)
			return
		}
		if m, ok := payload.(map[string]any); ok && m != nil {
			if targetPath != "" {
				if _, has := m["targetPath"]; !has {
					m["targetPath"] = targetPath
				}
			}
		}
		emit(name, payload)
	}
}

func (s *Service) liveRecordCallbacks(emit func(string, any), browseOnly bool, output, targetPath, recordSessionID, browserSessionID string) recorder.LiveCallbacks {
	eventPayload := func(extra map[string]any) map[string]any {
		m := map[string]any{
			"recordSessionId":  recordSessionID,
			"browserSessionId": browserSessionID,
		}
		if targetPath != "" {
			m["targetPath"] = targetPath
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	return recorder.LiveCallbacks{
		OnBrowserOpened: func() {
			if emit == nil {
				return
			}
			if browseOnly {
				emit("browser-opened", eventPayload(nil))
				return
			}
			emit("record-started", eventPayload(map[string]any{
				"resume": false,
				"output": output,
			}))
		},
		OnCaptureStart: func(resume bool) {
			if emit != nil {
				emit("record-started", eventPayload(map[string]any{"resume": resume}))
			}
		},
		OnCaptureStop: func(reason string) {
			if emit != nil {
				payload := eventPayload(map[string]any{"reason": reason})
				if reason == "idle" && s.recordIdleSeconds > 0 {
					payload["idleSeconds"] = s.recordIdleSeconds
				}
				emit("record-stopped", payload)
			}
		},
		OnPickerRequest: func() {
			if emit != nil {
				emit("toolbar-picker", eventPayload(nil))
			}
		},
		OnBrowserLost: func() {
			if emit != nil {
				emit("browser-lost", eventPayload(nil))
			}
		},
		OnStepRecorded: func(event recorder.RecordStepEvent) {
			if emit == nil {
				return
			}
			payload := eventPayload(map[string]any{"op": string(event.Op)})
			switch event.Op {
			case recorder.RecordStepUpsert:
				payload["index"] = event.Index
				payload["line"] = event.Line
			case recorder.RecordStepDelete:
				payload["index"] = event.Index
			case recorder.RecordStepSnapshot:
				payload["lines"] = event.Lines
			}
			emit("record-step", payload)
		},
	}
}

func (s *Service) HasLiveBrowser() bool {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	return session != nil && session.BrowserAlive()
}

func (s *Service) emitLiveRecordedSteps(session *recorder.LiveSession, emit func(string, any), targetPath, recordSessionID, browserSessionID string) {
	if emit == nil || session == nil {
		return
	}
	lines := session.SnapshotRecordedSteps()
	payload := map[string]any{
		"op":               string(recorder.RecordStepSnapshot),
		"recordSessionId":  recordSessionID,
		"browserSessionId": browserSessionID,
	}
	if targetPath != "" {
		payload["targetPath"] = targetPath
	}
	if len(lines) == 0 {
		payload["op"] = string(recorder.RecordStepReset)
	} else {
		payload["lines"] = lines
	}
	emit("record-step", payload)
}

func (s *Service) emitRecordStepReset() {
	s.mu.RLock()
	emit := s.recordEmit
	myGen := s.recordGen
	targetPath := s.recordTargetPath
	recordSessionID := s.recordSessionID
	browserSessionID := s.browserSessionID
	s.mu.RUnlock()
	if emit == nil {
		return
	}
	emit = s.guardedRecordEmit(myGen, targetPath, emit)
	emit("record-step", map[string]any{
		"op":               string(recorder.RecordStepReset),
		"recordSessionId":  recordSessionID,
		"browserSessionId": browserSessionID,
		"targetPath":       targetPath,
	})
}

func (s *Service) emitRecordStepDelete(index int) {
	s.mu.RLock()
	emit := s.recordEmit
	myGen := s.recordGen
	targetPath := s.recordTargetPath
	recordSessionID := s.recordSessionID
	browserSessionID := s.browserSessionID
	s.mu.RUnlock()
	if emit == nil || index < 0 {
		return
	}
	emit = s.guardedRecordEmit(myGen, targetPath, emit)
	emit("record-step", map[string]any{
		"op":               string(recorder.RecordStepDelete),
		"index":            index,
		"recordSessionId":  recordSessionID,
		"browserSessionId": browserSessionID,
		"targetPath":       targetPath,
	})
}

func (s *Service) startCaptureOnExistingSession(req RecordRequest, emit func(string, any)) RunResult {
	s.mu.RLock()
	session := s.liveSession
	myGen := s.recordGen
	recordSessionID := s.recordSessionID
	browserSessionID := s.browserSessionID
	targetPath := s.recordTargetPath
	s.mu.RUnlock()
	if session == nil || !session.BrowserAlive() {
		return RunResult{Error: "браузер не открыт"}
	}
	emit = s.guardedRecordEmit(myGen, targetPath, emit)
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
			emit("record-started", map[string]any{
				"append":           true,
				"sync":             true,
				"recordSessionId":  recordSessionID,
				"browserSessionId": browserSessionID,
				"targetPath":       targetPath,
			})
		}
		return RunResult{}
	}
	replay := recorder.ShouldSyncRecordedStepsOnCaptureStart(session)
	if err := session.BeginCapture(); err != nil {
		return RunResult{Error: err.Error()}
	}
	if emit != nil {
		emit("record-started", map[string]any{
			"append":           true,
			"resume":           !replay,
			"recordSessionId":  recordSessionID,
			"browserSessionId": browserSessionID,
			"targetPath":       targetPath,
		})
		if replay {
			s.emitLiveRecordedSteps(session, emit, targetPath, recordSessionID, browserSessionID)
		}
	}
	return RunResult{}
}

func (s *Service) RecordLive(req RecordRequest, emit func(string, any)) RunResult {
	s.activePlaywright.Add(1)
	defer s.activePlaywright.Done()

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
	} else {
		confined, err := paths.ConfineToProjectRoot(path, output)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		output = confined
	}
	appendTo := strings.TrimSpace(req.AppendTo)
	if appendTo != "" {
		confined, err := paths.ConfineToProjectRoot(path, appendTo)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		appendTo = confined
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
	s.recordGen++
	myGen := s.recordGen
	recordSessionID := fmt.Sprintf("record-%d", myGen)
	browserSessionID := fmt.Sprintf("browser-%d", myGen)
	targetPath := resolveRecordTargetPath(path, req)
	session := recorder.NewLiveSession()
	s.liveSession = session
	s.recordSessionID = recordSessionID
	s.browserSessionID = browserSessionID
	s.recordTargetPath = targetPath
	ctx, cancel := context.WithCancel(context.Background())
	s.recordCtx = ctx
	s.recordCancel = cancel
	s.recordEmit = emit
	s.recordIdleSeconds = idle
	s.mu.Unlock()

	emit = s.guardedRecordEmit(myGen, targetPath, emit)

	clearRecordSession := func() {
		s.mu.Lock()
		if s.recordGen == myGen {
			s.recordEmit = nil
		} else {
			logx.Debug("stale record emit cleanup skipped", "record_session_id", recordSessionID, "expected_gen", myGen, "current_gen", s.recordGen)
		}
		s.mu.Unlock()
	}
	defer clearRecordSession()

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
		AppendTo:          appendTo,
		FilterImportant:   req.FilterRecording,
		NavOnly:           req.NavOnlyRecording,
		HoverRecord:       req.HoverRecord,
		ScrollBeforeClick: appCfg.ScrollBeforeClick,
		HoverRecordMinMs:  appCfg.HoverRecordMinMs,
		TestClient:        testClient,
		HTTPCredentials:   httpCreds,
		BrowseOnly:        req.BrowseOnly,
		Callbacks:         s.liveRecordCallbacks(emit, req.BrowseOnly, output, targetPath, recordSessionID, browserSessionID),
	})

	s.mu.Lock()
	if s.recordGen == myGen {
		s.rememberClosedRecordSessionLocked(recordSessionID, browserSessionID)
		s.recordCancel = nil
		s.liveSession = nil
		s.recordCtx = nil
		s.recordSessionID = ""
		s.browserSessionID = ""
		s.recordTargetPath = ""
	} else {
		logx.Debug("stale record session cleanup skipped", "record_session_id", recordSessionID, "expected_gen", myGen, "current_gen", s.recordGen)
	}
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
	started, replay, err := s.recorderOps().BeginCapture(session)
	if err != nil {
		return false, err
	}
	if started && replay && emit != nil {
		s.mu.RLock()
		recordSessionID := s.recordSessionID
		browserSessionID := s.browserSessionID
		targetPath := s.recordTargetPath
		myGen := s.recordGen
		s.mu.RUnlock()
		s.emitLiveRecordedSteps(session, s.guardedRecordEmit(myGen, targetPath, emit), targetPath, recordSessionID, browserSessionID)
	}
	return started, nil
}

func (s *Service) RecordBaseline(req BaselineRecordRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	output := strings.TrimSpace(req.Output)
	if output == "" {
		output = filepath.Join(path, "recorded.feature")
	} else {
		confined, err := paths.ConfineToProjectRoot(path, output)
		if err != nil {
			return RunResult{Error: err.Error()}
		}
		output = confined
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
	out, err := s.cliRunner().Record(args)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

type BrowserSessionDTO struct {
	BrowserOpen      bool   `json:"browserOpen"`
	Recording        bool   `json:"recording"`
	Paused           bool   `json:"paused"`
	StepCount        int    `json:"stepCount"`
	BrowserSessionID string `json:"browserSessionId,omitempty"`
}

func (s *Service) PollBrowserSession() BrowserSessionDTO {
	s.mu.RLock()
	session := s.liveSession
	browserSessionID := s.browserSessionID
	s.mu.RUnlock()
	return s.recorderOps().PollBrowserSession(session, browserSessionID)
}

func (s *Service) PauseRecording() {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	s.recorderOps().PauseRecording(session)
}

func (s *Service) ResumeRecording() {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	s.recorderOps().ResumeRecording(session)
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

func (s *Service) CurrentRecordSessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.recordSessionID
}

func (s *Service) CurrentBrowserSessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.browserSessionID
}

func (s *Service) CurrentRecordTargetPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.recordTargetPath
}

func (s *Service) LastClosedBrowserSessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastClosedBrowserSessionID
}

func (s *Service) LastClosedRecordSessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastClosedRecordSessionID
}

func (s *Service) rememberClosedRecordSessionLocked(recordSessionID, browserSessionID string) {
	if recordSessionID == "" && browserSessionID == "" {
		return
	}
	s.lastClosedRecordSessionID = recordSessionID
	s.lastClosedBrowserSessionID = browserSessionID
}

func (s *Service) StopRecordingCapture() (bool, error) {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	stopped, err := s.recorderOps().StopCapture(session)
	if err != nil {
		return false, err
	}
	if stopped {
		s.emitRecordStepReset()
	}
	return stopped, nil
}

func (s *Service) CloseBrowser() {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session != nil && session.TestRunHeld() {
		return
	}
	s.closeBrowserLocked()
}

func (s *Service) closeBrowserForced() {
	s.closeBrowserLocked()
}

func (s *Service) closeBrowserLocked() {
	s.mu.Lock()
	cancel := s.recordCancel
	session := s.liveSession
	hadSession := session != nil || cancel != nil
	recordSessionID := s.recordSessionID
	browserSessionID := s.browserSessionID
	s.recordCancel = nil
	s.liveSession = nil
	s.recordEmit = nil
	s.recordCtx = nil
	s.recordSessionID = ""
	s.browserSessionID = ""
	s.recordTargetPath = ""
	if hadSession {
		s.rememberClosedRecordSessionLocked(recordSessionID, browserSessionID)
		s.recordGen++
	}
	s.mu.Unlock()
	if session != nil {
		session.Clear()
	}
	if cancel != nil {
		cancel()
	}
}

func (s *Service) CancelRecording() {
	s.closeBrowserForced()
}

func (s *Service) FocusBrowser() error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	return s.recorderOps().FocusBrowser(session)
}

func (s *Service) UpdateRecordingOptions(filter, navOnly, hover, headless, scrollBefore bool, hoverMinMs int) error {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	return s.recorderOps().UpdateRecordingOptions(session, filter, navOnly, hover, headless, scrollBefore, hoverMinMs)
}

func (s *Service) UndoRecordedStep() bool {
	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	index, ok := s.recorderOps().UndoLastRecordedStep(session)
	if !ok {
		return false
	}
	s.emitRecordStepDelete(index)
	return true
}
