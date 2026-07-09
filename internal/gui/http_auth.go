package gui

import (
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
	return s.settingOps().loadAppSettings()
}

func (s *Service) saveAppSettings(cfg *settings.AppSettings) error {
	return s.settingOps().saveAppSettings(cfg)
}

func (s *Service) ListHTTPAuthHosts() ([]string, error) {
	return s.settingOps().ListHTTPAuthHosts()
}

func (s *Service) HTTPAuthForHost(host string) (HTTPAuthCredentials, error) {
	return s.settingOps().HTTPAuthForHost(host)
}

func (s *Service) SaveHTTPAuth(req HTTPAuthRequest) error {
	return s.settingOps().SaveHTTPAuth(req)
}

func (s *Service) RemoveHTTPAuth(host string) error {
	return s.settingOps().RemoveHTTPAuth(host)
}

func (s *Service) PrepareRecordURL(url string) (string, error) {
	return s.settingOps().PrepareRecordURL(url)
}
