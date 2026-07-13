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

var (
	atomicCreateTemp = os.CreateTemp
	atomicRename     = os.Rename
	atomicRemove     = os.Remove
)

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
		return writeFileAtomic(abs, []byte(content), 0o644)
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
	if strings.TrimSpace(newName) != "" {
		fileName, err := normalizeFeatureFileName(newName)
		if err != nil {
			return "", err
		}
		target, err = s.ensureInsideProject(filepath.Join(dir, fileName))
		if err != nil {
			return "", err
		}
	} else {
		base := strings.TrimSuffix(filepath.Base(srcAbs), ext)
		target, err = nextAvailableFeaturePath(filepath.Join(dir, base+"-copy"+ext))
		if err != nil {
			return "", err
		}
		if target, err = s.ensureInsideProject(target); err != nil {
			return "", err
		}
	}
	if err := s.withWriteLock(func() error {
		return copyFeatureFileExclusive(srcAbs, target)
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
			fileName, err := normalizeFeatureFileName(filepath.Base(src))
			if err != nil {
				return fmt.Errorf("import %s: %w", filepath.Base(src), err)
			}
			target, err := nextAvailableFeaturePath(filepath.Join(destAbs, fileName))
			if err != nil {
				return fmt.Errorf("import %s: %w", filepath.Base(src), err)
			}
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
	if strings.TrimSpace(newName) == "" {
		return "", fmt.Errorf("new file name is required")
	}
	fileName, err := normalizeFeatureFileName(newName)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(srcAbs)
	target := filepath.Join(dir, fileName)
	if strings.EqualFold(filepath.Clean(srcAbs), filepath.Clean(target)) {
		return srcAbs, nil
	}
	if err := s.withWriteLock(func() error {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("file already exists: %s", fileName)
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
	return writeFileAtomic(path, append(data, '\n'), 0o644)
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
	return paths.PathGuard{Root: project}.ResolveExistingOrNew(target)
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

func normalizeFeatureFileName(name string) (string, error) {
	if name != strings.TrimSpace(name) {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	if name == "" {
		return "", fmt.Errorf("new feature name is required")
	}
	if name != strings.TrimRight(name, " .") {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	if filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	if filepath.Base(name) != name || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	if strings.ContainsAny(name, `<>:"|?*`) {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	for _, r := range name {
		if r < 32 {
			return "", fmt.Errorf("invalid file name %q", name)
		}
	}
	if !strings.EqualFold(filepath.Ext(name), ".feature") {
		name += ".feature"
	}
	if strings.TrimSuffix(name, filepath.Ext(name)) == "" {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	if isReservedWindowsFileName(name) {
		return "", fmt.Errorf("invalid file name %q", name)
	}
	return name, nil
}

func isReservedWindowsFileName(name string) bool {
	base := strings.TrimRight(name, " .")
	if idx := strings.IndexRune(base, '.'); idx >= 0 {
		base = base[:idx]
	}
	base = strings.ToUpper(base)
	switch base {
	case "CON", "PRN", "AUX", "NUL":
		return true
	}
	if len(base) == 4 {
		prefix := base[:3]
		suffix := base[3]
		if (prefix == "COM" || prefix == "LPT") && suffix >= '1' && suffix <= '9' {
			return true
		}
	}
	return false
}

func nextAvailableFeaturePath(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, nil
	} else if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	ext := filepath.Ext(path)
	for i := 2; i < 10000; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		} else if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("no available file name for %s", filepath.Base(path))
}

func copyFeatureFile(src, dest string) error {
	return copyFeatureFileExclusive(src, dest)
}

func copyFeatureFileExclusive(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	cleanup := true
	defer func() {
		_ = out.Close()
		if cleanup {
			_ = os.Remove(dest)
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	cleanup = false
	return nil
}

func restoreFileBackupByCopy(src, dest string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := atomicCreateTemp(filepath.Dir(dest), filepath.Base(dest)+".restore-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = atomicRemove(tmpPath)
		}
	}()
	if _, err := io.Copy(tmp, in); err != nil {
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := atomicRename(tmpPath, dest); err != nil {
		return err
	}
	cleanupTemp = false
	return nil
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}
	tmp, err := atomicCreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file %q: %w", path, err)
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = atomicRemove(tmpPath)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("write temp file %q: %w", path, err)
	}
	if err := tmp.Chmod(perm); err != nil {
		return fmt.Errorf("chmod temp file %q: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file %q: %w", path, err)
	}
	backupPath := ""
	backupCreated := false
	if _, err := os.Stat(path); err == nil {
		backupPath = path + ".bak"
		if err := atomicRemove(backupPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale backup %q: %w", backupPath, err)
		}
		if err := atomicRename(path, backupPath); err != nil {
			return fmt.Errorf("backup file %q: %w", path, err)
		}
		backupCreated = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat file %q: %w", path, err)
	}
	if err := atomicRename(tmpPath, path); err != nil {
		if backupCreated {
			if restoreErr := atomicRename(backupPath, path); restoreErr != nil {
				if copyErr := restoreFileBackupByCopy(backupPath, path, perm); copyErr != nil {
					return fmt.Errorf("replace file %q: %w (restore backup %q failed: %v, atomic copy restore failed: %w)", path, err, backupPath, restoreErr, copyErr)
				}
			}
			_ = atomicRemove(backupPath)
		}
		return fmt.Errorf("replace file %q: %w", path, err)
	}
	if backupCreated {
		_ = atomicRemove(backupPath)
	}
	cleanupTemp = false
	return nil
}
