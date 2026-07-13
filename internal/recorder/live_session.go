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
		notifySnapshot(opts.Callbacks.OnStepRecorded, *recorded)
	}

	pageClosed := make(chan struct{}, 1)
	page.OnClose(func(playwright.Page) {
		select {
		case pageClosed <- struct{}{}:
		default:
		}
	})

	lastURL := page.URL()
	pollState := session.pollStateForPage(page)
	var stepNotify StepNotifier
	if opts.Callbacks.OnStepRecorded != nil {
		stepNotify = func(event RecordStepEvent) {
			opts.Callbacks.OnStepRecorded(event)
		}
	}
	poll := time.NewTicker(100 * time.Millisecond)
	defer poll.Stop()
	nextEvaluateAt := time.Now()
	idlePolls := 0
	lastEventAt := time.Now()

	processToolbar := func() {
		if action := takeToolbarAction(page); action != "" {
			idlePolls = 0
			ProcessToolbarAction(action, session, page, opts, stepNotify, recorded)
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
		case <-poll.C:
		}
		if !session.BrowserAlive() {
			if opts.Callbacks.OnBrowserLost != nil {
				opts.Callbacks.OnBrowserLost()
			}
			return context.Canceled
		}
		// Poll toolbar at full rate; do not throttle behind recorder event backoff.
		if !session.TestRunHeld() {
			syncBrowserToolbar(page, session, opts.BrowseOnly)
			processToolbar()
		}
		now := time.Now()
		if now.Before(nextEvaluateAt) {
			continue
		}
		if session.TestRunHeld() {
			nextEvaluateAt = time.Now().Add(250 * time.Millisecond)
			continue
		}
		nextEvaluateAt = time.Now().Add(100 * time.Millisecond)
		if session.RelaunchPending() {
			if u := page.URL(); u != "" {
				session.SetResumeURL(u)
			}
			return ErrRelaunchHeadless
		}
		for session.IsPaused() {
			if !session.TestRunHeld() {
				syncBrowserToolbar(page, session, opts.BrowseOnly)
				processToolbar()
			}
			evaluateRecorderCleanup(page, `() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = true; }`)
			if session.RelaunchPending() {
				if u := page.URL(); u != "" {
					session.SetResumeURL(u)
				}
				return ErrRelaunchHeadless
			}
			if !session.BrowserAlive() {
				if opts.Callbacks.OnBrowserLost != nil {
					opts.Callbacks.OnBrowserLost()
				}
				return context.Canceled
			}
			if err := sleepContext(ctx, 100*time.Millisecond); err != nil {
				return err
			}
		}
		evaluateRecorderCleanup(page, `() => { if (window.__scenariaRecorder) window.__scenariaRecorder.paused = false; }`)
		if opts.IdleTimeout > 0 && session.CaptureEnabled() && time.Since(lastEventAt) >= opts.IdleTimeout {
			CompleteCaptureStop(session, page, stepNotify, "idle", opts.Callbacks.OnCaptureStop)
			if opts.Callbacks.OnStepRecorded != nil {
				continue
			}
			return nil
		}

		if !session.CaptureEnabled() {
			nextEvaluateAt = time.Now().Add(500 * time.Millisecond)
			continue
		}

		events, err := drainRecorderEvents(page)
		if err != nil {
			return err
		}
		hadEvents, urlChanged := session.applyRecorderPollBatch(pollState, events, page.URL(), now, stepNotify)
		if hadEvents || urlChanged {
			lastEventAt = time.Now()
		}
		lastURL = pollState.lastURL
		if hadEvents || urlChanged {
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
		_ = lastURL
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
	phase := browserToolbarPhase(session, browseOnly)
	script := fmt.Sprintf(`() => {
		if (!window.__scenariaToolbar) return;
		window.__scenariaToolbar.setState({
			recording: %v,
			paused: %v,
			browserOnly: %v,
			stepCount: %d,
			phase: %q,
		});
	}`, recording, paused, browserOnly, session.RecordedStepCount(), phase)
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
	_ = selector.ApplyLibraryHeuristicsFromSettings(page, appCfg)
}

func injectRecorderOnPage(page playwright.Page) error {
	if _, err := page.Evaluate(selector.HeuristicsJS); err != nil {
		return err
	}
	injectSelectorOrderFromSettings(page)
	_, err := page.Evaluate(selector.RecorderListenersJS)
	return err
}
