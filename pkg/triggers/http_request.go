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
	Header map[string][]string
	// The original body stream
	Body []byte
	// The original method
	Method string
	// The original path
	Path string
	// URL
	URL string
	// URL query parameters
	Query map[string][]string
	// Path parameters
	Params map[string]string
}

func (*HttpRequest) GetTriggerType() TriggerType {
	return TriggerType_Request
}

func HttpHeaders(rh *fasthttp.RequestHeader) map[string][]string {
	headerCopy := make(map[string][]string)

	rh.VisitAll(func(key []byte, val []byte) {
		keyString := string(key)

		if strings.ToLower(keyString) == "host" {
			// Don't copy the host header
			headerCopy["X-Forwarded-For"] = []string{string(val)}
		} else {
			headerCopy[string(key)] = append(headerCopy[string(key)], string(val))
		}
	})

	return headerCopy
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(rc *fasthttp.RequestCtx) *HttpRequest {
	headerCopy := HttpHeaders(&rc.Request.Header)
	queryArgs := make(map[string][]string)

	rc.QueryArgs().VisitAll(func(key []byte, val []byte) {
		k := string(key)

		if queryArgs[k] == nil {
			queryArgs[k] = make([]string, 0)
		}

		queryArgs[k] = append(queryArgs[k], string(val))
	})

	return &HttpRequest{
		Header: headerCopy,
		Body:   rc.Request.Body(),
		Method: string(rc.Method()),
		URL:    rc.URI().String(),
		Path:   string(rc.URI().PathOriginal()),
		Query:  queryArgs,
	}
}
