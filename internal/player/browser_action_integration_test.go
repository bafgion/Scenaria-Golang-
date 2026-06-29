//go:build integration

package player

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

func newTestBrowserSession(t *testing.T) (*browserSession, func()) {
	t.Helper()
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("run playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		_ = pw.Stop()
		t.Fatalf("launch: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		_ = browser.Close()
		_ = pw.Stop()
		t.Fatalf("new page: %v", err)
	}
	session := &browserSession{
		browser: browser,
		page:    page,
	}
	cleanup := func() {
		session.close()
		_ = browser.Close()
		_ = pw.Stop()
	}
	return session, cleanup
}

func TestExecuteGotoUsesDomContentLoaded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><h1>page</h1></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if !strings.Contains(session.page.URL(), srv.URL) {
		t.Fatalf("expected URL %q, got %q", srv.URL, session.page.URL())
	}
}

func TestExecutePressKeyboard(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><input id="q" /></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := session.page.Locator("#q").Focus(); err != nil {
		t.Fatalf("focus: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "press",
		Value1: "Enter",
	}, "", nil); err != nil {
		t.Fatalf("press: %v", err)
	}
}

func TestExecuteAssertHiddenVisibleFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><div class="modal">open</div></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "assert-hidden",
		Value1: ".modal",
	}, "", nil)
	if err == nil || !strings.Contains(err.Error(), "hidden") {
		t.Fatalf("expected hidden assertion error, got: %v", err)
	}
}

func TestExecuteAssertHiddenMissingElementPasses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "assert-hidden",
		Value1: ".modal",
	}, "", nil); err != nil {
		t.Fatalf("missing element should pass hidden check: %v", err)
	}
}

func TestExecuteScrollTo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body style="height:2000px"><div id="footer">end</div></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "scroll-to",
		Value1: "#footer",
	}, "", nil); err != nil {
		t.Fatalf("scroll-to: %v", err)
	}
}

func TestExecuteFillGenerated(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><input id="last" /></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	runCtx := NewRunContext(map[string]string{}, 42, "")
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "fill-generated",
		Value1: "last_name",
		Value2: "#last",
	}, "", runCtx); err != nil {
		t.Fatalf("fill-generated: %v", err)
	}
	value, err := session.page.Locator("#last").InputValue()
	if err != nil {
		t.Fatalf("read value: %v", err)
	}
	if strings.TrimSpace(value) == "" {
		t.Fatal("expected generated last name in input")
	}
}

func TestExecuteClick(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><button id="go">Go</button></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "click",
		Value1: "#go",
	}, "", nil); err != nil {
		t.Fatalf("click: %v", err)
	}
}

func TestExecuteFill(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><input id="name" /></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "fill",
		Value1: "alice",
		Value2: "#name",
	}, "", nil); err != nil {
		t.Fatalf("fill: %v", err)
	}
	value, err := session.page.Locator("#name").InputValue()
	if err != nil {
		t.Fatalf("read value: %v", err)
	}
	if value != "alice" {
		t.Fatalf("expected alice, got %q", value)
	}
}

func TestExecuteAssertVisible(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><h1 id="title">Hello</h1></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "assert-visible",
		Value1: "#title",
	}, "", nil); err != nil {
		t.Fatalf("assert-visible: %v", err)
	}
}

func TestExecuteAssertText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><p id="msg">Example Domain</p></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "assert-text",
		Value1: "Example Domain",
		Value2: "#msg",
	}, "", nil); err != nil {
		t.Fatalf("assert-text: %v", err)
	}
}

func TestExecuteWait(t *testing.T) {
	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	start := time.Now()
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "wait",
		Value1: "200ms",
	}, "", nil); err != nil {
		t.Fatalf("wait: %v", err)
	}
	if time.Since(start) < 150*time.Millisecond {
		t.Fatal("expected wait to pause execution")
	}
}

func TestExecuteHover(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<nav id="menu"><a id="item" href="#">Item</a></nav>
</body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "hover",
		Value1: "#menu",
	}, "", nil); err != nil {
		t.Fatalf("hover: %v", err)
	}
}

func TestExecuteSelect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<select id="lang"><option value="ru">RU</option><option value="en">EN</option></select>
</body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "select",
		Value1: "en",
		Value2: "#lang",
	}, "", nil); err != nil {
		t.Fatalf("select: %v", err)
	}
	value, err := session.page.Locator("#lang").InputValue()
	if err != nil {
		t.Fatalf("read value: %v", err)
	}
	if value != "en" {
		t.Fatalf("expected en, got %q", value)
	}
}

func TestExecuteCheck(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<input type="checkbox" id="agree" />
</body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "goto",
		Value1: srv.URL,
	}, "", nil); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := executeAction(context.Background(), session, stepdsl.Action{
		Kind:   "check",
		Value1: "#agree",
	}, "", nil); err != nil {
		t.Fatalf("check: %v", err)
	}
	checked, err := session.page.Locator("#agree").IsChecked()
	if err != nil {
		t.Fatalf("is checked: %v", err)
	}
	if !checked {
		t.Fatal("expected checkbox to be checked")
	}
}
