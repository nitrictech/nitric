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
