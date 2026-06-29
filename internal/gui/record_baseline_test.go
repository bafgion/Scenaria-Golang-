package gui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecordBaseline_WritesFeature(t *testing.T) {
	dir := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}

	old := cliRunRecord
	cliRunRecord = func(args []string) error {
		if len(args) < 2 || args[0] != "--output" {
			t.Fatalf("unexpected args: %v", args)
		}
		out := args[1]
		content := `# language: ru
Функционал: Demo
  Сценарий: Flow
    Когда открываю "https://example.com"
`
		return os.WriteFile(out, []byte(content), 0o644)
	}
	t.Cleanup(func() { cliRunRecord = old })

	result := svc.RecordBaseline(BaselineRecordRequest{
		Output:       "demo.feature",
		FeatureName:  "Demo",
		ScenarioName: "Flow",
		Steps:        []string{`открываю "https://example.com"`},
	})
	if result.Error != "" {
		t.Fatalf("unexpected error: %s", result.Error)
	}
	path := filepath.Join(dir, "demo.feature")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Функционал: Demo") {
		t.Fatalf("unexpected content: %s", data)
	}
}

func TestAppendRunPluginArgs(t *testing.T) {
	args := appendRunPluginArgs([]string{"/proj", "--dry-run"}, PluginRunRequest{Tag: "@smoke", Scenario: "Login"})
	want := []string{"/proj", "--dry-run", "--tag", "@smoke", "--scenario", "Login"}
	if len(args) != len(want) {
		t.Fatalf("got %v", args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, args[i], want[i])
		}
	}
}
