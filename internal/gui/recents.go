package gui

const maxRecents = 6

type RecentsDTO struct {
	Projects []string `json:"projects"`
	Features []string `json:"features"`
}

func (s *Service) LoadRecents() RecentsDTO {
	return s.settingOps().LoadRecents()
}

func (s *Service) RememberRecentProject(path string) error {
	return s.settingOps().RememberRecentProject(path)
}

func (s *Service) RememberRecentFeature(path string) error {
	return s.settingOps().RememberRecentFeature(path)
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
