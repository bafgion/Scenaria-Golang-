package gui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceInProjectStepsOnly(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.feature")
	second := filepath.Join(dir, "b.feature")
	if err := os.WriteFile(first, []byte("Функционал: UI\nСценарий: A\n\tДопустим открыт \"https://old.example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("Функционал: UI\nСценарий: B\n\tДопустим открыт \"https://old.example.com/page\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ReplaceInProject(ProjectReplaceRequest{
		Find:    "old.example.com",
		Replace: "new.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.FilesChanged != 2 || got.Replacements != 2 {
		t.Fatalf("got %+v", got)
	}
	a, _ := os.ReadFile(first)
	b, _ := os.ReadFile(second)
	if !strings.Contains(string(a), "new.example.com") || !strings.Contains(string(b), "new.example.com") {
		t.Fatalf("files not updated: %q %q", a, b)
	}
}

func TestReplaceInProjectSkipsFeatureTitle(t *testing.T) {
	dir := t.TempDir()
	feature := filepath.Join(dir, "demo.feature")
	payload := "Функционал: old title\nСценарий: old title\n\tДопустим открыт \"https://x.com\"\n"
	if err := os.WriteFile(feature, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old title", Replace: "new title"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Replacements != 0 {
		t.Fatalf("expected no replacements in headers, got %+v", got)
	}
}
