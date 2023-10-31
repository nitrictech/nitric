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

package base_http

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/span"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type (
	HttpMiddleware   func(*fasthttp.RequestCtx, pool.WorkerPool) bool
	EventConstructor func(topicName string, ctx *fasthttp.RequestCtx) v1.TriggerRequest
)

type RouteRegister func(*router.Router, pool.WorkerPool)

const (
	DefaultTopicRoute              = "/x-nitric-topic/{name}"
	DefaultScheduleRoute           = "/x-nitric-schedule/{name}"
	DefaultBucketNotificationRoute = "/x-nitric-notification/bucket/{name}"
)

type BaseHttpGatewayOptions struct {
	// Middleware for handling events
	// return bool will indicate whether to continue
	// to the next (default) behaviour or not...
	Middleware HttpMiddleware
	Router     RouteRegister
}

type BaseHttpGateway struct {
	address string
	server  *fasthttp.Server
	gateway.UnimplementedGatewayPlugin

	// Middleware for handling events
	// return bool will indicate whether to continue
	// to the next (default) behaviour or not...
	mw       HttpMiddleware
	routeReg RouteRegister
}

func HttpHeadersToMap(rh *fasthttp.RequestHeader) map[string][]string {
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

func (s *BaseHttpGateway) httpHandler(workerPool pool.WorkerPool) func(ctx *fasthttp.RequestCtx) {
	return func(rc *fasthttp.RequestCtx) {
		if s.mw != nil {
			if !s.mw(rc, workerPool) {
				// middleware has indicated that is has processed the request
				// so we can exit here
				return
			}
		}

		headerMap := HttpHeadersToMap(&rc.Request.Header)

		// httpTrigger := triggers.FromHttpRequest(rc)
		headers := map[string]*v1.HeaderValue{}
		for k, v := range headerMap {
			headers[k] = &v1.HeaderValue{Value: v}
		}

		query := map[string]*v1.QueryValue{}
		rc.QueryArgs().VisitAll(func(key []byte, val []byte) {
			k := string(key)

			if query[k] == nil {
				query[k] = &v1.QueryValue{}
			}

			query[k].Value = append(query[k].Value, string(val))
		})

		httpTrigger := &v1.TriggerRequest{
			Data: rc.Request.Body(),
			Context: &v1.TriggerRequest_Http{
				Http: &v1.HttpTriggerContext{
					Method:      string(rc.Request.Header.Method()),
					Path:        string(rc.URI().PathOriginal()),
					Headers:     headers,
					QueryParams: query,
				},
			},
		}

		wrkr, err := workerPool.GetWorker(&pool.GetWorkerOptions{
			Trigger: httpTrigger,
		})
		if err != nil {
			rc.Error("Unable to get worker to handle request", 500)
			return
		}

		response, err := wrkr.HandleTrigger(span.FromHeaders(context.TODO(), headerMap), httpTrigger)
		if err != nil {
			rc.Error(fmt.Sprintf("Error handling HTTP Request: %v", err), 500)
			return
		}

		if http := response.GetHttp(); http != nil {
			// Copy headers across
			for k, v := range http.Headers {
				for _, val := range v.Value {
					rc.Response.Header.Add(k, val)
				}
			}

			// Avoid content length header duplication
			rc.Response.Header.Del("Content-Length")
			rc.Response.SetStatusCode(int(http.Status))
			rc.Response.SetBody(response.Data)

			return
		}

		rc.Error("received invalid response type from worker", 500)
	}
}

func (s *BaseHttpGateway) Start(pool pool.WorkerPool) error {
	r := router.New()

	// Allow custom provider level routing for handling events/schedules etc.
	if s.routeReg != nil {
		s.routeReg(r, pool)
	}

	r.ANY("/{path?:*}", s.httpHandler(pool))

	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         r.Handler,
		ReadBufferSize:  8192,
	}

	return s.server.ListenAndServe(s.address)
}

func (s *BaseHttpGateway) Stop() error {
	tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if ok {
		_ = tp.ForceFlush(context.TODO())
		_ = tp.Shutdown(context.TODO())
	}

	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New(opts *BaseHttpGatewayOptions) (gateway.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &BaseHttpGateway{
		address:  address,
		mw:       opts.Middleware,
		routeReg: opts.Router,
	}, nil
}
