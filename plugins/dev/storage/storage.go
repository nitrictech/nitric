package storage_service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// StorageDriver - The interface used by the LocalStorage plugin to write/read files
// from the local file system
type StorageDriver interface {
	EnsureDirExists(string) error
	ExistsOrFail(string) error
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
	DeleteFile(string) error
}

// DefaultStorageDriver - The Storage Driver to be used when creating
// a new Local Storage Plugin using the New() method
type DefaultStorageDriver struct {
}

// Ensure interface methods are implemented.
var _ StorageDriver = (*DefaultStorageDriver)(nil)

// EnsureDirExists - Recursively creates directories for the given path
func (s *DefaultStorageDriver) EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// ExistsOrFail - Returns an error if the provided path doesn't exist in the file system
func (s *DefaultStorageDriver) ExistsOrFail(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	return nil
}

// WriteFile - Writes the given byte array to the given path
func (s *DefaultStorageDriver) WriteFile(file string, contents []byte, fileMode os.FileMode) error {
	return ioutil.WriteFile(file, contents, fileMode)
}

// ReadFile - Reads from the given path and returns the byte array
func (s *DefaultStorageDriver) ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

// DeleteFile - Deletes the file at the given path
func (s *DefaultStorageDriver) DeleteFile(file string) error {
	return os.Remove(file)
}

// LocalStorageService - The Nitric Storage Plugin for local development work
// Primarily used as part of the nitric run CLI function
type LocalStorageService struct {
	storageDriver StorageDriver
	storeDir      string
}

func (s *LocalStorageService) getBucketName(bucket string) string {
	return fmt.Sprintf("%s%s/", s.storeDir, bucket)
}

func (s *LocalStorageService) getFilename(bucket string, key string) string {
	return fmt.Sprintf("%s%s", s.getBucketName(bucket), key)
}

// Put will create a new item or overwrite an existing item in storage
func (s *LocalStorageService) Write(bucket string, key string, payload []byte) error {
	bucketName := s.getBucketName(bucket)

	if err := s.storageDriver.EnsureDirExists(bucketName); err == nil {
		fileName := s.getFilename(bucket, key)

		if err := s.storageDriver.WriteFile(fileName, payload, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Get will retrieve an item from Storage
func (s *LocalStorageService) Read(bucket string, key string) ([]byte, error) {
	bucketName := s.getBucketName(bucket)

	if err := s.storageDriver.EnsureDirExists(bucketName); err == nil {
		fileName := s.getFilename(bucket, key)

		if err := s.storageDriver.ExistsOrFail(fileName); err != nil {
			return nil, fmt.Errorf("object not found %s", fileName)
		}

		return s.storageDriver.ReadFile(fileName)
	} else {
		return nil, err
	}
}

// Delete will delete an item from Storage
func (s *LocalStorageService) Delete(bucket string, key string) error {
	return s.storageDriver.DeleteFile(s.getFilename(bucket, key))
}

// New creates a new default StorageService
func New() (sdk.StorageService, error) {
	storeDir := utils.GetEnv("LOCAL_BLOB_DIR", "/nitric/buckets/")
	defaultDriver := &DefaultStorageDriver{}

	return &LocalStorageService{
		storeDir:      storeDir,
		storageDriver: defaultDriver,
	}, nil
}

// NewWithStorageDriver creates a new StorageService with the given StorageDriver
// primarily used for mock testing
func NewWithStorageDriver(driver StorageDriver) (sdk.StorageService, error) {
	storeDir := utils.GetEnv("LOCAL_BLOB_DIR", "/nitric/buckets/")

	return &LocalStorageService{
		storeDir:      storeDir,
		storageDriver: driver,
	}, nil
}
