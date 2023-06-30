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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/fasthttp/router"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"

	"github.com/nitrictech/nitric/cloud/azure/runtime/core"
	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/worker"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type azMiddleware struct {
	provider core.AzProvider
}

func extractEvents(ctx *fasthttp.RequestCtx) ([]eventgrid.Event, error) {
	var eventgridEvents []eventgrid.Event
	bytes := ctx.Request.Body()
	// TODO: verify topic for validity
	if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
		return nil, errors.New("invalid event grid types")
	}

	return eventgridEvents, nil
}

func extractPayload(event eventgrid.Event) []byte {
	var payloadBytes []byte
	if stringData, ok := event.Data.(string); ok {
		payloadBytes = []byte(stringData)
	} else if byteData, ok := event.Data.([]byte); ok {
		payloadBytes = byteData
	} else {
		// Assume a json serializable struct for now...
		payloadBytes, _ = json.Marshal(event.Data)
	}
	return payloadBytes
}

func eventAuthorised(ctx *fasthttp.RequestCtx) bool {
	token := ctx.QueryArgs().Peek("token")
	evtToken := os.Getenv("EVENT_TOKEN")

	return string(token) == evtToken
}

func (a *azMiddleware) handleSubscriptionValidation(ctx *fasthttp.RequestCtx, events []eventgrid.Event) {
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

func (a *azMiddleware) handleSubscription(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}

		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
		}

		eventgridEvents, err := extractEvents(ctx)
		if err != nil {
			ctx.Error(err.Error(), 400)
			return
		}

		for _, event := range eventgridEvents {
			eventType := string(ctx.Request.Header.Peek("aeg-event-type"))
			if eventType == "SubscriptionValidation" {
				a.handleSubscriptionValidation(ctx, eventgridEvents)
				return
			}

			payloadBytes := extractPayload(event)

			topicName := ctx.UserValue("name").(string)

			evt := &v1.TriggerRequest{
				Data: payloadBytes,
				Context: &v1.TriggerRequest_Topic{
					Topic: &v1.TopicTriggerContext{
						Topic: topicName,
					},
				},
			}

			wrkr, err := process.GetWorker(&pool.GetWorkerOptions{
				Trigger: evt,
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

			_, err = wrkr.HandleTrigger(context.TODO(), evt)
			if err != nil {
				log.Default().Println("could not handle event: ", evt)
			}

			// TODO: event handling failure???
			ctx.SuccessString("text/plain", "success")
		}
	}
}

func (a *azMiddleware) handleSchedule(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
		}

		scheduleName := ctx.UserValue("name").(string)

		evt := &v1.TriggerRequest{
			// Send empty data for now (no reason to send data for schedules at the moment)
			Data: nil,
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{
					Topic: scheduleName,
				},
			},
		}

		wrkr, err := process.GetWorker(&pool.GetWorkerOptions{
			Trigger: evt,
			Filter: func(w worker.Worker) bool {
				_, isSchedule := w.(*worker.ScheduleWorker)
				return isSchedule
			},
		})
		if err != nil {
			log.Default().Println("could not get worker for schedule: ", scheduleName)
		}

		_, err = wrkr.HandleTrigger(context.TODO(), evt)
		if err != nil {
			log.Default().Println("could not handle event: ", evt)
		}

		ctx.SuccessString("text/plain", "success")
	}
}

// Converts the GCP event type to our abstract event type
func notificationEventToEventType(eventType *string) (v1.BucketNotificationType, error) {
	switch *eventType {
	case "Microsoft.Storage.BlobCreated":
		return v1.BucketNotificationType_Created, nil
	case "Microsoft.Storage.BlobDeleted":
		return v1.BucketNotificationType_Deleted, nil
	default:
		return v1.BucketNotificationType_All, fmt.Errorf("unsupported bucket notification event type %s", *eventType)
	}
}

func (a *azMiddleware) handleBucketNotification(process pool.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
		}

		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}

		eventgridEvents, err := extractEvents(ctx)
		if err != nil {
			ctx.Error(fmt.Sprintf("error occurred extracting events: %s", err.Error()), 400)
			return
		}

		for _, event := range eventgridEvents {
			azureEventType := string(ctx.Request.Header.Peek("aeg-event-type"))
			if azureEventType == "SubscriptionValidation" {
				a.handleSubscriptionValidation(ctx, eventgridEvents)
				return
			}

			bucketName := ctx.UserValue("name").(string)

			eventType, err := notificationEventToEventType(event.EventType)
			if err != nil {
				ctx.Error(err.Error(), 400)
				return
			}

			// Subject is in the form: "/blobServices/default/containers/test-container/blobs/new-file.txt"
			eventKeySeparated := strings.SplitN(*event.Subject, "/", 7)
			if len(eventKeySeparated) < 7 {
				ctx.Error("object key cannot be empty", 400)
				return
			}

			eventKey := eventKeySeparated[6]

			evt := &v1.TriggerRequest{
				Context: &v1.TriggerRequest_Notification{
					Notification: &v1.NotificationTriggerContext{
						Source: bucketName,
						Notification: &v1.NotificationTriggerContext_Bucket{
							Bucket: &v1.BucketNotification{
								Key:  eventKey,
								Type: eventType,
							},
						},
					},
				},
			}

			wrkr, err := process.GetWorker(&pool.GetWorkerOptions{
				Trigger: evt,
				Filter: func(w worker.Worker) bool {
					_, isNotification := w.(*worker.BucketNotificationWorker)
					return isNotification
				},
			})
			if err != nil {
				log.Default().Println("could not get worker for bucket notification: ", bucketName)
			}

			_, err = wrkr.HandleTrigger(context.TODO(), evt)
			if err != nil {
				log.Default().Println("could not handle event: ", evt)
			}

			ctx.SuccessString("text/plain", "success")
		}
	}
}

func (a *azMiddleware) router(r *router.Router, pool pool.WorkerPool) {
	r.ANY(base_http.DefaultTopicRoute, a.handleSubscription(pool))
	r.ANY(base_http.DefaultScheduleRoute, a.handleSchedule(pool))
	r.ANY(base_http.DefaultBucketNotificationRoute, a.handleBucketNotification(pool))
}

// Create a new HTTP Gateway plugin
func New(provider core.AzProvider) (gateway.GatewayService, error) {
	mw := &azMiddleware{
		provider: provider,
	}

	return base_http.New(&base_http.BaseHttpGatewayOptions{
		Router: mw.router,
	})
}
