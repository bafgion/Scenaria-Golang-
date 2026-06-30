package recorder

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

var ErrRelaunchHeadless = errors.New("recorder: relaunch browser for headless change")

// LiveSession controls an in-progress live recording (pause/resume/focus/undo).
type LiveSession struct {
	paused          atomic.Bool
	headless        atomic.Bool
	relaunch        atomic.Bool
	captureEnabled   atomic.Bool
	captureEver      atomic.Bool
	recorderInjected atomic.Bool
	recMu            sync.RWMutex
	filterImportant bool
	navOnly         bool
	hoverRecord     bool
	scrollBeforeClick bool
	hoverRecordMinMs  int
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
	s.SetRecorderOptions(filterImportant, navOnly, hoverRecord, false, 0)
}

func (s *LiveSession) SetRecorderOptions(filterImportant, navOnly, hoverRecord, scrollBeforeClick bool, hoverRecordMinMs int) {
	s.recMu.Lock()
	s.filterImportant = filterImportant
	s.navOnly = navOnly
	s.hoverRecord = hoverRecord
	s.scrollBeforeClick = scrollBeforeClick
	s.hoverRecordMinMs = hoverRecordMinMs
	s.recMu.Unlock()
}

func (s *LiveSession) RecorderPageConfig() PageRecorderConfig {
	s.recMu.RLock()
	defer s.recMu.RUnlock()
	return PageRecorderConfig{
		FilterImportant:   s.filterImportant,
		NavOnly:           s.navOnly,
		HoverRecord:       s.hoverRecord,
		ScrollBeforeClick: s.scrollBeforeClick,
		HoverRecordMinMs:  s.hoverRecordMinMs,
	}
}

func (s *LiveSession) RecorderFlags() (filterImportant, navOnly, hoverRecord bool) {
	cfg := s.RecorderPageConfig()
	return cfg.FilterImportant, cfg.NavOnly, cfg.HoverRecord
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

func (s *LiveSession) CaptureEnabled() bool {
	return s.captureEnabled.Load()
}

func (s *LiveSession) CaptureEverEnabled() bool {
	return s.captureEver.Load()
}

// InitBrowseMode opens the browser without capturing interactions until BeginCapture.
func (s *LiveSession) InitBrowseMode() {
	s.captureEnabled.Store(false)
	s.captureEver.Store(false)
	s.paused.Store(false)
	s.recorderInjected.Store(false)
}

func (s *LiveSession) InitRecordMode() {
	s.captureEnabled.Store(true)
	s.captureEver.Store(true)
	s.paused.Store(false)
	s.recorderInjected.Store(false)
}

// BeginCapture enables step capture and injects the recorder script on first use.
func (s *LiveSession) BeginCapture() error {
	s.mu.Lock()
	page := s.page
	steps := s.steps
	s.mu.Unlock()
	if page != nil && !s.recorderInjected.Load() {
		if err := page.Context().AddInitScript(playwright.Script{
			Content: playwright.String(selector.RecorderListenersJS),
		}); err != nil {
			return fmt.Errorf("register recorder init script: %w", err)
		}
		if _, err := page.Evaluate(selector.HeuristicsJS); err != nil {
			return fmt.Errorf("inject heuristics: %w", err)
		}
		if _, err := page.Evaluate(selector.RecorderListenersJS); err != nil {
			return fmt.Errorf("inject recorder script: %w", err)
		}
		if appCfg, err := settings.LoadDefaultAppSettings(); err == nil && appCfg != nil {
			_ = selector.ApplySelectorOrder(page, appCfg.SelectorClickStrategies, appCfg.SelectorInputStrategies)
		}
		if err := ApplyPageRecorderConfig(page, s.RecorderPageConfig()); err != nil {
			return fmt.Errorf("configure recorder: %w", err)
		}
		s.recorderInjected.Store(true)
	}
	s.captureEver.Store(true)
	s.captureEnabled.Store(true)
	s.paused.Store(false)
	if page != nil && steps != nil && len(*steps) == 0 {
		if u := strings.TrimSpace(page.URL()); u != "" && u != "about:blank" {
			*steps = append(*steps, RecordedStep{Action: "goto", Value: u})
		}
	}
	return nil
}

func (s *LiveSession) RecordedStepCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.steps == nil {
		return 0
	}
	return len(*s.steps)
}

func (s *LiveSession) BrowserAlive() bool {
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return false
	}
	return !page.IsClosed()
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
	s.recorderInjected.Store(false)
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
	return s.ApplyRecorderOptions(filterImportant, navOnly, hoverRecord, false, 0)
}

func (s *LiveSession) ApplyRecorderOptions(filterImportant, navOnly, hoverRecord, scrollBeforeClick bool, hoverRecordMinMs int) error {
	s.SetRecorderOptions(filterImportant, navOnly, hoverRecord, scrollBeforeClick, hoverRecordMinMs)
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	return ApplyPageRecorderConfig(page, s.RecorderPageConfig())
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

func (s *LiveSession) ExportTestClient(name string) (*settings.TestClient, error) {
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return nil, fmt.Errorf("браузер не открыт")
	}
	return player.CaptureTestClientFromPage(page, name)
}

func (s *LiveSession) PickSelector(ctx context.Context) (string, error) {
	s.mu.Lock()
	page := s.page
	s.mu.Unlock()
	if page == nil {
		return "", fmt.Errorf("браузер не открыт")
	}
	if s.CaptureEnabled() && !s.IsPaused() {
		return "", fmt.Errorf("поставьте запись на паузу")
	}
	return PickSelectorOnPage(ctx, page)
}
