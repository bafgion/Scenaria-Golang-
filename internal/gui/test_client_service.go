package gui

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

// TestClientService owns project-scoped test client JSON access.
type TestClientService struct {
	projectPath func() string
}

func NewTestClientService(projectPath func() string) *TestClientService {
	return &TestClientService{projectPath: projectPath}
}

func (s *TestClientService) requireProject() (string, error) {
	if s == nil || s.projectPath == nil {
		return "", fmt.Errorf("open a project folder first")
	}
	path := s.projectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return path, nil
}

func (s *TestClientService) List() ([]string, error) {
	path, err := s.requireProject()
	if err != nil {
		return nil, err
	}
	return settings.ListTestClientNames(path)
}

func (s *TestClientService) Details(name string) (string, error) {
	path, err := s.requireProject()
	if err != nil {
		return "", err
	}
	client, err := settings.LoadTestClientByName(path, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("name=%s base_url=%s cookies=%d local_storage=%d",
		client.Name, client.BaseURL, len(client.Cookies), len(client.LocalStorage)), nil
}

func (s *TestClientService) ReadJSON(name string) (string, error) {
	path, err := s.requireProject()
	if err != nil {
		return "", err
	}
	return settings.ReadTestClientJSON(path, name)
}

func (s *TestClientService) SaveJSON(name, content string) error {
	path, err := s.requireProject()
	if err != nil {
		return err
	}
	return settings.SaveTestClientFromJSON(path, name, content)
}

func (s *TestClientService) Delete(name string) error {
	path, err := s.requireProject()
	if err != nil {
		return err
	}
	return settings.DeleteTestClient(path, name)
}
