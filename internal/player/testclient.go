package player

import (
	"context"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

func applyTestClient(page playwright.Page, client *settings.TestClient) error {
	if client == nil {
		return nil
	}
	if len(client.Cookies) > 0 {
		cookies := make([]playwright.OptionalCookie, 0, len(client.Cookies))
		for _, cookie := range client.Cookies {
			cookies = append(cookies, playwright.OptionalCookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   playwright.String(cookie.Domain),
				Path:     playwright.String(cookie.Path),
				HttpOnly: playwright.Bool(cookie.HTTPOnly),
				Secure:   playwright.Bool(cookie.Secure),
			})
		}
		if err := page.Context().AddCookies(cookies); err != nil {
			return fmt.Errorf("apply test client cookies: %w", err)
		}
	}
	if len(client.LocalStorage) > 0 {
		script := "(() => {"
		for key, value := range client.LocalStorage {
			script += fmt.Sprintf("localStorage.setItem(%q, %q);", key, value)
		}
		script += "})()"
		if _, err := page.Evaluate(script); err != nil {
			return fmt.Errorf("apply test client local storage: %w", err)
		}
	}
	return nil
}

func loadTestClientForFeature(projectRoot, name string) (*settings.TestClient, error) {
	if name == "" || projectRoot == "" {
		return nil, nil
	}
	path, err := settings.TestClientPath(projectRoot, name)
	if err != nil {
		return nil, err
	}
	return settings.LoadTestClient(path)
}

func (e *PlaywrightExecutor) executeWithSession(ctx context.Context, input ScenarioInput, run func(context.Context, *browserSession) error) (ScenarioResult, error) {
	result := ScenarioResult{
		FeaturePath: input.FeaturePath,
		Scenario:    input.ScenarioName,
		Status:      "passed",
	}

	if e.options.AutoInstall {
		if err := playwright.Install(); err != nil {
			return ScenarioResult{}, fmt.Errorf("playwright install failed: %w", err)
		}
	}

	paths.ConfigurePlaywrightBrowsers()

	pw, err := playwright.Run()
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("start playwright: %w", err)
	}
	defer func() { _ = pw.Stop() }()

	session, err := newBrowserSession(pw, e.options)
	if err != nil {
		return ScenarioResult{}, err
	}

	failed := false
	if input.TestClient != nil {
		if err := applyTestClient(session.page, input.TestClient); err != nil {
			failed = true
			result.Status = "failed"
			result.Message = err.Error()
		}
	}
	if !failed {
		if err := run(ctx, session); err != nil {
			failed = true
			result.Status = "failed"
			result.Message = err.Error()
		}
	}
	if failed {
		result.ScreenshotPNG, result.TraceZIP, result.VideoWebM = captureFailureArtifacts(
			session, input, e.options.TraceDir, e.options.VideoDir,
		)
	}
	session.close()
	return result, nil
}
