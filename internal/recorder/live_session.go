package recorder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/browserconfig"
	"github.com/bafgion/scenaria-golang/internal/logx"
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
	defer closeRecorderResource("browser", func() error { return browser.Close() })

	contextOpts := browserconfig.NewContextOptions(headless, opts.HTTPCredentials)
	bctx, err := browser.NewContext(contextOpts)
	if err != nil {
		return fmt.Errorf("create browser context: %w", err)
	}
	defer ReleasePickerBinding(bctx)
	defer closeRecorderResource("context", func() error { return bctx.Close() })
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
	if opts.Callbacks.OnBrowserOpened != nil {
		opts.Callbacks.OnBrowserOpened()
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
	poll := time.NewTicker(100 * time.Millisecond)
	defer poll.Stop()
	nextEvaluateAt := time.Now()
	idlePolls := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pageClosed:
			if opts.Callbacks.OnBrowserLost != nil {
				opts.Callbacks.OnBrowserLost()
			}
			return context.Canceled
		case <-poll.C:
		}
		now := time.Now()
		if now.Before(nextEvaluateAt) {
			continue
		}
		if !session.BrowserAlive() {
			if opts.Callbacks.OnBrowserLost != nil {
				opts.Callbacks.OnBrowserLost()
			}
			return context.Canceled
		}
		if session.TestRunHeld() {
			nextEvaluateAt = time.Now().Add(250 * time.Millisecond)
			continue
		}
		nextEvaluateAt = time.Now().Add(100 * time.Millisecond)
		syncBrowserToolbar(page, session, opts.BrowseOnly)
		if action := takeToolbarAction(page); action != "" {
			idlePolls = 0
			switch action {
			case "stop":
				if session.CaptureEnabled() {
					session.EndCapture()
					if opts.Callbacks.OnCaptureStop != nil {
						opts.Callbacks.OnCaptureStop("manual")
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
			evaluateRecorderCleanup(page, `() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = true; }`)
			if err := sleepContext(ctx, 100*time.Millisecond); err != nil {
				return err
			}
			if session.RelaunchPending() {
				if u := page.URL(); u != "" {
					session.SetResumeURL(u)
				}
				return ErrRelaunchHeadless
			}
		}
		evaluateRecorderCleanup(page, `() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = false; }`)
		if opts.IdleTimeout > 0 && session.CaptureEnabled() && time.Since(lastEventAt) >= opts.IdleTimeout {
			session.captureEnabled.Store(false)
			session.paused.Store(false)
			if opts.Callbacks.OnCaptureStop != nil {
				opts.Callbacks.OnCaptureStop("idle")
			}
			if opts.Callbacks.OnStepRecorded != nil {
				continue
			}
			return nil
		}

		if !session.CaptureEnabled() {
			nextEvaluateAt = time.Now().Add(500 * time.Millisecond)
			continue
		}

		urlChanged := false
		if currentURL := page.URL(); currentURL != "" && currentURL != lastURL {
			session.AppendGotoStep(currentURL, stepNotify)
			lastURL = currentURL
			lastEventAt = time.Now()
			urlChanged = true
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
			session.AppendCoalescedStep(step, stepNotify)
		}
		if len(events) > 0 || urlChanged {
			idlePolls = 0
			nextEvaluateAt = time.Now().Add(100 * time.Millisecond)
		} else {
			idlePolls++
			if idlePolls >= 5 {
				nextEvaluateAt = time.Now().Add(750 * time.Millisecond)
			} else {
				nextEvaluateAt = time.Now().Add(300 * time.Millisecond)
			}
		}
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
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
	evaluateRecorderCleanup(page, script)
}

func closeRecorderResource(resource string, closeFn func() error) {
	if closeFn == nil {
		return
	}
	if err := closeFn(); err != nil {
		logx.Debug("recorder cleanup", "resource", resource, "error", err)
	}
}

func evaluateRecorderCleanup(page playwright.Page, script string) {
	if page == nil {
		return
	}
	if _, err := page.Evaluate(script); err != nil {
		logx.Debug("recorder evaluate", "error", err)
	}
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
