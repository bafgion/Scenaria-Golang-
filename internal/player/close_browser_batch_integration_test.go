//go:build integration

package player

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func closeBrowserCase(name string) RunCase {
	return RunCase{
		FeaturePath: "close-browser.feature",
		Name:        name,
		CaseID:      BuildCaseID("close-browser.feature", name, 0),
		Steps: []gherkin.Step{
			{Keyword: "Допустим", Text: `открыт "https://example.com"`, Line: 1},
			{Keyword: "И", Text: "закрываю браузер", Line: 2},
		},
	}
}

func openPageCase(name string) RunCase {
	return RunCase{
		FeaturePath: "follow-up.feature",
		Name:        name,
		CaseID:      BuildCaseID("follow-up.feature", name, 0),
		Steps: []gherkin.Step{
			{Keyword: "Допустим", Text: `открыт "https://example.com"`, Line: 1},
			{Keyword: "Тогда", Text: `вижу "h1"`, Line: 2},
		},
	}
}

func TestCloseBrowserMarksRemainingScenariosNotStarted(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			closeBrowserCase("First closes browser"),
			openPageCase("Second should not start"),
		},
	}
	runner := BrowserRunner{
		Executor: NewPlaywrightExecutor(PlaywrightExecutorOptions{
			BrowserName:   "chromium",
			Headless:      true,
			CloseAfterRun: true,
		}),
		ParallelWorkers: 1,
	}
	ctx := WithContinueOnFail(context.Background(), true)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	result, err := runner.Execute(ctx, plan)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(result.ScenarioResults) != 2 {
		t.Fatalf("expected 2 scenario results, got %d", len(result.ScenarioResults))
	}
	if result.ScenarioResults[0].Status != "passed" {
		t.Fatalf("first scenario status = %q message = %q", result.ScenarioResults[0].Status, result.ScenarioResults[0].Message)
	}
	second := result.ScenarioResults[1]
	if second.Status != "not-started" {
		t.Fatalf("second scenario status = %q, want not-started (message: %q)", second.Status, second.Message)
	}
	if second.Message != MsgScenarioNotStartedBrowserClosed {
		t.Fatalf("second scenario message = %q", second.Message)
	}
	if IsBrowserSessionClosed(errors.New(second.Message)) || second.Message == "browser session is closed" {
		t.Fatal("second scenario should not expose low-level browser session error")
	}
}
