package settings

import (
	"errors"
	"os"
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

func TestWriteJSONRestoresOriginalOnReplacementFailure(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "settings.json")
	if err := os.WriteFile(path, []byte("{\"browser\":\"old\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := atomicRename
	defer func() { atomicRename = origRename }()
	calls := 0
	atomicRename = func(oldpath, newpath string) error {
		calls++
		if calls == 2 {
			return errors.New("replace failed")
		}
		return os.Rename(oldpath, newpath)
	}
	err := SaveAppSettings(path, &AppSettings{Browser: "new"})
	if err == nil {
		t.Fatal("expected error")
	}
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(raw) != "{\"browser\":\"old\"}\n" {
		t.Fatalf("expected original content restored, got %q", string(raw))
	}
	if _, err := os.Stat(path + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup after restore, stat err=%v", err)
	}
}

func TestWriteJSONRetainsBackupWhenRestoreCannotReplaceTarget(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "settings.json")
	if err := os.WriteFile(path, []byte("{\"browser\":\"old\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	origRename := atomicRename
	defer func() { atomicRename = origRename }()
	calls := 0
	atomicRename = func(oldpath, newpath string) error {
		calls++
		switch calls {
		case 2, 3, 4:
			return errors.New("rename unavailable")
		default:
			return os.Rename(oldpath, newpath)
		}
	}
	err := SaveAppSettings(path, &AppSettings{Browser: "new"})
	if err == nil {
		t.Fatal("expected error")
	}
	if raw, readErr := os.ReadFile(path); readErr == nil && string(raw) != "" && string(raw) != "{\"browser\":\"old\"}\n" {
		t.Fatalf("unexpected replacement content after failed restore: %q", string(raw))
	}
	backup, readErr := os.ReadFile(path + ".bak")
	if readErr != nil {
		t.Fatalf("expected retained backup: %v", readErr)
	}
	if string(backup) != "{\"browser\":\"old\"}\n" {
		t.Fatalf("expected backup to retain original content, got %q", string(backup))
	}
}

func TestWriteJSONNewFileSave(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "settings.json")
	if err := SaveAppSettings(path, &AppSettings{Browser: "chromium"}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatal("expected written settings")
	}
	if _, err := os.Stat(path + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("did not expect backup for new file, stat err=%v", err)
	}
}
