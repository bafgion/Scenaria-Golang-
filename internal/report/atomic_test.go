package report

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAtomicReplacesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeAtomic(path, []byte("new")); err != nil {
		t.Fatalf("writeAtomic: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestWriteAtomicPreservesPreviousOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")
	original := []byte("original report")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	err := writeAtomicWithRename(path, []byte("new report"), func(_, _ string) error {
		return fmt.Errorf("simulated rename failure")
	})
	if err == nil {
		t.Fatal("expected rename failure")
	}

	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read preserved report: %v", readErr)
	}
	if string(got) != string(original) {
		t.Fatalf("previous report changed on failure: got %q want %q", got, original)
	}
}
