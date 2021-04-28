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

package storage_service_test

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	storage_plugin "github.com/nitric-dev/membrane/plugins/storage/dev"
)

type MockStorageDriver struct {
	ensureDirExistsError error
	existsOrFailError    error
	writeFileError       error
	readFileError        error
	deleteFileError      error
	directories          []string
	storedItems          map[string][]byte
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
	if m.existsOrFailError != nil {
		return m.existsOrFailError
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
	if m.readFileError != nil {
		return nil, m.readFileError
	}

	pathParts := strings.Split(file, "/")

	directory := strings.Join(pathParts[:len(pathParts)-1], "/") + "/"
	for _, dir := range m.directories {
		if dir == directory {
			if object, ok := m.storedItems[file]; ok {
				return object, nil
			}
			return nil, fmt.Errorf("file %s not found in directory %s", file, directory)
		}
	}

	return nil, fmt.Errorf("unable to retrieve file, directory %s does not exist", directory)
}

func (m *MockStorageDriver) DeleteFile(file string) error {
	if m.deleteFileError != nil {
		return m.deleteFileError
	}

	pathParts := strings.Split(file, "/")

	directory := strings.Join(pathParts[:len(pathParts)-1], "/") + "/"
	for _, dir := range m.directories {
		if dir == directory {
			if _, ok := m.storedItems[file]; ok {
				delete(m.storedItems, file)
			}
			return nil
		}
	}

	return fmt.Errorf("failed to delete item, bucket directory [%s] no found", directory)
}

var _ = Describe("Storage", func() {
	Context("Put", func() {
		// Test Put methods...
		Context("Given the storage driver is functioning without error", func() {
			workingDriver := &MockStorageDriver{}
			mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)
			It("Should store the provided byte array", func() {
				err := mockStoragePlugin.Write("test", "test", []byte("Test"))
				By("Not returning an error")
				Expect(err).To(BeNil())

				Expect(workingDriver.storedItems["/nitric/buckets/test/test"]).To(BeEquivalentTo([]byte("Test")))
			})
		})

		Context("Given the storage driver cannot create directories", func() {
			faultyDriver := &MockStorageDriver{
				ensureDirExistsError: fmt.Errorf("error creating directory"),
			}
			mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(faultyDriver)
			It("Should return an error", func() {

				err := mockStoragePlugin.Write("test", "test", []byte("Test"))
				By("By returning an error")
				Expect(err).ToNot(BeNil())

				By("Not creating any directories")
				Expect(faultyDriver.directories).To(BeNil())
			})
		})

		Context("Given the storage driver cannot create files", func() {
			faultyDriver := &MockStorageDriver{
				writeFileError: fmt.Errorf("error creating file"),
			}
			mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(faultyDriver)
			It("Should store the provided byte array", func() {
				err := mockStoragePlugin.Write("test", "test", []byte("Test"))
				By("By returning an error")
				Expect(err).ToNot(BeNil())

				By("By creating the requested directory")
				Expect(faultyDriver.directories).ToNot(BeNil())

				By("Not storing any items")
				Expect(faultyDriver.storedItems).To(BeNil())
			})
		})
	})

	Context("Get", func() {
		When("The storage driver is functioning without error", func() {
			When("The bucket directory exists", func() {
				When("The object file exists", func() {
					workingDriver := &MockStorageDriver{
						directories: []string{"/nitric/buckets/test-bucket/"},
						storedItems: map[string][]byte{"/nitric/buckets/test-bucket/test-key": []byte("Test")},
					}
					mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

					It("Should retrieve the object", func() {
						object, err := mockStoragePlugin.Read("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).To(BeNil())

						By("Returning the object")
						Expect(object).To(Equal([]byte("Test")))
					})
				})
				When("The object file doesn't exist", func() {
					workingDriver := &MockStorageDriver{
						directories: []string{"/nitric/buckets/test-bucket/"},
					}
					mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

					It("Should return an error", func() {
						object, err := mockStoragePlugin.Read("test-bucket", "test-key")
						By("Returning an error")
						Expect(err).ToNot(BeNil())

						By("Returning a nil object")
						Expect(object).To(BeNil())
					})
				})
			})
			When("The bucket directory doesn't exist", func() {
				workingDriver := &MockStorageDriver{}
				mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

				It("Should return an error", func() {
					object, err := mockStoragePlugin.Read("test-bucket", "test-key")
					By("Returning an error")
					Expect(err).ToNot(BeNil())

					By("Returning a nil object")
					Expect(object).To(BeNil())
				})
			})
		})

		When("The storage driver cannot read directories", func() {
			faultyDriver := &MockStorageDriver{
				ensureDirExistsError: fmt.Errorf("error creating directory"),
			}
			mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(faultyDriver)
			It("Should return an error", func() {

				object, err := mockStoragePlugin.Read("test", "test")
				By("By returning an error")
				Expect(err).ToNot(BeNil())

				By("Not return the object")
				Expect(object).To(BeNil())
			})
		})

		When("The storage driver cannot read files", func() {
			faultyDriver := &MockStorageDriver{
				readFileError: fmt.Errorf("error reading file"),
			}
			mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(faultyDriver)
			It("Should return an error", func() {

				object, err := mockStoragePlugin.Read("test", "test")
				By("By returning an error")
				Expect(err).ToNot(BeNil())

				By("Not return the object")
				Expect(object).To(BeNil())
			})
		})
	})

	Context("Delete", func() {
		When("The storage driver is functioning without error", func() {
			When("The bucket directory exists", func() {
				When("The object file exists", func() {
					workingDriver := &MockStorageDriver{
						directories: []string{"/nitric/buckets/test-bucket/"},
						storedItems: map[string][]byte{"/nitric/buckets/test-bucket/test-key": []byte("Test")},
					}
					mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

					It("Should delete the object", func() {
						err := mockStoragePlugin.Delete("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).To(BeNil())
					})
				})
				When("The object file doesn't exist", func() {
					workingDriver := &MockStorageDriver{
						directories: []string{"/nitric/buckets/test-bucket/"},
					}
					mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

					It("Should skip deleting the file", func() {
						err := mockStoragePlugin.Delete("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).To(BeNil())
					})
				})
			})
			When("The bucket directory doesn't exist", func() {
				workingDriver := &MockStorageDriver{
					directories: []string{},
				}
				mockStoragePlugin, _ := storage_plugin.NewWithStorageDriver(workingDriver)

				It("Should return an error", func() {
					err := mockStoragePlugin.Delete("test-bucket", "test-key")
					By("Returning an error")
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})

})
