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
	Username    string `json:"username"`
	HasPassword bool   `json:"hasPassword"`
}

func (s *Service) loadAppSettings() (*settings.AppSettings, error) {
	cfg, err := s.settingsStore.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &settings.AppSettings{Browser: "chromium"}
	}
	return cfg, nil
}

func (s *Service) saveAppSettings(cfg *settings.AppSettings) error {
	return s.settingsStore.Update(func(current *settings.AppSettings) error {
		*current = *cfg
		return nil
	})
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
	return HTTPAuthCredentials{Username: username, HasPassword: password != ""}, nil
}

func (s *Service) SaveHTTPAuth(req HTTPAuthRequest) error {
	host := strings.TrimSpace(req.Host)
	if host == "" {
		return fmt.Errorf("host is required")
	}
	return s.settingsStore.Update(func(cfg *settings.AppSettings) error {
		password := req.Password
		if strings.TrimSpace(password) == "" {
			if _, existing := httpauth.CredentialsForHost(host, cfg); existing != "" {
				password = existing
			}
		}
		httpauth.StoreHostCredentials(host, req.Username, password, cfg)
		return nil
	})
}

func (s *Service) RemoveHTTPAuth(host string) error {
	return s.settingsStore.Update(func(cfg *settings.AppSettings) error {
		httpauth.RemoveHostCredentials(host, cfg)
		return nil
	})
}

func (s *Service) PrepareRecordURL(url string) (string, error) {
	var clean string
	err := s.settingsStore.Update(func(cfg *settings.AppSettings) error {
		clean = httpauth.ApplyURLCredentials(url, cfg)
		return nil
	})
	if err != nil {
		return url, err
	}
	return clean, nil
}
