package settings

import (
	"path/filepath"
	"testing"
)

func TestAppSettingsRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")

	input := &AppSettings{
		Browser:             "chromium",
		Headless:            true,
		RecordingHoverMode:  true,
		RecordingFilterMode: false,
	}
	if err := SaveAppSettings(path, input); err != nil {
		t.Fatalf("SaveAppSettings failed: %v", err)
	}

	got, err := LoadAppSettings(path)
	if err != nil {
		t.Fatalf("LoadAppSettings failed: %v", err)
	}
	if got.Browser != input.Browser || got.Headless != input.Headless {
		t.Fatalf("unexpected settings loaded: %+v", *got)
	}
}

func TestTestClientRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "client.json")

	input := &TestClient{
		Name:    "local",
		BaseURL: "https://example.com",
		Cookies: []Cookie{
			{Name: "sid", Value: "abc", Domain: "example.com", Path: "/", HTTPOnly: true, Secure: true},
		},
		LocalStorage: map[string]string{"token": "123"},
	}
	if err := SaveTestClient(path, input); err != nil {
		t.Fatalf("SaveTestClient failed: %v", err)
	}

	got, err := LoadTestClient(path)
	if err != nil {
		t.Fatalf("LoadTestClient failed: %v", err)
	}
	if got.Name != input.Name || got.BaseURL != input.BaseURL {
		t.Fatalf("unexpected test client loaded: %+v", *got)
	}
	if len(got.Cookies) != 1 {
		t.Fatalf("unexpected cookies: %d", len(got.Cookies))
	}
}
