package gui

import "testing"

func TestAppendVanessaArgs(t *testing.T) {
	args := appendVanessaArgs([]string{"run", "--project", "/p"}, PluginRunRequest{
		Tag:         "@smoke",
		ExcludeTags: []string{"@wip", ""},
		Scenario:    "Login",
	})
	want := []string{"run", "--project", "/p", "--tag", "@smoke", "--exclude-tag", "@wip", "--scenario", "Login"}
	if len(args) != len(want) {
		t.Fatalf("got %v", args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("args[%d]=%q want %q", i, args[i], want[i])
		}
	}
}
