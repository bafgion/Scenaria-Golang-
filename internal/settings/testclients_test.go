package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListTestClientNames(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".scenaria", "test_clients")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DemoUser.json"), []byte(`{"name":"DemoUser"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DemoUser.json.example"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte(`x`), 0o644); err != nil {
		t.Fatal(err)
	}

	names, err := ListTestClientNames(root)
	if err != nil {
		t.Fatalf("ListTestClientNames: %v", err)
	}
	if len(names) != 1 || names[0] != "DemoUser" {
		t.Fatalf("unexpected names: %#v", names)
	}
}
