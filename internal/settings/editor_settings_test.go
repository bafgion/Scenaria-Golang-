package settings

import (
	"path/filepath"
	"testing"
)

func TestNormalizeEditorSettingsDefaults(t *testing.T) {
	got := NormalizeEditorSettings(EditorSettings{})
	def := DefaultEditorSettings()
	if got.FontSize != def.FontSize {
		t.Fatalf("fontSize=%d want %d", got.FontSize, def.FontSize)
	}
	if got.Theme != EditorThemeDark {
		t.Fatalf("theme=%q", got.Theme)
	}
}

func TestNormalizeEditorSettingsOverrides(t *testing.T) {
	minimap := false
	got := NormalizeEditorSettings(EditorSettings{
		FontSize: 18,
		WordWrap: EditorWordWrapOff,
		Minimap:  &minimap,
		Theme:    EditorThemeLight,
		TabSize:  2,
	})
	if got.FontSize != 18 || got.WordWrap != EditorWordWrapOff || got.Theme != EditorThemeLight {
		t.Fatalf("unexpected %+v", got)
	}
	if got.Minimap == nil || *got.Minimap {
		t.Fatalf("minimap should be false")
	}
	if got.TabSize != 2 {
		t.Fatalf("tabSize=%d", got.TabSize)
	}
}

func TestNormalizeEditorSettingsScenarioHints(t *testing.T) {
	disabled := false
	got := NormalizeEditorSettings(EditorSettings{
		ScenarioHintsEnabled:       &disabled,
		ScenarioHintsAutoFixOnSave: boolPtr(true),
	})
	if got.ScenarioHintsEnabled == nil || *got.ScenarioHintsEnabled {
		t.Fatalf("scenario hints should be disabled")
	}
	if got.ScenarioHintsAutoFixOnSave == nil || !*got.ScenarioHintsAutoFixOnSave {
		t.Fatalf("auto fix on save should be true")
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func TestAppSettingsPersistsEditor(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")
	minimap := true
	input := &AppSettings{
		Browser: "chromium",
		Editor: EditorSettings{
			FontSize: 17,
			Minimap:  &minimap,
			Theme:    EditorThemeLight,
		},
	}
	if err := SaveAppSettings(path, input); err != nil {
		t.Fatal(err)
	}
	got, err := LoadAppSettings(path)
	if err != nil {
		t.Fatal(err)
	}
	normalized := NormalizeEditorSettings(got.Editor)
	if normalized.FontSize != 17 || normalized.Theme != EditorThemeLight {
		t.Fatalf("editor not persisted: %+v", normalized)
	}
}
