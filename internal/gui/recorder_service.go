package gui

import (
	"context"
	"errors"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/recorder"
)

// RecorderService owns recorder session operations that can be delegated from gui.Service facade.
type RecorderService struct{}

func NewRecorderService() *RecorderService {
	return &RecorderService{}
}

func (s *RecorderService) PollBrowserSession(session *recorder.LiveSession, browserSessionID string) BrowserSessionDTO {
	if session == nil || !session.BrowserAlive() {
		return BrowserSessionDTO{}
	}
	return BrowserSessionDTO{
		BrowserOpen:      true,
		Recording:        session.CaptureEnabled(),
		Paused:           session.IsPaused(),
		StepCount:        session.RecordedStepCount(),
		BrowserSessionID: browserSessionID,
	}
}

func (s *RecorderService) BeginCapture(session *recorder.LiveSession) (started bool, replay bool, err error) {
	if session == nil {
		return false, false, fmt.Errorf("браузер не открыт")
	}
	if session.CaptureEnabled() {
		return false, false, nil
	}
	replay = recorder.ShouldSyncRecordedStepsOnCaptureStart(session)
	if err := session.BeginCapture(); err != nil {
		return false, false, err
	}
	return true, replay, nil
}

func (s *RecorderService) StopCapture(session *recorder.LiveSession) (bool, error) {
	if session == nil {
		return false, fmt.Errorf("браузер не открыт")
	}
	if !session.CaptureEnabled() {
		return false, nil
	}
	session.EndCapture()
	return true, nil
}

func (s *RecorderService) UndoLastRecordedStep(session *recorder.LiveSession) (removedIndex int, ok bool) {
	if session == nil {
		return -1, false
	}
	index := session.RecordedStepCount() - 1
	if index < 0 {
		return -1, false
	}
	if !session.UndoLastStep() {
		return -1, false
	}
	return index, true
}

func (s *RecorderService) PauseRecording(session *recorder.LiveSession) {
	if session != nil {
		session.Pause()
	}
}

func (s *RecorderService) ResumeRecording(session *recorder.LiveSession) {
	if session != nil {
		session.Resume()
	}
}

func (s *RecorderService) FocusBrowser(session *recorder.LiveSession) error {
	if session == nil {
		return fmt.Errorf("браузер не открыт")
	}
	return session.FocusBrowser()
}

func (s *RecorderService) UpdateRecordingOptions(session *recorder.LiveSession, filter, navOnly, hover, headless, scrollBefore bool, hoverMinMs int) error {
	if session == nil {
		return fmt.Errorf("запись не активна")
	}
	_ = session.ApplyRecorderOptions(filter, navOnly, hover, scrollBefore, hoverMinMs)
	session.RequestHeadless(headless)
	return nil
}

func (s *RecorderService) PickSelector(session *recorder.LiveSession, ctx context.Context) PickSelectorResult {
	if session == nil {
		return PickSelectorResult{Error: "браузер не открыт"}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	selector, err := session.PickSelector(ctx)
	if err != nil {
		if errors.Is(err, recorder.ErrPickerCancelled) {
			return PickSelectorResult{Error: "отменено"}
		}
		return PickSelectorResult{Error: err.Error()}
	}
	return PickSelectorResult{Selector: selector}
}
