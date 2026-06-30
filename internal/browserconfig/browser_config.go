package browserconfig

import (
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

var chromiumHeadedLaunchArgs = []string{
	"--disable-dev-shm-usage",
	"--no-first-run",
	"--no-default-browser-check",
	"--disable-extensions",
	"--start-maximized",
}

func NormalizeEngine(value string) string {
	engine := strings.ToLower(strings.TrimSpace(value))
	if engine == "" {
		return "chromium"
	}
	return engine
}

func LaunchOptions(browserName string, headless bool, slowMo float64) playwright.BrowserTypeLaunchOptions {
	opts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	}
	if slowMo > 0 {
		opts.SlowMo = playwright.Float(slowMo)
	}
	if NormalizeEngine(browserName) == "chromium" && !headless {
		opts.Args = chromiumHeadedLaunchArgs
	}
	return opts
}

func NewContextOptions(headless bool, httpCreds *playwright.HttpCredentials) playwright.BrowserNewContextOptions {
	opts := playwright.BrowserNewContextOptions{}
	if !headless {
		opts.NoViewport = playwright.Bool(true)
	}
	if httpCreds != nil {
		opts.HttpCredentials = httpCreds
	}
	return opts
}
