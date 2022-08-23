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
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/fasthttp/router"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"

	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/pkg/plugins/gateway/base_http"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"
)

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

func handleNotifications(ctx *fasthttp.RequestCtx, events []eventgrid.Event, pool worker.WorkerPool) {
	topicName := ctx.UserValue("name").(string)

	// TODO: As we are batch handling events
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

		evt := &triggers.Event{
			ID:      *event.ID,
			Topic:   topicName,
			Payload: payloadBytes,
		}

		wrkr, err := pool.GetWorker(&worker.GetWorkerOptions{
			Event: evt,
			Filter: func(w worker.Worker) bool {
				_, isSubscription := w.(*worker.SubscriptionWorker)
				return isSubscription
			},
		})
		if err != nil {
			log.Default().Println("could not get worker for topic: ", topicName)
			// TODO: Handle error
			continue
		}

		err = wrkr.HandleEvent(evt)
		if err != nil {
			log.Default().Println("could not handle event: ", evt)
			// TODO: Handle error
			continue
		}
	}

	// Return 200 OK (TODO: Determine how we could mark individual events for failure)
	// Or potentially requeue them here internally...
	ctx.SuccessString("text/plain", "success")
}

func handleSubscription(pool worker.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
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
			} else if eventType == "Notification" {
				// Handle notifications
				handleNotifications(ctx, eventgridEvents, pool)
			}
		}
	}
}

func handleSchedule(pool worker.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		topicName := ctx.UserValue("name").(string)

		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}

		evt := &triggers.Event{
			ID:      fmt.Sprintf("%s:%d", topicName, time.Now().UnixMilli()),
			Topic:   topicName,
			Payload: []byte{},
		}

		wrkr, err := pool.GetWorker(&worker.GetWorkerOptions{
			Event: evt,
			Filter: func(w worker.Worker) bool {
				_, isSchedule := w.(*worker.ScheduleWorker)
				return isSchedule
			},
		})

		if err != nil {
			ctx.Error(fmt.Sprintf("could not find handler for schedule: %s", topicName), 404)
			return
		}

		if err := wrkr.HandleEvent(evt); err != nil {
			ctx.Error("error handling schedule", 500)
		} else {
			ctx.SuccessString("text/plain", "success")
		}
	}
}

func routes(r *router.Router, p worker.WorkerPool) {
	r.GET("/dapr/config", func(ctx *fasthttp.RequestCtx) {
		ctx.Error("no config", 404)
	})
	r.ANY(base_http.DefaultScheduleRoute, handleSchedule(p))
	r.ANY(base_http.DefaultTopicRoute, handleSubscription(p))
}

// Create a new HTTP Gateway plugin
func New() (gateway.GatewayService, error) {
	return base_http.New(base_http.WithRouter(routes))
}
