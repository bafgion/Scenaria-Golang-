// Package brand holds product naming and release artifact constants (Python app/brand.py parity).
package brand

import "strings"

const (
	Name        = "Scenaria"
	Short       = "Scenaria"
	Tagline     = "Автотесты сайтов · Gherkin · Playwright"
	Description = "Запись и запуск автотестов сайтов в браузере."

	PackageName = "scenaria"
	ExeName     = "Scenaria.exe"
	GUIExeName  = "scenaria-gui.exe"
	AppDataDir  = "Scenaria"

	PortableDir = "Scenaria"
	PortableZip = "Scenaria-Portable.zip"
	SetupExe    = "Scenaria-Setup.exe"
	UpdateZip   = "Scenaria-update.zip"
	DistDir     = "Scenaria"

	DefaultGitHubRepo = "bafgion/Scenaria-Golang-"
	GitHubRepoEnv     = "SCENARIA_GITHUB_REPO"
)

// AboutText matches Python app.qt.branding.about_text().
func AboutText() string {
	return Description + "\n" + Tagline
}

// UserAgent is the HTTP User-Agent product token (Python: f"{BRAND_NAME}/{app_version()}").
func UserAgent(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return Name
	}
	return Name + "/" + version
}

// OverlayTitle is the draggable browser toolbar caption suffix.
func OverlayTitle() string {
	return Name + " — перетащите панель"
}
