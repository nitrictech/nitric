package storage_plugin

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// StorageDriver - The interface used by the LocalStorage plugin to write/read files
// from the local file system
type StorageDriver interface {
	EnsureDirExists(string) error
	WriteFile(string, []byte, os.FileMode) error
}

// DefaultStorageDriver - The Storage Driver to be used when creating
// a new Local Storage Plugin using the New() method
type DefaultStorageDriver struct {
	StorageDriver
}

// EnsureDirExists - Recurively creates directories for the given path
func (s *DefaultStorageDriver) EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// WriteFile - Writes the given byte array to the given path
func (s *DefaultStorageDriver) WriteFile(file string, contents []byte, fileMode os.FileMode) error {
	return ioutil.WriteFile(file, contents, fileMode)
}

// LocalStoragePlugin - The Nitric Storage Plugin for local development work
// Primarily used as part of the nitric run CLI function
type LocalStoragePlugin struct {
	sdk.UnimplementedStoragePlugin
	storageDriver StorageDriver
	storeDir      string
}

// Put will create a new item or overwrite an existing item in storage
func (s *LocalStoragePlugin) Put(bucket string, key string, payload []byte) error {
	bucketName := fmt.Sprintf("%s%s/", s.storeDir, bucket)

	if err := s.storageDriver.EnsureDirExists(bucketName); err == nil {
		fileName := fmt.Sprintf("%s%s", bucketName, key)

		if err := s.storageDriver.WriteFile(fileName, payload, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Get will retrieve an item from Storage
func (s *LocalStoragePlugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// New creates a new default StoragePlugin
func New() (sdk.StoragePlugin, error) {
	storeDir := utils.GetEnv("LOCAL_BLOB_DIR", "/nitric/buckets/")
	defaultDriver := &DefaultStorageDriver{}

	return &LocalStoragePlugin{
		storeDir:      storeDir,
		storageDriver: defaultDriver,
	}, nil
}

// NewWithStorageDriver creates a new StoragePlugin with the given StorageDriver
// primarily used for mock testing
func NewWithStorageDriver(driver StorageDriver) (sdk.StoragePlugin, error) {
	storeDir := utils.GetEnv("LOCAL_BLOB_DIR", "/nitric/buckets/")

	return &LocalStoragePlugin{
		storeDir:      storeDir,
		storageDriver: driver,
	}, nil
}
