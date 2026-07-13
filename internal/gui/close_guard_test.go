package gui

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestServiceCloseGuardReasonsDetectDirtyTabs(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(feature, []byte("Feature: Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	svc.projectPath = root
	svc.settingsStore = settings.NewStore(filepath.Join(root, "settings.json"))
	svc.settingsService = NewSettingsService(svc.settingsStore)

	if err := svc.SaveFeatureDraft(feature, "Feature: Demo\n"); err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveSettings(AppSettingsDTO{OpenTabs: []string{feature}}); err != nil {
		t.Fatal(err)
	}
	reasons := svc.CloseGuardReasons()
	if len(reasons) == 0 {
		t.Fatal("expected close guard reasons")
	}
	found := false
	for _, reason := range reasons {
		if reason == "unsaved_tabs" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unsaved tabs reason, got %v", reasons)
	}
	if err := svc.ClearFeatureDraft(feature); err != nil {
		t.Fatal(err)
	}
	if svc.HasDirtyTabs() {
		t.Fatal("expected dirty tabs to clear after draft removal")
	}
}

func TestServiceCloseGuardReasonsDetectLiveDirtyTabs(t *testing.T) {
	svc := NewService()
	svc.settingsStore = settings.NewStore(filepath.Join(t.TempDir(), "settings.json"))
	svc.settingsService = NewSettingsService(svc.settingsStore)
	if svc.HasDirtyTabs() {
		t.Fatal("expected no dirty tabs initially")
	}
	svc.UpdateDirtyTabsState(true)
	if !svc.HasDirtyTabs() {
		t.Fatal("expected live dirty tab state to trigger close guard")
	}
	reasons := svc.CloseGuardReasons()
	found := false
	for _, reason := range reasons {
		if reason == "unsaved_tabs" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unsaved tabs reason, got %v", reasons)
	}
	svc.UpdateDirtyTabsState(false)
	if svc.HasDirtyTabs() {
		t.Fatal("expected live dirty tab state to clear")
	}
}

func TestRunServiceHasActiveRun(t *testing.T) {
	svc := NewService()
	if svc.HasActiveRun() {
		t.Fatal("expected no active run initially")
	}
	begin, err := svc.runner().TryBegin(1, RunRequest{}, nil, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer begin.Cancel()
	if !svc.HasActiveRun() {
		t.Fatal("expected active run to be detected")
	}
	svc.runner().Finish(begin.Gen)
	if svc.HasActiveRun() {
		t.Fatal("expected active run to clear after finish")
	}
}
