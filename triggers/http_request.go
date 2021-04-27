package triggers

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// HttpRequest - Storage information that captures a HTTP Request
type HttpRequest struct {
	// The original Headers
	// Header *fasthttp.RequestHeader
	Header map[string]string
	// The original body stream
	Body []byte
	// The original method
	Method string
	// The original path
	Path string
	// URL query parameters
	Query map[string]string
}

func (*HttpRequest) GetTriggerType() TriggerType {
	return TriggerType_Request
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(ctx *fasthttp.RequestCtx) *HttpRequest {
	headerCopy := make(map[string]string)
	queryArgs := make(map[string]string)

	ctx.Request.Header.VisitAll(func(key []byte, val []byte) {
		keyString := string(key)

		if strings.ToLower(keyString) == "host" {
			// Don't copy the host header
			headerCopy["X-Forwarded-For"] = string(val)
		} else {
			headerCopy[string(key)] = string(val)
		}
	})

	ctx.QueryArgs().VisitAll(func(key []byte, val []byte) {
		queryArgs[string(key)] = string(val)
	})

	return &HttpRequest{
		Header: headerCopy,
		Body:   ctx.PostBody(),
		Method: string(ctx.Method()),
		Path:   string(ctx.Path()),
		Query:  queryArgs,
	}
}
