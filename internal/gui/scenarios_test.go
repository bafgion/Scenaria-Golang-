package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectScenarioTitles(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "demo.feature")
	content := `Функционал: Demo
  Сценарий: Login
    Когда шаг
  Сценарий: Logout
    Когда шаг
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	titles, err := svc.ListScenarioTitles()
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 2 || titles[0] != "Login" || titles[1] != "Logout" {
		t.Fatalf("unexpected titles: %v", titles)
	}
}
