package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AppSettings struct {
	Browser             string `json:"browser"`
	Headless            bool   `json:"headless"`
	RecordingHoverMode  bool   `json:"recording_hover_mode"`
	RecordingFilterMode bool   `json:"recording_filter_mode"`
}

type TestClient struct {
	Name         string            `json:"name"`
	BaseURL      string            `json:"base_url"`
	Cookies      []Cookie          `json:"cookies"`
	LocalStorage map[string]string `json:"local_storage"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	HTTPOnly bool   `json:"http_only"`
	Secure   bool   `json:"secure"`
}

func LoadAppSettings(path string) (*AppSettings, error) {
	var cfg AppSettings
	if err := readJSON(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveAppSettings(path string, cfg *AppSettings) error {
	return writeJSON(path, cfg)
}

func LoadTestClient(path string) (*TestClient, error) {
	var client TestClient
	if err := readJSON(path, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

func SaveTestClient(path string, client *TestClient) error {
	return writeJSON(path, client)
}

func TestClientPath(projectRoot, name string) (string, error) {
	if projectRoot == "" || name == "" {
		return "", fmt.Errorf("project root and test client name are required")
	}
	return filepath.Join(projectRoot, ".scenaria", "test_clients", name+".json"), nil
}

func readJSON(path string, dst any) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read json file %q: %w", path, err)
	}
	if err := json.Unmarshal(payload, dst); err != nil {
		return fmt.Errorf("decode json file %q: %w", path, err)
	}
	return nil
}

func writeJSON(path string, src any) error {
	payload, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return fmt.Errorf("encode json file %q: %w", path, err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write json file %q: %w", path, err)
	}
	return nil
}
