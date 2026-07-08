//go:build integration

package player

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

func TestCanceledGotoReturnsPromptly(t *testing.T) {
	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- pageGoto(ctx, session.page, "http://10.255.255.1/", session.navigationWaitUntil())
	}()

	time.Sleep(100 * time.Millisecond)
	cancelStart := time.Now()
	cancel()
	err := <-done
	cancelElapsed := time.Since(cancelStart)
	if err == nil {
		t.Fatal("expected context error")
	}
	if cancelElapsed > 2*time.Second {
		t.Fatalf("cancel took too long: %v", cancelElapsed)
	}
	drainPendingAsync(3 * time.Second)
}

func TestCanceledWaitForReturnsPromptly(t *testing.T) {
	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	locator := session.page.Locator("#missing-element")
	done := make(chan error, 1)
	go func() {
		done <- waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(30000),
		})
	}()

	time.Sleep(50 * time.Millisecond)
	cancelStart := time.Now()
	cancel()
	err := <-done
	if err == nil {
		t.Fatal("expected context error")
	}
	if elapsed := time.Since(cancelStart); elapsed > 2*time.Second {
		t.Fatalf("cancel took too long: %v", elapsed)
	}
	drainPendingAsync(3 * time.Second)
}

func TestCanceledPressReturnsPromptly(t *testing.T) {
	session, cleanup := newTestBrowserSession(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- pressKey(ctx, session.page, "Enter")
	}()
	cancel()
	err := <-done
	if err == nil {
		t.Fatal("expected context error")
	}
	drainPendingAsync(3 * time.Second)
}

func TestCanceledExpectDownloadReturnsPromptly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slow" {
			time.Sleep(30 * time.Second)
			return
		}
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><a href="/slow" download="file.txt">download</a></body></html>`))
	}))
	defer srv.Close()

	session, cleanup := newTestBrowserSession(t)
	defer cleanup()
	if err := pageGoto(context.Background(), session.page, srv.URL, session.navigationWaitUntil()); err != nil {
		t.Fatalf("goto: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := expectDownload(ctx, session.page, func() error {
			_, err := session.page.Evaluate(`() => document.querySelector('a')?.click()`)
			return err
		})
		done <- err
	}()

	time.Sleep(50 * time.Millisecond)
	cancelStart := time.Now()
	cancel()
	err := <-done
	if err == nil {
		t.Fatal("expected context error")
	}
	if elapsed := time.Since(cancelStart); elapsed > 2*time.Second {
		t.Fatalf("cancel took too long: %v", elapsed)
	}
	drainPendingAsync(3 * time.Second)
}

func TestBrowserCloseFreesStuckGoto(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	var once sync.Once
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			once.Do(func() {})
			select {}
		}),
	}
	go func() { _ = server.Serve(listener) }()
	defer func() {
		_ = server.Close()
		_ = listener.Close()
	}()

	session, cleanup := newTestBrowserSession(t)
	ctx := context.Background()
	go func() {
		_ = pageGoto(ctx, session.page, "http://"+listener.Addr().String()+"/", session.navigationWaitUntil())
	}()
	time.Sleep(100 * time.Millisecond)

	beforePending := PendingAsyncCount()
	session.close()
	cleanup()
	drainPendingAsync(5 * time.Second)

	if pending := PendingAsyncCount(); pending > beforePending {
		t.Fatalf("pending async drains increased after close: before=%d after=%d", beforePending, pending)
	}
}

func TestRepeatedCancelBoundedGoroutinesIntegration(t *testing.T) {
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("run playwright: %v", err)
	}
	defer func() { _ = pw.Stop() }()

	const repeats = 20
	const slack = 12

	before := runtime.NumGoroutine()
	for i := 0; i < repeats; i++ {
		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		})
		if err != nil {
			t.Fatalf("launch browser %d: %v", i, err)
		}
		page, err := browser.NewPage()
		if err != nil {
			_ = browser.Close()
			t.Fatalf("new page %d: %v", i, err)
		}
		session := &browserSession{
			browser: browser,
			page:    page,
		}
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			_ = pageGoto(ctx, page, "http://10.255.255.1/", session.navigationWaitUntil())
		}()
		time.Sleep(30 * time.Millisecond)
		cancel()
		session.close()
		_ = browser.Close()
		drainPendingAsync(3 * time.Second)
	}
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	after := runtime.NumGoroutine()
	if after > before+slack {
		t.Fatalf("goroutines grew from %d to %d after %d cancels (slack %d)", before, after, repeats, slack)
	}
}

func TestBrowserPoolShutdownLatencyWithWorkers(t *testing.T) {
	ctx := context.Background()
	pool, err := newBrowserPool(ctx, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}, 4)
	if err != nil {
		t.Skip("playwright not available:", err)
	}

	start := time.Now()
	pool.Close()
	elapsed := time.Since(start)
	if elapsed > 20*time.Second {
		t.Fatalf("pool close took too long: %v", elapsed)
	}
}
