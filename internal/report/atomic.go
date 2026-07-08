package report

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type renameFunc func(oldpath, newpath string) error

func writeAtomic(path string, data []byte) error {
	return writeAtomicWithRename(path, data, os.Rename)
}

func writeAtomicWithRename(path string, data []byte, rename renameFunc) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create report dir %q: %w", dir, err)
		}
	}

	tmp, err := os.CreateTemp(dir, ".scenaria-write-*")
	if err != nil {
		return fmt.Errorf("create temp report %q: %w", path, err)
	}
	tmpName := tmp.Name()
	removeTemp := true
	defer func() {
		_ = tmp.Close()
		if removeTemp {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("write temp report %q: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp report %q: %w", path, err)
	}
	if err := rename(tmpName, path); err != nil {
		return fmt.Errorf("replace report %q: %w", path, err)
	}
	removeTemp = false
	return nil
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
