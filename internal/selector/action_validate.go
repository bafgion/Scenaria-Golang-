package selector

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

type actionValidationProfile string

const (
	profileVisible   actionValidationProfile = "visible"
	profileHidden    actionValidationProfile = "hidden"
	profileClickable actionValidationProfile = "clickable"
	profileEditable  actionValidationProfile = "editable"
	profileSelect    actionValidationProfile = "select"
	profileCheckbox  actionValidationProfile = "checkbox"
	profileUpload    actionValidationProfile = "upload"
	profileEnabled   actionValidationProfile = "enabled"
	profileDisabled  actionValidationProfile = "disabled"
)

type ActionValidationResult struct {
	OK         bool
	Message    string
	MatchCount int
	Warnings   []string
}

func validationProfile(kind string) actionValidationProfile {
	switch kind {
	case "click", "double-click", "download-click", "scroll-to", "draw-signature":
		return profileClickable
	case "hover":
		return profileVisible
	case "fill", "fill-generated", "clear", "press-in", "prompt-email-code":
		return profileEditable
	case "select":
		return profileSelect
	case "check", "uncheck":
		return profileCheckbox
	case "upload":
		return profileUpload
	case "assert-enabled", "wait-enabled":
		return profileEnabled
	case "assert-disabled", "wait-disabled":
		return profileDisabled
	case "assert-hidden", "wait-hidden":
		return profileHidden
	default:
		return profileVisible
	}
}

func (v Validator) ValidateActionTarget(ctx context.Context, page playwright.Page, action stepdsl.Action, selector string, timeout time.Duration) ActionValidationResult {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return ActionValidationResult{Message: err.Error()}
	}
	if err := ValidateSyntax(selector); err != nil {
		return ActionValidationResult{Message: err.Error()}
	}
	if timeout <= 0 {
		timeout = 8 * time.Second
	}

	matchCount, matchWarn := selectorMatchDiagnostics(page, selector, action.Kind)
	warnings := make([]string, 0)
	if matchWarn != "" {
		warnings = append(warnings, matchWarn)
	}

	profile := validationProfile(action.Kind)
	var locator playwright.Locator
	if action.Kind == "hover" {
		locator = ResolveHoverLocator(page, selector)
		if hoverWarn := hoverAmbiguityWarning(page, selector); hoverWarn != "" {
			warnings = append(warnings, hoverWarn)
		}
	} else {
		locator = ResolveChainedLocator(page, selector)
	}

	var err error
	switch profile {
	case profileHidden:
		err = v.validateHidden(ctx, page, selector, timeout)
	case profileClickable:
		err = validateClickable(ctx, locator, selector, timeout)
	case profileEditable:
		err = validateEditable(ctx, locator, selector, timeout)
	case profileSelect:
		err = validateSelectTarget(ctx, locator, selector, timeout)
	case profileCheckbox:
		err = validateCheckboxTarget(ctx, locator, selector, timeout)
	case profileUpload:
		err = validateUploadTarget(ctx, locator, selector, timeout)
	case profileEnabled:
		err = validateEnabledState(ctx, locator, selector, true, timeout)
	case profileDisabled:
		err = validateEnabledState(ctx, locator, selector, false, timeout)
	default:
		err = v.ValidateVisible(ctx, page, selector, timeout)
	}

	if err != nil {
		return ActionValidationResult{
			Message:    err.Error(),
			MatchCount: matchCount,
			Warnings:   warnings,
		}
	}

	msg := actionSuccessMessage(profile)
	return ActionValidationResult{
		OK:         true,
		Message:    msg,
		MatchCount: matchCount,
		Warnings:   warnings,
	}
}

func actionSuccessMessage(profile actionValidationProfile) string {
	switch profile {
	case profileHidden:
		return "элемент скрыт"
	case profileClickable:
		return "элемент виден и доступен для клика"
	case profileEditable:
		return "поле доступно для ввода"
	case profileSelect:
		return "элемент доступен для выбора"
	case profileCheckbox:
		return "checkbox/radio доступен"
	case profileUpload:
		return "поле загрузки файла найдено"
	case profileEnabled:
		return "элемент включён"
	case profileDisabled:
		return "элемент отключён"
	default:
		return "элемент найден"
	}
}

func selectorMatchDiagnostics(page playwright.Page, selector string, actionKind string) (int, string) {
	if strings.Contains(selector, " >> ") {
		parts := strings.SplitN(selector, " >> ", 2)
		containerSel := strings.TrimSpace(parts[0])
		targetSel := strings.TrimSpace(parts[1])
		containerCount, _ := page.Locator(containerSel).Count()
		if containerCount != 1 {
			return containerCount, fmt.Sprintf("контейнер %q: %d совпадений", containerSel, containerCount)
		}
		targetCount, _ := page.Locator(selector).Count()
		if targetCount != 1 {
			return targetCount, fmt.Sprintf("цель %q: %d совпадений", targetSel, targetCount)
		}
		return 1, ""
	}
	count, _ := page.Locator(selector).Count()
	if count > 1 && actionKind != "hover" {
		return count, fmt.Sprintf("неоднозначный селектор: %d совпадений", count)
	}
	return count, ""
}

func hoverAmbiguityWarning(page playwright.Page, selector string) string {
	for _, candidate := range HoverLocatorCandidates(selector) {
		count, err := page.Locator(candidate).Count()
		if err == nil && count > 1 {
			return fmt.Sprintf("hover candidate %q: %d совпадений", candidate, count)
		}
	}
	return ""
}

func validateClickable(ctx context.Context, locator playwright.Locator, selector string, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	enabled, err := locator.IsEnabled()
	if err != nil {
		return fmt.Errorf("selector %q enabled check: %w", selector, err)
	}
	if !enabled {
		return fmt.Errorf("selector %q is not enabled for click", selector)
	}
	return ctx.Err()
}

func validateEditable(ctx context.Context, locator playwright.Locator, selector string, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	ok, err := locator.Evaluate(`el => {
		if (!el) return false;
		const tag = (el.tagName || '').toUpperCase();
		if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') {
			return !el.disabled && !el.readOnly;
		}
		if (el.isContentEditable) return true;
		const role = (el.getAttribute('role') || '').toLowerCase();
		return ['textbox','combobox','searchbox','spinbutton'].includes(role);
	}`, nil)
	if err != nil {
		return fmt.Errorf("selector %q editable check: %w", selector, err)
	}
	if editable, _ := ok.(bool); !editable {
		return fmt.Errorf("selector %q is not an editable target", selector)
	}
	return ctx.Err()
}

func validateSelectTarget(ctx context.Context, locator playwright.Locator, selector string, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	ok, err := locator.Evaluate(`el => {
		if (!el) return false;
		const tag = (el.tagName || '').toUpperCase();
		if (tag === 'SELECT') return !el.disabled;
		const role = (el.getAttribute('role') || '').toLowerCase();
		return role === 'combobox' || role === 'listbox';
	}`, nil)
	if err != nil {
		return fmt.Errorf("selector %q select check: %w", selector, err)
	}
	if valid, _ := ok.(bool); !valid {
		return fmt.Errorf("selector %q is not a select/combobox target", selector)
	}
	return ctx.Err()
}

func validateCheckboxTarget(ctx context.Context, locator playwright.Locator, selector string, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	ok, err := locator.Evaluate(`el => {
		if (!el) return false;
		const tag = (el.tagName || '').toUpperCase();
		if (tag === 'INPUT') {
			const type = (el.type || '').toLowerCase();
			return type === 'checkbox' || type === 'radio';
		}
		const role = (el.getAttribute('role') || '').toLowerCase();
		return role === 'checkbox' || role === 'radio';
	}`, nil)
	if err != nil {
		return fmt.Errorf("selector %q checkbox check: %w", selector, err)
	}
	if valid, _ := ok.(bool); !valid {
		return fmt.Errorf("selector %q is not a checkbox/radio target", selector)
	}
	return ctx.Err()
}

func validateUploadTarget(ctx context.Context, locator playwright.Locator, selector string, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	ok, err := locator.Evaluate(`el => {
		if (!el) return false;
		return el.tagName === 'INPUT' && (el.type || '').toLowerCase() === 'file';
	}`, nil)
	if err != nil {
		return fmt.Errorf("selector %q upload check: %w", selector, err)
	}
	if valid, _ := ok.(bool); !valid {
		return fmt.Errorf("selector %q is not input[type=file]", selector)
	}
	return ctx.Err()
}

func validateEnabledState(ctx context.Context, locator playwright.Locator, selector string, wantEnabled bool, timeout time.Duration) error {
	if err := waitLocatorVisible(ctx, locator, timeout); err != nil {
		return fmt.Errorf("selector %q is not visible: %w", selector, err)
	}
	enabled, err := locator.IsEnabled()
	if err != nil {
		return fmt.Errorf("selector %q enabled check: %w", selector, err)
	}
	if wantEnabled && !enabled {
		return fmt.Errorf("selector %q is not enabled", selector)
	}
	if !wantEnabled && enabled {
		return fmt.Errorf("selector %q is not disabled", selector)
	}
	return ctx.Err()
}

func waitLocatorVisible(ctx context.Context, locator playwright.Locator, timeout time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
}

func replaySafeStep(ctx context.Context, page playwright.Page, action stepdsl.Action, baseURL string, timeout time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	ms := float64(timeout.Milliseconds())
	switch action.Kind {
	case "goto":
		url := stepdsl.ResolveURL(action.Value1, baseURL)
		_, err := page.Goto(url, playwright.PageGotoOptions{Timeout: playwright.Float(ms)})
		return err
	case "click", "double-click":
		loc := ResolveChainedLocator(page, action.Value1)
		if action.Kind == "double-click" {
			return loc.Dblclick(playwright.LocatorDblclickOptions{Timeout: playwright.Float(ms)})
		}
		return loc.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(ms)})
	case "hover":
		return ResolveHoverLocator(page, action.Value1).Hover(playwright.LocatorHoverOptions{Timeout: playwright.Float(ms)})
	case "scroll-to":
		return ResolveChainedLocator(page, action.Value1).ScrollIntoViewIfNeeded(playwright.LocatorScrollIntoViewIfNeededOptions{Timeout: playwright.Float(ms)})
	default:
		return nil
	}
}
