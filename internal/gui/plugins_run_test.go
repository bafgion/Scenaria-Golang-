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
		PlatformExe:       `C:\1cv8\bin\1cv8.exe`,
		EPFPath:           `C:\vanessa\vanessa-automation.epf`,
		IBConnection:      `Srvr="localhost";Ref="base";`,
		ReportAllure:      true,
		VaDir:             `features\smoke`,
		VaFiles:           `a.feature,b.feature`,
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
		"--platform-exe", `C:\1cv8\bin\1cv8.exe`,
		"--epf", `C:\vanessa\vanessa-automation.epf`,
		"--ib", `Srvr="localhost";Ref="base";`,
		"--allure",
		"--dir", `features\smoke`,
		"--files", `a.feature,b.feature`,
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
