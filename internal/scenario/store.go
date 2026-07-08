package scenario

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type featureCacheEntry struct {
	modTime time.Time
	size    int64
	feature *gherkin.Feature
}

type FeatureStore struct {
	mu    sync.RWMutex
	cache map[string]featureCacheEntry
}

func NewFeatureStore() *FeatureStore {
	return &FeatureStore{cache: make(map[string]featureCacheEntry)}
}

func (s *FeatureStore) Load(path string) (*gherkin.Feature, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat feature %q: %w", path, err)
	}
	s.mu.RLock()
	if cached, ok := s.cache[path]; ok && cached.modTime.Equal(stat.ModTime()) && cached.size == stat.Size() {
		feature := cached.feature
		s.mu.RUnlock()
		return feature, nil
	}
	s.mu.RUnlock()

	feature, err := gherkin.ParseFeatureFile(path)
	if err != nil {
		return nil, fmt.Errorf("load feature %q: %w", path, err)
	}
	s.mu.Lock()
	if s.cache == nil {
		s.cache = make(map[string]featureCacheEntry)
	}
	s.cache[path] = featureCacheEntry{
		modTime: stat.ModTime(),
		size:    stat.Size(),
		feature: feature,
	}
	s.mu.Unlock()
	return feature, nil
}

func (s *FeatureStore) Invalidate(path string) {
	s.mu.Lock()
	delete(s.cache, path)
	s.mu.Unlock()
}

func (s *FeatureStore) ClearCache() {
	s.mu.Lock()
	s.cache = make(map[string]featureCacheEntry)
	s.mu.Unlock()
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
