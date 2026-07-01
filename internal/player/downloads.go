package player

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (c *RunContext) allocateDownloadPath(suggested string) string {
	c.downloadSeq++
	name := strings.TrimSpace(suggested)
	if name == "" {
		name = "download"
	}
	name = filepath.Base(name)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	if base == "" {
		base = "download"
	}
	base = sanitizeDownloadName(base)
	unique := fmt.Sprintf("%s-%d-%d%s", base, c.runSeed, c.downloadSeq, ext)
	return filepath.Join(c.DownloadDir(), unique)
}

func sanitizeDownloadName(name string) string {
	name = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '.', r == '-', r == '_':
			return r
		default:
			return '_'
		}
	}, name)
	name = strings.Trim(name, "._-")
	if name == "" {
		return "download"
	}
	return name
}
