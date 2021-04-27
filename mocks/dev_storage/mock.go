package mock_dev_storage

import (
	"fmt"
	"os"
	"strings"
)

type MockStorageDriverOptions struct {
	EnsureDirExistsError error
	WriteFileError       error
	ReadFileError        error
	Directories          []string
	StoredItems          map[string][]byte
}

type MockStorageDriver struct {
	ensureDirExistsError error
	writeFileError       error
	readFileError        error
	directories          []string
	storedItems          map[string][]byte
}

func (m *MockStorageDriver) GetStoredItems() map[string][]byte {
	return m.storedItems
}

func (m *MockStorageDriver) EnsureDirExists(dir string) error {
	if m.ensureDirExistsError != nil {
		return m.ensureDirExistsError
	}

	if m.directories == nil {
		m.directories = make([]string, 0)
	}

	m.directories = append(m.directories, dir)

	return nil
}

func (m *MockStorageDriver) ExistsOrFail(path string) error {
	if m.storedItems == nil {
		return fmt.Errorf("%s does not exist", path)
	}

	if m.storedItems[path] == nil {
		return fmt.Errorf("%s does not exist", path)
	}

	return nil
}

func (m *MockStorageDriver) WriteFile(file string, contents []byte, fileMode os.FileMode) error {
	if m.writeFileError != nil {
		return m.writeFileError
	}

	// Capture for later eval
	if m.storedItems == nil {
		m.storedItems = make(map[string][]byte)
	}

	pathParts := strings.Split(file, "/")

	directory := strings.Join(pathParts[:len(pathParts)-1], "/") + "/"
	for _, dir := range m.directories {
		if dir == directory {
			m.storedItems[file] = contents
			return nil
		}
	}

	return fmt.Errorf("Cannot create file as directory does not exist")
}

func (m *MockStorageDriver) ReadFile(file string) ([]byte, error) {
	if m.storedItems == nil {
		m.storedItems = make(map[string][]byte)
	}

	if m.storedItems[file] == nil {
		return nil, fmt.Errorf("No such file: %s", file)
	}

	return m.storedItems[file], nil
}

func (m *MockStorageDriver) DeleteFile(file string) error {
	return nil
}

func NewMockStorageDriver(opts *MockStorageDriverOptions) *MockStorageDriver {
	return &MockStorageDriver{
		ensureDirExistsError: opts.EnsureDirExistsError,
		writeFileError:       opts.WriteFileError,
		readFileError:        opts.ReadFileError,
		directories:          opts.Directories,
		storedItems:          opts.StoredItems,
	}
}
