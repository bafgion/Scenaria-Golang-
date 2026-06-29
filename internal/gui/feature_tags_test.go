package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func TestCollectFeatureTags(t *testing.T) {
	root := t.TempDir()
	featurePath := filepath.Join(root, "smoke.feature")
	content := `@smoke
Функционал: Демо

  @ui
  Сценарий: Один
    Допустим пустой шаг
`
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	store := scenario.NewFeatureStore()
	tags := collectFeatureTags(store, []string{featurePath})
	if len(tags[featurePath]) == 0 {
		t.Fatalf("expected tags for %s, got %v", featurePath, tags)
	}
	if !gherkin.TagsInclude(tags[featurePath], "smoke") {
		t.Fatalf("missing smoke in %v", tags[featurePath])
	}
	if !gherkin.TagsInclude(tags[featurePath], "ui") {
		t.Fatalf("missing ui in %v", tags[featurePath])
	}
}
