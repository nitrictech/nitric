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
	triggers2 "github.com/nitric-dev/membrane/pkg/triggers"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	worker2 "github.com/nitric-dev/membrane/pkg/worker"

	"github.com/nitric-dev/membrane/pkg/sdk"
	"github.com/valyala/fasthttp"
)

type HttpGateway struct {
	address string
	server  *fasthttp.Server
	sdk.UnimplementedGatewayPlugin
}

func httpHandler(pool worker2.WorkerPool) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		wrkr, err := pool.GetWorker()

		if err != nil {
			ctx.Error("Unable to get worker to handle request", 500)
			return
		}

		httpTrigger := triggers2.FromHttpRequest(ctx)
		response, err := wrkr.HandleHttpRequest(httpTrigger)

		if err != nil {
			ctx.Error(fmt.Sprintf("Error handling HTTP Request: %v", err), 500)
			return
		}

		if response.Header != nil {
			response.Header.CopyTo(&ctx.Response.Header)
		}

		// Avoid content length header duplication
		ctx.Response.Header.Del("Content-Length")
		ctx.Response.SetStatusCode(response.StatusCode)
		ctx.Response.SetBody(response.Body)
	}
}

func (s *HttpGateway) Start(pool worker2.WorkerPool) error {
	s.server = &fasthttp.Server{
		Handler: httpHandler(pool),
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
	address := utils2.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &HttpGateway{
		address: address,
	}, nil
}
