package queue_service

import "os"

// StorageDriver - The interface used by the LocalStorage plugin to write/read files
// from the local file system
type StorageDriver interface {
	EnsureDirExists(string) error
	ExistsOrFail(string) error
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
	DeleteFile(string) error
}
