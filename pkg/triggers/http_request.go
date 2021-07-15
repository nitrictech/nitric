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
		Body:   ctx.Request.Body(),
		Method: string(ctx.Method()),
		Path:   string(ctx.Path()),
		Query:  queryArgs,
	}
}
