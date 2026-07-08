package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFeatureStoreLoadUsesMtimeCache(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.feature")
	if err := os.WriteFile(path, []byte(`Функционал: Demo
  Сценарий: S
    Допустим открыт "https://example.com"
`), 0o644); err != nil {
		t.Fatal(err)
	}

	store := NewFeatureStore()
	first, err := store.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("expected cached feature pointer on unchanged file")
	}

	if err := os.WriteFile(path, []byte(`Функционал: Demo
  Сценарий: S2
    Допустим открыт "https://example.com"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	third, err := store.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if third == first {
		t.Fatal("expected cache miss after file change")
	}
}
