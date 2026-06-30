package paths

import (
	"os"
	"path/filepath"
)

// AppInstallDir is the directory containing the running scenaria-gui / scenaria binary.
func AppInstallDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}
