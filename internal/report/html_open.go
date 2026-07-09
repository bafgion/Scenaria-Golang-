package report

import (
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

// ResolveHTMLReportPath picks the HTML report to open or list in the UI.
// When runs/latest.json exists, its htmlPath wins over a legacy flat .scenaria/report.html.
func ResolveHTMLReportPath(projectRoot, path string) (string, error) {
	path = strings.TrimSpace(path)
	projectRoot = strings.TrimSpace(projectRoot)
	if projectRoot == "" {
		return path, nil
	}
	latest, err := ReadLatestRunPointer(projectRoot)
	if err != nil {
		return "", err
	}
	if latest == nil || strings.TrimSpace(latest.HTMLPath) == "" {
		return path, nil
	}
	latestHTML := strings.TrimSpace(latest.HTMLPath)
	if path == "" {
		return latestHTML, nil
	}
	legacy, err := legacyFlatHTMLReportPath(projectRoot)
	if err != nil {
		return "", err
	}
	if samePath(path, legacy) {
		return latestHTML, nil
	}
	return path, nil
}

func legacyFlatHTMLReportPath(projectRoot string) (string, error) {
	scenariaDir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(scenariaDir, "report.html"), nil
}

func samePath(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == "" || b == "" {
		return false
	}
	aAbs, errA := filepath.Abs(a)
	bAbs, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return aAbs == bAbs
}

// PreferredHTMLReportPath picks full or light HTML report file based on openMode ("full" or "light").
// Falls back to the other variant or the original path when the preferred file is missing.
func PreferredHTMLReportPath(reportPath, openMode string) string {
	reportPath = strings.TrimSpace(reportPath)
	if reportPath == "" {
		return ""
	}
	mode := settings.NormalizeHTMLReportOpenMode(openMode)
	fullPath, lightPath := htmlModePairPaths(reportPath, false)
	switch mode {
	case "light":
		if fileExists(lightPath) {
			return lightPath
		}
		if fileExists(fullPath) {
			return fullPath
		}
	default:
		if fileExists(fullPath) {
			return fullPath
		}
		if fileExists(lightPath) {
			return lightPath
		}
	}
	return reportPath
}
