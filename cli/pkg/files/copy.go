package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// CopyDir copies the contents of src directory to dst directory using Afero.
// It recursively copies all files and subdirectories while preserving
// file permissions and directory structure.
func CopyDir(fs afero.Fs, src, dst string) error {
	srcInfo, err := fs.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	err = fs.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return afero.Walk(fs, src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return fs.MkdirAll(dstPath, info.Mode())
		} else {
			return copyFile(fs, path, dstPath, info.Mode())
		}
	})
}

func copyFile(fs afero.Fs, src, dst string, mode os.FileMode) error {
	srcFile, err := fs.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := fs.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return fs.Chmod(dst, mode)
}
