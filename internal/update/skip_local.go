package update

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SkipRemoteReleaseCheck reports whether GitHub release comparison should be skipped.
// Local dev builds and desktop smoke always run the latest sources from the repo.
func SkipRemoteReleaseCheck(installDir string) bool {
	if os.Getenv("SCENARIA_DESKTOP_SMOKE") != "" {
		return true
	}
	if os.Getenv("SCENARIA_SKIP_UPDATE_CHECK") != "" {
		return true
	}
	return isLocalDevInstall(installDir)
}

func isLocalDevInstall(installDir string) bool {
	dir := strings.TrimSpace(installDir)
	if dir == "" {
		return false
	}
	norm := strings.ToLower(filepath.ToSlash(dir))
	markers := []string{
		"/build/bin",
		"/build/",
		"scenaria_go/build",
		"scenaria-go/build",
	}
	for _, marker := range markers {
		if strings.Contains(norm, marker) {
			return true
		}
	}
	return false
}

func localUpToDateInfo(currentVersion string) *Info {
	v := DisplayVersion(currentVersion)
	return &Info{
		CurrentVersion:  v,
		LatestVersion:   v,
		UpdateAvailable: false,
		Message:         fmt.Sprintf("Установлена актуальная версия (%s)", v),
	}
}
