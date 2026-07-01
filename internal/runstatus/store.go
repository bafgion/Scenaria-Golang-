package runstatus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

type Entry struct {
	Path         string `json:"path"`
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	DurationMS   int    `json:"duration_ms"`
	FailedStep   *int   `json:"failed_step,omitempty"`
	Runner       string `json:"runner"`
	At           string `json:"at"`
	ExampleIndex *int   `json:"example_index,omitempty"`
}

type Store struct {
	path string
}

func Open(projectRoot string) (*Store, error) {
	dir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("create run status dir: %w", err)
	}
	return &Store{path: filepath.Join(dir, "run_status.json")}, nil
}

func (s *Store) Record(entry Entry) error {
	entries, err := s.load()
	if err != nil {
		return err
	}
	if entry.At == "" {
		entry.At = time.Now().UTC().Format(time.RFC3339)
	}
	entries = append([]Entry{entry}, entries...)
	payload, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("encode run status: %w", err)
	}
	return os.WriteFile(s.path, append(payload, '\n'), 0o644)
}

func (s *Store) Latest(path string) (*Entry, error) {
	entries, err := s.load()
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
	entries, err := s.load()
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (s *Store) load() ([]Entry, error) {
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
