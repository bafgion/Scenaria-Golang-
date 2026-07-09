package gui

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

// RunService owns run history reads and the active run session lifecycle.
type RunService struct {
	projectPath projectPathProvider
	mu          sync.RWMutex
	session     *RunSession
	runCtx      context.Context
	runCancel   context.CancelFunc
	runGen      uint64
}

func NewRunService(projectPath projectPathProvider) *RunService {
	return &RunService{projectPath: projectPath}
}

func (s *RunService) ListRunResults(limit int) ([]RunResultEntry, error) {
	if s == nil || s.projectPath == nil {
		return []RunResultEntry{}, nil
	}
	path := s.projectPath()
	if path == "" {
		return []RunResultEntry{}, nil
	}
	store, err := runstatus.Open(path)
	if err != nil {
		return nil, err
	}
	entries, err := store.List(limit)
	if err != nil {
		return nil, err
	}
	out := make([]RunResultEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, RunResultEntry{
			Path:       e.Path,
			Success:    e.Success,
			Message:    e.Message,
			Runner:     e.Runner,
			At:         e.At,
			FailedStep: e.FailedStep,
		})
	}
	return out, nil
}

func (s *RunService) FlakyMetrics(historyLimit int) (FlakyMetricsDTO, error) {
	if s == nil || s.projectPath == nil {
		return FlakyMetricsDTO{}, nil
	}
	path := s.projectPath()
	if path == "" {
		return FlakyMetricsDTO{}, nil
	}
	store, err := runstatus.Open(path)
	if err != nil {
		return FlakyMetricsDTO{}, err
	}
	if historyLimit <= 0 {
		historyLimit = 200
	}
	entries, err := store.List(historyLimit)
	if err != nil {
		return FlakyMetricsDTO{}, err
	}
	scenarios, steps := runstatus.FlakyStats(entries)
	out := FlakyMetricsDTO{
		Scenarios: make([]FlakyScenarioDTO, 0, len(scenarios)),
		Steps:     make([]FlakyStepDTO, 0, len(steps)),
	}
	for _, item := range scenarios {
		out.Scenarios = append(out.Scenarios, FlakyScenarioDTO{
			Path:       item.Path,
			Failures:   item.Failures,
			Passes:     item.Passes,
			Total:      item.Total,
			Flaky:      item.Flaky,
			LastFailed: item.LastFailed,
		})
	}
	for _, item := range steps {
		out.Steps = append(out.Steps, FlakyStepDTO{
			Path:       item.Path,
			Step:       item.Step,
			Failures:   item.Failures,
			LastFailed: item.LastFailed,
		})
	}
	return out, nil
}

type RunBeginResult struct {
	Context       context.Context
	Cancel        context.CancelFunc
	RunID         string
	Gen           uint64
	TempResources []string
}

func (s *RunService) CurrentSession() *RunSession {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.session == nil {
		return nil
	}
	copy := *s.session
	copy.RequestSnapshot = cloneRunRequest(copy.RequestSnapshot)
	copy.TempResources = append([]string(nil), copy.TempResources...)
	return &copy
}

func (s *RunService) TryBegin(projectVersion uint64, req RunRequest, tempResources []string, timeout time.Duration) (RunBeginResult, error) {
	if s == nil {
		return RunBeginResult{}, fmt.Errorf("run service is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runCancel != nil && s.runCtx != nil && s.runCtx.Err() == nil {
		cancel()
		return RunBeginResult{}, fmt.Errorf("запуск уже выполняется — нажмите «Стоп» и попробуйте снова")
	}
	s.runGen++
	myGen := s.runGen
	s.runCtx = ctx
	s.runCancel = cancel
	s.session = &RunSession{
		RunID:           runID,
		ProjectVersion:  projectVersion,
		RequestSnapshot: cloneRunRequest(req),
		Context:         ctx,
		TempResources:   append([]string(nil), tempResources...),
	}
	return RunBeginResult{
		Context:       ctx,
		Cancel:        cancel,
		RunID:         runID,
		Gen:           myGen,
		TempResources: append([]string(nil), s.session.TempResources...),
	}, nil
}

func (s *RunService) Finish(gen uint64) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runGen != gen {
		return
	}
	s.runCtx = nil
	s.runCancel = nil
	s.session = nil
}

func (s *RunService) Cancel() {
	if s == nil {
		return
	}
	s.mu.Lock()
	cancel := s.runCancel
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (s *RunService) CancelIfActive() {
	if s == nil {
		return
	}
	s.mu.Lock()
	if s.runCancel != nil {
		s.runCancel()
		s.runCancel = nil
		s.runCtx = nil
		s.session = nil
	}
	s.mu.Unlock()
}
