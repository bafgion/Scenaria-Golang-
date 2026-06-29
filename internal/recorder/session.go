package recorder

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	playwright "github.com/mxschmitt/playwright-go"
)

// LiveSession controls an in-progress live recording (pause/resume/focus/undo).
type LiveSession struct {
	paused atomic.Bool
	mu     sync.Mutex
	page   playwright.Page
	steps  *[]RecordedStep
}

func NewLiveSession() *LiveSession {
	return &LiveSession{}
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
