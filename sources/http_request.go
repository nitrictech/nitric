package sources

import (
	"io"
	"net/http"
)

// HttpRequest - Storage information that captures a HTTP Request
type HttpRequest struct {
	// The original Headers
	Header http.Header
	// The original body stread
	Body io.ReadCloser
	// The original method
	Method string
	// The original path
	Path string
}

func (*HttpRequest) GetSourceType() SourceType {
	return SourceType_Request
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(r *http.Request) *HttpRequest {
	return &HttpRequest{
		Header: r.Header,
		Body:   r.Body,
		Method: r.Method,
		Path:   r.URL.Path,
	}
}
