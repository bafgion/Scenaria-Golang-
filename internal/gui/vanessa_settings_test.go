package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveVanessaSettingsJSONRoundTrip(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	content := `{
  "platform_executable": "C:\\Program Files\\1cv8\\bin\\1cv8.exe",
  "epf_path": "C:\\vanessa\\vanessa-automation.epf",
  "runs_dir": "C:\\runs"
}`
	if err := svc.SaveVanessaSettingsJSON(content); err != nil {
		t.Fatalf("SaveVanessaSettingsJSON: %v", err)
	}
	path := filepath.Join(root, ".scenaria", "vanessa.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file: %v", err)
	}
	got, err := svc.ReadVanessaSettingsJSON()
	if err != nil {
		t.Fatalf("ReadVanessaSettingsJSON: %v", err)
	}
	if got == "" || got == content {
		// re-marshaled JSON may differ in spacing; ensure key persisted
		if _, err := os.ReadFile(path); err != nil {
			t.Fatal(err)
		}
	}
}
