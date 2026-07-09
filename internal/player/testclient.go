package player

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
		script, err := localStorageInitScript(client.BaseURL, client.LocalStorage)
		if err != nil {
			return err
		}
		if err := page.Context().AddInitScript(playwright.Script{Content: playwright.String(script)}); err != nil {
			return fmt.Errorf("register test client local storage: %w", err)
		}
		if sameOrigin(page.URL(), client.BaseURL) || strings.TrimSpace(client.BaseURL) == "" {
			if _, err := page.Evaluate(script); err != nil {
				return fmt.Errorf("apply test client local storage: %w", err)
			}
		}
	}
	return nil
}

func localStorageInitScript(baseURL string, values map[string]string) (string, error) {
	payload, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("encode test client local storage: %w", err)
	}
	origin := originFromURL(baseURL)
	originJSON, err := json.Marshal(origin)
	if err != nil {
		return "", fmt.Errorf("encode test client origin: %w", err)
	}
	return fmt.Sprintf(`(() => {
  const expectedOrigin = %s;
  if (expectedOrigin && window.location.origin !== expectedOrigin) return;
  const values = %s;
  for (const [key, value] of Object.entries(values)) {
    localStorage.setItem(key, value);
  }
})()`, string(originJSON), string(payload)), nil
}

func sameOrigin(currentURL, baseURL string) bool {
	current := originFromURL(currentURL)
	target := originFromURL(baseURL)
	return current != "" && target != "" && current == target
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
	defer func() {
		if e.options.CloseAfterRun {
			stopPW()
		}
	}()

	session, err := newBrowserSession(pw, e.options)
	if err != nil {
		return ScenarioResult{}, err
	}
	stopWatch := session.watchContext(ctx)
	defer stopWatch()
	if e.options.CloseAfterRun {
		defer session.close()
	}

	return e.runScenarioOnSession(ctx, session, input, run)
}

func (e *PlaywrightExecutor) runScenarioOnSession(
	ctx context.Context,
	session *browserSession,
	input ScenarioInput,
	run ScenarioSessionRun,
) (ScenarioResult, error) {
	started := time.Now()
	result := ScenarioResult{
		FeaturePath:  input.FeaturePath,
		Scenario:     input.ScenarioName,
		CaseID:       input.CaseID,
		ExampleIndex: input.ExampleIndex,
		Status:       "passed",
	}

	if err := ctx.Err(); err != nil {
		return ScenarioResult{}, err
	}

	failed := false
	var runCtx *RunContext
	if input.TestClient != nil {
		page, err := session.currentPage()
		if err != nil {
			failed = true
			result.Status = "failed"
			result.Message = UserFacingBrowserError(err)
			result.FailedStep = failedStepIndex(0)
		} else if err := ApplyTestClient(page, input.TestClient); err != nil {
			failed = true
			result.Status = "failed"
			result.Message = err.Error()
		}
	}
	if !failed {
		var err error
		runCtx, err = run(ctx, session)
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
		result.ScreenshotPath, result.TraceZIPPath, result.VideoWebMPath = captureFailureArtifacts(
			session, input, e.options.TraceDir, e.options.VideoDir,
		)
		restartTraceRecording(session)
	} else {
		discardTraceRecording(session)
	}
	if runCtx != nil {
		result.StepRecords = runCtx.StepRecords()
	}
	result.DurationMS = time.Since(started).Milliseconds()
	return result, nil
}
