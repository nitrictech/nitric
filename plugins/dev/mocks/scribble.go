package mocks

import (
	"encoding/json"
	"internal/oserror"
	"os"
)

type MockScribble struct {
	// ifaces.ScribbleIface
	// Read(string, string, interface{}) error
	// ReadAll(string) ([]string, error)
	// Write(string, string, interface{}) error
	// Delete(string, string) error
	store      map[string]map[string][]byte
	readErr    error
	readAllErr error
	writeErr   error
	deleteErr  error
}

func (m *MockScribbleDriver) ensureCollectionExists(collection string) {
	if _, ok := m.store[collection]; !ok {
		m.store[collection] = make(map[string][]byte)
	}
}

func (m *MockScribbleDriver) clearStore() {
	m.store = make(map[string]map[string][]byte)
}

func (m *MockScribbleDriver) Read(collection string, key string, v interface{}) error {
	if m[collection] == nil {
		return oserror.ErrNotExist
	}

	if item, ok := m.store[collection][key]; ok {
		json.Unmarshal(item, &v)

		return nil
	}

	// TODO: This should produce the same error as stat()
	// in the case of a file not existing
	return oserror.ErrNotExist
}

func (m *MockScribbleDriver) Write(collection string, key string, v interface{}) error {
	m.ensureCollectionExists(collection)

	os.IsNotExist()

	bytes, _ := json.Marshal(v)

	m.store[collection][key] = bytes

	return nil
}

// Delete
func (m *MockScribbleDriver) Delete(collection string, key string) error {
	// m.ensureCollectionExists(collection)

	if m.store[collection] == nil {
		return oserror.ErrNotExist
	}

	if _, ok := m.store[collection][key]; ok {
		m.store[collection][key] = nil
		return nil
	}

	return oserror.ErrNotExist
}

// ReadAll
func (m *MockScribble) ReadAll(collection string) ([]string, error) {
	if m.store[collection] == nil {
		return nil, oserror.ErrNotExist
	}

	vals := make([]string, 0)
	for _, value := range m.store[collection] {
		vals = append(vals, string(value))
	}

	return vals, nil
}
