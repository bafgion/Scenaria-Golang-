package plugin

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MaxPluginDownloadBytes caps remote plugin zip downloads (100 MiB).
const MaxPluginDownloadBytes = 100 << 20

var (
	installMkdirAll  = os.MkdirAll
	installMkdirTemp = os.MkdirTemp
	installRemoveAll = os.RemoveAll
	installRename    = os.Rename
)

func FetchAndInstall(projectRoot, name, source string) error {
	source = strings.TrimSpace(source)
	if source == "" {
		return fmt.Errorf("plugin source is required")
	}
	dest, err := addonPath(projectRoot, name)
	if err != nil {
		return err
	}

	staging, err := createPluginStagingDir(projectRoot, name)
	if err != nil {
		return err
	}
	stagingActive := true
	defer func() {
		if stagingActive {
			_ = installRemoveAll(staging)
		}
	}()

	switch {
	case strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://"):
		tmp, err := downloadToTemp(source)
		if err != nil {
			return err
		}
		defer os.Remove(tmp)
		if strings.HasSuffix(strings.ToLower(tmp), ".zip") || isZipFile(tmp) {
			if err := extractZip(tmp, staging); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("downloaded plugin is not a zip archive")
		}
	case strings.HasSuffix(strings.ToLower(source), ".zip"):
		if err := extractZip(source, staging); err != nil {
			return err
		}
	default:
		info, err := os.Stat(source)
		if err != nil {
			return fmt.Errorf("plugin source not found: %w", err)
		}
		if info.IsDir() {
			if err := copyDir(source, staging); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("plugin source must be a directory or .zip archive")
		}
	}
	if err := validateStagedPlugin(staging, name); err != nil {
		return err
	}
	if err := commitPluginInstall(projectRoot, name, source, staging, dest); err != nil {
		return err
	}
	stagingActive = false
	return nil
}

func downloadToTemp(url string) (string, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download plugin: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download plugin: HTTP %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp("", "scenaria-plugin-*.zip")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	limited := io.LimitReader(resp.Body, MaxPluginDownloadBytes+1)
	written, err := io.Copy(tmp, limited)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	if written > MaxPluginDownloadBytes {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("download plugin: exceeds %d byte limit", MaxPluginDownloadBytes)
	}
	return tmp.Name(), nil
}

func isZipFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	buf := make([]byte, 2)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 'P' && buf[1] == 'K'
}

func extractZip(zipPath, dest string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open plugin zip: %w", err)
	}
	defer reader.Close()
	for _, file := range reader.File {
		target := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal zip path: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := installMkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := installMkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			src.Close()
			return err
		}
		_, copyErr := io.Copy(out, src)
		out.Close()
		src.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return installMkdirAll(target, 0o755)
		}
		in, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, in, 0o644)
	})
}

func createPluginStagingDir(projectRoot, name string) (string, error) {
	parent := filepath.Join(projectRoot, ".scenaria", "plugin-staging")
	if err := installMkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("create plugin staging dir: %w", err)
	}
	dir, err := installMkdirTemp(parent, name+"-*")
	if err != nil {
		return "", fmt.Errorf("create plugin staging dir: %w", err)
	}
	return dir, nil
}

func createPluginBackupPath(projectRoot, name string) (string, error) {
	parent := filepath.Join(projectRoot, ".scenaria", "plugin-backups")
	if err := installMkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("create plugin backup dir: %w", err)
	}
	dir, err := installMkdirTemp(parent, name+"-*")
	if err != nil {
		return "", fmt.Errorf("create plugin backup dir: %w", err)
	}
	if err := installRemoveAll(dir); err != nil {
		return "", fmt.Errorf("prepare plugin backup dir: %w", err)
	}
	return dir, nil
}

func validateStagedPlugin(staging, expectedID string) error {
	payload, err := os.ReadFile(filepath.Join(staging, "plugin.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin descriptor plugin.json is required")
		}
		return fmt.Errorf("read staged plugin descriptor: %w", err)
	}
	var desc Descriptor
	if err := json.Unmarshal(payload, &desc); err != nil {
		return fmt.Errorf("decode staged plugin descriptor: %w", err)
	}
	if strings.TrimSpace(desc.ID) != "" {
		if desc.ID != strings.TrimSpace(desc.ID) {
			return fmt.Errorf("plugin descriptor id has leading or trailing whitespace")
		}
		if err := ValidatePluginID(desc.ID); err != nil {
			return fmt.Errorf("invalid plugin descriptor id: %w", err)
		}
		if desc.ID != expectedID {
			return fmt.Errorf("plugin descriptor id %q does not match requested plugin id %q", desc.ID, expectedID)
		}
	}
	return nil
}

func commitPluginInstall(projectRoot, name, source, staging, dest string) error {
	backup := ""
	hadExisting := false
	if _, err := os.Stat(dest); err == nil {
		var backupErr error
		backup, backupErr = createPluginBackupPath(projectRoot, name)
		if backupErr != nil {
			return backupErr
		}
		if err := installRename(dest, backup); err != nil {
			_ = installRemoveAll(backup)
			return fmt.Errorf("backup existing plugin: %w", err)
		}
		hadExisting = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat existing plugin: %w", err)
	}

	if err := installMkdirAll(filepath.Dir(dest), 0o755); err != nil {
		if rollbackErr := rollbackPluginInstall(dest, backup, hadExisting); rollbackErr != nil {
			return fmt.Errorf("create addons dir: %w (rollback failed: %v)", err, rollbackErr)
		}
		return fmt.Errorf("create addons dir: %w", err)
	}
	if err := installRename(staging, dest); err != nil {
		if rollbackErr := rollbackPluginInstall(dest, backup, hadExisting); rollbackErr != nil {
			return fmt.Errorf("commit plugin files: %w (rollback failed: %v)", err, rollbackErr)
		}
		return fmt.Errorf("commit plugin files: %w", err)
	}

	if err := Install(projectRoot, name, source); err != nil {
		if rollbackErr := rollbackPluginInstall(dest, backup, hadExisting); rollbackErr != nil {
			return fmt.Errorf("update plugin registry: %w (rollback failed: %v)", err, rollbackErr)
		}
		return err
	}
	if hadExisting {
		_ = installRemoveAll(backup)
	}
	return nil
}

func rollbackPluginInstall(dest, backup string, hadExisting bool) error {
	if err := installRemoveAll(dest); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove failed plugin: %w", err)
	}
	if hadExisting {
		if err := installRename(backup, dest); err != nil {
			return fmt.Errorf("restore previous plugin: %w", err)
		}
	}
	return nil
}
