package gui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

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

func (s *Service) SaveFeatureDraft(featurePath, content string) error {
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

func (s *Service) LoadFeatureDraft(featurePath string) (string, error) {
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

func (s *Service) ClearFeatureDraft(featurePath string) error {
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
