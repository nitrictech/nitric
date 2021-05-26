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

// The Digital Ocean App Platform HTTP gateway plugin
package appplatform_service

import (
	"fmt"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
	"github.com/valyala/fasthttp"
)

type HttpGateway struct {
	address string
	server  *fasthttp.Server
	sdk.UnimplementedGatewayPlugin
}

func httpHandler(handler handler.TriggerHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		httpTrigger := triggers.FromHttpRequest(ctx)
		response, err := handler.HandleHttpRequest(httpTrigger)

		if err != nil {
			ctx.Error(fmt.Sprintf("Error handling HTTP Request: %v", err), 500)
			return
		}
		// responseBody, _ := ioutil.ReadAll(response.Body)
		if response.Header != nil {
			// Set headers...
			response.Header.VisitAll(func(key []byte, val []byte) {
				ctx.Response.Header.AddBytesKV(key, val)
			})
		}

		// Avoid content length header duplication
		ctx.Response.Header.Del("Content-Length")
		ctx.Response.SetStatusCode(response.StatusCode)
		ctx.Response.SetBody(response.Body)
	}
}

func (s *HttpGateway) Start(handler handler.TriggerHandler) error {
	s.server = &fasthttp.Server{
		Handler: httpHandler(handler),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *HttpGateway) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &HttpGateway{
		address: address,
	}, nil
}
