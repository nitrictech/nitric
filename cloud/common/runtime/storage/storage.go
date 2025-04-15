// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"mime"
	"net/http"
	"path/filepath"
)

// DetectContentType - detects content type using the file's extension or content.
// Extension matching is performed first, if the content type can't be determined by the file extension then the file contents will be inspected to attempt to determine the content type.
// If the content type cannot be determined it returns "application/octet-stream"
func DetectContentType(filename string, content []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}

	return http.DetectContentType(content)
}
