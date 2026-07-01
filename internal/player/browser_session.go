package player

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/browserconfig"
	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type PlaywrightExecutorOptions struct {
	BrowserName       string
	Headless          bool
	BaseURL           string
	AutoInstall       bool
	SlowMo            float64
	TraceDir          string
	VideoDir          string
	HTTPCredentials   *playwright.HttpCredentials
	MaxLoopIterations int
	NavWaitUntil      string
	PromptEmailCode   EmailCodePrompter
}

type PlaywrightExecutor struct {
	options PlaywrightExecutorOptions
}

func NewPlaywrightExecutor(options PlaywrightExecutorOptions) *PlaywrightExecutor {
	return &PlaywrightExecutor{options: options}
}

type browserSession struct {
	mu           sync.Mutex
	browser      playwright.Browser
	context      playwright.BrowserContext
	page         playwright.Page
	closed       bool
	traceEnabled bool
	traceStopped bool
	videoEnabled bool
	videoRetained bool
	navWaitUntil  *playwright.WaitUntilState
}

func newBrowserSession(pw *playwright.Playwright, options PlaywrightExecutorOptions) (*browserSession, error) {
	browser, err := launchBrowser(pw, options)
	if err != nil {
		return nil, err
	}
	ctxOpts := browserconfig.NewContextOptions(options.Headless, options.HTTPCredentials)
	if strings.TrimSpace(options.VideoDir) != "" {
		if err := os.MkdirAll(options.VideoDir, 0o755); err != nil {
			_ = browser.Close()
			return nil, fmt.Errorf("create video dir: %w", err)
		}
		ctxOpts.RecordVideo = &playwright.RecordVideo{Dir: playwright.String(options.VideoDir)}
	}
	bctx, err := browser.NewContext(ctxOpts)
	if err != nil {
		_ = browser.Close()
		return nil, fmt.Errorf("create browser context: %w", err)
	}
	page, err := bctx.NewPage()
	if err != nil {
		_ = bctx.Close()
		_ = browser.Close()
		return nil, fmt.Errorf("create browser page: %w", err)
	}
	navWaitUntil, err := ParseNavWaitUntil(options.NavWaitUntil)
	if err != nil {
		_ = bctx.Close()
		_ = browser.Close()
		return nil, err
	}
	session := &browserSession{
		browser:      browser,
		context:      bctx,
		page:         page,
		videoEnabled: strings.TrimSpace(options.VideoDir) != "",
		navWaitUntil: navWaitUntil,
	}
	if strings.TrimSpace(options.TraceDir) != "" {
		if err := bctx.Tracing().Start(playwright.TracingStartOptions{
			Screenshots: playwright.Bool(true),
			Snapshots:   playwright.Bool(true),
		}); err != nil {
			session.close()
			return nil, fmt.Errorf("start playwright trace: %w", err)
		}
		session.traceEnabled = true
	}
	return session, nil
}

func (s *browserSession) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeLocked()
}

func (s *browserSession) closeLocked() {
	if s == nil || s.closed {
		return
	}
	if s.traceEnabled && !s.traceStopped && s.context != nil {
		closeBrowserResource("trace", func() error { return s.context.Tracing().Stop() })
		s.traceStopped = true
	}
	if s.page != nil {
		if s.videoEnabled && !s.videoRetained {
			recorder := s.page.Video()
			closeBrowserResource("page", func() error { return s.page.Close() })
			s.page = nil
			if recorder != nil {
				path, err := recorder.Path()
				if err == nil && strings.TrimSpace(path) != "" {
					if rmErr := os.Remove(path); rmErr != nil {
						logx.Debug("browser cleanup", "resource", "video-temp", "error", rmErr)
					}
				}
			}
		} else {
			closeBrowserResource("page", func() error { return s.page.Close() })
			s.page = nil
		}
	}
	if s.context != nil {
		closeBrowserResource("context", func() error { return s.context.Close() })
		s.context = nil
	}
	if s.browser != nil {
		closeBrowserResource("browser", func() error { return s.browser.Close() })
		s.browser = nil
	}
	s.closed = true
}

func (s *browserSession) navigationWaitUntil() *playwright.WaitUntilState {
	if s != nil && s.navWaitUntil != nil {
		return s.navWaitUntil
	}
	return playwright.WaitUntilStateDomcontentloaded
}

func (s *browserSession) resetForScenario() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.context == nil {
		return fmt.Errorf("browser session is closed")
	}
	if s.page != nil {
		_ = s.page.Close()
		s.page = nil
	}
	page, err := s.context.NewPage()
	if err != nil {
		return fmt.Errorf("reset browser page: %w", err)
	}
	s.page = page
	return nil
}

// finalizeVideoRecording closes page and context so Playwright flushes the webm, then reads it.
func (s *browserSession) finalizeVideoRecording(videoDir string) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s == nil || s.closed || !s.videoEnabled || s.page == nil {
		return nil
	}
	recorder := s.page.Video()
	_ = s.page.Close()
	s.page = nil
	if s.context != nil {
		_ = s.context.Close()
		s.context = nil
	}
	s.videoRetained = true
	if recorder == nil {
		return nil
	}
	return readVideoRecording(recorder, videoDir)
}

func (s *browserSession) setPage(page playwright.Page) {
	s.page = page
	if page != nil {
		_ = page.BringToFront()
	}
}

func (s *browserSession) openPages() []playwright.Page {
	if s.context == nil {
		return nil
	}
	return s.context.Pages()
}

func launchBrowser(pw *playwright.Playwright, options PlaywrightExecutorOptions) (playwright.Browser, error) {
	name := browserconfig.NormalizeEngine(options.BrowserName)
	launchOpts := browserconfig.LaunchOptions(name, options.Headless, options.SlowMo)

	switch name {
	case "chromium":
		return pw.Chromium.Launch(launchOpts)
	case "firefox":
		return pw.Firefox.Launch(launchOpts)
	case "webkit":
		return pw.WebKit.Launch(launchOpts)
	default:
		return nil, fmt.Errorf("unsupported browser %q (supported: chromium, firefox, webkit)", options.BrowserName)
	}
}

func waitAllMatchesHidden(ctx context.Context, page playwright.Page, sel, failFmt string) error {
	locator := selector.ResolveChainedLocator(page, sel)
	count, err := locator.Count()
	if err != nil {
		return fmt.Errorf("hidden check failed: %w", err)
	}
	if count == 0 {
		return nil
	}
	for i := 0; i < count; i++ {
		item := locator.Nth(i)
		visible, err := item.IsVisible()
		if err != nil {
			return fmt.Errorf("hidden check failed: %w", err)
		}
		if !visible {
			continue
		}
		if err := waitForLocator(ctx, item, playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateHidden,
		}); err != nil {
			return fmt.Errorf(failFmt+": %w", sel, err)
		}
	}
	return nil
}

func executeAction(ctx context.Context, session *browserSession, action stepdsl.Action, baseURL string, runCtx *RunContext) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.closed {
		return fmt.Errorf("browser session is closed")
	}
	page := session.page

	switch action.Kind {
	case "goto":
		url := stepdsl.ResolveURL(action.Value1, baseURL)
		if UrlsMatch(page.URL(), url) {
			return nil
		}
		_, err := page.Goto(url, playwright.PageGotoOptions{
			WaitUntil: session.navigationWaitUntil(),
			Timeout:   timeoutMs(ctx, NavTimeoutMs),
		})
		if err != nil {
			return fmt.Errorf("goto failed: %w", err)
		}
		return nil
	case "click":
		if err := clickWithFallback(ctx, page, action.Value1); err != nil {
			return err
		}
		return nil
	case "double-click":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.Dblclick(playwright.LocatorDblclickOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("double click failed: %w", err)
		}
		return nil
	case "hover":
		if err := hoverTarget(ctx, page, action.Value1); err != nil {
			return fmt.Errorf("hover failed: %w", err)
		}
		return nil
	case "fill":
		locator := selector.ResolveChainedLocator(page, action.Value2)
		if err := locator.Fill(action.Value1, playwright.LocatorFillOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("fill failed: %w", err)
		}
		return nil
	case "fill-generated":
		if runCtx == nil {
			return fmt.Errorf("fill-generated requires run context")
		}
		value, err := runCtx.GenerateByKind(action.Value1)
		if err != nil {
			return err
		}
		if err := selector.ResolveChainedLocator(page, action.Value2).Fill(value, playwright.LocatorFillOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("fill generated failed: %w", err)
		}
		return nil
	case "select":
		locator := selector.ResolveChainedLocator(page, action.Value2)
		if _, err := locator.SelectOption(playwright.SelectOptionValues{Values: &[]string{action.Value1}}, playwright.LocatorSelectOptionOptions{
			Timeout: timeoutMs(ctx, ActionTimeoutMs),
		}); err != nil {
			return fmt.Errorf("select failed: %w", err)
		}
		return nil
	case "upload":
		path := action.Value1
		if runCtx != nil {
			resolved, err := runCtx.ResolveText(path)
			if err != nil {
				return err
			}
			path = resolved
		}
		if !filepath.IsAbs(path) && runCtx != nil && strings.TrimSpace(runCtx.projectRoot) != "" {
			path = filepath.Join(runCtx.projectRoot, path)
		}
		if err := selector.ResolveChainedLocator(page, action.Value2).SetInputFiles(path); err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}
		return nil
	case "clear":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.Clear(playwright.LocatorClearOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("clear failed: %w", err)
		}
		return nil
	case "check":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.Check(playwright.LocatorCheckOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}
		return nil
	case "uncheck":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.Uncheck(playwright.LocatorUncheckOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("uncheck failed: %w", err)
		}
		return nil
	case "press":
		if err := pressKey(ctx, page, action.Value1); err != nil {
			return fmt.Errorf("key press failed: %w", err)
		}
		return nil
	case "press-in":
		locator := selector.ResolveChainedLocator(page, action.Value2)
		if err := locator.Press(action.Value1, playwright.LocatorPressOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("key press in element failed: %w", err)
		}
		return nil
	case "assert-visible":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("expected element %q to be visible: %w", action.Value1, err)
		}
		return nil
	case "assert-hidden":
		if err := waitAllMatchesHidden(ctx, page, action.Value1, "expected element %q to be hidden"); err != nil {
			return err
		}
		return nil
	case "assert-text":
		locator := selector.ResolveChainedLocator(page, action.Value2)
		if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("assert text failed: %w", err)
		}
		text, err := locator.InnerText()
		if err != nil {
			return fmt.Errorf("read text failed: %w", err)
		}
		if !strings.Contains(text, action.Value1) {
			return fmt.Errorf("expected text %q in %q, got %q", action.Value1, action.Value2, text)
		}
		return nil
	case "assert-url":
		if err := waitForURL(ctx, page, action.Value1, false); err != nil {
			return err
		}
		return nil
	case "assert-url-contains":
		if err := waitForURL(ctx, page, action.Value1, true); err != nil {
			return err
		}
		return nil
	case "scroll-to":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.ScrollIntoViewIfNeeded(playwright.LocatorScrollIntoViewIfNeededOptions{
			Timeout: timeoutMs(ctx, ActionTimeoutMs),
		}); err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}
		return nil
	case "drag-drop":
		src := selector.ResolveChainedLocator(page, action.Value1)
		dst := selector.ResolveChainedLocator(page, action.Value2)
		if err := src.DragTo(dst, playwright.LocatorDragToOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)}); err != nil {
			return fmt.Errorf("drag failed: %w", err)
		}
		return nil
	case "wait-visible":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("wait for visible failed: %w", err)
		}
		return nil
	case "wait-hidden":
		if err := waitAllMatchesHidden(ctx, page, action.Value1, "wait for hidden on %q"); err != nil {
			return err
		}
		return nil
	case "wait":
		duration, err := stepdsl.ParseWaitDuration(action.Value1)
		if err != nil {
			return err
		}
		duration = capWaitDuration(ctx, duration)
		if duration <= 0 {
			return ctx.Err()
		}
		timer := time.NewTimer(duration)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	case "reload":
		if _, err := page.Reload(playwright.PageReloadOptions{
			WaitUntil: session.navigationWaitUntil(),
			Timeout:   timeoutMs(ctx, NavTimeoutMs),
		}); err != nil {
			return fmt.Errorf("reload failed: %w", err)
		}
		return nil
	case "go-back":
		if _, err := page.GoBack(playwright.PageGoBackOptions{
			WaitUntil: session.navigationWaitUntil(),
			Timeout:   timeoutMs(ctx, NavTimeoutMs),
		}); err != nil {
			return fmt.Errorf("go back failed: %w", err)
		}
		return nil
	case "close-browser":
		session.closeLocked()
		return nil
	case "remember-text":
		if runCtx == nil {
			return fmt.Errorf("remember-text requires run context")
		}
		runCtx.Remember(action.Value2, action.Value1)
		return nil
	case "remember-field":
		if runCtx == nil {
			return fmt.Errorf("remember-field requires run context")
		}
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("remember field failed: %w", err)
		}
		value, err := locator.InputValue()
		if err != nil || strings.TrimSpace(value) == "" {
			value, err = locator.InnerText()
		}
		if err != nil {
			return fmt.Errorf("remember field failed: %w", err)
		}
		runCtx.Remember(action.Value2, strings.TrimSpace(value))
		return nil
	case "remember-url":
		if runCtx == nil {
			return fmt.Errorf("remember-url requires run context")
		}
		runCtx.Remember(action.Value1, page.URL())
		return nil
	case "download-click":
		if runCtx == nil {
			return fmt.Errorf("download-click requires run context")
		}
		return downloadByClick(ctx, page, action.Value1, runCtx)
	case "assert-download-contains":
		if runCtx == nil {
			return fmt.Errorf("assert-download-contains requires run context")
		}
		path := runCtx.LastDownload()
		if path == "" {
			return fmt.Errorf("no downloaded file recorded")
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read downloaded file: %w", err)
		}
		if !strings.Contains(string(content), action.Value1) {
			return fmt.Errorf("downloaded file does not contain %q", action.Value1)
		}
		return nil
	case "draw-signature":
		return drawSignature(ctx, page, action.Value1)
	case "prompt-email-code":
		if runCtx == nil {
			return fmt.Errorf("prompt-email-code requires run context")
		}
		return runEmailCode(ctx, page, action, runCtx)
	case "switch-tab":
		return switchTab(session, action, runCtx)
	case "close-tab":
		return closeCurrentTab(session, runCtx)
	case "assert-tab-count":
		if len(session.openPages()) != action.IntVal {
			return fmt.Errorf("expected %d tab(s), got %d", action.IntVal, len(session.openPages()))
		}
		return nil
	default:
		return fmt.Errorf("unsupported action kind %q", action.Kind)
	}
}

func downloadByClick(ctx context.Context, page playwright.Page, sel string, runCtx *RunContext) error {
	download, err := expectDownload(ctx, page, func() error {
		locator := selector.ResolveChainedLocator(page, sel)
		return locator.Click(playwright.LocatorClickOptions{Timeout: timeoutMs(ctx, ActionTimeoutMs)})
	})
	if err != nil {
		return fmt.Errorf("download click failed: %w", err)
	}
	filename := download.SuggestedFilename()
	path := runCtx.allocateDownloadPath(filename)
	if err := download.SaveAs(path); err != nil {
		return fmt.Errorf("save download failed: %w", err)
	}
	runCtx.SetLastDownload(path)
	return nil
}

func clickWithFallback(ctx context.Context, page playwright.Page, sel string) error {
	locator := selector.ResolveChainedLocator(page, sel)
	timeout := timeoutMs(ctx, ActionTimeoutMs)
	if err := locator.Click(playwright.LocatorClickOptions{Timeout: timeout}); err != nil {
		if revealHoverMenu(ctx, page, sel) == nil {
			locator = selector.ResolveChainedLocator(page, sel)
			if err2 := locator.Click(playwright.LocatorClickOptions{Timeout: timeout}); err2 == nil {
				_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
					State: playwright.LoadStateDomcontentloaded,
				})
				return nil
			}
		}
		return fmt.Errorf("click failed: %w", err)
	}
	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	return nil
}

func hoverTarget(ctx context.Context, page playwright.Page, sel string) error {
	locator := selector.ResolveHoverLocator(page, sel)
	if err := locator.ScrollIntoViewIfNeeded(playwright.LocatorScrollIntoViewIfNeededOptions{
		Timeout: timeoutMs(ctx, ActionTimeoutMs),
	}); err != nil {
		return err
	}
	if err := locator.Hover(playwright.LocatorHoverOptions{
		Force:   playwright.Bool(true),
		Timeout: timeoutMs(ctx, ActionTimeoutMs),
	}); err != nil {
		return err
	}
	return sleepContext(ctx, 400*time.Millisecond)
}

func revealHoverMenu(ctx context.Context, page playwright.Page, sel string) error {
	for _, candidate := range selector.HoverLocatorCandidates(sel) {
		locator := page.Locator(candidate).First()
		count, err := locator.Count()
		if err != nil || count == 0 {
			continue
		}
		if err := locator.Hover(playwright.LocatorHoverOptions{
			Force:   playwright.Bool(true),
			Timeout: timeoutMs(ctx, ActionTimeoutMs),
		}); err != nil {
			continue
		}
		if err := sleepContext(ctx, 400*time.Millisecond); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("could not reveal hover menu for %q", sel)
}

func switchTab(session *browserSession, action stepdsl.Action, runCtx *RunContext) error {
	pages := session.openPages()
	if len(pages) == 0 {
		return fmt.Errorf("no open tabs")
	}
	var target playwright.Page
	switch action.Mode {
	case "index":
		if action.IntVal < 0 {
			return fmt.Errorf("tab number must be at least 1")
		}
		if action.IntVal >= len(pages) {
			return fmt.Errorf("tab index %d out of range (1..%d)", action.IntVal+1, len(pages))
		}
		target = pages[action.IntVal]
	case "title":
		for _, p := range pages {
			title, _ := p.Title()
			if strings.Contains(title, action.Value1) {
				target = p
				break
			}
		}
		if target == nil {
			return fmt.Errorf("no tab with title containing %q", action.Value1)
		}
	case "url":
		for _, p := range pages {
			if strings.Contains(p.URL(), action.Value1) {
				target = p
				break
			}
		}
		if target == nil {
			return fmt.Errorf("no tab with URL containing %q", action.Value1)
		}
	case "first":
		target = pages[0]
	case "new":
		target = pages[len(pages)-1]
	default:
		return fmt.Errorf("unsupported tab switch mode %q", action.Mode)
	}
	session.setPage(target)
	if runCtx != nil {
		runCtx.SetPage(target)
	}
	return nil
}

func closeCurrentTab(session *browserSession, runCtx *RunContext) error {
	current := session.page
	if current == nil {
		return fmt.Errorf("no current tab")
	}
	pages := session.openPages()
	if len(pages) <= 1 {
		return fmt.Errorf("cannot close the only open tab")
	}
	if err := current.Close(); err != nil {
		return fmt.Errorf("close tab failed: %w", err)
	}
	remaining := session.openPages()
	if len(remaining) == 0 {
		return fmt.Errorf("no tabs left after close")
	}
	fallback := remaining[len(remaining)-1]
	session.setPage(fallback)
	if runCtx != nil {
		runCtx.SetPage(fallback)
	}
	return nil
}
