package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type HTTPAuthRequest struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type HTTPAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Service) loadAppSettings() (*settings.AppSettings, error) {
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &settings.AppSettings{Browser: "chromium"}
	}
	return cfg, nil
}

func (s *Service) saveAppSettings(cfg *settings.AppSettings) error {
	path := settings.DefaultAppSettingsPath()
	if path == "" {
		return fmt.Errorf("settings path is not configured")
	}
	return settings.SaveAppSettings(path, cfg)
}

func (s *Service) ListHTTPAuthHosts() ([]string, error) {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return nil, err
	}
	return httpauth.ListHosts(cfg), nil
}

func (s *Service) HTTPAuthForHost(host string) (HTTPAuthCredentials, error) {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return HTTPAuthCredentials{}, err
	}
	username, password := httpauth.CredentialsForHost(host, cfg)
	return HTTPAuthCredentials{Username: username, Password: password}, nil
}

func (s *Service) SaveHTTPAuth(req HTTPAuthRequest) error {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return err
	}
	host := strings.TrimSpace(req.Host)
	if host == "" {
		return fmt.Errorf("host is required")
	}
	httpauth.StoreHostCredentials(host, req.Username, req.Password, cfg)
	return s.saveAppSettings(cfg)
}

func (s *Service) RemoveHTTPAuth(host string) error {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return err
	}
	httpauth.RemoveHostCredentials(host, cfg)
	return s.saveAppSettings(cfg)
}

func (s *Service) PrepareRecordURL(url string) (string, error) {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return url, err
	}
	clean := httpauth.ApplyURLCredentials(url, cfg)
	if err := s.saveAppSettings(cfg); err != nil {
		return clean, err
	}
	return clean, nil
}
