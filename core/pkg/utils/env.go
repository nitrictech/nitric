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

package utils

import (
	"os"
	"path/filepath"
)

// GetEnv - Retrieve an environment variable with a fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetDevVolumePath - Returns the default directory to be used for local development plugins
// this directory points at a docker volume, used to share data between running containers.
func GetDevVolumePath() string {
	return GetEnv("NITRIC_DEV_VOLUME", "nitric/")
}

// GetRelativeDevPath - create a path, relative to the Dev Volume
func GetRelativeDevPath(relativePath string) string {
	return filepath.Join(GetDevVolumePath(), relativePath)
}
