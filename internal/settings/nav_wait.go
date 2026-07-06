package settings

import "strings"

// ResolveNavWaitUntil returns navigation wait policy: project config overrides app settings.
func ResolveNavWaitUntil(projectRoot string, app *AppSettings) string {
	if projectRoot != "" {
		if proj, err := LoadProjectConfig(projectRoot); err == nil {
			if wait := strings.TrimSpace(proj.NavWaitUntil); wait != "" {
				return wait
			}
		}
	}
	if app != nil {
		return strings.TrimSpace(app.NavWaitUntil)
	}
	return ""
}
