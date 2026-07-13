package recorder

import (
	"fmt"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

// ToolbarActionResult describes side effects from a browser-toolbar command.
type ToolbarActionResult struct {
	CloseBrowser bool
}

func (s *LiveSession) pollStateForPage(page playwright.Page) *recorderPollState {
	s.pollMu.Lock()
	defer s.pollMu.Unlock()
	if s.pollState == nil {
		url := ""
		if page != nil {
			url = page.URL()
		}
		state := newRecorderPollState(url, s.RecordURLWaitAfterClick())
		s.pollState = &state
	}
	return s.pollState
}

// FlushPendingRecorderEvents drains browser-side recorder events into the Go step buffer.
func FlushPendingRecorderEvents(session *LiveSession, page playwright.Page, notify StepNotifier) error {
	if session == nil || page == nil || !session.CaptureEnabled() {
		return nil
	}
	events, err := drainRecorderEvents(page)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	state := session.pollStateForPage(page)
	session.applyRecorderPollBatch(state, events, page.URL(), time.Now(), notify)
	return nil
}

// CompleteCaptureStop drains pending events, emits a snapshot, disables capture, and resets the segment buffer.
func CompleteCaptureStop(
	session *LiveSession,
	page playwright.Page,
	notify StepNotifier,
	reason string,
	onStop func(string),
) bool {
	if session == nil || !session.CaptureEnabled() {
		return false
	}
	_ = FlushPendingRecorderEvents(session, page, notify)
	shouldResetSegment := notify != nil
	if notify != nil {
		session.mu.Lock()
		var steps []RecordedStep
		if session.steps != nil {
			steps = append([]RecordedStep(nil), (*session.steps)...)
		}
		session.mu.Unlock()
		notifySnapshot(notify, steps)
	}
	session.StopCapturePreserveBuffer()
	if shouldResetSegment {
		session.ResetCaptureSegment()
	}
	if onStop != nil {
		onStop(reason)
	}
	return true
}

// ProcessToolbarAction handles one browser-toolbar command.
func ProcessToolbarAction(
	action string,
	session *LiveSession,
	page playwright.Page,
	opts LiveOptions,
	notify StepNotifier,
	recorded *[]RecordedStep,
) ToolbarActionResult {
	if session == nil {
		return ToolbarActionResult{}
	}
	switch action {
	case "stop":
		if session.CaptureEnabled() {
			CompleteCaptureStop(session, page, notify, "manual", opts.Callbacks.OnCaptureStop)
		}
		return ToolbarActionResult{}
	case "pause":
		if session.CaptureEnabled() {
			session.Pause()
		}
	case "resume":
		if session.CaptureEnabled() && session.IsPaused() {
			session.Resume()
		}
	case "record":
		if !session.CaptureEnabled() {
			replay := ShouldSyncRecordedStepsOnCaptureStart(session)
			if err := session.BeginCapture(); err != nil {
				setBrowserToolbarError(page, err.Error())
				return ToolbarActionResult{}
			}
			clearBrowserToolbarError(page)
			if opts.Callbacks.OnCaptureStart != nil {
				opts.Callbacks.OnCaptureStart(!replay)
			}
			if replay && notify != nil && recorded != nil {
				notifySnapshot(notify, *recorded)
			}
		}
	case "picker":
		if opts.Callbacks.OnPickerRequest != nil {
			opts.Callbacks.OnPickerRequest()
		}
	default:
		return ToolbarActionResult{}
	}
	return ToolbarActionResult{}
}

func setBrowserToolbarError(page playwright.Page, message string) {
	if page == nil || message == "" {
		return
	}
	script := fmt.Sprintf(`() => {
		if (!window.__scenariaToolbar) return;
		window.__scenariaToolbar.setState({ phase: 'error', error: %q });
	}`, message)
	evaluateRecorderCleanup(page, script)
}

func clearBrowserToolbarError(page playwright.Page) {
	evaluateRecorderCleanup(page, `() => {
		if (!window.__scenariaToolbar) return;
		window.__scenariaToolbar.setState({ phase: '', error: '' });
	}`)
}

func browserToolbarPhase(session *LiveSession, browseOnly bool) string {
	if session == nil {
		return "idle"
	}
	if session.CaptureEnabled() {
		if session.IsPaused() {
			return "paused"
		}
		return "recording"
	}
	if browseOnly {
		return "idle"
	}
	return "idle"
}
