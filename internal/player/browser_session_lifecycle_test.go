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
