//go:build integration

package gui

import (
	"context"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	playwright "github.com/mxschmitt/playwright-go"
)

func TestRunOnLiveBrowserPreservesRecorderPage(t *testing.T) {
	page, cleanup := openLiveBrowserPage(t)
	defer cleanup()

	live := recorder.NewLiveSession()
	steps := []recorder.RecordedStep{}
	live.Bind(page, &steps)

	runSvc := &RunService{}
	plan := player.ExecutionPlan{
		Cases: []player.RunCase{{
			FeaturePath: "live.feature",
			Name:        "reuse",
			Steps: []gherkin.Step{{
				Keyword: "Допустим",
				Text:    `жду "100ms"`,
				Line:    1,
			}},
		}},
	}
	exec := player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
		BrowserName:   "chromium",
		CloseAfterRun: false,
	})

	_, err := runSvc.RunOnLiveBrowser(context.Background(), exec, plan, "domcontentloaded", live)
	if err != nil {
		t.Fatalf("RunOnLiveBrowser: %v", err)
	}
	if !live.BrowserAlive() {
		t.Fatal("live session browser not alive after run")
	}
	if live.TestRunHeld() {
		t.Fatal("test hold should be released after run")
	}
}

func TestRunOnLiveBrowserCancelPreservesRecorderPage(t *testing.T) {
	page, cleanup := openLiveBrowserPage(t)
	defer cleanup()

	live := recorder.NewLiveSession()
	steps := []recorder.RecordedStep{}
	live.Bind(page, &steps)

	runSvc := &RunService{}
	plan := player.ExecutionPlan{
		Cases: []player.RunCase{{
			FeaturePath: "live.feature",
			Name:        "reuse",
			Steps: []gherkin.Step{{
				Keyword: "Допустим",
				Text:    `жду "60s"`,
				Line:    1,
			}},
		}},
	}
	exec := player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
		BrowserName:   "chromium",
		CloseAfterRun: false,
	})

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		_, err := runSvc.RunOnLiveBrowser(ctx, exec, plan, "domcontentloaded", live)
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
		t.Fatal("RunOnLiveBrowser did not finish after cancel")
	}
	if !live.BrowserAlive() {
		t.Fatal("live session browser not alive after canceled run")
	}
	if live.TestRunHeld() {
		t.Fatal("test hold should be released after canceled run")
	}
}

func openLiveBrowserPage(t *testing.T) (playwright.Page, func()) {
	t.Helper()
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("run playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(true)})
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
	return page, func() {
		_ = page.Close()
		_ = browser.Close()
		_ = pw.Stop()
	}
}
