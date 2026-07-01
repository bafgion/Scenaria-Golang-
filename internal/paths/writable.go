package paths

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WritableScenariaDir returns a writable .scenaria directory for projectRoot.
// When the project folder is read-only (e.g. Program Files), artifacts go to AppData.
func WritableScenariaDir(projectRoot string) (string, error) {
	projectRoot = strings.TrimSpace(projectRoot)
	if projectRoot == "" {
		return "", fmt.Errorf("project root is required")
	}
	local := ScenariaProjectDir(projectRoot)
	if dirWritable(local) {
		return local, nil
	}
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		abs = projectRoot
	}
	slug := projectScenariaSlug(abs)
	if legacy := LegacyProjectMirrorDir(slug); legacy != "" && dirWritable(legacy) {
		return legacy, nil
	}
	fallback := filepath.Join(AppDataDir(), "projects", slug)
	if err := os.MkdirAll(fallback, 0o755); err != nil {
		return "", fmt.Errorf("create scenaria dir: %w", err)
	}
	if !dirWritable(fallback) {
		return "", fmt.Errorf("scenaria dir is not writable: %s", fallback)
	}
	return fallback, nil
}

// ScenariaArtifactPath joins sub (e.g. "report.html") with the writable .scenaria dir.
func ScenariaArtifactPath(projectRoot, sub string) (string, error) {
	sub = strings.TrimSpace(sub)
	if sub == "" {
		return "", fmt.Errorf("artifact name is required")
	}
	dir, err := WritableScenariaDir(projectRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, sub), nil
}

// RemapScenariaArtifact keeps artifactPath when project .scenaria is writable; otherwise
// maps files under projectRoot/.scenaria into the AppData mirror.
func RemapScenariaArtifact(projectRoot, artifactPath string) string {
	artifactPath = strings.TrimSpace(artifactPath)
	if artifactPath == "" || strings.TrimSpace(projectRoot) == "" {
		return artifactPath
	}
	local := ScenariaProjectDir(projectRoot)
	if dirWritable(local) {
		return artifactPath
	}
	rel, ok := scenariaRelative(local, artifactPath)
	if !ok {
		return artifactPath
	}
	writable, err := WritableScenariaDir(projectRoot)
	if err != nil {
		return artifactPath
	}
	return filepath.Join(writable, rel)
}

func dirWritable(dir string) bool {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	probe := filepath.Join(dir, ".write_probe")
	file, err := os.Create(probe)
	if err != nil {
		return false
	}
	_ = file.Close()
	_ = os.Remove(probe)
	return true
}

func scenariaRelative(scenariaDir, artifactPath string) (string, bool) {
	absScenaria, err := filepath.Abs(scenariaDir)
	if err != nil {
		return "", false
	}
	absArtifact, err := filepath.Abs(artifactPath)
	if err != nil {
		return "", false
	}
	rel, err := filepath.Rel(absScenaria, absArtifact)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", false
	}
	return rel, true
}

func projectScenariaSlug(absProjectRoot string) string {
	clean := strings.ToLower(filepath.Clean(absProjectRoot))
	sum := sha256.Sum256([]byte(clean))
	hex := fmt.Sprintf("%x", sum[:6])
	base := filepath.Base(absProjectRoot)
	safe := strings.NewReplacer(" ", "_", ":", "", "\\", "_", "/", "_").Replace(base)
	if safe == "" || safe == "." {
		safe = "project"
	}
	return safe + "-" + hex
}
