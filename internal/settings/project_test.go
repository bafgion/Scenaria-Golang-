package settings

import "testing"

func TestLoadProjectConfigDefaults(t *testing.T) {
	cfg, err := LoadProjectConfig(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultRunner != "playwright" || cfg.FeaturesRoot != "features" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestSaveAndLoadProjectConfig(t *testing.T) {
	root := t.TempDir()
	if err := SaveProjectConfig(root, ProjectConfig{DefaultRunner: "playwright", BaseURL: "https://app.local"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadProjectConfig(root)
	if err != nil || cfg.BaseURL != "https://app.local" {
		t.Fatalf("unexpected config: %+v err=%v", cfg, err)
	}
}
