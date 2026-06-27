package scenario

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type FeatureStore struct{}

func NewFeatureStore() *FeatureStore {
	return &FeatureStore{}
}

func (s *FeatureStore) Load(path string) (*gherkin.Feature, error) {
	feature, err := gherkin.ParseFeatureFile(path)
	if err != nil {
		return nil, fmt.Errorf("load feature %q: %w", path, err)
	}
	return feature, nil
}

func (s *FeatureStore) Save(path string, feature *gherkin.Feature) error {
	if err := gherkin.SaveFeatureFile(path, feature); err != nil {
		return fmt.Errorf("save feature %q: %w", path, err)
	}
	return nil
}

func (s *FeatureStore) Discover(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".feature") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover features in %q: %w", root, err)
	}
	sort.Strings(files)
	return files, nil
}
