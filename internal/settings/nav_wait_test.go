package settings

import "testing"

func TestResolveNavWaitUntilProjectOverridesApp(t *testing.T) {
	root := t.TempDir()
	if err := SaveProjectConfig(root, ProjectConfig{NavWaitUntil: "networkidle"}); err != nil {
		t.Fatal(err)
	}
	app := &AppSettings{NavWaitUntil: "load"}
	got := ResolveNavWaitUntil(root, app)
	if got != "networkidle" {
		t.Fatalf("project override: got %q", got)
	}
}

func TestResolveNavWaitUntilFallsBackToApp(t *testing.T) {
	app := &AppSettings{NavWaitUntil: "load"}
	got := ResolveNavWaitUntil("", app)
	if got != "load" {
		t.Fatalf("app fallback: got %q", got)
	}
}
