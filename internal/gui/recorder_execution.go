package gui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

// RecorderExecutionHost supplies project-scoped dependencies for live record orchestration.
type RecorderExecutionHost struct {
	ProjectPath        func() string
	LoadAppSettings    func() (*settings.AppSettings, error)
	SaveAppSettings    func(*settings.AppSettings) error
	WithActivePlaywright func(func())
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

func (s *RecorderService) liveRecordCallbacks(
	emit func(string, any),
	browseOnly bool,
	output, targetPath, recordSessionID, browserSessionID string,
	idleSeconds int,
) recorder.LiveCallbacks {
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
				if reason == "idle" && idleSeconds > 0 {
					payload["idleSeconds"] = idleSeconds
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

func (s *RecorderService) emitLiveRecordedSteps(
	session *recorder.LiveSession,
	emit func(string, any),
	targetPath, recordSessionID, browserSessionID string,
) {
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

func (s *RecorderService) emitRecordStepReset() {
	snap := s.session.Snapshot()
	if snap.Emit == nil {
		return
	}
	emit := s.session.GuardedEmit(snap.Gen, snap.TargetPath, snap.Emit)
	emit("record-step", map[string]any{
		"op":               string(recorder.RecordStepReset),
		"recordSessionId":  snap.RecordSessionID,
		"browserSessionId": snap.BrowserSessionID,
		"targetPath":       snap.TargetPath,
	})
}

func (s *RecorderService) emitRecordStepDelete(index int) {
	snap := s.session.Snapshot()
	if snap.Emit == nil || index < 0 {
		return
	}
	emit := s.session.GuardedEmit(snap.Gen, snap.TargetPath, snap.Emit)
	emit("record-step", map[string]any{
		"op":               string(recorder.RecordStepDelete),
		"index":            index,
		"recordSessionId":  snap.RecordSessionID,
		"browserSessionId": snap.BrowserSessionID,
		"targetPath":       snap.TargetPath,
	})
}

func (s *RecorderService) startCaptureOnExistingSession(
	req RecordRequest,
	emit func(string, any),
	host RecorderExecutionHost,
) RunResult {
	snap := s.session.Snapshot()
	session := snap.Session
	if session == nil || !session.BrowserAlive() {
		return RunResult{Error: "браузер не открыт"}
	}
	emit = s.session.GuardedEmit(snap.Gen, snap.TargetPath, emit)
	appCfg, err := host.LoadAppSettings()
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
				"recordSessionId":  snap.RecordSessionID,
				"browserSessionId": snap.BrowserSessionID,
				"targetPath":       snap.TargetPath,
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
			"recordSessionId":  snap.RecordSessionID,
			"browserSessionId": snap.BrowserSessionID,
			"targetPath":       snap.TargetPath,
		})
		if replay {
			s.emitLiveRecordedSteps(session, emit, snap.TargetPath, snap.RecordSessionID, snap.BrowserSessionID)
		}
	}
	return RunResult{}
}

func (s *RecorderService) RecordLive(req RecordRequest, emit func(string, any), host RecorderExecutionHost) RunResult {
	if s == nil {
		return RunResult{Error: "recorder service is not configured"}
	}
	if host.WithActivePlaywright == nil {
		return RunResult{Error: "recorder host is not configured"}
	}
	var result RunResult
	host.WithActivePlaywright(func() {
		result = s.recordLive(req, emit, host)
	})
	return result
}

func (s *RecorderService) recordLive(req RecordRequest, emit func(string, any), host RecorderExecutionHost) RunResult {
	path := host.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	if strings.TrimSpace(req.URL) != "" && !strings.HasPrefix(strings.TrimSpace(req.URL), "http") {
		return RunResult{Error: "укажите корректный URL (https://…) или оставьте поле пустым"}
	}

	existing := s.LiveSession()
	if existing != nil && existing.BrowserAlive() && !req.BrowseOnly {
		return s.startCaptureOnExistingSession(req, emit, host)
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

	appCfg, err := host.LoadAppSettings()
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	cleanURL := strings.TrimSpace(req.URL)
	if cleanURL != "" {
		cleanURL = httpauth.ApplyURLCredentials(cleanURL, appCfg)
		if err := host.SaveAppSettings(appCfg); err != nil {
			return RunResult{Error: err.Error()}
		}
	}
	httpCreds := httpauth.PlaywrightHTTPCredentials(cleanURL, appCfg)

	session := recorder.NewLiveSession()
	targetPath := resolveRecordTargetPath(path, req)
	ctx, cancel := context.WithCancel(context.Background())
	begin := s.session.BeginSession(BeginRecorderSessionInput{
		TargetPath:   targetPath,
		Session:      session,
		RecordCtx:    ctx,
		RecordCancel: cancel,
		Emit:         emit,
		IdleSeconds:  idle,
	})
	myGen := begin.Gen
	recordSessionID := begin.RecordSessionID
	browserSessionID := begin.BrowserSessionID

	emit = s.session.GuardedEmit(myGen, targetPath, emit)

	defer s.session.ClearEmitIfGen(myGen)

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
		Callbacks:         s.liveRecordCallbacks(emit, req.BrowseOnly, output, targetPath, recordSessionID, browserSessionID, idle),
	})

	s.session.FinishSession(myGen, recordSessionID, browserSessionID)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return RunResult{Output: "Браузер закрыт."}
		}
		return RunResult{Error: fmt.Errorf("record: %w", err).Error()}
	}
	return RunResult{Output: fmt.Sprintf("Запись сохранена: %s\n", output)}
}

func (s *RecorderService) BeginRecordingCaptureWithReplay() (bool, error) {
	snap := s.session.Snapshot()
	started, replay, err := s.BeginCapture(snap.Session)
	if err != nil {
		return false, err
	}
	if started && replay && snap.Emit != nil {
		s.emitLiveRecordedSteps(
			snap.Session,
			s.session.GuardedEmit(snap.Gen, snap.TargetPath, snap.Emit),
			snap.TargetPath,
			snap.RecordSessionID,
			snap.BrowserSessionID,
		)
	}
	return started, nil
}

func (s *RecorderService) StopRecordingCaptureWithReset() (bool, error) {
	stopped, err := s.StopCapture(s.LiveSession())
	if err != nil {
		return false, err
	}
	if stopped {
		s.emitRecordStepReset()
	}
	return stopped, nil
}

func (s *RecorderService) UndoRecordedStepWithEmit() bool {
	index, ok := s.UndoLastRecordedStep(s.LiveSession())
	if !ok {
		return false
	}
	s.emitRecordStepDelete(index)
	return true
}
