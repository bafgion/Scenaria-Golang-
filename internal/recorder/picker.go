package recorder

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/selector"
	playwright "github.com/mxschmitt/playwright-go"
)

var ErrPickerCancelled = errors.New("выбор элемента отменён")

var pickerBindingMu sync.Mutex
var pickerResult chan string

func ensurePickerBindings(ctx playwright.BrowserContext) error {
	pickerBindingMu.Lock()
	defer pickerBindingMu.Unlock()
	if pickerResult != nil {
		return nil
	}
	pickerResult = make(chan string, 1)
	if err := ctx.ExposeBinding("pickSelectorDone", func(_ *playwright.BindingSource, args ...any) any {
		if len(args) == 0 {
			return nil
		}
		selector, _ := args[0].(string)
		select {
		case pickerResult <- selector:
		default:
		}
		return nil
	}); err != nil {
		pickerResult = nil
		return err
	}
	return ctx.ExposeBinding("pickSelectorCancel", func(_ *playwright.BindingSource, _ ...any) any {
		select {
		case pickerResult <- "":
		default:
		}
		return nil
	})
}

func drainPickerResults() {
	for {
		select {
		case <-pickerResult:
		default:
			return
		}
	}
}

func PickSelectorOnPage(ctx context.Context, page playwright.Page) (string, error) {
	if page == nil {
		return "", fmt.Errorf("страница браузера недоступна")
	}
	if err := ensurePickerBindings(page.Context()); err != nil {
		return "", fmt.Errorf("picker bindings: %w", err)
	}
	drainPickerResults()

	if _, err := page.Evaluate(selector.PickerInstallScript); err != nil {
		uninstallPicker(page)
		return "", fmt.Errorf("install picker: %w", err)
	}

	timeout := time.NewTimer(5 * time.Minute)
	defer timeout.Stop()

	select {
	case <-ctx.Done():
		uninstallPicker(page)
		return "", ctx.Err()
	case <-timeout.C:
		uninstallPicker(page)
		return "", fmt.Errorf("время выбора элемента истекло")
	case value := <-pickerResult:
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
