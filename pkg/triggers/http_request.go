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
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/pkg/utils"
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
	// Extracted params (if configured) from the path
	Params map[string]string
	// URL query parameters
	Query map[string][]string
}

func (*HttpRequest) GetTriggerType() TriggerType {
	return TriggerType_Request
}

const paramToken = ":"

func parsePathParams(exp string, path string) (map[string]string, error) {
	pathParts := strings.Split(path, "/")
	expPathParts := strings.Split(exp, "/")
	params := make(map[string]string)

	for i, s := range expPathParts {
		if strings.HasPrefix(s, paramToken) {
			paramName := strings.Replace(s, paramToken, "", -1)
			params[paramName] = pathParts[i]
		} else if pathParts[i] != expPathParts[i] {
			return nil, fmt.Errorf("unable to match path")
		}
	}

	return params, nil
}

// FromHttpRequest (constructs a HttpRequest source type from a HttpRequest)
func FromHttpRequest(ctx *fasthttp.RequestCtx) *HttpRequest {
	headerCopy := make(map[string][]string)
	queryArgs := make(map[string][]string)

	ctx.Request.Header.VisitAll(func(key []byte, val []byte) {
		keyString := string(key)

		if strings.ToLower(keyString) == "host" {
			// Don't copy the host header
			headerCopy["X-Forwarded-For"] = []string{string(val)}
		} else {
			headerCopy[string(key)] = []string{string(val)}
		}
	})

	ctx.Request.Header.VisitAllCookie(func(key []byte, val []byte) {
		headerCopy[string(key)] = append(headerCopy[string(key)], string(val))
	})

	ctx.QueryArgs().VisitAll(func(key []byte, val []byte) {
		k := string(key)

		if queryArgs[k] == nil {
			queryArgs[k] = make([]string, 0)
		}

		queryArgs[k] = append(queryArgs[k], string(val))
	})

	// Get gateway path if one is configured to parse on
	gwPath := utils.GetEnv("GW_PATH", "")

	var params = make(map[string]string)
	if gwPath != "" {
		// Check if we should parse path params
		params, _ = parsePathParams(gwPath, string(ctx.Path()))
	}

	return &HttpRequest{
		Header: headerCopy,
		Body:   ctx.Request.Body(),
		Method: string(ctx.Method()),
		Path:   string(ctx.Path()),
		Query:  queryArgs,
		Params: params,
	}
}
