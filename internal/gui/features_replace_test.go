package gui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceInProjectStepsOnly(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.feature")
	second := filepath.Join(dir, "b.feature")
	if err := os.WriteFile(first, []byte("Функционал: UI\nСценарий: A\n\tДопустим открыт \"https://old.example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("Функционал: UI\nСценарий: B\n\tДопустим открыт \"https://old.example.com/page\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ReplaceInProject(ProjectReplaceRequest{
		Find:    "old.example.com",
		Replace: "new.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.FilesChanged != 2 || got.Replacements != 2 {
		t.Fatalf("got %+v", got)
	}
	a, _ := os.ReadFile(first)
	b, _ := os.ReadFile(second)
	if !strings.Contains(string(a), "new.example.com") || !strings.Contains(string(b), "new.example.com") {
		t.Fatalf("files not updated: %q %q", a, b)
	}
}

func TestReplaceInProjectSkipsFeatureTitle(t *testing.T) {
	dir := t.TempDir()
	feature := filepath.Join(dir, "demo.feature")
	payload := "Функционал: old title\nСценарий: old title\n\tДопустим открыт \"https://x.com\"\n"
	if err := os.WriteFile(feature, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old title", Replace: "new title"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Replacements != 0 {
		t.Fatalf("expected no replacements in headers, got %+v", got)
	}
}

func TestReplaceInProjectReadFailureDoesNotModifyAnyFile(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	second := writeReplaceFeature(t, dir, "b.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origRead := projectReplaceReadFile
	origWrite := projectReplaceWriteFile
	defer func() {
		projectReplaceReadFile = origRead
		projectReplaceWriteFile = origWrite
	}()
	writes := 0
	projectReplaceReadFile = func(path string) ([]byte, error) {
		if filepath.Clean(path) == filepath.Clean(second) {
			return nil, errors.New("simulated read failure")
		}
		return os.ReadFile(path)
	}
	projectReplaceWriteFile = func(path string, data []byte, perm os.FileMode) error {
		writes++
		return writeFileAtomic(path, data, perm)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "read b.feature") {
		t.Fatalf("expected read failure, got %v", err)
	}
	if writes != 0 {
		t.Fatalf("expected no writes before all files are read, got %d", writes)
	}
	assertReplaceFileContains(t, first, "old.example.com")
	assertReplaceFileContains(t, second, "old.example.com")
}

func TestReplaceInProjectWriteFailureRollsBackCommittedFiles(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	second := writeReplaceFeature(t, dir, "b.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origWrite := projectReplaceWriteFile
	defer func() { projectReplaceWriteFile = origWrite }()
	projectReplaceWriteFile = func(path string, data []byte, perm os.FileMode) error {
		if filepath.Clean(path) == filepath.Clean(second) {
			return errors.New("simulated write failure")
		}
		return writeFileAtomic(path, data, perm)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "write b.feature") {
		t.Fatalf("expected write failure, got %v", err)
	}
	assertReplaceFileContains(t, first, "old.example.com")
	assertReplaceFileContains(t, second, "old.example.com")
}

func TestReplaceInProjectExternalChangeBeforeFirstCommitAbortsWithoutOverwrite(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	second := writeReplaceFeature(t, dir, "b.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origRead := projectReplaceReadFile
	origWrite := projectReplaceWriteFile
	defer func() {
		projectReplaceReadFile = origRead
		projectReplaceWriteFile = origWrite
	}()
	reads := map[string]int{}
	writes := 0
	projectReplaceReadFile = func(path string) ([]byte, error) {
		clean := filepath.Clean(path)
		reads[clean]++
		if clean == filepath.Clean(first) && reads[clean] == 2 {
			if err := os.WriteFile(first, []byte("external same size edit old.example.com\n"), 0o644); err != nil {
				return nil, err
			}
		}
		return os.ReadFile(path)
	}
	projectReplaceWriteFile = func(path string, data []byte, perm os.FileMode) error {
		writes++
		return writeFileAtomic(path, data, perm)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "file changed externally: a.feature") {
		t.Fatalf("expected external modification error, got %v", err)
	}
	if writes != 0 {
		t.Fatalf("expected no writes after first-file conflict, got %d", writes)
	}
	assertReplaceFileContains(t, first, "external same size edit old.example.com")
	assertReplaceFileContains(t, second, "old.example.com")
}

func TestReplaceInProjectExternalChangeBeforeLaterCommitRollsBackCommittedFiles(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	second := writeReplaceFeature(t, dir, "b.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origRead := projectReplaceReadFile
	defer func() { projectReplaceReadFile = origRead }()
	reads := map[string]int{}
	projectReplaceReadFile = func(path string) ([]byte, error) {
		clean := filepath.Clean(path)
		reads[clean]++
		if clean == filepath.Clean(second) && reads[clean] == 2 {
			if err := os.WriteFile(second, []byte("external later edit old.example.com\n"), 0o644); err != nil {
				return nil, err
			}
		}
		return os.ReadFile(path)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "file changed externally: b.feature") {
		t.Fatalf("expected external modification error, got %v", err)
	}
	assertReplaceFileContains(t, first, "old.example.com")
	assertReplaceFileContains(t, second, "external later edit old.example.com")
}

func TestReplaceInProjectExternalSameSizeEditIsDetected(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origRead := projectReplaceReadFile
	defer func() { projectReplaceReadFile = origRead }()
	reads := map[string]int{}
	projectReplaceReadFile = func(path string) ([]byte, error) {
		clean := filepath.Clean(path)
		reads[clean]++
		if clean == filepath.Clean(first) && reads[clean] == 2 {
			raw, err := os.ReadFile(first)
			if err != nil {
				return nil, err
			}
			changed := strings.Replace(string(raw), "old.example.com", "alt.example.com", 1)
			if len(changed) != len(string(raw)) {
				return nil, errors.New("test setup expected same-size edit")
			}
			if err := os.WriteFile(first, []byte(changed), 0o644); err != nil {
				return nil, err
			}
		}
		return os.ReadFile(path)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "file changed externally: a.feature") {
		t.Fatalf("expected same-size external modification error, got %v", err)
	}
	assertReplaceFileContains(t, first, "alt.example.com")
}

func TestReplaceInProjectExternalDeleteBeforeCommitAborts(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origRead := projectReplaceReadFile
	defer func() { projectReplaceReadFile = origRead }()
	reads := map[string]int{}
	projectReplaceReadFile = func(path string) ([]byte, error) {
		clean := filepath.Clean(path)
		reads[clean]++
		if clean == filepath.Clean(first) && reads[clean] == 2 {
			if err := os.Remove(first); err != nil {
				return nil, err
			}
		}
		return os.ReadFile(path)
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "verify a.feature before replace") {
		t.Fatalf("expected verify failure after external delete, got %v", err)
	}
	if _, statErr := os.Stat(first); !os.IsNotExist(statErr) {
		t.Fatalf("externally deleted file should remain deleted, stat err=%v", statErr)
	}
}

func TestReplaceInProjectRollbackFailureIsReported(t *testing.T) {
	dir := t.TempDir()
	first := writeReplaceFeature(t, dir, "a.feature", "old.example.com")
	second := writeReplaceFeature(t, dir, "b.feature", "old.example.com")
	svc := openReplaceTestProject(t, dir)

	origWrite := projectReplaceWriteFile
	defer func() { projectReplaceWriteFile = origWrite }()
	projectReplaceWriteFile = func(path string, data []byte, perm os.FileMode) error {
		switch {
		case filepath.Clean(path) == filepath.Clean(second):
			return errors.New("simulated write failure")
		case filepath.Clean(path) == filepath.Clean(first) && strings.Contains(string(data), "old.example.com"):
			return errors.New("simulated rollback failure")
		default:
			return writeFileAtomic(path, data, perm)
		}
	}

	_, err := svc.ReplaceInProject(ProjectReplaceRequest{Find: "old.example.com", Replace: "new.example.com"})
	if err == nil || !strings.Contains(err.Error(), "rollback failed") {
		t.Fatalf("expected rollback failure, got %v", err)
	}
	assertReplaceFileContains(t, first, "new.example.com")
	assertReplaceFileContains(t, second, "old.example.com")
}

func writeReplaceFeature(t *testing.T, dir, name, url string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	content := "Р¤СѓРЅРєС†РёРѕРЅР°Р»: UI\nРЎС†РµРЅР°СЂРёР№: A\n\tР”РѕРїСѓСЃС‚РёРј РѕС‚РєСЂС‹С‚ \"https://" + url + "\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func openReplaceTestProject(t *testing.T, dir string) *Service {
	t.Helper()
	svc := NewService()
	if _, err := svc.OpenProject(dir); err != nil {
		t.Fatal(err)
	}
	return svc
}

func assertReplaceFileContains(t *testing.T, path, want string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), want) {
		t.Fatalf("%s does not contain %q: %q", path, want, string(raw))
	}
}
