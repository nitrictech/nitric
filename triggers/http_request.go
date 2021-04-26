package triggers

import (
	"github.com/valyala/fasthttp"
)

// HttpRequest - Storage information that captures a HTTP Request
type HttpRequest struct {
	// The original Headers
	Header *fasthttp.RequestHeader
	// The original body stream
	Body []byte
	// The original method
	Method string
	// The original path
	Path string
	// URL query parameters
	Query *fasthttp.Args
}

func (*HttpRequest) GetTriggerType() TriggerType {
	return TriggerType_Request
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(ctx *fasthttp.RequestCtx) *HttpRequest {
	return &HttpRequest{
		Header: &ctx.Request.Header,
		Body:   ctx.PostBody(),
		Method: string(ctx.Method()),
		Path:   string(ctx.Path()),
		Query:  ctx.QueryArgs(),
	}
}
