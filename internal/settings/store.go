package settings

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Store serializes settings read-modify-write operations.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: strings.TrimSpace(path)}
}

func DefaultStore() *Store {
	return NewStore(DefaultAppSettingsPath())
}

func (s *Store) Load() (*AppSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) Update(mutator func(*AppSettings) error) error {
	if mutator == nil {
		return fmt.Errorf("settings mutator is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, err := s.loadLocked()
	if err != nil {
		return err
	}
	if err := mutator(cfg); err != nil {
		return err
	}
	return s.saveLocked(cfg)
}

func (s *Store) loadLocked() (*AppSettings, error) {
	if s == nil || s.path == "" {
		return &AppSettings{Browser: "chromium"}, nil
	}
	cfg, err := LoadAppSettings(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &AppSettings{Browser: "chromium"}, nil
		}
		return nil, err
	}
	if strings.TrimSpace(cfg.Browser) == "" {
		cfg.Browser = "chromium"
	}
	return cfg, nil
}

func (s *Store) saveLocked(cfg *AppSettings) error {
	if s == nil || s.path == "" {
		return fmt.Errorf("settings path unavailable")
	}
	return SaveAppSettings(s.path, cfg)
}
