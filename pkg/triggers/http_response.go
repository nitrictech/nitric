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

	"github.com/valyala/fasthttp"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
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

// FromTriggerResponse (constructs a HttpResponse from a FaaS TriggerResponse)
func FromTriggerResponse(triggerResponse *pb.TriggerResponse) (*HttpResponse, error) {
	// FIXME: This will panic if the incorrect response type is provided
	httpContext := triggerResponse.GetHttp()
	if httpContext != nil {
		fasthttpHeader := &fasthttp.ResponseHeader{}

		if httpContext.GetHeaders() != nil {
			for key, val := range httpContext.GetHeaders() {
				for _, v := range val.Value {
					fasthttpHeader.Add(key, v)
				}
			}
		}

		return &HttpResponse{
			Header:     fasthttpHeader,
			StatusCode: int(httpContext.Status),
			Body:       triggerResponse.GetData(),
		}, nil
	}

	return nil, fmt.Errorf("TriggerResponse does not container HTTP Context")
}
