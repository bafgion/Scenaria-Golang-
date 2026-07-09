package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

type featurePathConfiner func(path string) (string, error)
type projectReadLocker func(func() error) error
type projectWriteLocker func(func() error) error
type projectRootGetter func() string

// FileOperationService owns filesystem read/write operations for feature files.
type FileOperationService struct {
	confineFeaturePath featurePathConfiner
	projectRoot        projectRootGetter
	withReadLock       projectReadLocker
	withWriteLock      projectWriteLocker
}

func NewFileOperationService(
	confiner featurePathConfiner,
	projectRoot projectRootGetter,
	readLock projectReadLocker,
	writeLock projectWriteLocker,
) *FileOperationService {
	return &FileOperationService{
		confineFeaturePath: confiner,
		projectRoot:        projectRoot,
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

func (s *FileOperationService) DeleteFeature(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("feature path is required")
	}
	abs, err := s.ensureInsideProject(path)
	if err != nil {
		return err
	}
	if err := s.withWriteLock(func() error { return os.Remove(abs) }); err != nil {
		return fmt.Errorf("delete feature: %w", err)
	}
	return nil
}

func (s *FileOperationService) DuplicateFeature(path, newName string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("feature path is required")
	}
	srcAbs, err := s.confineFeaturePath(path)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(srcAbs)
	ext := filepath.Ext(srcAbs)
	target := ""
	if name := strings.TrimSpace(newName); name != "" {
		name = strings.TrimSuffix(name, ext)
		name = strings.TrimSuffix(name, ".feature")
		if name == "" {
			return "", fmt.Errorf("new feature name is required")
		}
		target = filepath.Join(dir, name+ext)
		if _, err := os.Stat(target); err == nil {
			return "", fmt.Errorf("file already exists: %s", filepath.Base(target))
		} else if !os.IsNotExist(err) {
			return "", err
		}
	} else {
		base := strings.TrimSuffix(filepath.Base(srcAbs), ext)
		target = filepath.Join(dir, base+"-copy"+ext)
		for i := 2; i < 100; i++ {
			if _, err := os.Stat(target); os.IsNotExist(err) {
				break
			}
			target = filepath.Join(dir, fmt.Sprintf("%s-copy-%d%s", base, i, ext))
		}
	}
	if err := s.withWriteLock(func() error {
		payload, err := os.ReadFile(srcAbs)
		if err != nil {
			return fmt.Errorf("read feature: %w", err)
		}
		return os.WriteFile(target, payload, 0o644)
	}); err != nil {
		return "", fmt.Errorf("write duplicate: %w", err)
	}
	return target, nil
}

func (s *FileOperationService) MoveFeature(src, destDir string) (string, error) {
	srcAbs, err := s.ensureInsideProject(src)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(filepath.Ext(srcAbs), ".feature") {
		return "", fmt.Errorf("only .feature files can be moved")
	}
	destAbs, err := s.ensureInsideProject(destDir)
	if err != nil {
		return "", err
	}
	target := filepath.Join(destAbs, filepath.Base(srcAbs))
	if strings.EqualFold(filepath.Clean(srcAbs), filepath.Clean(target)) {
		return srcAbs, nil
	}
	if err := s.withWriteLock(func() error {
		info, err := os.Stat(destAbs)
		if err != nil {
			return fmt.Errorf("destination folder: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("destination must be a folder")
		}
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("file already exists: %s", filepath.Base(target))
		} else if !os.IsNotExist(err) {
			return err
		}
		return os.Rename(srcAbs, target)
	}); err != nil {
		return "", fmt.Errorf("move feature: %w", err)
	}
	return target, nil
}

func (s *FileOperationService) ImportFeatures(destDir string, paths []string) ([]string, error) {
	destAbs, err := s.ensureInsideProject(destDir)
	if err != nil {
		return nil, err
	}
	imported := make([]string, 0, len(paths))
	if err := s.withWriteLock(func() error {
		if info, err := os.Stat(destAbs); err != nil || !info.IsDir() {
			if err != nil {
				return fmt.Errorf("destination folder: %w", err)
			}
			return fmt.Errorf("destination must be a folder")
		}
		for _, src := range paths {
			src = strings.TrimSpace(src)
			if src == "" || !strings.EqualFold(filepath.Ext(src), ".feature") {
				continue
			}
			target := uniqueFeaturePath(filepath.Join(destAbs, filepath.Base(src)))
			if err := copyFeatureFile(src, target); err != nil {
				return fmt.Errorf("import %s: %w", filepath.Base(src), err)
			}
			imported = append(imported, target)
		}
		return nil
	}); err != nil {
		return imported, err
	}
	if len(imported) == 0 {
		return nil, fmt.Errorf("no .feature files to import")
	}
	return imported, nil
}

func (s *FileOperationService) RenameFeature(path, newName string) (string, error) {
	srcAbs, err := s.ensureInsideProject(path)
	if err != nil {
		return "", err
	}
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return "", fmt.Errorf("new file name is required")
	}
	if strings.ContainsAny(newName, `/\`) {
		return "", fmt.Errorf("invalid file name %q", newName)
	}
	if !strings.EqualFold(filepath.Ext(newName), ".feature") {
		newName += ".feature"
	}
	dir := filepath.Dir(srcAbs)
	target := filepath.Join(dir, newName)
	if strings.EqualFold(filepath.Clean(srcAbs), filepath.Clean(target)) {
		return srcAbs, nil
	}
	if err := s.withWriteLock(func() error {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("file already exists: %s", newName)
		} else if !os.IsNotExist(err) {
			return err
		}
		return os.Rename(srcAbs, target)
	}); err != nil {
		return "", fmt.Errorf("rename feature: %w", err)
	}
	return target, nil
}

func (s *FileOperationService) SaveFeatureDraft(featurePath, content string) error {
	path, err := featureDraftPath(featurePath)
	if err != nil {
		return err
	}
	payload := featureDraft{
		FeaturePath: featurePath,
		Content:     content,
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func (s *FileOperationService) LoadFeatureDraft(featurePath string) (string, error) {
	path, err := featureDraftPath(featurePath)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var draft featureDraft
	if err := json.Unmarshal(data, &draft); err != nil {
		return "", err
	}
	return draft.Content, nil
}

func (s *FileOperationService) ClearFeatureDraft(featurePath string) error {
	path, err := featureDraftPath(featurePath)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *FileOperationService) ensureInsideProject(target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", fmt.Errorf("path is required")
	}
	project := ""
	if s != nil && s.projectRoot != nil {
		project = s.projectRoot()
	}
	if project == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	projectAbs, err := filepath.Abs(project)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(projectAbs, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path is outside the project")
	}
	return abs, nil
}

type featureDraft struct {
	FeaturePath string `json:"featurePath"`
	Content     string `json:"content"`
	UpdatedAt   string `json:"updatedAt"`
}

func featureDraftPath(featurePath string) (string, error) {
	featurePath = strings.TrimSpace(featurePath)
	if featurePath == "" {
		return "", fmt.Errorf("feature path is required")
	}
	root := paths.InferProjectRoot([]string{featurePath})
	if root == "" {
		return "", fmt.Errorf("project root not found for %q", featurePath)
	}
	rel, err := filepath.Rel(root, featurePath)
	if err != nil {
		rel = filepath.Base(featurePath)
	}
	safe := strings.NewReplacer(string(os.PathSeparator), "_", ":", "_").Replace(rel)
	dir := filepath.Join(root, ".scenaria", "drafts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create drafts dir: %w", err)
	}
	return filepath.Join(dir, safe+".json"), nil
}

func uniqueFeaturePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	ext := filepath.Ext(path)
	for i := 2; i < 100; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
	return filepath.Join(dir, base+"-copy"+ext)
}

func copyFeatureFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
