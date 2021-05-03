// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mock_dev_storage

import (
	"encoding/json"
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

func (m *MockScribble) ensureCollectionExists(collection string) {
	if _, ok := m.store[collection]; !ok {
		m.store[collection] = make(map[string][]byte)
	}
}

func (m *MockScribble) SetCollection(collection string, items map[string]interface{}) {
	m.ensureCollectionExists(collection)

	for key, value := range items {
		itemBytes, _ := json.Marshal(value)
		m.store[collection][key] = itemBytes
	}
}

func (m *MockScribble) GetCollection(collection string) map[string]interface{} {
	if m.store[collection] == nil {
		return nil
	}

	tmpMap := make(map[string]interface{})
	for key, value := range m.store[collection] {
		var item interface{}
		_ = json.Unmarshal(value, &item)
		tmpMap[key] = item
	}

	return tmpMap
}

func (m *MockScribble) ClearStore() {
	m.store = make(map[string]map[string][]byte)
}

func (m *MockScribble) Read(collection string, key string, v interface{}) error {
	if m.store[collection] == nil {
		return os.ErrNotExist
	}

	if item, ok := m.store[collection][key]; ok {
		json.Unmarshal(item, &v)

		return nil
	}

	// TODO: This should produce the same error as stat()
	// in the case of a file not existing
	return os.ErrNotExist
}

func (m *MockScribble) Write(collection string, key string, v interface{}) error {
	m.ensureCollectionExists(collection)

	bytes, _ := json.Marshal(v)

	m.store[collection][key] = bytes

	return nil
}

// Delete
func (m *MockScribble) Delete(collection string, key string) error {
	// m.ensureCollectionExists(collection)

	if m.store[collection] == nil {
		return os.ErrNotExist
	}

	if _, ok := m.store[collection][key]; ok {
		m.store[collection][key] = nil
		return nil
	}

	return os.ErrNotExist
}

// ReadAll
func (m *MockScribble) ReadAll(collection string) ([]string, error) {
	if m.store[collection] == nil {
		return nil, os.ErrNotExist
	}

	vals := make([]string, 0)
	for _, value := range m.store[collection] {
		vals = append(vals, string(value))
	}

	return vals, nil
}

func NewMockScribble() *MockScribble {
	return &MockScribble{
		store: make(map[string]map[string][]byte),
	}
}
