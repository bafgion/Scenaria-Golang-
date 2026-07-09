package gui

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
)

func emptyPlan() player.ExecutionPlan {
	return player.ExecutionPlan{}
}

func closeBrowserPlan() player.ExecutionPlan {
	return player.ExecutionPlan{Cases: []player.RunCase{{
		Steps: []gherkin.Step{{Keyword: "И", Text: "закрываю браузер", Line: 1}},
	}}}
}

func multiScenarioPlan() player.ExecutionPlan {
	return player.ExecutionPlan{Cases: []player.RunCase{
		{Name: "A"},
		{Name: "B"},
	}}
}

func TestCanReuseLiveBrowserRequiresOptIn(t *testing.T) {
	runSvc := &RunService{}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: false}, true, emptyPlan()) {
		t.Fatal("expected false without opt-in")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Workers: 2}, true, emptyPlan()) {
		t.Fatal("expected false with workers>1")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, DryRun: true}, true, emptyPlan()) {
		t.Fatal("expected false for dry-run")
	}
}

func TestCanReuseLiveBrowserAllowsSequentialChromium(t *testing.T) {
	runSvc := &RunService{}
	if !runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, true, emptyPlan()) {
		t.Fatal("expected true with live browser and opt-in")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "firefox", Workers: 1}, true, emptyPlan()) {
		t.Fatal("expected false for non-chromium live browser reuse")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, false, emptyPlan()) {
		t.Fatal("expected false without live browser")
	}
}

func TestCanReuseLiveBrowserDisabledForMultiScenario(t *testing.T) {
	runSvc := &RunService{}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, true, multiScenarioPlan()) {
		t.Fatal("expected false for multi-scenario plan")
	}
}

func TestCanReuseLiveBrowserDisabledWhenPlanClosesBrowser(t *testing.T) {
	runSvc := &RunService{}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, true, closeBrowserPlan()) {
		t.Fatal("expected false when plan contains close-browser")
	}
}
