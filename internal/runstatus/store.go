package runstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

type Entry struct {
	Path          string `json:"path"`
	CaseID        string `json:"case_id,omitempty"`
	RunID         string `json:"run_id,omitempty"`
	Success       bool   `json:"success"`
	Status        string `json:"status,omitempty"`
	Message       string `json:"message"`
	DurationMS    int    `json:"duration_ms"`
	FailedStep    *int   `json:"failed_step,omitempty"`
	StepDurations []int  `json:"step_durations,omitempty"`
	Runner        string `json:"runner"`
	At            string `json:"at"`
	ExampleIndex  *int   `json:"example_index,omitempty"`
}

type Store struct {
	path string
	mu   sync.Mutex
}

var (
	runStatusCreateTemp = os.CreateTemp
	runStatusRename     = os.Rename
	runStatusRemove     = os.Remove
)

func Open(projectRoot string) (*Store, error) {
	dir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("create run status dir: %w", err)
	}
	return &Store{path: filepath.Join(dir, "run_status.json")}, nil
}

func (s *Store) Record(entry Entry) error {
	return s.RecordBatch([]Entry{entry})
}

// RecordBatch prepends multiple entries in one atomic read-modify-write.
// Batch order is completion order; stored newest-first (last completed first).
func (s *Store) RecordBatch(batch []Entry) error {
	if len(batch) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := s.loadLocked()
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	prepended := make([]Entry, len(batch))
	for i, entry := range batch {
		if entry.At == "" {
			entry.At = now
		}
		prepended[len(batch)-1-i] = entry
	}
	entries = append(prepended, entries...)
	payload, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("encode run status: %w", err)
	}
	return writeAtomicLocked(s.path, append(payload, '\n'))
}

func (s *Store) Latest(path string) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := s.loadLocked()
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Path == path {
			copy := entry
			return &copy, nil
		}
	}
	return nil, nil
}

// List returns the most recent run entries (newest first).
func (s *Store) List(limit int) ([]Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := s.loadLocked()
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (s *Store) load() ([]Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) loadLocked() ([]Entry, error) {
	payload, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("read run status: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(payload, &entries); err != nil {
		return nil, fmt.Errorf("decode run status: %w", err)
	}
	return entries, nil
}

func writeAtomicLocked(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create run status dir: %w", err)
		}
	}
	tmp, err := runStatusCreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create run status temp: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = runStatusRemove(tmpPath)
		return fmt.Errorf("write run status temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = runStatusRemove(tmpPath)
		return fmt.Errorf("close run status temp: %w", err)
	}
	backupPath := ""
	backupCreated := false
	if _, err := os.Stat(path); err == nil {
		backupPath = path + ".bak"
		if err := runStatusRemove(backupPath); err != nil && !os.IsNotExist(err) {
			_ = runStatusRemove(tmpPath)
			return fmt.Errorf("remove stale run status backup: %w", err)
		}
		if err := runStatusRename(path, backupPath); err != nil {
			_ = runStatusRemove(tmpPath)
			return fmt.Errorf("backup run status: %w", err)
		}
		backupCreated = true
	} else if !os.IsNotExist(err) {
		_ = runStatusRemove(tmpPath)
		return fmt.Errorf("stat run status: %w", err)
	}
	if err := runStatusRename(tmpPath, path); err != nil {
		if backupCreated {
			if restoreErr := runStatusRename(backupPath, path); restoreErr != nil {
				if copyErr := restoreRunStatusBackupByCopy(backupPath, path); copyErr != nil {
					_ = runStatusRemove(tmpPath)
					return fmt.Errorf("replace run status: %w (restore backup %q failed: %v, atomic copy restore failed: %w)", err, backupPath, restoreErr, copyErr)
				}
			}
			_ = runStatusRemove(backupPath)
		}
		_ = runStatusRemove(tmpPath)
		return fmt.Errorf("replace run status: %w", err)
	}
	if backupCreated {
		_ = runStatusRemove(backupPath)
	}
	return nil
}

func restoreRunStatusBackupByCopy(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := runStatusCreateTemp(filepath.Dir(dest), filepath.Base(dest)+".restore-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanupTemp := true
	defer func() {
		_ = tmp.Close()
		if cleanupTemp {
			_ = runStatusRemove(tmpPath)
		}
	}()
	if _, err := io.Copy(tmp, in); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := runStatusRename(tmpPath, dest); err != nil {
		return err
	}
	cleanupTemp = false
	return nil
}
