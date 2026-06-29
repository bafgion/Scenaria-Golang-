package player

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
	playwright "github.com/mxschmitt/playwright-go"
)

func ResolveRunHTTPCredentials(baseURL string, plan ExecutionPlan, cfg *settings.AppSettings) *playwright.HttpCredentials {
	return httpauth.PlaywrightHTTPCredentials(ResolveAuthURL(baseURL, plan), cfg)
}

func ResolveAuthURL(baseURL string, plan ExecutionPlan) string {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL != "" {
		return baseURL
	}
	for _, runCase := range plan.Cases {
		if runCase.TestClient != nil && strings.TrimSpace(runCase.TestClient.BaseURL) != "" {
			return runCase.TestClient.BaseURL
		}
	}
	return FirstNavigationURL(plan)
}

func FirstNavigationURL(plan ExecutionPlan) string {
	for _, runCase := range plan.Cases {
		for _, step := range runCase.Steps {
			action, err := stepdsl.Parse(step)
			if err != nil || action.Kind != "goto" {
				continue
			}
			if url := strings.TrimSpace(action.Value1); url != "" {
				return url
			}
		}
	}
	return ""
}
