package paths

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// OpenExternalURL opens an http(s) link in the default browser.
func OpenExternalURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("url is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("unsupported url scheme")
	}
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
