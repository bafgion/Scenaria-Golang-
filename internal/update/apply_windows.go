//go:build windows

package update

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

const (
	updateStagingDir = "_update_staging"
	updateBatName    = "_apply_update.bat"
	updateVBSName    = "_apply_update.vbs"
)

// ApplyDownloaded launches the correct auto-update flow for the downloaded artifact.
func ApplyDownloaded(assetPath string, kind ApplyKind, installDir string, parentPID int, report Reporter) error {
	assetPath = strings.TrimSpace(assetPath)
	installDir = strings.TrimSpace(installDir)
	if assetPath == "" || installDir == "" {
		return fmt.Errorf("update paths are required")
	}
	installDir = filepath.Clean(installDir)
	switch kind {
	case ApplyKindPortable:
		return applyPortableZip(assetPath, installDir, parentPID, report)
	case ApplyKindSetup:
		report.report("apply", "Запуск установщика…", 92)
		err := applySetupExe(assetPath)
		if err == nil {
			report.report("restart", "Перезапуск приложения…", 100)
		}
		return err
	default:
		return fmt.Errorf("unsupported update kind %q", kind)
	}
}

const createNoWindow = 0x08000000

func applySetupExe(setupPath string) error {
	cmd := exec.Command(setupPath,
		"/VERYSILENT",
		"/SUPPRESSMSGBOXES",
		"/CLOSEAPPLICATIONS",
		"/RESTARTAPPLICATIONS",
		"/NORESTART",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | createNoWindow,
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start installer: %w", err)
	}
	return nil
}

func applyPortableZip(zipPath, installDir string, parentPID int, report Reporter) error {
	report.report("prepare", "Распаковка обновления…", 84)
	tempRoot, err := os.MkdirTemp("", "scenaria-update-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempRoot)

	extracted := filepath.Join(tempRoot, "extracted")
	if err := extractZip(zipPath, extracted); err != nil {
		return err
	}
	payloadRoot, err := findPayloadRoot(extracted)
	if err != nil {
		return err
	}
	report.report("prepare", "Копирование файлов…", 88)
	localStaging := filepath.Join(installDir, updateStagingDir)
	_ = os.RemoveAll(localStaging)
	if err := copyTree(payloadRoot, localStaging); err != nil {
		return err
	}
	scriptPath := filepath.Join(installDir, updateBatName)
	if err := os.WriteFile(scriptPath, []byte(portableUpdateScript(installDir, localStaging, parentPID)), 0o644); err != nil {
		return err
	}
	report.report("apply", "Установка обновления…", 95)
	if err := launchHiddenBatch(scriptPath); err != nil {
		return err
	}
	report.report("restart", "Перезапуск приложения…", 100)
	return nil
}

func extractZip(zipPath, targetDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open update zip: %w", err)
	}
	defer reader.Close()
	for _, file := range reader.File {
		dest := filepath.Join(targetDir, file.Name)
		if !strings.HasPrefix(filepath.Clean(dest), filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("zip path escapes target: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
		if err != nil {
			_ = src.Close()
			return err
		}
		_, copyErr := io.Copy(out, src)
		closeErr := errorsClose(src, out)
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

func errorsClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func findPayloadRoot(extracted string) (string, error) {
	var matches []string
	_ = filepath.WalkDir(extracted, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.EqualFold(d.Name(), brand.GUIExeName) {
			matches = append(matches, filepath.Dir(path))
		}
		return nil
	})
	if len(matches) == 0 {
		return "", fmt.Errorf("%s not found in update archive", brand.GUIExeName)
	}
	return matches[0], nil
}

func copyTree(source, target string) error {
	return filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(target, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}
		return copyFile(path, dest)
	})
}

func copyFile(source, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func launchHiddenBatch(scriptPath string) error {
	if launchHiddenBatchHook != nil {
		return launchHiddenBatchHook(scriptPath)
	}
	vbsPath := filepath.Join(filepath.Dir(scriptPath), updateVBSName)
	bat := strings.ReplaceAll(filepath.Clean(scriptPath), `"`, `""`)
	content := fmt.Sprintf("CreateObject(\"WScript.Shell\").Run \"cmd /c \"\"%s\"\"\", 0, False\r\n", bat)
	if err := os.WriteFile(vbsPath, []byte(content), 0o644); err != nil {
		return err
	}
	cmd := exec.Command("wscript.exe", "//Nologo", vbsPath)
	cmd.Dir = filepath.Dir(scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | createNoWindow | 0x01000000,
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch update script: %w", err)
	}
	time.Sleep(400 * time.Millisecond)
	return nil
}
