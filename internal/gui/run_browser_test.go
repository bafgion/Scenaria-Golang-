package gui

import "testing"

func TestCanReuseLiveBrowserRequiresOptIn(t *testing.T) {
	svc := NewService()
	if svc.canReuseLiveBrowser(RunRequest{ReuseLiveBrowser: false}) {
		t.Fatal("expected false without opt-in")
	}
	if svc.canReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, Workers: 2}) {
		t.Fatal("expected false with workers>1")
	}
	if svc.canReuseLiveBrowser(RunRequest{ReuseLiveBrowser: true, DryRun: true}) {
		t.Fatal("expected false for dry-run")
	}
}
