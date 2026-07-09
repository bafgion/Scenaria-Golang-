package player

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBrowserSessionCloseIsIdempotent(t *testing.T) {
	session := &browserSession{videoEnabled: true}
	session.close()
	session.close()
	if !session.isClosed() {
		t.Fatal("expected session to be closed")
	}
}

func TestExternalSessionCloseDetachesWithoutClosing(t *testing.T) {
	session := &browserSession{external: true}
	session.closeLocked(true)
	if !session.isClosed() {
		t.Fatal("expected external session to be marked closed")
	}
}

func TestResetForScenarioClearsBrowserStorage(t *testing.T) {
	ctx := context.Background()
	pw, stopPW, err := startPlaywright(ctx)
	if err != nil {
		t.Skip("playwright not available:", err)
	}
	defer stopPW()
	session, err := newBrowserSession(pw, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true})
	if err != nil {
		t.Skip("browser not available:", err)
	}
	defer session.close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><body>ok</body></html>`))
	}))
	defer srv.Close()

	page, err := session.currentPage()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := page.Goto(srv.URL); err != nil {
		t.Fatal(err)
	}
	if _, err := page.Evaluate(`() => {
		localStorage.setItem('scenaria-reset-test', 'leaked')
		sessionStorage.setItem('scenaria-reset-test', 'leaked')
		document.cookie = 'scenaria_reset_test=leaked; path=/'
	}`); err != nil {
		t.Fatal(err)
	}

	if err := session.resetForScenario(); err != nil {
		t.Fatal(err)
	}
	page, err = session.currentPage()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := page.Goto(srv.URL); err != nil {
		t.Fatal(err)
	}
	state, err := page.Evaluate(`() => ({
		local: localStorage.getItem('scenaria-reset-test') || '',
		session: sessionStorage.getItem('scenaria-reset-test') || '',
		cookie: document.cookie || '',
	})`)
	if err != nil {
		t.Fatal(err)
	}
	got := state.(map[string]any)
	if got["local"] != "" || got["session"] != "" || strings.Contains(got["cookie"].(string), "scenaria_reset_test") {
		t.Fatalf("browser storage leaked after reset: %+v", got)
	}
}
