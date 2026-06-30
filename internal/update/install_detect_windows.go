//go:build windows

package update

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// InnoAppUninstallKey matches installer/scenaria-setup.iss AppId.
const InnoAppUninstallKey = `Software\Microsoft\Windows\CurrentVersion\Uninstall\{A4B8C2D1-5E6F-4A7B-9C0D-1E2F3A4B5C6D}_is1`

func DetectInstallMode(installDir string) InstallMode {
	installDir = filepath.Clean(strings.TrimSpace(installDir))
	if installDir == "" {
		return InstallModePortable
	}
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, InnoAppUninstallKey, registry.QUERY_VALUE)
	if err != nil {
		return InstallModePortable
	}
	defer key.Close()
	location, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		return InstallModePortable
	}
	location = filepath.Clean(strings.TrimSpace(location))
	if location == "" {
		return InstallModePortable
	}
	if strings.EqualFold(location, installDir) {
		return InstallModeSetup
	}
	return InstallModePortable
}
