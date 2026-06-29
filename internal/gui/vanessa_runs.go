package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

func (s *Service) ListVanessaRunDirs(limit int) ([]string, error) {
	projectRoot := s.ProjectPath()
	if projectRoot == "" {
		return nil, fmt.Errorf("open a project folder first")
	}
	cfg, err := vanessa.LoadSettings(projectRoot)
	if err != nil {
		return nil, err
	}
	runsRoot := vanessa.ResolveRunsDir(cfg)
	entries, err := os.ReadDir(runsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	type runDir struct {
		path string
		mod  int64
	}
	dirs := make([]runDir, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "run-") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		dirs = append(dirs, runDir{
			path: filepath.Join(runsRoot, name),
			mod:  info.ModTime().UnixNano(),
		})
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].mod > dirs[j].mod })
	if limit <= 0 {
		limit = 20
	}
	if len(dirs) > limit {
		dirs = dirs[:limit]
	}
	out := make([]string, 0, len(dirs))
	for _, item := range dirs {
		out = append(out, item.path)
	}
	return out, nil
}
