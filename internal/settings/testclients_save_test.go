package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveTestClientFromJSON(t *testing.T) {
	root := t.TempDir()
	content := `{
  "name": "ignored",
  "base_url": "https://example.com",
  "cookies": [{"name": "sid", "value": "1", "domain": "example.com", "path": "/"}],
  "local_storage": {"k": "v"}
}`
	if err := SaveTestClientFromJSON(root, "demo", content); err != nil {
		t.Fatalf("SaveTestClientFromJSON: %v", err)
	}
	path := filepath.Join(TestClientsDir(root), "demo.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file: %v", err)
	}
	client, err := LoadTestClientByName(root, "demo")
	if err != nil {
		t.Fatalf("LoadTestClientByName: %v", err)
	}
	if client.Name != "demo" {
		t.Fatalf("name = %q, want demo", client.Name)
	}
	if client.BaseURL != "https://example.com" {
		t.Fatalf("base_url = %q", client.BaseURL)
	}
	if len(client.Cookies) != 1 || client.Cookies[0].Name != "sid" {
		t.Fatalf("cookies = %+v", client.Cookies)
	}
	if client.LocalStorage["k"] != "v" {
		t.Fatalf("local_storage = %+v", client.LocalStorage)
	}

	raw, err := ReadTestClientJSON(root, "demo")
	if err != nil {
		t.Fatalf("ReadTestClientJSON: %v", err)
	}
	if !strings.Contains(raw, `"base_url"`) {
		t.Fatalf("json = %q", raw)
	}

	if err := DeleteTestClient(root, "demo"); err != nil {
		t.Fatalf("DeleteTestClient: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected deleted file, stat err = %v", err)
	}
}

func TestSaveTestClientFromJSONRejectsInvalidName(t *testing.T) {
	root := t.TempDir()
	err := SaveTestClientFromJSON(root, "../evil", `{}`)
	if err == nil {
		t.Fatal("expected error for invalid name")
	}
}
