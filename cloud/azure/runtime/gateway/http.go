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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/fasthttp/router"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type azMiddleware struct {
	provider resource.AzProvider
}

func extractEvents(ctx *fasthttp.RequestCtx) ([]eventgrid.Event, error) {
	var eventgridEvents []eventgrid.Event
	bytes := ctx.Request.Body()
	if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
		return nil, errors.New("invalid event grid types")
	}

	return eventgridEvents, nil
}

func extractMessage(event eventgrid.Event) (*topicpb.TopicMessage, error) {
	var payloadBytes []byte
	var err error
	if stringData, ok := event.Data.(string); ok {
		// byte payloads are automatically base64 encoded when publishing to event grid
		payloadBytes, err = base64.StdEncoding.DecodeString(stringData)
		if err != nil {
			return nil, err
		}
	} else if byteData, ok := event.Data.([]byte); ok {
		payloadBytes = byteData
	} else {
		return nil, fmt.Errorf("invalid event data type: %T", event.Data)
	}

	var message topicpb.TopicMessage

	if err = proto.Unmarshal(payloadBytes, &message); err != nil {
		return nil, err
	}

	return &message, nil
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

func (a *azMiddleware) handleSubscription(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}

		eventgridEvents, err := extractEvents(ctx)
		if err != nil {
			ctx.Error(err.Error(), 400)
			return
		}

		for _, event := range eventgridEvents {
			eventType := string(ctx.Request.Header.Peek("aeg-event-type"))
			if eventType == "SubscriptionValidation" {
				fmt.Println("handling validation event")
				a.handleSubscriptionValidation(ctx, eventgridEvents)
				return
			}

			message, err := extractMessage(event)
			if err != nil {
				fmt.Println("error extracting message", err)
				ctx.Error(err.Error(), 500)
				return
			}

			topicName := ctx.UserValue("name").(string)

			evt := &topicspb.ServerMessage{
				Content: &topicspb.ServerMessage_MessageRequest{
					MessageRequest: &topicspb.MessageRequest{
						TopicName: topicName,
						Message:   message,
					},
				},
			}

			resp, err := opts.TopicsListenerPlugin.HandleRequest(evt)
			if err != nil {
				logger.Errorf("error handling event from topic %s: %s", topicName, err.Error())
				ctx.Error("failed handling event, error returned from subscriber function", 500)
				return
			}

			if !resp.GetMessageResponse().Success {
				logger.Errorf("event handling failed %s, subscriber returned success=false", topicName)
				ctx.Error("subscriber returned non-success response", 500)
				return
			}
		}
		ctx.SuccessString("text/plain", "success")
	}
}

func (a *azMiddleware) handleSchedule(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}

		scheduleName := ctx.UserValue("name").(string)

		evt := &schedulespb.ServerMessage{
			Content: &schedulespb.ServerMessage_IntervalRequest{
				IntervalRequest: &schedulespb.IntervalRequest{
					ScheduleName: scheduleName,
				},
			},
		}

		_, err := opts.SchedulesPlugin.HandleRequest(evt)
		if err != nil {
			ctx.Error(fmt.Sprintf("failed handling schedule %s", scheduleName), 500)
		}

		ctx.SuccessString("text/plain", "success")
	}
}

// Converts the GCP event type to our abstract event type
func notificationEventToEventType(eventType *string) (*storagepb.BlobEventType, error) {
	switch *eventType {
	case "Microsoft.Storage.BlobCreated":
		return storagepb.BlobEventType_Created.Enum(), nil
	case "Microsoft.Storage.BlobDeleted":
		return storagepb.BlobEventType_Deleted.Enum(), nil
	default:
		return nil, fmt.Errorf("unsupported bucket notification event type %s", *eventType)
	}
}

func (a *azMiddleware) handleBucketNotification(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
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

			evt := &storagepb.ServerMessage{
				Content: &storagepb.ServerMessage_BlobEventRequest{
					BlobEventRequest: &storagepb.BlobEventRequest{
						BucketName: bucketName,
						Event: &storagepb.BlobEventRequest_BlobEvent{
							BlobEvent: &storagepb.BlobEvent{
								Key:  eventKey,
								Type: *eventType,
							},
						},
					},
				},
			}

			resp, err := opts.StorageListenerPlugin.HandleRequest(evt)
			if err != nil {
				logger.Errorf("could not handle event: %s", err)
				ctx.Error("failed handling event", 500)
				return
			}

			if !resp.GetBlobEventResponse().Success {
				logger.Errorf("failed handling event: %s", evt)
				ctx.Error("failed handling event", 500)
				return
			}

			ctx.SuccessString("text/plain", "success")
		}
	}
}

func (a *azMiddleware) router(r *router.Router, opts *gateway.GatewayStartOpts) {
	evtToken := os.Getenv("EVENT_TOKEN")

	r.ANY("/"+evtToken+base_http.DefaultTopicRoute, a.handleSubscription(opts))
	r.ANY("/"+evtToken+base_http.DefaultScheduleRoute, a.handleSchedule(opts))
	r.ANY("/"+evtToken+base_http.DefaultBucketNotificationRoute, a.handleBucketNotification(opts))
}

// Create a new HTTP Gateway plugin
func New(provider resource.AzProvider) (gateway.GatewayService, error) {
	mw := &azMiddleware{
		provider: provider,
	}

	return base_http.NewHttpGateway(&base_http.HttpGatewayOptions{
		RouteRegistrationHook: mw.router,
	})
}
