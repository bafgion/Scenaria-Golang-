package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	InstalledAt string `json:"installed_at"`
}

type Manifest struct {
	Plugins []Entry `json:"plugins"`
}

func ManifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".scenaria", "plugins.json")
}

func LoadManifest(projectRoot string) (Manifest, error) {
	path := ManifestPath(projectRoot)
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Manifest{Plugins: []Entry{}}, nil
		}
		return Manifest{}, fmt.Errorf("read plugin manifest %q: %w", path, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(payload, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode plugin manifest %q: %w", path, err)
	}
	if manifest.Plugins == nil {
		manifest.Plugins = []Entry{}
	}
	return manifest, nil
}

func SaveManifest(projectRoot string, manifest Manifest) error {
	path := ManifestPath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create plugin dir %q: %w", filepath.Dir(path), err)
	}
	payload, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode plugin manifest %q: %w", path, err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write plugin manifest %q: %w", path, err)
	}
	return nil
}

func List(projectRoot string) ([]Entry, error) {
	manifest, err := LoadManifest(projectRoot)
	if err != nil {
		return nil, err
	}
	return manifest.Plugins, nil
}

func Install(projectRoot string, name string, source string) error {
	manifest, err := LoadManifest(projectRoot)
	if err != nil {
		return err
	}

	updated := false
	for i := range manifest.Plugins {
		if manifest.Plugins[i].Name == name {
			manifest.Plugins[i].Source = source
			manifest.Plugins[i].InstalledAt = time.Now().UTC().Format(time.RFC3339)
			updated = true
			break
		}
	}
	if !updated {
		manifest.Plugins = append(manifest.Plugins, Entry{
			Name:        name,
			Source:      source,
			InstalledAt: time.Now().UTC().Format(time.RFC3339),
		})
	}
	return SaveManifest(projectRoot, manifest)
}

func Uninstall(projectRoot string, name string) (bool, error) {
	manifest, err := LoadManifest(projectRoot)
	if err != nil {
		return false, err
	}
	out := make([]Entry, 0, len(manifest.Plugins))
	removed := false
	for _, entry := range manifest.Plugins {
		if entry.Name == name {
			removed = true
			continue
		}
		out = append(out, entry)
	}
	manifest.Plugins = out
	if err := SaveManifest(projectRoot, manifest); err != nil {
		return false, err
	}
	return removed, nil
}
