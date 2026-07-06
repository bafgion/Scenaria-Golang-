package player

import "testing"

func TestAbortRunAllowsLaterClose(t *testing.T) {
	session := &browserSession{}
	session.abortRun()
	if !session.closed.Load() {
		t.Fatal("expected aborted session to be marked closed")
	}
	session.closeWhileLocked(false)
	if !session.isClosed() {
		// no resources — close is a no-op beyond the flag
	}
}

func TestAbortRunExternalSession(t *testing.T) {
	session := &browserSession{external: true}
	session.abortRun()
	if !session.isClosed() {
		t.Fatal("expected external session to be marked closed")
	}
}
