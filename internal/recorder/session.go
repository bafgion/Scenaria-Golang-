package recorder

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	playwright "github.com/mxschmitt/playwright-go"
)

var ErrRelaunchHeadless = errors.New("recorder: relaunch browser for headless change")

// LiveSession controls an in-progress live recording (pause/resume/focus/undo).
type LiveSession struct {
	paused   atomic.Bool
	headless atomic.Bool
	relaunch atomic.Bool
	recMu    sync.RWMutex
	filterImportant bool
	navOnly         bool
	hoverRecord     bool
	resumeURL       string
	mu       sync.Mutex
	page     playwright.Page
	steps    *[]RecordedStep
}

func NewLiveSession() *LiveSession {
	return &LiveSession{}
}

func (s *LiveSession) InitHeadless(headless bool) {
	s.headless.Store(headless)
}

func (s *LiveSession) SetRecorderFlags(filterImportant, navOnly, hoverRecord bool) {
	s.recMu.Lock()
	s.filterImportant = filterImportant
	s.navOnly = navOnly
	s.hoverRecord = hoverRecord
	s.recMu.Unlock()
}

func (s *LiveSession) SetResumeURL(url string) {
	s.recMu.Lock()
	s.resumeURL = strings.TrimSpace(url)
	s.recMu.Unlock()
}

func (s *LiveSession) ResumeURL(fallback string) string {
	s.recMu.RLock()
	defer s.recMu.RUnlock()
	if s.resumeURL != "" {
		return s.resumeURL
	}
	return fallback
}

func (s *LiveSession) RecorderFlags() (bool, bool, bool) {
	s.recMu.RLock()
	defer s.recMu.RUnlock()
	return s.filterImportant, s.navOnly, s.hoverRecord
}

func (s *LiveSession) Headless() bool {
	return s.headless.Load()
}

func (s *LiveSession) RequestHeadless(headless bool) {
	if s.headless.Load() == headless {
		return
	}
	s.headless.Store(headless)
	s.relaunch.Store(true)
}

func (s *LiveSession) RelaunchPending() bool {
	return s.relaunch.Load()
}

func (s *LiveSession) ClearRelaunch() {
	s.relaunch.Store(false)
}

func (s *LiveSession) Pause()  { s.paused.Store(true) }
func (s *LiveSession) Resume() { s.paused.Store(false) }
func (s *LiveSession) IsPaused() bool {
	return s.paused.Load()
}

func (s *LiveSession) Bind(page playwright.Page, steps *[]RecordedStep) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.page = page
	s.steps = steps
}

func (s *LiveSession) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.page = nil
	s.steps = nil
}

func (s *LiveSession) FocusBrowser() error {
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return fmt.Errorf("браузер не открыт")
	}
	return page.BringToFront()
}

func (s *LiveSession) ApplyRecorderConfig(filterImportant, navOnly, hoverRecord bool) error {
	s.SetRecorderFlags(filterImportant, navOnly, hoverRecord)
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return nil
	}
	script := fmt.Sprintf(`() => {
		if (!window.__scenariaRecorder) return;
		window.__scenariaRecorder.filterImportant = %v;
		window.__scenariaRecorder.navOnly = %v;
		window.__scenariaRecorder.hoverRecord = %v;
	}`, filterImportant, navOnly, hoverRecord)
	_, err := page.Evaluate(script)
	return err
}

func (s *LiveSession) UndoLastStep() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.steps == nil || len(*s.steps) <= 1 {
		return false
	}
	*s.steps = (*s.steps)[:len(*s.steps)-1]
	if s.page != nil {
		_, _ = s.page.Evaluate(`() => {
			const r = window.__scenariaRecorder;
			if (r && r.events.length) r.events.pop();
		}`, nil)
	}
	return true
}

func (s *LiveSession) PickSelector(ctx context.Context) (string, error) {
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return "", fmt.Errorf("браузер не открыт")
	}
	if !s.IsPaused() {
		return "", fmt.Errorf("поставьте запись на паузу")
	}
	return PickSelectorOnPage(ctx, page)
}
