package plugin

import "testing"

func TestResolveRunVanessa(t *testing.T) {
	desc := Descriptor{ID: "vanessa", Commands: []string{"va run"}}
	target, err := ResolveRun("/proj", desc, true)
	if err != nil {
		t.Fatal(err)
	}
	if target.Runner != "va" {
		t.Fatalf("runner=%q", target.Runner)
	}
	if len(target.Args) != 4 || target.Args[0] != "run" || target.Args[2] != "/proj" || target.Args[3] != "--dry-run" {
		t.Fatalf("args=%v", target.Args)
	}
}

func TestResolveRunPlaywrightPlugin(t *testing.T) {
	desc := Descriptor{ID: "pw", Commands: []string{"run"}}
	target, err := ResolveRun("/proj", desc, false)
	if err != nil {
		t.Fatal(err)
	}
	if target.Runner != "run" || target.Args[0] != "/proj" {
		t.Fatalf("got %+v", target)
	}
}
