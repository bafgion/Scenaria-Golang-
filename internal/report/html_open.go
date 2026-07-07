package report

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

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
