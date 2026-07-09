package gui_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gui"
)

// Facade smoke tests: gui.Service should delegate domain work to extracted services,
// not reimplement filesystem / project index logic inline.
func TestServiceFacadeDelegatesProjectAndFileOps(t *testing.T) {
	root := t.TempDir()
	featurePath := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: Демо\n\n  Сценарий: A\n    Допустим пустой шаг\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := gui.NewService()
	info, err := svc.OpenProject(root)
	if err != nil {
		t.Fatalf("OpenProject: %v", err)
	}
	if len(info.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(info.Features))
	}

	text, err := svc.ReadFeature(featurePath)
	if err != nil {
		t.Fatalf("ReadFeature: %v", err)
	}
	if text == "" {
		t.Fatal("expected feature text")
	}

	if err := svc.SaveFeature(featurePath, text+"\n"); err != nil {
		t.Fatalf("SaveFeature: %v", err)
	}

	results, err := svc.ListRunResults(5)
	if err != nil {
		t.Fatalf("ListRunResults: %v", err)
	}
	_ = results

	settings, err := svc.LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	_ = settings
}

func TestServiceFacadeEditorAnalysis(t *testing.T) {
	svc := gui.NewService()
	issues := svc.ValidateFeature("Функционал: X\n\n  Сценарий: A\n    Неизвестный шаг\n")
	if len(issues) == 0 {
		t.Fatal("expected validation issues")
	}
	steps := svc.ParseEditorSteps("Функционал: X\n\n  Сценарий: A\n    Допустим пустой шаг\n")
	if len(steps) == 0 {
		t.Fatal("expected parsed steps")
	}
}
