package ifaces

import (
	"fmt"
	"os"
)

type StorageDriver interface {
	EnsureDirExists(string) error
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
}

type UnimplementedStorageDriver struct {
	StorageDriver
}

func (s *UnimplementedStorageDriver) EnsureDirExists(string) error {
	return fmt.Errorf("Unimplemented!")
}

func (s *UnimplementedStorageDriver) WriteFile(string, []byte, os.FileMode) error {
	return fmt.Errorf("Unimplemented!")
}

func (s *UnimplementedStorageDriver) ReadFile(string) ([]byte, error) {
	return nil, fmt.Errorf("Unimplemented!")
}
