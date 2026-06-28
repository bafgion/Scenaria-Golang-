package player

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

const otpKeyboardDelayMs = 80

var (
	emailRE = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	codeNormalizeRE = regexp.MustCompile(`[\s\-–—]`)
	inputEventsJS = `(el) => {
 el.dispatchEvent(new Event('input', { bubbles: true }));
 el.dispatchEvent(new Event('change', { bubbles: true }));
}`
)

// ResolveEmailForCode finds display email from step, variables, or page content.
func (c *RunContext) ResolveEmailForCode(stepEmail string, priorSteps []gherkin.Step) (string, error) {
	if stepEmail = strings.TrimSpace(stepEmail); stepEmail != "" {
		resolved, err := c.ResolveText(stepEmail)
		if err != nil {
			return "", err
		}
		return resolved, nil
	}
	for _, key := range []string{"email", "user_email", "login_email"} {
		if value := strings.TrimSpace(c.Variables[key]); value != "" {
			return value, nil
		}
	}
	for _, step := range priorSteps {
		action, err := stepdsl.Parse(step)
		if err != nil || action.Kind != "fill" {
			continue
		}
		if strings.Contains(action.Value1, "@") {
			return action.Value1, nil
		}
	}
	if c.page != nil {
		content, err := c.page.Content()
		if err == nil {
			if match := emailRE.FindString(content); match != "" {
				return match, nil
			}
		}
	}
	return "", nil
}

func normalizeVerificationCode(code string) string {
	return codeNormalizeRE.ReplaceAllString(strings.TrimSpace(code), "")
}

func otpFieldsVisible(locator playwright.Locator, minCount int) bool {
	count, err := locator.Count()
	if err != nil || count < minCount {
		return false
	}
	visible, err := locator.First().IsVisible()
	return err == nil && visible
}

func segmentedInputValues(locator playwright.Locator, fieldCount int) string {
	values := make([]string, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		value, err := locator.Nth(i).InputValue()
		if err != nil {
			values = append(values, "")
			continue
		}
		values = append(values, value)
	}
	return strings.Join(values, "")
}

func otpAutoSubmitted(locator playwright.Locator, fieldCount int, expected string) bool {
	if expected == "" {
		return false
	}
	count, err := locator.Count()
	if err != nil {
		return false
	}
	if count == 0 {
		return true
	}
	filled := segmentedInputValues(locator, min(fieldCount, count))
	if filled == expected {
		return true
	}
	for i := 0; i < min(count, fieldCount); i++ {
		visible, err := locator.Nth(i).IsVisible()
		if err == nil && visible {
			return false
		}
	}
	return true
}

func fillVerificationCode(page playwright.Page, locator playwright.Locator, code string, digits int, method string) error {
	value := normalizeVerificationCode(code)
	if value == "" {
		return fmt.Errorf("verification code is empty")
	}
	method = strings.ToLower(strings.TrimSpace(method))
	if method == "" {
		method = "fill"
	}

	count, err := locator.Count()
	if err != nil || count == 0 {
		return fmt.Errorf("code fields not found")
	}

	segmented := digits > 1 || count > 1
	if segmented {
		fieldCount := digits
		if fieldCount < 2 {
			fieldCount = count
		}
		if count < fieldCount {
			return fmt.Errorf("found %d code field(s), need %d", count, fieldCount)
		}
		if len(value) > fieldCount {
			value = value[:fieldCount]
		}
		return fillSegmentedCode(page, locator, value, fieldCount, method)
	}

	target := locator.First()
	if err := target.Click(); err != nil {
		return err
	}
	if method == "keyboard" {
		return page.Keyboard().Type(value, playwright.KeyboardTypeOptions{
			Delay: playwright.Float(otpKeyboardDelayMs),
		})
	}
	if err := target.Fill(value); err != nil {
		return err
	}
	_, _ = target.Evaluate(inputEventsJS, nil)
	current, err := target.InputValue()
	if err == nil && current != value && method != "keyboard" {
		if err := target.Click(); err != nil {
			return err
		}
		return page.Keyboard().Type(value, playwright.KeyboardTypeOptions{
			Delay: playwright.Float(otpKeyboardDelayMs),
		})
	}
	return tryAutoSubmit(page, locator)
}

func fillSegmentedCode(page playwright.Page, locator playwright.Locator, value string, fieldCount int, method string) error {
	tryFill := func(stage string, fillFn func() error) error {
		if err := fillFn(); err != nil {
			if otpAutoSubmitted(locator, fieldCount, value) {
				return nil
			}
			return err
		}
		if segmentedInputValues(locator, fieldCount) == value || otpAutoSubmitted(locator, fieldCount, value) {
			return nil
		}
		return fmt.Errorf("segmented code fill failed at %s", stage)
	}

	if method == "fill" {
		return tryFill("fill", func() error {
			for i, ch := range value {
				cell := locator.Nth(i)
				if err := cell.Click(); err != nil {
					return err
				}
				if err := cell.Fill(string(ch)); err != nil {
					return err
				}
				_, _ = cell.Evaluate(inputEventsJS, nil)
			}
			return nil
		})
	}

	if err := tryFill("keyboard-batch", func() error {
		first := locator.Nth(0)
		if err := first.Click(); err != nil {
			return err
		}
		time.Sleep(120 * time.Millisecond)
		return page.Keyboard().Type(value, playwright.KeyboardTypeOptions{
			Delay: playwright.Float(otpKeyboardDelayMs),
		})
	}); err == nil {
		return nil
	}

	if !otpFieldsVisible(locator, 1) {
		return nil
	}

	return tryFill("keyboard-char", func() error {
		for i, ch := range value {
			cell := locator.Nth(i)
			if err := cell.Click(); err != nil {
				return err
			}
			time.Sleep(60 * time.Millisecond)
			if err := page.Keyboard().Type(string(ch), playwright.KeyboardTypeOptions{
				Delay: playwright.Float(otpKeyboardDelayMs),
			}); err != nil {
				return err
			}
			time.Sleep(40 * time.Millisecond)
		}
		return nil
	})
}

func runEmailCode(page playwright.Page, action stepdsl.Action, runCtx *RunContext) error {
	locator := page.Locator(action.Value2)
	if err := locator.First().WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60000),
	}); err != nil {
		return fmt.Errorf("email code field not visible: %w", err)
	}
	if !otpFieldsVisible(locator, 1) {
		return nil
	}
	code, err := runCtx.EmailCode()
	if err != nil {
		return err
	}
	_, _ = runCtx.ResolveEmailForCode(action.Value1, runCtx.PriorSteps())
	digits := action.IntVal
	if digits <= 0 {
		digits = len(normalizeVerificationCode(code))
	}
	method := action.Value3
	if method == "" {
		method = "fill"
	}
	if err := fillVerificationCode(page, locator, code, digits, method); err != nil {
		return err
	}
	return tryAutoSubmit(page, locator)
}

func tryAutoSubmit(page playwright.Page, locator playwright.Locator) error {
	submit := page.Locator(`button[type="submit"], input[type="submit"], button:has-text("Подтвердить"), button:has-text("Войти")`)
	count, err := submit.Count()
	if err != nil || count == 0 {
		return nil
	}
	first := submit.First()
	visible, _ := first.IsVisible()
	if !visible {
		return nil
	}
	return first.Click()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
