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

// The GCP HTTP gateway plugin for CloudRun
package cloudrun_plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/utils"
	"github.com/nitric-dev/membrane/pkg/worker"

	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	"github.com/valyala/fasthttp"
)

type HttpProxyGateway struct {
	address string
	server  *fasthttp.Server
}

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func httpHandler(pool worker.WorkerPool) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		wrkr, err := pool.GetWorker()
		if err != nil {
			ctx.Error("Unable to get worker to handle request", 500)
			return
		}

		bodyBytes := ctx.Request.Body()

		// Check if the payload contains a pubsub event
		// TODO: We probably want to use a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
			// We have an event from pubsub here...
			event := &triggers.Event{
				ID: pubsubEvent.Message.ID,
				// Set the topic
				Topic: pubsubEvent.Message.Attributes["x-nitric-topic"],
				// Set the payload
				Payload: pubsubEvent.Message.Data,
			}

			if err := wrkr.HandleEvent(event); err == nil {
				// return a successful response
				ctx.SuccessString("text/plain", "success")
			} else {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
			}

			return
		}

		httpTrigger := triggers.FromHttpRequest(ctx)
		response, err := wrkr.HandleHttpRequest(httpTrigger)

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

func (s *HttpProxyGateway) Start(pool worker.WorkerPool) error {
	// Start the fasthttp server
	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         httpHandler(pool),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *HttpProxyGateway) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (gateway.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpProxyGateway{
		address: address,
	}, nil
}
