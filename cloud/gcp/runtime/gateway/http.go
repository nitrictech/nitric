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
	"log"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"

	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/core"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/worker"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type gcpMiddleware struct {
	provider core.GcpProvider
}

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (g *gcpMiddleware) handleSubscription(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		bodyBytes := ctx.Request.Body()

		// Check if the payload contains a pubsub event
		// TODO: We probably want to use a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
			// We have an event from pubsub here...
			topicName := ctx.UserValue("name").(string)
			if topicName == "" {
				ctx.Error("Can not handle event for empty topic", 400)
			}

			event := &v1.TriggerRequest{
				Data: pubsubEvent.Message.Data,
				Context: &v1.TriggerRequest_Topic{
					Topic: &v1.TopicTriggerContext{
						Topic: topicName,
					},
				},
			}

			worker, err := process.GetWorker(&pool.GetWorkerOptions{
				Trigger: event,
			})
			if err != nil {
				ctx.Error("Could not find handle for event", 500)
			}

			traceKey := propagator.CloudTraceFormatPropagator{}.Fields()[0]
			traceCtx := context.TODO()

			if pubsubEvent.Message.Attributes[traceKey] != "" {
				var mc propagation.MapCarrier = pubsubEvent.Message.Attributes
				traceCtx = propagator.CloudTraceFormatPropagator{}.Extract(traceCtx, mc)
			} else {
				var hc propagation.HeaderCarrier = base_http.HttpHeadersToMap(&ctx.Request.Header)
				traceCtx = propagator.CloudTraceFormatPropagator{}.Extract(traceCtx, hc)
			}

			if _, err := worker.HandleTrigger(traceCtx, event); err == nil {
				// return a successful response
				ctx.SuccessString("text/plain", "success")
			} else {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
			}
		}
	}
}

func (g *gcpMiddleware) handleSchedule(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		scheduleName := ctx.UserValue("name").(string)
		if scheduleName == "" {
			ctx.Error("Can not handle event for empty schedule", 400)
		}

		evt := &v1.TriggerRequest{
			// Send empty data for now (no reason to send data for schedules at the moment)
			Data: nil,
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{
					Topic: scheduleName,
				},
			},
		}

		worker, err := process.GetWorker(&pool.GetWorkerOptions{
			Trigger: evt,
			Filter: func(w worker.Worker) bool {
				_, isSchedule := w.(*worker.ScheduleWorker)
				return isSchedule
			},
		})
		if err != nil {
			log.Default().Println("could not get worker for schedule: ", scheduleName)
		}

		var hc propagation.HeaderCarrier = base_http.HttpHeadersToMap(&ctx.Request.Header)
		traceCtx := propagator.CloudTraceFormatPropagator{}.Extract(context.TODO(), hc)

		_, err = worker.HandleTrigger(traceCtx, evt)
		if err != nil {
			log.Default().Println("could not handle event: ", evt)
		}

		ctx.SuccessString("text/plain", "success")
	}
}

// Converts the GCP event type to our abstract event type
func notificationEventToEventType(eventType string) (v1.BucketNotificationType, error) {
	switch eventType {
	case "OBJECT_FINALIZE":
		return v1.BucketNotificationType_Created, nil
	case "OBJECT_DELETE":
		return v1.BucketNotificationType_Deleted, nil
	default:
		return v1.BucketNotificationType_All, fmt.Errorf("unsupported bucket notification event type %s", eventType)
	}
}

func (g *gcpMiddleware) handleBucketNotification(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		bodyBytes := ctx.Request.Body()

		// Check if the payload contains a pubsub event
		// TODO: We probably want to use a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
			bucketName := ctx.UserValue("name").(string)

			key := pubsubEvent.Message.Attributes["objectId"]
			eventType, err := notificationEventToEventType(pubsubEvent.Message.Attributes["eventType"])
			if err != nil {
				ctx.Error(err.Error(), 400)
				return
			}

			evt := &v1.TriggerRequest{
				Context: &v1.TriggerRequest_Notification{
					Notification: &v1.NotificationTriggerContext{
						Source: bucketName,
						Notification: &v1.NotificationTriggerContext_Bucket{
							Bucket: &v1.BucketNotification{
								Key:  key,
								Type: eventType,
							},
						},
					},
				},
			}

			worker, err := process.GetWorker(&pool.GetWorkerOptions{
				Trigger: evt,
				Filter: func(w worker.Worker) bool {
					_, ok := w.(*worker.BucketNotificationWorker)
					return ok
				},
			})
			if err != nil {
				ctx.Error("Could not find handle for event", 500)
			}

			traceKey := propagator.CloudTraceFormatPropagator{}.Fields()[0]
			traceCtx := context.TODO()

			if pubsubEvent.Message.Attributes[traceKey] != "" {
				var mc propagation.MapCarrier = pubsubEvent.Message.Attributes
				traceCtx = propagator.CloudTraceFormatPropagator{}.Extract(traceCtx, mc)
			} else {
				var hc propagation.HeaderCarrier = base_http.HttpHeadersToMap(&ctx.Request.Header)
				traceCtx = propagator.CloudTraceFormatPropagator{}.Extract(traceCtx, hc)
			}

			if _, err := worker.HandleTrigger(traceCtx, evt); err == nil {
				// return a successful response
				ctx.SuccessString("text/plain", "success")
			} else {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
			}
		}
	}
}

func (g *gcpMiddleware) router(r *router.Router, pool pool.WorkerPool) {
	r.ANY(base_http.DefaultTopicRoute, g.handleSubscription(pool))
	r.ANY(base_http.DefaultScheduleRoute, g.handleSchedule(pool))
	r.ANY(base_http.DefaultBucketNotificationRoute, g.handleBucketNotification(pool))
}

// New - Create a New cloudrun gateway plugin
func New(provider core.GcpProvider) (gateway.GatewayService, error) {
	mw := &gcpMiddleware{
		provider: provider,
	}

	return base_http.New(&base_http.BaseHttpGatewayOptions{
		Router: mw.router,
	})
}
