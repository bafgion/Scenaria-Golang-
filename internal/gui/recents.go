package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

const maxRecents = 6

type RecentsDTO struct {
	Projects []string `json:"projects"`
	Features []string `json:"features"`
}

func (s *Service) LoadRecents() RecentsDTO {
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil || cfg == nil {
		return RecentsDTO{}
	}
	return RecentsDTO{
		Projects: trimRecents(cfg.RecentProjects),
		Features: trimRecents(cfg.RecentFeatures),
	}
}

func (s *Service) RememberRecentProject(path string) error {
	return rememberRecent(path, true)
}

func (s *Service) RememberRecentFeature(path string) error {
	return rememberRecent(path, false)
}

func rememberRecent(itemPath string, project bool) error {
	itemPath = strings.TrimSpace(itemPath)
	if itemPath == "" {
		return nil
	}
	settingsPath := settings.DefaultAppSettingsPath()
	if settingsPath == "" {
		return fmt.Errorf("settings path unavailable")
	}
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil {
		return err
	}
	if cfg == nil {
		cfg = &settings.AppSettings{Browser: "chromium"}
	}
	if project {
		cfg.RecentProjects = pushRecent(cfg.RecentProjects, itemPath)
	} else {
		cfg.RecentFeatures = pushRecent(cfg.RecentFeatures, itemPath)
	}
	return settings.SaveAppSettings(settingsPath, cfg)
}

func pushRecent(list []string, item string) []string {
	out := []string{item}
	for _, existing := range list {
		if existing != item {
			out = append(out, existing)
		}
		if len(out) >= maxRecents {
			break
		}
	}
	return out
}

func trimRecents(list []string) []string {
	if len(list) > maxRecents {
		return list[:maxRecents]
	}
	return list
}

func clampSidebarWidth(w int) int {
	if w < 120 {
		return 260
	}
	if w > 480 {
		return 480
	}
	return w
}
