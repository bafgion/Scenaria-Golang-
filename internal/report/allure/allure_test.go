package allure

import "testing"

func TestWritePlaceholder(t *testing.T) {
	dir := t.TempDir()
	if err := WritePlaceholder(Options{OutputDir: dir}); err != nil {
		t.Fatalf("WritePlaceholder: %v", err)
	}
}
