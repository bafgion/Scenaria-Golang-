package gui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

func (s *Service) ReadVanessaSettingsJSON() (string, error) {
	projectRoot := s.ProjectPath()
	if projectRoot == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	path := filepath.Join(projectRoot, ".scenaria", "vanessa.json")
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := vanessa.DefaultSettings()
			out, merr := json.MarshalIndent(cfg, "", "  ")
			if merr != nil {
				return "", merr
			}
			return string(out) + "\n", nil
		}
		return "", err
	}
	return string(payload), nil
}

func (s *Service) SaveVanessaSettingsJSON(content string) error {
	projectRoot := s.ProjectPath()
	if projectRoot == "" {
		return fmt.Errorf("open a project folder first")
	}
	var cfg vanessa.Settings
	if err := json.Unmarshal([]byte(content), &cfg); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return vanessa.SaveProjectSettings(projectRoot, cfg)
}
