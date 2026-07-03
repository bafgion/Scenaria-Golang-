package player

import "testing"

func TestBrowserSessionCloseIsIdempotent(t *testing.T) {
	session := &browserSession{videoEnabled: true}
	session.close()
	session.close()
	if !session.closed {
		t.Fatal("expected session to be closed")
	}
}

func TestExternalSessionCloseDetachesWithoutClosing(t *testing.T) {
	session := &browserSession{external: true}
	session.closeLocked(true)
	if !session.closed {
		t.Fatal("expected external session to be marked closed")
	}
}
