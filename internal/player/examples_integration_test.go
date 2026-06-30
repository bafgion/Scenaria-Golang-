//go:build integration

package player_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/mxschmitt/playwright-go"
)

func TestBundledExamplesPlaywright(t *testing.T) {
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}

	examplesDir := filepath.Join("..", "..", "examples")
	if _, err := os.Stat(examplesDir); err != nil {
		t.Fatalf("examples dir: %v", err)
	}

	store := scenario.NewFeatureStore()
	files, err := store.Discover(examplesDir)
	if err != nil {
		t.Fatalf("discover examples: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no example feature files found")
	}

	inputs := make([]player.FeatureInput, 0, len(files))
	for _, path := range files {
		feature, loadErr := store.Load(path)
		if loadErr != nil {
			t.Fatalf("load %s: %v", path, loadErr)
		}
		if issues := gherkin.ValidateFeature(feature); len(issues) > 0 {
			t.Fatalf("%s validation: %+v", path, issues)
		}
		inputs = append(inputs, player.FeatureInput{Path: path, Feature: feature})
	}

	plan := player.BuildExecutionPlanWithTestClient(inputs, "", "", nil, "")
	if len(plan.Cases) == 0 {
		t.Fatal("execution plan has no runnable scenarios")
	}
	if len(plan.Cases) < 6 {
		t.Fatalf("expected at least 6 runnable scenarios (outline expands), got %d", len(plan.Cases))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	runner := player.BrowserRunner{
		Executor: player.NewPlaywrightExecutor(player.PlaywrightExecutorOptions{
			BrowserName: "chromium",
			Headless:    true,
		}),
	}
	result, err := runner.Execute(ctx, plan)
	if err != nil {
		t.Fatalf("execute examples: %v", err)
	}

	failures := make([]string, 0)
	for _, scenarioResult := range result.ScenarioResults {
		if scenarioResult.Status == "passed" {
			continue
		}
		failures = append(failures, strings.TrimSpace(scenarioResult.Scenario)+": "+scenarioResult.Message)
	}
	if len(failures) > 0 {
		t.Fatalf("example scenarios failed:\n%s", strings.Join(failures, "\n"))
	}
	if len(result.ScenarioResults) != len(plan.Cases) {
		t.Fatalf("expected %d scenario results, got %d", len(plan.Cases), len(result.ScenarioResults))
	}
}
