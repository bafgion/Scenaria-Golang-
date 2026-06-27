package scenario

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestFeatureStore_LoadSave(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "demo.feature")

	store := NewFeatureStore()
	feature := &gherkin.Feature{
		Title: "Demo",
		Scenarios: []gherkin.Scenario{
			{
				Title: "S1",
				Steps: []gherkin.Step{
					{Keyword: "Когда", Text: "выполняю шаг"},
				},
			},
		},
	}

	if err := store.Save(path, feature); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := store.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if got.Title != feature.Title {
		t.Fatalf("unexpected loaded feature title: %q", got.Title)
	}
}

func TestFeatureStore_Discover(t *testing.T) {
	tmp := t.TempDir()
	dirA := filepath.Join(tmp, "a")
	dirB := filepath.Join(tmp, "b")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatalf("mkdir a failed: %v", err)
	}
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatalf("mkdir b failed: %v", err)
	}

	featureA := filepath.Join(dirA, "first.feature")
	featureB := filepath.Join(dirB, "second.feature")
	other := filepath.Join(tmp, "README.txt")

	if err := os.WriteFile(featureA, []byte("Функционал: A\nСценарий: S\nКогда шаг\n"), 0o644); err != nil {
		t.Fatalf("write feature a failed: %v", err)
	}
	if err := os.WriteFile(featureB, []byte("Функционал: B\nСценарий: S\nКогда шаг\n"), 0o644); err != nil {
		t.Fatalf("write feature b failed: %v", err)
	}
	if err := os.WriteFile(other, []byte("ignore"), 0o644); err != nil {
		t.Fatalf("write other file failed: %v", err)
	}

	store := NewFeatureStore()
	files, err := store.Discover(tmp)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("unexpected discovered files: %#v", files)
	}
}
