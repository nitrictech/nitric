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

	fasthttprouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nitrictech/nitric/cloud/common/runtime/env"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

type (
	HttpMiddleware func(*fasthttp.RequestCtx, *gateway.GatewayStartOpts) bool
	// EventConstructor func(topicName string, ctx *fasthttp.RequestCtx) v1.TriggerRequest
)

// A callback function that allows for custom routing configuration to be setup by consumer packages.
type RouterRegistrationCallback func(*fasthttprouter.Router, *gateway.GatewayStartOpts)

const (
	DefaultTopicRoute              = "/x-nitric-topic/{name}"
	DefaultScheduleRoute           = "/x-nitric-schedule/{name}"
	DefaultBucketNotificationRoute = "/x-nitric-notification/bucket/{name}"
)

type HttpGatewayOptions struct {
	RouteRegistrationHook RouterRegistrationCallback
}

type HttpGateway struct {
	address string
	server  *fasthttp.Server
	gateway.UnimplementedGatewayPlugin
	routeRegistrationHook RouterRegistrationCallback
}

func HttpHeadersToMap(rh *fasthttp.RequestHeader) map[string][]string {
	headerCopy := make(map[string][]string)

	rh.VisitAll(func(key []byte, val []byte) {
		keyString := string(key)

		if strings.ToLower(keyString) == "host" {
			// Don't copy the host header
			headerCopy["X-Forwarded-For"] = []string{string(val)}
		} else if strings.ToLower(keyString) == "x-forwarded-authorization" {
			// Forward original authorization header
			headerCopy["Authorization"] = []string{string(val)}
		} else {
			headerCopy[string(key)] = append(headerCopy[string(key)], string(val))
		}
	})

	return headerCopy
}

func (s *HttpGateway) newApiHandler(opts *gateway.GatewayStartOpts, apiNameParam string, originalPathParam string) func(ctx *fasthttp.RequestCtx) {
	return func(rc *fasthttp.RequestCtx) {
		// The API name is captured in the path using a path rewrite at the cloud API Gateway layer, and used to route the request to the correct workers
		// the path is extracted here for routing. The original path is captured and passed to the workers, removing the rewrite.
		apiName, apiOk := rc.UserValue(apiNameParam).(string)
		originalPath, pathOk := rc.UserValue(originalPathParam).(string)
		if !apiOk || !pathOk {
			rc.Error("invalid path", 400)
			return
		}

		headerMap := HttpHeadersToMap(&rc.Request.Header)

		// httpTrigger := triggers.FromHttpRequest(rc)
		headers := map[string]*apispb.HeaderValue{}
		for k, v := range headerMap {
			headers[k] = &apispb.HeaderValue{Value: v}
		}

		query := map[string]*apispb.QueryValue{}
		rc.QueryArgs().VisitAll(func(key []byte, val []byte) {
			k := string(key)

			if query[k] == nil {
				query[k] = &apispb.QueryValue{}
			}

			query[k].Value = append(query[k].Value, string(val))
		})

		httpTrigger := &apispb.ServerMessage{
			Content: &apispb.ServerMessage_HttpRequest{
				HttpRequest: &apispb.HttpRequest{
					Method:      string(rc.Request.Header.Method()),
					Path:        originalPath,
					Headers:     headers,
					QueryParams: query,
					Body:        rc.Request.Body(),
				},
			},
		}

		resp, err := opts.ApiPlugin.HandleRequest(apiName, httpTrigger)
		if err != nil {
			rc.Error("Unable to get worker to handle request", 500)
			return
		}

		if http := resp.GetHttpResponse(); http != nil {
			// Copy headers across
			for k, v := range http.Headers {
				for _, val := range v.Value {
					rc.Response.Header.Add(k, val)
				}
			}

			// Avoid content length header duplication
			rc.Response.Header.Del("Content-Length")
			rc.Response.SetStatusCode(int(http.Status))
			rc.Response.SetBody(resp.GetHttpResponse().Body)

			return
		}

		rc.Error("received invalid response type from worker", 500)
	}
}

func (s *HttpGateway) newHttpProxyHandler(opts *gateway.GatewayStartOpts) func(ctx *fasthttp.RequestCtx) {
	return func(rc *fasthttp.RequestCtx) {
		logger.Debugf("handling HTTP request: %s", rc.Request.URI())

		// Copy the cloud provider authorization header to the X-Platform-Authorization header
		// This will preserve the Bearer token used to communicate with the compute platform in case needed in future
		if auth := rc.Request.Header.Peek("Authorization"); len(auth) > 0 {
			rc.Request.Header.Set("X-Platform-Authorization", string(auth))
		}

		// Copy the X-Forwarded-Authorization header to the Authorization header
		// In cloud environments, the Authorization header is usually stripped by the cloud provider
		// at the api gateway layer, and forwarded as a custom header in order to authenticate with the compute platform its forwarding to.
		if auth := rc.Request.Header.Peek("X-Forwarded-Authorization"); len(auth) > 0 {
			rc.Request.Header.Set("Authorization", string(auth))
		}

		resp, err := opts.HttpPlugin.HandleRequest(&rc.Request)
		if err != nil {
			logger.Errorf("error handling request: %s", err)
			rc.Error("Internal Server Error", 500)
			return
		}

		// Copy the given response back
		resp.CopyTo(&rc.Response)
	}
}

func (s *HttpGateway) newDaprConfigHandler() func(ctx *fasthttp.RequestCtx) {
	return func(rc *fasthttp.RequestCtx) {
		rc.Error("No config available", 404)
	}
}

// Start the HTTP server and listen for requests, then route them to the appropriate handler(s
func (s *HttpGateway) Start(opts *gateway.GatewayStartOpts) error {
	r := fasthttprouter.New()

	// Allow custom provider level routing for handling events/schedules etc.
	if s.routeRegistrationHook != nil {
		s.routeRegistrationHook(r, opts)
	}

	// Handle Dapr config request
	r.ANY("/dapr/config", s.newDaprConfigHandler())

	// if opts.ApiPlugin.WorkerCount() > 0 {
	// Capture the API Name to allow for accurate worker routing.
	// also capture the original path so it can be passed to the worker without the name prefix.
	const apiNameParam = "apiName"
	const originalPathParam = "path"
	r.ANY(fmt.Sprintf("/x-nitric-api/{%s}/{%s?:*}", apiNameParam, originalPathParam), s.newApiHandler(opts, apiNameParam, originalPathParam))
	// }

	// proxy to http if available
	// if opts.HttpPlugin.WorkerCount() > 0 {
	if opts.HttpPlugin != nil {
		r.ANY("/{path?:*}", s.newHttpProxyHandler(opts))
	}
	// }

	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         r.Handler,
		ReadBufferSize:  8192,
	}

	return s.server.ListenAndServe(s.address)
}

// Stop the HTTP gateway server
func (s *HttpGateway) Stop() error {
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
func NewHttpGateway(opts *HttpGatewayOptions) (gateway.GatewayService, error) {
	address := env.GATEWAY_ADDRESS.String()

	return &HttpGateway{
		address:               address,
		routeRegistrationHook: opts.RouteRegistrationHook,
	}, nil
}
