package gui

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type projectReadLock func(func() error) error

// ProjectService owns project discovery/index read operations.
type ProjectService struct {
	withReadLock projectReadLock
}

func NewProjectService(withReadLock projectReadLock) *ProjectService {
	return &ProjectService{withReadLock: withReadLock}
}

func (s *ProjectService) ProjectInfo(path string, version uint64) (ProjectInfo, error) {
	if path == "" {
		return ProjectInfo{}, fmt.Errorf("no project opened")
	}
	if s == nil || s.withReadLock == nil {
		return ProjectInfo{}, fmt.Errorf("project service is not configured")
	}
	var info ProjectInfo
	err := s.withReadLock(func() error {
		store := scenario.NewFeatureStore()
		files, err := store.Discover(path)
		if err != nil {
			return err
		}
		tags := collectProjectTags(store, files)
		featureTags := collectFeatureTags(store, files)
		info = ProjectInfo{
			Path:        path,
			Features:    files,
			Tags:        tags,
			FeatureTags: featureTags,
			Version:     version,
		}
		return nil
	})
	if err != nil {
		return ProjectInfo{}, err
	}
	return info, nil
}
