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

package http_service

import (
	"encoding/json"
	"fmt"
	triggers2 "github.com/nitric-dev/membrane/pkg/triggers"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	worker2 "github.com/nitric-dev/membrane/pkg/worker"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/mitchellh/mapstructure"
	"github.com/nitric-dev/membrane/pkg/sdk"
	"github.com/valyala/fasthttp"
)

// HttpService - The HTTP gateway plugin for Azure
type HttpService struct {
	address string
	server  *fasthttp.Server
}

func handleSubscriptionValidation(ctx *fasthttp.RequestCtx, events []eventgrid.Event) {
	subPayload := events[0]
	var validateData eventgrid.SubscriptionValidationEventData
	if err := mapstructure.Decode(subPayload.Data, &validateData); err != nil {
		ctx.Error("Invalid subscription event data", 400)
		return
	}

	response := eventgrid.SubscriptionValidationResponse{
		ValidationResponse: validateData.ValidationCode,
	}

	responseBody, _ := json.Marshal(response)

	ctx.Success("application/json", responseBody)
}

func handleNotifications(ctx *fasthttp.RequestCtx, events []eventgrid.Event, wrkr worker2.Worker) {
	// FIXME: As we are batch handling events in azure
	// how do we notify of failed event handling?
	for _, event := range events {
		// XXX: Assume we have a nitric event for now
		// We have a valid nitric event
		// Decode and pass to our function
		var payloadBytes []byte
		if stringData, ok := event.Data.(string); ok {
			payloadBytes = []byte(stringData)
		} else if byteData, ok := event.Data.([]byte); ok {
			payloadBytes = byteData
		} else {
			// Assume a json serializable struct for now...
			payloadBytes, _ = json.Marshal(event.Data)
		}

		// FIXME: Handle error
		wrkr.HandleEvent(&triggers2.Event{
			// FIXME: Check if ID is nil
			ID:      *event.ID,
			Topic:   *event.Topic,
			Payload: payloadBytes,
		})
	}

	// Return 200 OK (TODO: Determine how we could mark individual events for failure)
	// Or potentially requeue them here internally...
	ctx.SuccessString("text/plain", "success")
}

func handleRequest(ctx *fasthttp.RequestCtx, wrkr worker2.Worker) {
	response, err := wrkr.HandleHttpRequest(triggers2.FromHttpRequest(ctx))

	if err != nil {
		ctx.Error(fmt.Sprintf("Error Handling Request: %v", err), 500)
		return
	}
	if response.Header != nil {
		response.Header.VisitAll(func(key []byte, val []byte) {
			ctx.Response.Header.AddBytesKV(key, val)
		})
	}

	// Avoid content length header duplication
	ctx.Response.Header.Del("Content-Length")
	ctx.Response.SetStatusCode(response.StatusCode)
	ctx.Response.SetBody(response.Body)
}

func httpHandler(pool worker2.WorkerPool) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		wrkr, err := pool.GetWorker()

		if err != nil {
			ctx.Error("Unable to get worker to handle request", 500)
			return
		}
		// Handle Event/Subscription Request Types
		eventType := string(ctx.Request.Header.Peek("aeg-event-type"))

		// Handle an eventgrid webhook event
		if eventType != "" {
			var eventgridEvents []eventgrid.Event
			bytes := ctx.Request.Body()
			// TODO: verify topic for validity
			if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
				ctx.Error("Invalid event grid types", 400)
				return
			}

			// Handle Eventgrid event
			if eventType == "SubscriptionValidation" {
				// Validate a subscription
				handleSubscriptionValidation(ctx, eventgridEvents)
				return
			} else if eventType == "Notification" {
				// Handle notifications
				handleNotifications(ctx, eventgridEvents, wrkr)
				return
			}
		}

		// Handle a standard HTTP request
		handleRequest(ctx, wrkr)
	}
}

func (s *HttpService) Start(pool worker2.WorkerPool) error {
	// Start the fasthttp server
	s.server = &fasthttp.Server{
		Handler: httpHandler(pool),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *HttpService) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// Create a new HTTP Gateway plugin
func New() (sdk.GatewayService, error) {
	address := utils2.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpService{
		address: address,
	}, nil
}
