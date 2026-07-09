package gui

import "testing"

func TestCanReuseLiveBrowserRequiresOptIn(t *testing.T) {
	runSvc := &RunService{}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: false}, true) {
		t.Fatal("expected false without opt-in")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Workers: 2}, true) {
		t.Fatal("expected false with workers>1")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, DryRun: true}, true) {
		t.Fatal("expected false for dry-run")
	}
}

func TestCanReuseLiveBrowserAllowsSequentialChromium(t *testing.T) {
	runSvc := &RunService{}
	if !runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, true) {
		t.Fatal("expected true with live browser and opt-in")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "firefox", Workers: 1}, true) {
		t.Fatal("expected false for non-chromium live browser reuse")
	}
	if runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, false) {
		t.Fatal("expected false without live browser")
	}
}

func TestCanReuseLiveBrowserAllowsMultiScenarioPlan(t *testing.T) {
	runSvc := &RunService{}
	if !runSvc.CanReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Browser: "chromium", Workers: 1}, true) {
		t.Fatal("expected multi-scenario live reuse to stay opt-in enabled; per-scenario close-browser decides reuse")
	}
}
