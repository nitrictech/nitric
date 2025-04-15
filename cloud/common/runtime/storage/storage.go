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
