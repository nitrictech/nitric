package triggers

import (
	"github.com/valyala/fasthttp"
)

// HttpRequest - Storage information that captures a HTTP Request
type HttpResponse struct {
	// The original Headers
	Header *fasthttp.ResponseHeader
	// The original body stream
	Body []byte
	// The original method
	StatusCode int
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpResponse(resp *fasthttp.Response) *HttpResponse {
	return &HttpResponse{
		Header:     &resp.Header,
		Body:       resp.Body(),
		StatusCode: resp.StatusCode(),
	}
}
