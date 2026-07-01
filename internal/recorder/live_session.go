package recorder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/browserconfig"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
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
	browser, err := pw.Chromium.Launch(browserconfig.LaunchOptions("chromium", headless, 0))
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer browser.Close()

	contextOpts := browserconfig.NewContextOptions(headless, opts.HTTPCredentials)
	bctx, err := browser.NewContext(contextOpts)
	if err != nil {
		return fmt.Errorf("create browser context: %w", err)
	}
	defer bctx.Close()
	if err := registerBrowserInitScripts(bctx, !opts.BrowseOnly); err != nil {
		return fmt.Errorf("register browser init scripts: %w", err)
	}

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
	if strings.TrimSpace(startURL) != "" {
		if _, err := page.Goto(startURL); err != nil {
			return fmt.Errorf("goto start URL: %w", err)
		}
	}
	if !opts.BrowseOnly {
		if err := injectRecorderOnPage(page); err != nil {
			return fmt.Errorf("inject recorder script: %w", err)
		}
		session.recorderInjected.Store(true)
	}
	if _, err := page.Evaluate(selector.BrowserToolbarJS); err != nil {
		return fmt.Errorf("inject browser toolbar: %w", err)
	}
	if !opts.BrowseOnly {
		if err := applyRecorderConfig(page, opts, session); err != nil {
			return fmt.Errorf("configure recorder: %w", err)
		}
	}

	session.Bind(page, recorded)
	defer func() {
		session.mu.Lock()
		session.page = nil
		session.mu.Unlock()
	}()

	if session.CaptureEnabled() && opts.Callbacks.OnStepRecorded != nil {
		for i, st := range *recorded {
			if line, ok := RecordedStepToLine(st); ok {
				opts.Callbacks.OnStepRecorded(i, line)
			}
		}
	}

	pageClosed := make(chan struct{}, 1)
	page.OnClose(func(playwright.Page) {
		select {
		case pageClosed <- struct{}{}:
		default:
		}
	})

	lastURL := page.URL()
	lastEventAt := time.Now()
	stepNotify := func(index int, line string) {
		if opts.Callbacks.OnStepRecorded != nil {
			opts.Callbacks.OnStepRecorded(index, line)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pageClosed:
			if opts.Callbacks.OnBrowserLost != nil {
				opts.Callbacks.OnBrowserLost()
			}
			return context.Canceled
		default:
		}
		if !session.BrowserAlive() {
			if opts.Callbacks.OnBrowserLost != nil {
				opts.Callbacks.OnBrowserLost()
			}
			return context.Canceled
		}
		syncBrowserToolbar(page, session, opts.BrowseOnly)
		if action := takeToolbarAction(page); action != "" {
			switch action {
			case "stop":
				if session.CaptureEnabled() {
					session.EndCapture()
					if opts.Callbacks.OnCaptureStop != nil {
						opts.Callbacks.OnCaptureStop()
					}
					if opts.Callbacks.OnStepRecorded != nil {
						continue
					}
					return context.Canceled
				}
				return context.Canceled
			case "pause":
				session.Pause()
			case "resume":
				session.Resume()
			case "record":
				if !session.CaptureEnabled() {
					replay := ShouldSyncRecordedStepsOnCaptureStart(session)
					if err := session.BeginCapture(); err != nil {
						// Keep the browser alive; user can retry from the IDE.
						continue
					}
					if opts.Callbacks.OnCaptureStart != nil {
						opts.Callbacks.OnCaptureStart(!replay)
					}
					if replay {
						for i, st := range *recorded {
							if line, ok := RecordedStepToLine(st); ok {
								stepNotify(i, line)
							}
						}
					}
				}
			case "picker":
				if opts.Callbacks.OnPickerRequest != nil {
					opts.Callbacks.OnPickerRequest()
				}
			}
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
		if opts.IdleTimeout > 0 && session.CaptureEnabled() && time.Since(lastEventAt) >= opts.IdleTimeout {
			session.captureEnabled.Store(false)
			session.paused.Store(false)
			if opts.Callbacks.OnCaptureStop != nil {
				opts.Callbacks.OnCaptureStop()
			}
			if opts.Callbacks.OnStepRecorded != nil {
				continue
			}
			return nil
		}

		if !session.CaptureEnabled() {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if currentURL := page.URL(); currentURL != "" && currentURL != lastURL {
			appendGotoStep(recorded, currentURL, stepNotify)
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
			appendCoalescedStep(recorded, step, stepNotify)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func syncBrowserToolbar(page playwright.Page, session *LiveSession, browseOnly bool) {
	recording := session.CaptureEnabled()
	paused := session.IsPaused()
	browserOnly := browseOnly && !recording
	script := fmt.Sprintf(`() => {
		if (!window.__scenariaToolbar) return;
		window.__scenariaToolbar.setState({
			recording: %v,
			paused: %v,
			browserOnly: %v,
			stepCount: %d,
		});
	}`, recording, paused, browserOnly, session.RecordedStepCount())
	_, _ = page.Evaluate(script)
}

func takeToolbarAction(page playwright.Page) string {
	raw, err := page.Evaluate(`() => window.__scenariaToolbar?.takeAction?.() || null`)
	if err != nil || raw == nil {
		return ""
	}
	if action, ok := raw.(string); ok {
		return action
	}
	return ""
}

func registerBrowserInitScripts(bctx playwright.BrowserContext, includeRecorder bool) error {
	scripts := []string{selector.HeuristicsJS, selector.BrowserToolbarJS}
	if includeRecorder {
		scripts = append(scripts, selector.RecorderListenersJS)
	}
	for _, script := range scripts {
		if err := bctx.AddInitScript(playwright.Script{Content: playwright.String(script)}); err != nil {
			return err
		}
	}
	return nil
}

func injectSelectorOrderFromSettings(page playwright.Page) {
	appCfg, err := settings.LoadDefaultAppSettings()
	if err != nil || appCfg == nil {
		return
	}
	_ = selector.ApplySelectorOrder(page, appCfg.SelectorClickStrategies, appCfg.SelectorInputStrategies)
}

func injectRecorderOnPage(page playwright.Page) error {
	if _, err := page.Evaluate(selector.HeuristicsJS); err != nil {
		return err
	}
	injectSelectorOrderFromSettings(page)
	_, err := page.Evaluate(selector.RecorderListenersJS)
	return err
}
