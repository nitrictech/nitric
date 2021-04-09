package triggers

import (
	"io"
	"net/http"
	"strings"
)

// HttpRequest - Storage information that captures a HTTP Request
type HttpRequest struct {
	// The original Headers
	Header http.Header
	// The original body stream
	Body io.ReadCloser
	// The original method
	Method string
	// The original path
	Path string
}

func (*HttpRequest) GetTriggerType() TriggerType {
	return TriggerType_Request
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(r *http.Request) *HttpRequest {
	r.URL.RawPath
	return &HttpRequest{
		Header: r.Header,
		Body:   r.Body,
		Method: r.Method,
		Path:   strings.Join(r.URL.RawPath, "?", r.URL.RawQuery),
	}
}
