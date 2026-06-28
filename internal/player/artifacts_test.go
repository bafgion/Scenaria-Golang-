package player

import (
	"path/filepath"
	"testing"
)

func TestArtifactBaseName(t *testing.T) {
	got := artifactBaseName(ScenarioInput{
		FeaturePath:  filepath.Join("proj", "login.feature"),
		ScenarioName: "Успешный вход",
	})
	if got == "" || got == "scenario" {
		t.Fatalf("unexpected artifact name: %q", got)
	}
}
