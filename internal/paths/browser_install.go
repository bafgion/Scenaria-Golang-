package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

var BrowserEngineLabels = map[string]string{
	"chromium": "Chromium",
	"firefox":  "Firefox",
	"webkit":   "WebKit",
}

func NormalizeBrowserEngine(engine string) string {
	switch strings.ToLower(strings.TrimSpace(engine)) {
	case "firefox":
		return "firefox"
	case "webkit":
		return "webkit"
	default:
		return "chromium"
	}
}

func browserDirPrefix(engine string) string {
	switch NormalizeBrowserEngine(engine) {
	case "firefox":
		return "firefox-"
	case "webkit":
		return "webkit-"
	default:
		return "chromium-"
	}
}

func msPlaywrightDir() string {
	if custom := os.Getenv("PLAYWRIGHT_BROWSERS_PATH"); custom != "" {
		return custom
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		return filepath.Join(local, "ms-playwright")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache", "ms-playwright")
}

func directoryHasEngine(dir, engine string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	prefix := browserDirPrefix(engine)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			return true
		}
	}
	return false
}

func BrowserCacheCandidates(engine string) []string {
	engine = NormalizeBrowserEngine(engine)
	bundled := BundledBrowsersDir()
	local := msPlaywrightDir()

	candidates := make([]string, 0, 2)
	if bundled != "" {
		if info, err := os.Stat(bundled); err == nil && info.IsDir() && directoryHasEngine(bundled, "chromium") {
			candidates = append(candidates, bundled)
		}
	}
	if local != "" {
		candidates = append(candidates, local)
	}

	preferred := make([]string, 0, len(candidates))
	for _, path := range candidates {
		if directoryHasEngine(path, engine) {
			preferred = append(preferred, path)
		}
	}
	if len(preferred) > 0 {
		return preferred
	}
	return candidates
}

func InstallPlaywrightBrowsersPath(engine string) string {
	engine = NormalizeBrowserEngine(engine)
	bundled := BundledBrowsersDir()
	if bundled != "" {
		if info, err := os.Stat(bundled); err == nil && info.IsDir() && directoryHasEngine(bundled, "chromium") {
			_ = os.MkdirAll(bundled, 0o755)
			return bundled
		}
	}
	local := msPlaywrightDir()
	_ = os.MkdirAll(local, 0o755)
	return local
}

func ResolvePlaywrightBrowsersPath(engine string) string {
	candidates := BrowserCacheCandidates(engine)
	if len(candidates) > 0 {
		return candidates[0]
	}
	return msPlaywrightDir()
}

func ConfigurePlaywrightBrowsersForEngine(engine string) string {
	path := ResolvePlaywrightBrowsersPath(engine)
	if path == "" {
		return ConfigurePlaywrightBrowsers()
	}
	_ = os.MkdirAll(path, 0o755)
	_ = os.Setenv("PLAYWRIGHT_BROWSERS_PATH", path)
	return path
}

func FindBrowserExecutable(cache, engine string) (string, bool) {
	engine = NormalizeBrowserEngine(engine)
	info, err := os.Stat(cache)
	if err != nil || !info.IsDir() {
		return "", false
	}
	prefix := browserDirPrefix(engine)
	entries, err := os.ReadDir(cache)
	if err != nil {
		return "", false
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			names = append(names, entry.Name())
		}
	}
	sort.Slice(names, func(i, j int) bool { return names[i] > names[j] })

	for _, name := range names {
		folder := filepath.Join(cache, name)
		for _, candidate := range browserExecutableCandidates(engine, folder) {
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate, true
			}
		}
	}
	return "", false
}

func browserExecutableCandidates(engine, folder string) []string {
	switch NormalizeBrowserEngine(engine) {
	case "chromium":
		return []string{
			filepath.Join(folder, "chrome-win64", "chrome.exe"),
			filepath.Join(folder, "chrome-linux", "chrome"),
			filepath.Join(folder, "chrome-mac", "Chromium.app", "Contents", "MacOS", "Chromium"),
		}
	case "firefox":
		return []string{
			filepath.Join(folder, "firefox", "firefox.exe"),
			filepath.Join(folder, "firefox", "firefox"),
		}
	default:
		return []string{
			filepath.Join(folder, "Playwright.exe"),
			filepath.Join(folder, "pw_run.sh"),
		}
	}
}

func BrowserInstallStatus(engine string) (bool, string) {
	engine = NormalizeBrowserEngine(engine)
	label := BrowserEngineLabels[engine]
	ConfigurePlaywrightBrowsersForEngine(engine)

	for _, cache := range BrowserCacheCandidates(engine) {
		if executable, ok := FindBrowserExecutable(cache, engine); ok {
			return true, executable
		}
	}
	searched := strings.Join(BrowserCacheCandidates(engine), ", ")
	if searched == "" {
		searched = msPlaywrightDir()
	}
	return false, fmt.Sprintf("%s не установлен (каталоги: %s)", label, searched)
}

// EnsurePlaywrightEngine configures browser paths and installs the engine only when missing.
// Already-installed browsers are not re-downloaded on every recorder/run session.
func EnsurePlaywrightEngine(engine string) error {
	engine = NormalizeBrowserEngine(engine)
	ConfigurePlaywrightBrowsersForEngine(engine)
	if installed, _ := BrowserInstallStatus(engine); installed {
		return nil
	}
	_, err := InstallBrowserEngine(engine, nil)
	return err
}

func InstallBrowserEngine(engine string, onLine func(string)) (string, error) {
	engine = NormalizeBrowserEngine(engine)
	label := BrowserEngineLabels[engine]
	cache := InstallPlaywrightBrowsersPath(engine)
	ConfigurePlaywrightBrowsersForEngine(engine)
	if onLine != nil {
		onLine(fmt.Sprintf("Установка %s в %s…", label, cache))
	}

	var log strings.Builder
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{engine},
		Stdout:   &log,
		Stderr:   &log,
		Verbose:  true,
	}); err != nil {
		tail := strings.TrimSpace(log.String())
		if tail == "" {
			return "", fmt.Errorf("не удалось установить %s: %w", label, err)
		}
		return "", fmt.Errorf("не удалось установить %s:\n%s", label, tail)
	}
	if onLine != nil && log.Len() > 0 {
		for _, line := range strings.Split(strings.TrimSpace(log.String()), "\n") {
			if strings.TrimSpace(line) != "" {
				onLine(line)
			}
		}
	}

	installed, detail := BrowserInstallStatus(engine)
	if !installed {
		return "", fmt.Errorf("после установки %s недоступен: %s\nКаталог браузеров: %s", label, detail, cache)
	}
	return detail, nil
}
