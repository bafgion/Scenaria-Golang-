package player

import (
	"context"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

// ApplyTestClient applies cookies and localStorage from a test client profile to a page.
func ApplyTestClient(page playwright.Page, client *settings.TestClient) error {
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

type ScenarioSessionRun func(context.Context, *browserSession) (*RunContext, error)

func (e *PlaywrightExecutor) executeWithSession(ctx context.Context, input ScenarioInput, run ScenarioSessionRun) (ScenarioResult, error) {
	if err := ctx.Err(); err != nil {
		return ScenarioResult{}, err
	}

	if e.options.AutoInstall {
		if err := paths.EnsurePlaywrightEngine(e.options.BrowserName); err != nil {
			return ScenarioResult{}, fmt.Errorf("playwright install failed: %w", err)
		}
	} else {
		paths.ConfigurePlaywrightBrowsersForEngine(e.options.BrowserName)
	}

	pw, stopPW, err := startPlaywright(ctx)
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("start playwright: %w", err)
	}
	defer stopPW()

	session, err := newBrowserSession(pw, e.options)
	if err != nil {
		return ScenarioResult{}, err
	}
	stopWatch := session.watchContext(ctx)
	defer stopWatch()
	defer session.close()

	return e.runScenarioOnSession(ctx, session, input, run)
}

func (e *PlaywrightExecutor) runScenarioOnSession(
	ctx context.Context,
	session *browserSession,
	input ScenarioInput,
	run ScenarioSessionRun,
) (ScenarioResult, error) {
	result := ScenarioResult{
		FeaturePath: input.FeaturePath,
		Scenario:    input.ScenarioName,
		Status:      "passed",
	}

	if err := ctx.Err(); err != nil {
		return ScenarioResult{}, err
	}

	failed := false
	if input.TestClient != nil {
		if err := ApplyTestClient(session.page, input.TestClient); err != nil {
			failed = true
			result.Status = "failed"
			result.Message = err.Error()
		}
	}
	if !failed {
		runCtx, err := run(ctx, session)
		if err != nil {
			failed = true
			result.Status = "failed"
			result.Message = err.Error()
			if runCtx != nil {
				if idx := runCtx.FailedLeafStep(); idx >= 0 {
					result.FailedStep = &idx
				}
			}
		}
	}
	if failed {
		result.ScreenshotPNG, result.TraceZIP, result.VideoWebM = captureFailureArtifacts(
			session, input, e.options.TraceDir, e.options.VideoDir,
		)
	}
	return result, nil
}

