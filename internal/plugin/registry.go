package plugin

import (
	"encoding/json"
	"fmt"
	"io"
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

var (
	registryMkdirAll    = os.MkdirAll
	registryCreateTemp  = os.CreateTemp
	registryRename      = os.Rename
	registryRemove      = os.Remove
	registryRemoveAll   = os.RemoveAll
	registryStat        = os.Stat
	registryReadFile    = os.ReadFile
	registryWriteTemp   = func(file *os.File, data []byte) (int, error) { return file.Write(data) }
	registryWriteBackup = restoreRegistryBackupByCopy
)

func ManifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".scenaria", "plugins.json")
}

func LoadManifest(projectRoot string) (Manifest, error) {
	path := ManifestPath(projectRoot)
	payload, err := registryReadFile(path)
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
	if err := validateManifestPluginIDs(manifest); err != nil {
		return err
	}
	path := ManifestPath(projectRoot)
	if err := registryMkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create plugin dir %q: %w", filepath.Dir(path), err)
	}
	payload, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode plugin manifest %q: %w", path, err)
	}
	if err := writeRegistryAtomic(path, append(payload, '\n'), 0o644); err != nil {
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
	if err := ValidatePluginID(name); err != nil {
		return err
	}
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
	if err := ValidatePluginID(name); err != nil {
		return false, err
	}
	dest, err := addonPath(projectRoot, name)
	if err != nil {
		return false, err
	}
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

	backup := ""
	hadFiles := false
	if _, err := registryStat(dest); err == nil {
		var backupErr error
		backup, backupErr = createPluginBackupPath(projectRoot, name)
		if backupErr != nil {
			return false, backupErr
		}
		if err := registryRename(dest, backup); err != nil {
			_ = registryRemoveAll(backup)
			return false, fmt.Errorf("backup plugin files: %w", err)
		}
		hadFiles = true
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("stat plugin files: %w", err)
	}

	if err := SaveManifest(projectRoot, manifest); err != nil {
		if hadFiles {
			if rollbackErr := registryRename(backup, dest); rollbackErr != nil {
				return false, fmt.Errorf("update plugin registry: %w (restore plugin files failed: %v)", err, rollbackErr)
			}
		}
		return false, err
	}
	if hadFiles {
		if err := registryRemoveAll(backup); err != nil && !os.IsNotExist(err) {
			return true, fmt.Errorf("plugin uninstalled but cleanup failed: %w", err)
		}
	}
	return removed || hadFiles, nil
}

func validateManifestPluginIDs(manifest Manifest) error {
	for _, entry := range manifest.Plugins {
		if err := ValidatePluginID(entry.Name); err != nil {
			return fmt.Errorf("invalid plugin manifest entry %q: %w", entry.Name, err)
		}
	}
	return nil
}

func writeRegistryAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := registryCreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = registryRemove(tmpPath)
		}
	}()
	if _, err := registryWriteTemp(tmp, data); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Chmod(perm); err != nil {
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	backupPath := ""
	backupCreated := false
	if _, err := registryStat(path); err == nil {
		backupPath = path + ".bak"
		if err := registryRemove(backupPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale backup: %w", err)
		}
		if err := registryRename(path, backupPath); err != nil {
			return fmt.Errorf("backup manifest: %w", err)
		}
		backupCreated = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat manifest: %w", err)
	}

	if err := registryRename(tmpPath, path); err != nil {
		if backupCreated {
			if restoreErr := registryRename(backupPath, path); restoreErr != nil {
				if copyErr := registryWriteBackup(backupPath, path, perm); copyErr != nil {
					return fmt.Errorf("replace manifest: %w (restore backup failed: %v, copy restore failed: %w)", err, restoreErr, copyErr)
				}
			}
		}
		return fmt.Errorf("replace manifest: %w", err)
	}
	if backupCreated {
		_ = registryRemove(backupPath)
	}
	cleanupTemp = false
	return nil
}

func restoreRegistryBackupByCopy(src, dest string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := registryCreateTemp(filepath.Dir(dest), filepath.Base(dest)+".restore-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = registryRemove(tmpPath)
		}
	}()
	if _, err := io.Copy(tmp, in); err != nil {
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := registryRename(tmpPath, dest); err != nil {
		return err
	}
	cleanupTemp = false
	return nil
}
