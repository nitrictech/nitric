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

package storage

import "fmt"

type Operation int

const (
	READ  Operation = iota
	WRITE           = iota
)

func (op Operation) String() string {
	// The order of this array must match the iota order above.
	return [2]string{"READ", "WRITE"}[op]
}

type FileInfo struct {
	Key string
}

type StorageService interface {
	Read(bucket string, key string) ([]byte, error)
	Write(bucket string, key string, object []byte) error
	Delete(bucket string, key string) error
	ListFiles(bucket string) ([]*FileInfo, error)
	PreSignUrl(bucket string, key string, operation Operation, expiry uint32) (string, error)
}

type UnimplementedStoragePlugin struct{}

var _ StorageService = (*UnimplementedStoragePlugin)(nil)

func (*UnimplementedStoragePlugin) Read(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Write(bucket string, key string, object []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Delete(bucket string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) ListFiles(bucket string) ([]*FileInfo, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) PreSignUrl(bucket string, key string, operation Operation, expiry uint32) (string, error) {
	return "", fmt.Errorf("UNIMPLEMENTED")
}
