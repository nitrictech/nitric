package files_test

import (
	"os"
	"testing"

	"github.com/nitrictech/nitric/cli/pkg/files"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyDir(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(fs afero.Fs) error
		src         string
		dst         string
		expectError bool
		validate    func(t *testing.T, fs afero.Fs)
	}{
		{
			name: "successfully copy directory with files and subdirectories",
			setup: func(fs afero.Fs) error {
				if err := fs.MkdirAll("src", 0755); err != nil {
					return err
				}
				if err := fs.MkdirAll("src/subdir", 0755); err != nil {
					return err
				}
				if err := fs.MkdirAll("src/subdir/nested", 0755); err != nil {
					return err
				}

				if err := afero.WriteFile(fs, "src/file1.txt", []byte("content1"), 0644); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/file2.txt", []byte("content2"), 0755); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/subdir/file3.txt", []byte("content3"), 0600); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/subdir/nested/file4.txt", []byte("content4"), 0666); err != nil {
					return err
				}
				return nil
			},
			src:         "src",
			dst:         "dst",
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs) {
				exists, err := afero.Exists(fs, "dst")
				require.NoError(t, err)
				assert.True(t, exists)

				content, err := afero.ReadFile(fs, "dst/file1.txt")
				require.NoError(t, err)
				assert.Equal(t, "content1", string(content))

				content, err = afero.ReadFile(fs, "dst/file2.txt")
				require.NoError(t, err)
				assert.Equal(t, "content2", string(content))

				content, err = afero.ReadFile(fs, "dst/subdir/file3.txt")
				require.NoError(t, err)
				assert.Equal(t, "content3", string(content))

				content, err = afero.ReadFile(fs, "dst/subdir/nested/file4.txt")
				require.NoError(t, err)
				assert.Equal(t, "content4", string(content))

				info, err := fs.Stat("dst/file1.txt")
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0644), info.Mode().Perm())

				info, err = fs.Stat("dst/file2.txt")
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

				info, err = fs.Stat("dst/subdir/file3.txt")
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

				info, err = fs.Stat("dst/subdir/nested/file4.txt")
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0666), info.Mode().Perm())
			},
		},
		{
			name: "copy to existing destination directory",
			setup: func(fs afero.Fs) error {
				if err := fs.MkdirAll("src", 0755); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/file.txt", []byte("content"), 0644); err != nil {
					return err
				}
				if err := fs.MkdirAll("dst", 0755); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "dst/existing.txt", []byte("existing"), 0644); err != nil {
					return err
				}
				return nil
			},
			src:         "src",
			dst:         "dst",
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "dst/file.txt")
				require.NoError(t, err)
				assert.Equal(t, "content", string(content))

				content, err = afero.ReadFile(fs, "dst/existing.txt")
				require.NoError(t, err)
				assert.Equal(t, "existing", string(content))
			},
		},
		{
			name: "copy empty directory",
			setup: func(fs afero.Fs) error {
				return fs.MkdirAll("src", 0755)
			},
			src:         "src",
			dst:         "dst",
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs) {
				exists, err := afero.Exists(fs, "dst")
				require.NoError(t, err)
				assert.True(t, exists)

				isDir, err := afero.IsDir(fs, "dst")
				require.NoError(t, err)
				assert.True(t, isDir)
			},
		},
		{
			name:        "source directory does not exist",
			setup:       func(fs afero.Fs) error { return nil },
			src:         "nonexistent",
			dst:         "dst",
			expectError: true,
			validate:    func(t *testing.T, fs afero.Fs) {},
		},
		{
			name: "source is a file, not a directory",
			setup: func(fs afero.Fs) error {
				return afero.WriteFile(fs, "src", []byte("content"), 0644)
			},
			src:         "src",
			dst:         "dst",
			expectError: true,
			validate:    func(t *testing.T, fs afero.Fs) {},
		},
		{
			name: "copy directory with symbolic links (should copy as regular files)",
			setup: func(fs afero.Fs) error {
				if err := fs.MkdirAll("src", 0755); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/target.txt", []byte("target content"), 0644); err != nil {
					return err
				}
				// Note: Afero's MemMapFs doesn't support symlinks, so we'll test with regular files
				if err := afero.WriteFile(fs, "src/link.txt", []byte("link content"), 0644); err != nil {
					return err
				}
				return nil
			},
			src:         "src",
			dst:         "dst",
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "dst/target.txt")
				require.NoError(t, err)
				assert.Equal(t, "target content", string(content))

				content, err = afero.ReadFile(fs, "dst/link.txt")
				require.NoError(t, err)
				assert.Equal(t, "link content", string(content))
			},
		},
		{
			name: "copy directory with special characters in names",
			setup: func(fs afero.Fs) error {
				if err := fs.MkdirAll("src", 0755); err != nil {
					return err
				}
				if err := fs.MkdirAll("src/dir with spaces", 0755); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/file with spaces.txt", []byte("content"), 0644); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, "src/dir with spaces/file.txt", []byte("nested content"), 0644); err != nil {
					return err
				}
				return nil
			},
			src:         "src",
			dst:         "dst",
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "dst/file with spaces.txt")
				require.NoError(t, err)
				assert.Equal(t, "content", string(content))

				content, err = afero.ReadFile(fs, "dst/dir with spaces/file.txt")
				require.NoError(t, err)
				assert.Equal(t, "nested content", string(content))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			if tt.setup != nil {
				err := tt.setup(fs)
				require.NoError(t, err)
			}

			err := files.CopyDir(fs, tt.src, tt.dst)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, fs)
		})
	}
}

func TestCopyDir_Integration(t *testing.T) {
	fs := afero.NewMemMapFs()

	setupComplexStructure := func() error {
		dirs := []string{
			"src",
			"src/a",
			"src/a/b",
			"src/a/b/c",
			"src/x",
			"src/x/y",
		}

		for _, dir := range dirs {
			if err := fs.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		files := map[string]string{
			"src/root.txt":        "root content",
			"src/a/file1.txt":     "file1 content",
			"src/a/b/file2.txt":   "file2 content",
			"src/a/b/c/file3.txt": "file3 content",
			"src/x/file4.txt":     "file4 content",
			"src/x/y/file5.txt":   "file5 content",
		}

		for path, content := range files {
			if err := afero.WriteFile(fs, path, []byte(content), 0644); err != nil {
				return err
			}
		}

		return nil
	}

	require.NoError(t, setupComplexStructure())

	err := files.CopyDir(fs, "src", "dst")
	require.NoError(t, err)

	expectedFiles := map[string]string{
		"dst/root.txt":        "root content",
		"dst/a/file1.txt":     "file1 content",
		"dst/a/b/file2.txt":   "file2 content",
		"dst/a/b/c/file3.txt": "file3 content",
		"dst/x/file4.txt":     "file4 content",
		"dst/x/y/file5.txt":   "file5 content",
	}

	for path, expectedContent := range expectedFiles {
		content, err := afero.ReadFile(fs, path)
		require.NoError(t, err, "Failed to read %s", path)
		assert.Equal(t, expectedContent, string(content), "Content mismatch for %s", path)
	}

	expectedDirs := []string{
		"dst",
		"dst/a",
		"dst/a/b",
		"dst/a/b/c",
		"dst/x",
		"dst/x/y",
	}

	for _, dir := range expectedDirs {
		exists, err := afero.Exists(fs, dir)
		require.NoError(t, err)
		assert.True(t, exists, "Directory %s should exist", dir)

		isDir, err := afero.IsDir(fs, dir)
		require.NoError(t, err)
		assert.True(t, isDir, "%s should be a directory", dir)
	}
}

func TestCopyDir_FilePermissions(t *testing.T) {
	fs := afero.NewMemMapFs()

	require.NoError(t, fs.MkdirAll("src", 0755))
	require.NoError(t, afero.WriteFile(fs, "src/readonly.txt", []byte("readonly"), 0444))
	require.NoError(t, afero.WriteFile(fs, "src/executable.txt", []byte("executable"), 0755))
	require.NoError(t, afero.WriteFile(fs, "src/normal.txt", []byte("normal"), 0644))

	err := files.CopyDir(fs, "src", "dst")
	require.NoError(t, err)

	info, err := fs.Stat("dst/readonly.txt")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0444), info.Mode().Perm())

	info, err = fs.Stat("dst/executable.txt")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

	info, err = fs.Stat("dst/normal.txt")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
}
