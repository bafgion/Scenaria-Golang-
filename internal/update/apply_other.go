//go:build !windows

package update

import "fmt"

func ApplyDownloaded(assetPath string, kind ApplyKind, installDir string, parentPID int, report Reporter) error {
	return fmt.Errorf("auto-update is only supported on Windows")
}
