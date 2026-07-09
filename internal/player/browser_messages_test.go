package player

import (
	"errors"
	"testing"
)

func TestUserFacingBrowserError(t *testing.T) {
	err := errors.New("browser session is closed")
	if got := UserFacingBrowserError(err); got != MsgBrowserClosed {
		t.Fatalf("got %q", got)
	}
	if got := UserFacingBrowserError(errors.New("element not found")); got != "element not found" {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeScenarioMessage(t *testing.T) {
	if got := NormalizeScenarioMessage("browser page is not available"); got != MsgBrowserClosed {
		t.Fatalf("got %q", got)
	}
}
