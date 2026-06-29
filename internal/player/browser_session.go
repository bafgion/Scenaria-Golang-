package player

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
}

type PlaywrightExecutor struct {
	options PlaywrightExecutorOptions
}

func NewPlaywrightExecutor(options PlaywrightExecutorOptions) *PlaywrightExecutor {
	return &PlaywrightExecutor{options: options}
}

type browserSession struct {
	browser      playwright.Browser
	context      playwright.BrowserContext
	page         playwright.Page
	closed       bool
	traceEnabled bool
	traceStopped bool
	videoEnabled bool
}

func newBrowserSession(pw *playwright.Playwright, options PlaywrightExecutorOptions) (*browserSession, error) {
	browser, err := launchBrowser(pw, options)
	if err != nil {
		return nil, err
	}
	ctxOpts := playwright.BrowserNewContextOptions{}
	if options.HTTPCredentials != nil {
		ctxOpts.HttpCredentials = options.HTTPCredentials
	}
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
	session := &browserSession{
		browser:      browser,
		context:      bctx,
		page:         page,
		videoEnabled: strings.TrimSpace(options.VideoDir) != "",
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
	if s.closed {
		return
	}
	if s.traceEnabled && !s.traceStopped && s.context != nil {
		_ = s.context.Tracing().Stop()
		s.traceStopped = true
	}
	if s.page != nil {
		_ = s.page.Close()
	}
	if s.context != nil {
		_ = s.context.Close()
	}
	if s.browser != nil {
		_ = s.browser.Close()
	}
	s.closed = true
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
	name := strings.ToLower(strings.TrimSpace(options.BrowserName))
	if name == "" {
		name = "chromium"
	}

	launchOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(options.Headless),
	}
	if options.SlowMo > 0 {
		launchOpts.SlowMo = playwright.Float(options.SlowMo)
	}

	switch name {
	case "chromium":
		return pw.Chromium.Launch(launchOpts)
	case "firefox":
		return pw.Firefox.Launch(launchOpts)
	case "webkit":
		return pw.WebKit.Launch(launchOpts)
	default:
		return nil, fmt.Errorf("unsupported browser %q (supported: chromium, firefox, webkit)", name)
	}
}

func executeAction(ctx context.Context, session *browserSession, action stepdsl.Action, baseURL string, runCtx *RunContext) error {
	if session == nil || session.closed {
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
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   playwright.Float(NavTimeoutMs),
		})
		if err != nil {
			return fmt.Errorf("goto failed: %w", err)
		}
		return nil
	case "click":
		if err := clickWithFallback(page, action.Value1); err != nil {
			return err
		}
		return nil
	case "double-click":
		locator := selector.ResolveChainedLocator(page, action.Value1)
		if err := locator.Dblclick(playwright.LocatorDblclickOptions{Timeout: playwright.Float(9000)}); err != nil {
			return fmt.Errorf("double click failed: %w", err)
		}
		return nil
	case "hover":
		if err := hoverTarget(page, action.Value1); err != nil {
			return fmt.Errorf("hover failed: %w", err)
		}
		return nil
	case "fill":
		if err := page.Fill(action.Value2, action.Value1); err != nil {
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
		if err := page.Fill(action.Value2, value); err != nil {
			return fmt.Errorf("fill generated failed: %w", err)
		}
		return nil
	case "select":
		if _, err := page.SelectOption(action.Value2, playwright.SelectOptionValues{Values: &[]string{action.Value1}}); err != nil {
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
		if err := page.SetInputFiles(action.Value2, path); err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}
		return nil
	case "clear":
		if err := page.Locator(action.Value1).Clear(); err != nil {
			return fmt.Errorf("clear failed: %w", err)
		}
		return nil
	case "check":
		if err := page.Check(action.Value1); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}
		return nil
	case "uncheck":
		if err := page.Uncheck(action.Value1); err != nil {
			return fmt.Errorf("uncheck failed: %w", err)
		}
		return nil
	case "press":
		if err := page.Keyboard().Press(action.Value1); err != nil {
			return fmt.Errorf("key press failed: %w", err)
		}
		return nil
	case "press-in":
		if err := page.Press(action.Value2, action.Value1); err != nil {
			return fmt.Errorf("key press in element failed: %w", err)
		}
		return nil
	case "assert-visible":
		if err := page.Locator(action.Value1).WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("expected element %q to be visible: %w", action.Value1, err)
		}
		return nil
	case "assert-hidden":
		locator := page.Locator(action.Value1)
		count, err := locator.Count()
		if err != nil {
			return fmt.Errorf("hidden check failed: %w", err)
		}
		if count == 0 {
			return nil
		}
		visible, err := locator.First().IsVisible()
		if err != nil {
			return fmt.Errorf("hidden check failed: %w", err)
		}
		if visible {
			return fmt.Errorf("expected element %q to be hidden", action.Value1)
		}
		return nil
	case "assert-text":
		locator := page.Locator(action.Value2)
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
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
		if !UrlsMatch(page.URL(), action.Value1) {
			return fmt.Errorf("expected URL %q, got %q", action.Value1, page.URL())
		}
		return nil
	case "assert-url-contains":
		if !strings.Contains(page.URL(), action.Value1) {
			return fmt.Errorf("expected URL to contain %q, got %q", action.Value1, page.URL())
		}
		return nil
	case "scroll-to":
		if err := page.Locator(action.Value1).ScrollIntoViewIfNeeded(); err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}
		return nil
	case "drag-drop":
		if err := page.Locator(action.Value1).DragTo(page.Locator(action.Value2)); err != nil {
			return fmt.Errorf("drag failed: %w", err)
		}
		return nil
	case "wait-visible":
		if err := page.Locator(action.Value1).WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}); err != nil {
			return fmt.Errorf("wait for visible failed: %w", err)
		}
		return nil
	case "wait-hidden":
		if err := page.Locator(action.Value1).WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateHidden,
		}); err != nil {
			return fmt.Errorf("wait for hidden failed: %w", err)
		}
		return nil
	case "wait":
		duration, err := stepdsl.ParseWaitDuration(action.Value1)
		if err != nil {
			return err
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
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   playwright.Float(NavTimeoutMs),
		}); err != nil {
			return fmt.Errorf("reload failed: %w", err)
		}
		return nil
	case "go-back":
		if _, err := page.GoBack(playwright.PageGoBackOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   playwright.Float(NavTimeoutMs),
		}); err != nil {
			return fmt.Errorf("go back failed: %w", err)
		}
		return nil
	case "close-browser":
		session.close()
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
		locator := page.Locator(action.Value1)
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
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
		return downloadByClick(page, action.Value1, runCtx)
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
		return drawSignature(page, action.Value1)
	case "prompt-email-code":
		if runCtx == nil {
			return fmt.Errorf("prompt-email-code requires run context")
		}
		return runEmailCode(page, action, runCtx)
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

func downloadByClick(page playwright.Page, selector string, runCtx *RunContext) error {
	download, err := page.ExpectDownload(func() error {
		return page.Click(selector)
	})
	if err != nil {
		return fmt.Errorf("download click failed: %w", err)
	}
	filename := download.SuggestedFilename()
	if strings.TrimSpace(filename) == "" {
		filename = fmt.Sprintf("download-%d", time.Now().UnixNano())
	}
	path := filepath.Join(runCtx.DownloadDir(), filename)
	if err := download.SaveAs(path); err != nil {
		return fmt.Errorf("save download failed: %w", err)
	}
	runCtx.SetLastDownload(path)
	return nil
}

func drawSignature(page playwright.Page, selector string) error {
	locator := page.Locator(selector)
	if err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		return fmt.Errorf("draw signature failed: %w", err)
	}
	box, err := locator.BoundingBox()
	if err != nil || box == nil {
		return fmt.Errorf("draw signature: element has no bounding box")
	}
	startX := box.X + box.Width*0.2
	startY := box.Y + box.Height*0.6
	endX := box.X + box.Width*0.8
	endY := box.Y + box.Height*0.4
	midX := box.X + box.Width*0.5
	midY := box.Y + box.Height*0.7
	mouse := page.Mouse()
	if err := mouse.Move(startX, startY); err != nil {
		return err
	}
	if err := mouse.Down(); err != nil {
		return err
	}
	if err := mouse.Move(midX, midY); err != nil {
		return err
	}
	if err := mouse.Move(endX, endY); err != nil {
		return err
	}
	return mouse.Up()
}

func clickWithFallback(page playwright.Page, sel string) error {
	locator := selector.ResolveChainedLocator(page, sel)
	timeout := playwright.Float(9000)
	if err := locator.Click(playwright.LocatorClickOptions{Timeout: timeout}); err != nil {
		if revealHoverMenu(page, sel) == nil {
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

func hoverTarget(page playwright.Page, sel string) error {
	locator := selector.ResolveHoverLocator(page, sel)
	if err := locator.ScrollIntoViewIfNeeded(); err != nil {
		return err
	}
	if err := locator.Hover(playwright.LocatorHoverOptions{Force: playwright.Bool(true)}); err != nil {
		return err
	}
	time.Sleep(400 * time.Millisecond)
	return nil
}

func revealHoverMenu(page playwright.Page, sel string) error {
	for _, candidate := range selector.HoverLocatorCandidates(sel) {
		locator := page.Locator(candidate).First()
		count, err := locator.Count()
		if err != nil || count == 0 {
			continue
		}
		if err := locator.Hover(playwright.LocatorHoverOptions{Force: playwright.Bool(true)}); err != nil {
			continue
		}
		time.Sleep(400 * time.Millisecond)
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
