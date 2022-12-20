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
package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"

	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	ep "github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/triggers"
	"github.com/nitrictech/nitric/core/pkg/worker"
)

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func middleware(rc *fasthttp.RequestCtx, pool worker.WorkerPool) bool {
	bodyBytes := rc.Request.Body()

	// Check if the payload contains a pubsub event
	// TODO: We probably want to use a simpler method than this
	// like reading off the request origin to ensure it is from pubsub
	var pubsubEvent PubSubMessage
	if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
		// We have an event from pubsub here...

		topic := pubsubEvent.Message.Attributes["x-nitric-topic"]

		// need to determine if the underlying data is a nitric event
		var event *triggers.Event
		messageJson := &ep.NitricEvent{}
		// Check if it's a nitric event
		if err := json.Unmarshal(pubsubEvent.Message.Data, messageJson); err == nil && messageJson.ID != "" {
			// reserialize the nitric event payload
			payload, _ := json.Marshal(messageJson.Payload)

			event = &triggers.Event{
				ID:         messageJson.ID,
				Topic:      topic,
				Payload:    payload,
				Attributes: pubsubEvent.Message.Attributes,
			}
		} else {
			event = &triggers.Event{
				ID: pubsubEvent.Message.ID,
				// Set the topic
				Topic: topic,
				// Set the original full payload payload
				Payload:    pubsubEvent.Message.Data,
				Attributes: pubsubEvent.Message.Attributes,
			}
		}

		wrkr, err := pool.GetWorker(&worker.GetWorkerOptions{
			Event: event,
		})
		if err != nil {
			rc.Error("Could not find handle for event", 500)
			return false
		}

		traceKey := propagator.CloudTraceFormatPropagator{}.Fields()[0]
		ctx := context.TODO()

		if pubsubEvent.Message.Attributes[traceKey] != "" {
			var mc propagation.MapCarrier = pubsubEvent.Message.Attributes
			ctx = propagator.CloudTraceFormatPropagator{}.Extract(ctx, mc)
		} else {
			var hc propagation.HeaderCarrier = triggers.HttpHeaders(&rc.Request.Header)
			ctx = propagator.CloudTraceFormatPropagator{}.Extract(ctx, hc)
		}

		if err := wrkr.HandleEvent(ctx, event); err == nil {
			// return a successful response
			rc.SuccessString("text/plain", "success")
		} else {
			rc.Error(fmt.Sprintf("Error handling event %v", err), 500)
		}

		// We've already handled the request
		// do not continue processing
		return false
	}

	// Let the base plugin handle the request
	return true
}

// New - Create a New cloudrun gateway plugin
func New() (gateway.GatewayService, error) {
	// plugin is derived from base http plugin
	return base_http.New(middleware)
}
