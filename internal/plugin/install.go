package plugin

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FetchAndInstall(projectRoot, name, source string) error {
	source = strings.TrimSpace(source)
	if source == "" {
		return fmt.Errorf("plugin source is required")
	}
	dest := filepath.Join(projectRoot, "addons", name)
	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("clean addon dir: %w", err)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	switch {
	case strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://"):
		tmp, err := downloadToTemp(source)
		if err != nil {
			return err
		}
		defer os.Remove(tmp)
		if strings.HasSuffix(strings.ToLower(tmp), ".zip") || isZipFile(tmp) {
			if err := extractZip(tmp, dest); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("downloaded plugin is not a zip archive")
		}
	case strings.HasSuffix(strings.ToLower(source), ".zip"):
		if err := extractZip(source, dest); err != nil {
			return err
		}
	default:
		info, err := os.Stat(source)
		if err != nil {
			return fmt.Errorf("plugin source not found: %w", err)
		}
		if info.IsDir() {
			if err := copyDir(source, dest); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("plugin source must be a directory or .zip archive")
		}
	}
	return Install(projectRoot, name, source)
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
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
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
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
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
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, in, 0o644)
	})
}
