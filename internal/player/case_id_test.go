package player

import "testing"

func TestBuildCaseIDIncludesExampleIndex(t *testing.T) {
	got := BuildCaseID("features/login.feature", "Login", 2)
	want := "features/login.feature::Login#2"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBuildCaseIDOmitExampleZero(t *testing.T) {
	got := BuildCaseID("a.feature", "Smoke", 0)
	if got != "a.feature::Smoke" {
		t.Fatalf("got %q", got)
	}
}
