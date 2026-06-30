//go:build !windows

package update

func DetectInstallMode(installDir string) InstallMode {
	return InstallModePortable
}
