package settings

import (
	"path/filepath"
	"testing"
)

func TestStoreLoadMissingFileReturnsDefaults(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "settings.json")
	cfg, err := NewStore(path).Load()
	if err != nil {
		t.Fatalf("load missing settings: %v", err)
	}
	if cfg == nil || cfg.Browser != "chromium" {
		t.Fatalf("got %+v", cfg)
	}
}

func TestStoreUpdateCreatesMissingFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "settings.json")
	store := NewStore(path)
	if err := store.Update(func(cfg *AppSettings) error {
		cfg.RunDialogConfirmed = true
		return nil
	}); err != nil {
		t.Fatalf("update missing settings: %v", err)
	}
	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("load saved settings: %v", err)
	}
	if cfg == nil || !cfg.RunDialogConfirmed {
		t.Fatalf("got %+v", cfg)
	}
}
