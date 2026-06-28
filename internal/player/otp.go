package player

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

var emailRE = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

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

func fillVerificationCode(page playwright.Page, locator playwright.Locator, code string, digits int, method string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return fmt.Errorf("verification code is empty")
	}
	if digits <= 0 {
		digits = len(code)
	}
	if digits > 1 {
		inputs, err := locator.Locator("input").All()
		if err == nil && len(inputs) >= digits {
			for i := 0; i < digits; i++ {
				digit := string(code[i])
				field := inputs[i]
				if method == "keyboard" {
					if err := field.Click(); err != nil {
						return err
					}
					if err := page.Keyboard().Type(digit); err != nil {
						return err
					}
				} else if err := field.Fill(digit); err != nil {
					return err
				}
			}
			return nil
		}
	}
	if method == "keyboard" {
		if err := locator.First().Click(); err != nil {
			return err
		}
		return page.Keyboard().Type(code)
	}
	if err := locator.First().Fill(code); err != nil {
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
