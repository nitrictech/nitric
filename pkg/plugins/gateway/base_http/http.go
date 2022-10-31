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
	"context"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/utils"
	"github.com/nitrictech/nitric/pkg/worker"
)

type HttpMiddleware func(*fasthttp.RequestCtx, worker.WorkerPool) bool

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
	return func(rc *fasthttp.RequestCtx) {
		if s.mw != nil {
			if !s.mw(rc, pool) {
				// middleware has indicated that is has processed the request
				// so we can exit here
				return
			}
		}

		ctx, httpTrigger := triggers.FromHttpRequest(rc)

		wrkr, err := pool.GetWorker(&worker.GetWorkerOptions{
			Http: httpTrigger,
		})
		if err != nil {
			rc.Error("Unable to get worker to handle request", 500)
			return
		}

		response, err := wrkr.HandleHttpRequest(ctx, httpTrigger)
		if err != nil {
			rc.Error(fmt.Sprintf("Error handling HTTP Request: %v", err), 500)
			return
		}

		if response.Header != nil {
			response.Header.CopyTo(&rc.Response.Header)
		}

		// Avoid content length header duplication
		rc.Response.Header.Del("Content-Length")
		rc.Response.SetStatusCode(response.StatusCode)
		rc.Response.SetBody(response.Body)
	}
}

func (s *BaseHttpGateway) Start(pool worker.WorkerPool) error {
	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         s.httpHandler(pool),
		ReadBufferSize:  8192,
	}

	return s.server.ListenAndServe(s.address)
}

func (s *BaseHttpGateway) Stop() error {
	tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if ok {
		_ = tp.ForceFlush(context.TODO())
	}

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
