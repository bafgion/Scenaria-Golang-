package recorder

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

var ErrPickerCancelled = errors.New("выбор элемента отменён")

type pickerBinding struct {
	mu            sync.Mutex
	result        chan string
	allowedPage   playwright.Page
	allowedOrigin string
}

func pageOrigin(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func (pb *pickerBinding) setAllowedPage(page playwright.Page) {
	origin := ""
	if page != nil {
		origin = pageOrigin(page.URL())
	}
	pb.mu.Lock()
	pb.allowedPage = page
	pb.allowedOrigin = origin
	pb.mu.Unlock()
}

func (pb *pickerBinding) clearAllowedPage() {
	pb.mu.Lock()
	pb.allowedPage = nil
	pb.allowedOrigin = ""
	pb.mu.Unlock()
}

func (pb *pickerBinding) bindingAllowed(source *playwright.BindingSource) bool {
	if source == nil || source.Page == nil {
		return false
	}
	pb.mu.Lock()
	allowedPage := pb.allowedPage
	allowedOrigin := pb.allowedOrigin
	pb.mu.Unlock()
	if allowedPage == nil || source.Page != allowedPage {
		return false
	}
	if allowedOrigin == "" {
		return true
	}
	return pageOrigin(source.Page.URL()) == allowedOrigin
}

var (
	pickerMu       sync.Mutex
	pickerBindings = map[playwright.BrowserContext]*pickerBinding{}
)

// ReleasePickerBinding drops per-context picker channels when a browser context closes.
func ReleasePickerBinding(ctx playwright.BrowserContext) {
	if ctx == nil {
		return
	}
	pickerMu.Lock()
	delete(pickerBindings, ctx)
	pickerMu.Unlock()
}

func ensurePickerBindings(ctx playwright.BrowserContext) (*pickerBinding, error) {
	pickerMu.Lock()
	if pb, ok := pickerBindings[ctx]; ok {
		pickerMu.Unlock()
		return pb, nil
	}
	pb := &pickerBinding{result: make(chan string, 1)}
	pickerBindings[ctx] = pb
	pickerMu.Unlock()

	if err := ctx.ExposeBinding("pickSelectorDone", func(source *playwright.BindingSource, args ...any) any {
		if !pb.bindingAllowed(source) {
			return nil
		}
		if len(args) == 0 {
			return nil
		}
		selector, _ := args[0].(string)
		select {
		case pb.result <- selector:
		default:
		}
		return nil
	}); err != nil {
		ReleasePickerBinding(ctx)
		return nil, err
	}
	if err := ctx.ExposeBinding("pickSelectorCancel", func(source *playwright.BindingSource, _ ...any) any {
		if !pb.bindingAllowed(source) {
			return nil
		}
		select {
		case pb.result <- "":
		default:
		}
		return nil
	}); err != nil {
		ReleasePickerBinding(ctx)
		return nil, err
	}
	return pb, nil
}

func drainPickerResults(pb *pickerBinding) {
	for {
		select {
		case <-pb.result:
		default:
			return
		}
	}
}

func PickSelectorOnPage(ctx context.Context, page playwright.Page) (string, error) {
	if page == nil {
		return "", fmt.Errorf("страница браузера недоступна")
	}
	bctx := page.Context()
	pb, err := ensurePickerBindings(bctx)
	if err != nil {
		return "", fmt.Errorf("picker bindings: %w", err)
	}
	drainPickerResults(pb)

	if _, err := page.Evaluate(selector.HeuristicsJS); err != nil {
		return "", fmt.Errorf("inject heuristics: %w", err)
	}
	if appCfg, err := settings.LoadDefaultAppSettings(); err == nil && appCfg != nil {
		_ = selector.ApplySelectorOrder(page, appCfg.SelectorClickStrategies, appCfg.SelectorInputStrategies)
	}

	if _, err := page.Evaluate(selector.PickerInstallScript); err != nil {
		uninstallPicker(page)
		return "", fmt.Errorf("install picker: %w", err)
	}
	pb.setAllowedPage(page)
	defer pb.clearAllowedPage()

	timeout := time.NewTimer(5 * time.Minute)
	defer timeout.Stop()

	select {
	case <-ctx.Done():
		uninstallPicker(page)
		return "", ctx.Err()
	case <-timeout.C:
		uninstallPicker(page)
		return "", fmt.Errorf("время выбора элемента истекло")
	case value := <-pb.result:
		uninstallPicker(page)
		if value == "" {
			return "", ErrPickerCancelled
		}
		return value, nil
	}
}

func uninstallPicker(page playwright.Page) {
	if page == nil {
		return
	}
	_, _ = page.Evaluate(selector.PickerUninstallScript)
}
