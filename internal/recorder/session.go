package recorder

import "sync/atomic"

// LiveSession controls an in-progress live recording (pause/resume).
type LiveSession struct {
	paused atomic.Bool
}

func NewLiveSession() *LiveSession {
	return &LiveSession{}
}

func (s *LiveSession) Pause()  { s.paused.Store(true) }
func (s *LiveSession) Resume() { s.paused.Store(false) }
func (s *LiveSession) IsPaused() bool {
	return s.paused.Load()
}

func (s *LiveSession) waitIfPaused() {
	for s.IsPaused() {
		// busy-wait with short sleep in caller loop
	}
}
