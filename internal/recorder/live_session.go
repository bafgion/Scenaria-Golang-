package recorder

import (
	"context"
	"fmt"
	"time"

	"github.com/bafgion/scenaria-golang/internal/player"
	playwright "github.com/mxschmitt/playwright-go"
)

func runLiveBrowserSession(
	ctx context.Context,
	pw *playwright.Playwright,
	opts LiveOptions,
	session *LiveSession,
	recorded *[]RecordedStep,
	headless bool,
) error {
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer browser.Close()

	contextOpts := playwright.BrowserNewContextOptions{}
	if opts.HTTPCredentials != nil {
		contextOpts.HttpCredentials = opts.HTTPCredentials
	}
	bctx, err := browser.NewContext(contextOpts)
	if err != nil {
		return fmt.Errorf("create browser context: %w", err)
	}
	defer bctx.Close()

	page, err := bctx.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}
	if opts.TestClient != nil {
		if err := player.ApplyTestClient(page, opts.TestClient); err != nil {
			return fmt.Errorf("apply test client: %w", err)
		}
	}

	startURL := session.ResumeURL(opts.StartURL)
	if len(*recorded) > 0 {
		for i := len(*recorded) - 1; i >= 0; i-- {
			if (*recorded)[i].Action == "goto" && (*recorded)[i].Value != "" {
				startURL = (*recorded)[i].Value
				break
			}
		}
	}
	if _, err := page.Goto(startURL); err != nil {
		return fmt.Errorf("goto start URL: %w", err)
	}
	if _, err := page.Evaluate(RecorderScript); err != nil {
		return fmt.Errorf("inject recorder script: %w", err)
	}
	if err := applyRecorderConfig(page, opts, session); err != nil {
		return fmt.Errorf("configure recorder: %w", err)
	}

	session.Bind(page, recorded)
	defer func() {
		session.mu.Lock()
		session.page = nil
		session.mu.Unlock()
	}()

	lastURL := page.URL()
	lastEventAt := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if session.RelaunchPending() {
			if u := page.URL(); u != "" {
				session.SetResumeURL(u)
			}
			return ErrRelaunchHeadless
		}
		for session.IsPaused() {
			_, _ = page.Evaluate(`() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = true; }`, nil)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
			}
			if session.RelaunchPending() {
				if u := page.URL(); u != "" {
					session.SetResumeURL(u)
				}
				return ErrRelaunchHeadless
			}
		}
		_, _ = page.Evaluate(`() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = false; }`, nil)
		if time.Since(lastEventAt) >= opts.IdleTimeout {
			return nil
		}

		if currentURL := page.URL(); currentURL != "" && currentURL != lastURL {
			*recorded = append(*recorded, RecordedStep{Action: "goto", Value: currentURL})
			lastURL = currentURL
			lastEventAt = time.Now()
		}

		raw, err := page.Evaluate(`() => {
			const r = window.__scenariaRecorder;
			if (!r || !r.events.length) return [];
			const out = r.events.splice(0, r.events.length);
			return out;
		}`)
		if err != nil {
			return fmt.Errorf("read recorder events: %w", err)
		}
		events, err := decodeEvents(raw)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			lastEventAt = time.Now()
		}
		for _, event := range events {
			detail := normalizeDetail(event.Detail)
			step, ok := EventToRecordedStep(event.Type, detail)
			if !ok {
				continue
			}
			*recorded = append(*recorded, step)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
