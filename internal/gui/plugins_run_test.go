package gui

import "testing"

func TestAppendVanessaArgs(t *testing.T) {
	args := appendVanessaArgs([]string{"run", "--project", "/p"}, PluginRunRequest{
		Tag:               "@smoke",
		ExcludeTags:       []string{"@wip", ""},
		Scenario:          "Login",
		RerunFailedRunDir: "/runs/run-1",
		InstallEPF:        true,
		EPFURL:            "https://example.com/va.epf",
		EPFDest:           `C:\va.epf`,
	})
	want := []string{
		"run", "--project", "/p",
		"--tag", "@smoke",
		"--exclude-tag", "@wip",
		"--scenario", "Login",
		"--rerun-failed", "/runs/run-1",
		"--epf-install",
		"--epf-url", "https://example.com/va.epf",
		"--epf-dest", `C:\va.epf`,
	}
	if len(args) != len(want) {
		t.Fatalf("got %v", args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("args[%d]=%q want %q", i, args[i], want[i])
		}
	}
}
