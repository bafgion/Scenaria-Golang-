//go:build windows

package edge

import "os"

// desktopSmokeBrowserArgs enables Playwright CDP when scenaria-gui runs desktop smoke.
// Set SCENARIA_DESKTOP_SMOKE=1 (and optionally SCENARIA_DESKTOP_SMOKE_PORT) before launch.
func desktopSmokeBrowserArgs() []string {
	if os.Getenv("SCENARIA_DESKTOP_SMOKE") == "" {
		return nil
	}
	port := os.Getenv("SCENARIA_DESKTOP_SMOKE_PORT")
	if port == "" {
		port = "9333"
	}
	return []string{"--remote-debugging-port=" + port}
}
