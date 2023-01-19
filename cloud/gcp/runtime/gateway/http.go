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
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/triggers"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func middleware(rc *fasthttp.RequestCtx, workerPool pool.WorkerPool) bool {
	bodyBytes := rc.Request.Body()

	// Check if the payload contains a pubsub event
	// TODO: We probably want to use a simpler method than this
	// like reading off the request origin to ensure it is from pubsub
	var pubsubEvent PubSubMessage
	if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
		// We have an event from pubsub here...
		topic := pubsubEvent.Message.Attributes["x-nitric-topic"]

		event := &v1.TriggerRequest{
			Data: pubsubEvent.Message.Data,
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{
					Topic: topic,
				},
			},
		}

		wrkr, err := workerPool.GetWorker(&pool.GetWorkerOptions{
			Trigger: event,
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

		if _, err := wrkr.HandleTrigger(ctx, event); err == nil {
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
