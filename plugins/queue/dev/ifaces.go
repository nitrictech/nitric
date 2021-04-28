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
