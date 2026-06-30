package gui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

// CaptureBrowserSession saves cookies and localStorage from the open browser as a TestClient profile.
func (s *Service) CaptureBrowserSession(name string) (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("test client name is required")
	}

	s.mu.RLock()
	session := s.liveSession
	s.mu.RUnlock()
	if session == nil || !session.BrowserAlive() {
		return "", fmt.Errorf("браузер не открыт — откройте браузер или запись, войдите на сайт и повторите")
	}

	client, err := session.ExportTestClient(name)
	if err != nil {
		return "", err
	}

	payload, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode test client: %w", err)
	}
	if err := settings.SaveTestClientFromJSON(path, name, string(payload)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Сохранено: %s (cookies: %d, localStorage: %d, base_url: %s)",
		name, len(client.Cookies), len(client.LocalStorage), client.BaseURL), nil
}
