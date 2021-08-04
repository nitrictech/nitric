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
package base_http

import (
	"fmt"
	"time"

	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/utils"
	"github.com/nitric-dev/membrane/pkg/worker"

	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	"github.com/valyala/fasthttp"
)

type HttpMiddleware func(*fasthttp.RequestCtx, worker.Worker) bool

type BaseHttpGateway struct {
	address string
	server  *fasthttp.Server
	gateway.UnimplementedGatewayPlugin

	// Middleware for handling events
	// return bool will indicate whether to continue
	// to the next (default) behaviour or not...
	mw HttpMiddleware
}

func (s *BaseHttpGateway) httpHandler(pool worker.WorkerPool) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		wrkr, err := pool.GetWorker()

		if err != nil {
			ctx.Error("Unable to get worker to handle request", 500)
			return
		}

		if s.mw != nil {
			if !s.mw(ctx, wrkr) {
				// middleware has indicated that is has processed the request
				// so we can exit here
				return
			}
		}

		httpTrigger := triggers.FromHttpRequest(ctx)
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

func (s *BaseHttpGateway) Start(pool worker.WorkerPool) error {
	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         s.httpHandler(pool),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *BaseHttpGateway) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New(mw HttpMiddleware) (gateway.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &BaseHttpGateway{
		address: address,
		mw:      mw,
	}, nil
}
