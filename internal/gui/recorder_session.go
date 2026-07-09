package gui

import (
	"context"
	"fmt"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/recorder"
)

// RecorderSessionSnapshot is a consistent read of recorder session identity fields.
type RecorderSessionSnapshot struct {
	Gen              uint64
	RecordSessionID  string
	BrowserSessionID string
	TargetPath       string
	Emit             func(string, any)
	Session          *recorder.LiveSession
	RecordCtx        context.Context
	IdleSeconds      int
}

// RecorderSessionManager owns the active recorder/browser session state.
type RecorderSessionManager struct {
	mu                         sync.RWMutex
	liveSession                *recorder.LiveSession
	recordCtx                  context.Context
	recordCancel               context.CancelFunc
	recordGen                  uint64
	recordSessionID            string
	browserSessionID           string
	recordTargetPath           string
	recordEmit                 func(string, any)
	lastClosedRecordSessionID  string
	lastClosedBrowserSessionID string
	recordIdleSeconds          int
}

func (m *RecorderSessionManager) HasLiveBrowser() bool {
	if m == nil {
		return false
	}
	m.mu.RLock()
	session := m.liveSession
	m.mu.RUnlock()
	return session != nil && session.BrowserAlive()
}

func (m *RecorderSessionManager) Snapshot() RecorderSessionSnapshot {
	if m == nil {
		return RecorderSessionSnapshot{}
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return RecorderSessionSnapshot{
		Gen:              m.recordGen,
		RecordSessionID:  m.recordSessionID,
		BrowserSessionID: m.browserSessionID,
		TargetPath:       m.recordTargetPath,
		Emit:             m.recordEmit,
		Session:          m.liveSession,
		RecordCtx:        m.recordCtx,
		IdleSeconds:      m.recordIdleSeconds,
	}
}

func (m *RecorderSessionManager) LiveSession() *recorder.LiveSession {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.liveSession
}

func (m *RecorderSessionManager) RecordSessionID() string {
	if m == nil {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recordSessionID
}

func (m *RecorderSessionManager) BrowserSessionID() string {
	if m == nil {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.browserSessionID
}

func (m *RecorderSessionManager) TargetPath() string {
	if m == nil {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recordTargetPath
}

func (m *RecorderSessionManager) LastClosedBrowserSessionID() string {
	if m == nil {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastClosedBrowserSessionID
}

func (m *RecorderSessionManager) LastClosedRecordSessionID() string {
	if m == nil {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastClosedRecordSessionID
}

func (m *RecorderSessionManager) GuardedEmit(gen uint64, targetPath string, emit func(string, any)) func(string, any) {
	if emit == nil {
		return nil
	}
	return func(name string, payload any) {
		if m == nil {
			return
		}
		m.mu.RLock()
		active := m.recordGen == gen
		m.mu.RUnlock()
		if !active {
			logx.Debug("stale record event skipped", "event", name, "expected_gen", gen)
			return
		}
		if payloadMap, ok := payload.(map[string]any); ok && payloadMap != nil {
			if targetPath != "" {
				if _, has := payloadMap["targetPath"]; !has {
					payloadMap["targetPath"] = targetPath
				}
			}
		}
		emit(name, payload)
	}
}

type BeginRecorderSessionInput struct {
	TargetPath   string
	Session      *recorder.LiveSession
	RecordCtx    context.Context
	RecordCancel context.CancelFunc
	Emit         func(string, any)
	IdleSeconds  int
}

type BeginRecorderSessionResult struct {
	Gen              uint64
	RecordSessionID  string
	BrowserSessionID string
}

func (m *RecorderSessionManager) BeginSession(input BeginRecorderSessionInput) BeginRecorderSessionResult {
	if m == nil {
		return BeginRecorderSessionResult{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.recordCancel != nil {
		m.recordCancel()
	}
	m.recordGen++
	myGen := m.recordGen
	recordSessionID := fmt.Sprintf("record-%d", myGen)
	browserSessionID := fmt.Sprintf("browser-%d", myGen)
	m.liveSession = input.Session
	m.recordSessionID = recordSessionID
	m.browserSessionID = browserSessionID
	m.recordTargetPath = input.TargetPath
	m.recordCtx = input.RecordCtx
	m.recordCancel = input.RecordCancel
	m.recordEmit = input.Emit
	m.recordIdleSeconds = input.IdleSeconds
	return BeginRecorderSessionResult{
		Gen:              myGen,
		RecordSessionID:  recordSessionID,
		BrowserSessionID: browserSessionID,
	}
}

func (m *RecorderSessionManager) ClearEmitIfGen(gen uint64) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.recordGen == gen {
		m.recordEmit = nil
	} else {
		logx.Debug("stale record emit cleanup skipped", "expected_gen", gen, "current_gen", m.recordGen)
	}
}

func (m *RecorderSessionManager) FinishSession(gen uint64, recordSessionID, browserSessionID string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.recordGen != gen {
		logx.Debug("stale record session cleanup skipped", "record_session_id", recordSessionID, "expected_gen", gen, "current_gen", m.recordGen)
		return
	}
	m.rememberClosedLocked(recordSessionID, browserSessionID)
	m.recordCancel = nil
	m.liveSession = nil
	m.recordCtx = nil
	m.recordSessionID = ""
	m.browserSessionID = ""
	m.recordTargetPath = ""
}

func (m *RecorderSessionManager) SetEmit(emit func(string, any)) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.recordEmit = emit
	m.mu.Unlock()
}

func (m *RecorderSessionManager) Generation() uint64 {
	if m == nil {
		return 0
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recordGen
}

func (m *RecorderSessionManager) CancelContextOnProjectSwitch() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.recordCancel != nil {
		m.recordCancel()
		m.recordCancel = nil
		m.recordCtx = nil
	}
}

func (m *RecorderSessionManager) CloseBrowser(force bool) {
	if m == nil {
		return
	}
	m.mu.RLock()
	session := m.liveSession
	m.mu.RUnlock()
	if !force && session != nil && session.TestRunHeld() {
		return
	}
	m.closeLocked()
}

func (m *RecorderSessionManager) closeLocked() {
	m.mu.Lock()
	cancel := m.recordCancel
	session := m.liveSession
	hadSession := session != nil || cancel != nil
	recordSessionID := m.recordSessionID
	browserSessionID := m.browserSessionID
	m.recordCancel = nil
	m.liveSession = nil
	m.recordEmit = nil
	m.recordCtx = nil
	m.recordSessionID = ""
	m.browserSessionID = ""
	m.recordTargetPath = ""
	if hadSession {
		m.rememberClosedLocked(recordSessionID, browserSessionID)
		m.recordGen++
	}
	m.mu.Unlock()
	if session != nil {
		session.Clear()
	}
	if cancel != nil {
		cancel()
	}
}

func (m *RecorderSessionManager) rememberClosedLocked(recordSessionID, browserSessionID string) {
	if recordSessionID == "" && browserSessionID == "" {
		return
	}
	m.lastClosedRecordSessionID = recordSessionID
	m.lastClosedBrowserSessionID = browserSessionID
}
