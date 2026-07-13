package plugin

import (
	"reflect"
	"testing"
)

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

func TestResolveRunLegacyCommandPreservesQuotedArguments(t *testing.T) {
	desc := Descriptor{ID: "pw", Commands: []string{`run --dir "C:\Program Files\Scenaria Tests" --scenario "Scenario With Spaces"`}}
	target, err := ResolveRun(`C:\Project Root`, desc, false)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{`C:\Project Root`, "--dir", `C:\Program Files\Scenaria Tests`, "--scenario", "Scenario With Spaces"}
	if target.Runner != "run" || !reflect.DeepEqual(target.Args, want) {
		t.Fatalf("got %+v want runner=run args=%#v", target, want)
	}
}

func TestResolveRunStructuredCommandPreservesEmptyAndSpacedArgs(t *testing.T) {
	desc := Descriptor{
		ID: "pw",
		StructuredRuns: []CommandSpec{{
			Runner: "run",
			Args:   []string{"--label", "value with spaces", "--empty", ""},
		}},
	}
	target, err := ResolveRun("/proj", desc, true)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/proj", "--label", "value with spaces", "--empty", "", "--dry-run"}
	if target.Runner != "run" || !reflect.DeepEqual(target.Args, want) {
		t.Fatalf("got %+v want runner=run args=%#v", target, want)
	}
}

func TestResolveRunLegacyCommandRejectsUnterminatedQuote(t *testing.T) {
	desc := Descriptor{ID: "pw", Commands: []string{`run --scenario "missing end`}}
	if _, err := ResolveRun("/proj", desc, false); err == nil {
		t.Fatal("expected unterminated quote error")
	}
}

func TestResolveRunStructuredVACommandUsesExplicitArgs(t *testing.T) {
	desc := Descriptor{
		ID: "custom-va",
		CommandSpecs: []CommandSpec{{
			Runner: "va",
			Args:   []string{"run", "--project", "/explicit", "--scenario", "Smoke"},
		}},
	}
	target, err := ResolveRun("/proj", desc, true)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"run", "--project", "/explicit", "--scenario", "Smoke", "--dry-run"}
	if target.Runner != "va" || !reflect.DeepEqual(target.Args, want) {
		t.Fatalf("got %+v want runner=va args=%#v", target, want)
	}
}
