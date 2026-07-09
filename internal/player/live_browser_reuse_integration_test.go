//go:build integration

package player

import (
	"context"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestLiveBrowserReusePageStaysOpenAfterRun(t *testing.T) {
	base, cleanup := newTestBrowserSession(t)
	defer cleanup()

	attached, err := AttachToPage(base.page, "domcontentloaded")
	if err != nil {
		t.Fatalf("attach: %v", err)
	}

	exec := NewPlaywrightExecutor(PlaywrightExecutorOptions{
		BrowserName:   "chromium",
		CloseAfterRun: false,
	})
	runner := BrowserRunner{}
	_, err = runner.ExecuteSequentialAttached(context.Background(), exec, waitPlan("100ms"), attached)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if base.page.IsClosed() {
		t.Fatal("page closed after attached run")
	}
}

func TestLiveBrowserReuseSurvivesCancelAndSecondRun(t *testing.T) {
	base, cleanup := newTestBrowserSession(t)
	defer cleanup()

	exec := NewPlaywrightExecutor(PlaywrightExecutorOptions{
		BrowserName:   "chromium",
		CloseAfterRun: false,
	})
	runner := BrowserRunner{}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		attached, err := AttachToPage(base.page, "domcontentloaded")
		if err != nil {
			errCh <- err
			return
		}
		_, err = runner.ExecuteSequentialAttached(ctx, exec, waitPlan("60s"), attached)
		errCh <- err
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected canceled run error")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("canceled attached run did not finish promptly")
	}
	if base.page.IsClosed() {
		t.Fatal("page closed after canceled attached run")
	}

	attached, err := AttachToPage(base.page, "domcontentloaded")
	if err != nil {
		t.Fatalf("reattach: %v", err)
	}
	_, err = runner.ExecuteSequentialAttached(context.Background(), exec, waitPlan("100ms"), attached)
	if err != nil {
		t.Fatalf("second run after cancel: %v", err)
	}
	if base.page.IsClosed() {
		t.Fatal("page closed after second attached run")
	}
}

func waitPlan(waitText string) ExecutionPlan {
	return ExecutionPlan{
		Cases: []RunCase{{
			FeaturePath: "live.feature",
			Name:        "reuse",
			Steps: []gherkin.Step{{
				Keyword: "Допустим",
				Text:    `жду "` + waitText + `"`,
				Line:    1,
			}},
		}},
	}
}
