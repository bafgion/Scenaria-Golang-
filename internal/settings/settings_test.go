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

func TestAppSettingsUntitledTabsRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")

	input := &AppSettings{
		Browser:        "chromium",
		SessionProject: "C:/proj",
		OpenTabs:       []string{"C:/proj/a.feature", "__untitled__:1/demo.feature"},
		ActiveTab:      "__untitled__:1/demo.feature",
		UntitledTabs: []UntitledTabSession{
			{Path: "__untitled__:1/demo.feature", Content: "Функционал: X\nСценарий: Y\n"},
		},
	}
	if err := SaveAppSettings(path, input); err != nil {
		t.Fatalf("SaveAppSettings failed: %v", err)
	}
	got, err := LoadAppSettings(path)
	if err != nil {
		t.Fatalf("LoadAppSettings failed: %v", err)
	}
	if len(got.UntitledTabs) != 1 || got.UntitledTabs[0].Path != input.UntitledTabs[0].Path {
		t.Fatalf("unexpected untitled tabs: %+v", got.UntitledTabs)
	}
	if got.UntitledTabs[0].Content != input.UntitledTabs[0].Content {
		t.Fatalf("unexpected untitled content: %q", got.UntitledTabs[0].Content)
	}
}
