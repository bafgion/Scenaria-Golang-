package gui

import "testing"

func TestCloneRunRequestCopiesTargetsAndVars(t *testing.T) {
	orig := RunRequest{
		Tag:      "@smoke",
		Scenario: "Login",
		Targets:  []string{"C:/proj/a.feature", "C:/proj/b.feature"},
		Vars: map[string]string{
			"BASE_URL": "https://example.com",
		},
	}

	cloned := cloneRunRequest(orig)

	orig.Targets[0] = "C:/proj/changed.feature"
	orig.Vars["BASE_URL"] = "https://changed.example.com"

	if cloned.Targets[0] != "C:/proj/a.feature" || cloned.Targets[1] != "C:/proj/b.feature" {
		t.Fatalf("targets were not copied: %#v", cloned.Targets)
	}
	if cloned.Vars["BASE_URL"] != "https://example.com" {
		t.Fatalf("vars were not copied: %#v", cloned.Vars)
	}
	if cloned.Tag != orig.Tag || cloned.Scenario != orig.Scenario {
		t.Fatalf("non-collection fields should be preserved: %#v", cloned)
	}
}
