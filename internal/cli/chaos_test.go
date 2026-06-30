package cli

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestChaosDryRunManyFeatures(t *testing.T) {
	tmp := t.TempDir()
	const fileCount = 40
	for i := 0; i < fileCount; i++ {
		path := filepath.Join(tmp, fmt.Sprintf("f-%03d.feature", i))
		content := fmt.Sprintf("Функционал: Batch %d\nСценарий: S\n\tДопустим открыт \"https://example.com/%d\"\n", i, i)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write feature: %v", err)
		}
	}
	if err := RunRun([]string{tmp, "--dry-run"}); err != nil {
		t.Fatalf("dry-run over %d features: %v", fileCount, err)
	}
}

func TestChaosDryRunRandomFeatureNames(t *testing.T) {
	tmp := t.TempDir()
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("%s-%d.feature", randomName(rng), i)
		path := filepath.Join(tmp, name)
		content := "Функционал: X\nСценарий: Y\n\tКогда нажимаю \"#ok\"\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	if err := RunRun([]string{tmp, "--dry-run"}); err != nil {
		t.Fatalf("dry-run: %v", err)
	}
}

func randomName(rng *rand.Rand) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	n := 3 + rng.Intn(8)
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}
