package gui

import (
	"fmt"
	"os"
)

type featurePathConfiner func(path string) (string, error)
type projectReadLocker func(func() error) error
type projectWriteLocker func(func() error) error

// FileOperationService owns filesystem read/write operations for feature files.
type FileOperationService struct {
	confineFeaturePath featurePathConfiner
	withReadLock       projectReadLocker
	withWriteLock      projectWriteLocker
}

func NewFileOperationService(
	confiner featurePathConfiner,
	readLock projectReadLocker,
	writeLock projectWriteLocker,
) *FileOperationService {
	return &FileOperationService{
		confineFeaturePath: confiner,
		withReadLock:       readLock,
		withWriteLock:      writeLock,
	}
}

func (s *FileOperationService) ReadFeature(path string) (string, error) {
	if s == nil || s.confineFeaturePath == nil || s.withReadLock == nil {
		return "", fmt.Errorf("file operation service is not configured")
	}
	abs, err := s.confineFeaturePath(path)
	if err != nil {
		return "", err
	}
	var payload []byte
	if err := s.withReadLock(func() error {
		var readErr error
		payload, readErr = os.ReadFile(abs)
		return readErr
	}); err != nil {
		return "", fmt.Errorf("read feature: %w", err)
	}
	return string(payload), nil
}

func (s *FileOperationService) SaveFeature(path, content string) error {
	if s == nil || s.confineFeaturePath == nil || s.withWriteLock == nil {
		return fmt.Errorf("file operation service is not configured")
	}
	abs, err := s.confineFeaturePath(path)
	if err != nil {
		return err
	}
	if err := s.withWriteLock(func() error {
		return os.WriteFile(abs, []byte(content), 0o644)
	}); err != nil {
		return fmt.Errorf("save feature: %w", err)
	}
	return nil
}
