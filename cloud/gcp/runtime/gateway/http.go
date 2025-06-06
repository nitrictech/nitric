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
	"encoding/json"
	"fmt"
	"os"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type gcpMiddleware struct {
	provider resource.GcpResourceResolver
}

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func eventAuthorised(ctx *fasthttp.RequestCtx) bool {
	token := ctx.QueryArgs().Peek("token")
	evtToken := os.Getenv("EVENT_TOKEN")

	fmt.Println("checking:", string(token), evtToken)

	return string(token) == evtToken
}

func (g *gcpMiddleware) handleSubscription(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
			return
		}

		bodyBytes := ctx.Request.Body()
		// Check if the payload contains a pubsub event
		// Consider using a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
			// We have an event from pubsub here...
			topicName := ctx.UserValue("name").(string)
			if topicName == "" {
				ctx.Error("Can not handle event for empty topic", 400)
			}

			var message topicpb.TopicMessage
			err := proto.Unmarshal(pubsubEvent.Message.Data, &message)
			if err != nil {
				fmt.Println("could not parse message as a nitric event attempting to parse as generic json payload")

				messageData := map[string]any{}
				err := json.Unmarshal(pubsubEvent.Message.Data, &messageData)
				if err != nil {
					ctx.Error("could not unmarshal event data", 500)
					return
				}

				structPayload, err := structpb.NewStruct(messageData)
				if err != nil {
					ctx.Error("could not convert message data to struct", 500)
					return
				}

				message = topicpb.TopicMessage{
					Content: &topicspb.TopicMessage_StructPayload{
						StructPayload: structPayload,
					},
				}
			}

			event := &topicspb.ServerMessage{
				Content: &topicspb.ServerMessage_MessageRequest{
					MessageRequest: &topicspb.MessageRequest{
						TopicName: topicName,
						Message:   &message,
					},
				},
			}

			response, err := opts.TopicsListenerPlugin.HandleRequest(event)
			if err != nil {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
				return
			}

			if !response.GetMessageResponse().Success {
				ctx.Error("Event handler returned success false", 500)
				return
			}

			ctx.SuccessString("text/plain", "success")
		}
	}
}

func (g *gcpMiddleware) handleSchedule(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
			return
		}

		scheduleName := ctx.UserValue("name").(string)
		if scheduleName == "" {
			ctx.Error("Can not handle event for empty schedule", 400)
		}

		_, err := opts.SchedulesPlugin.HandleRequest(&schedulespb.ServerMessage{
			Content: &schedulespb.ServerMessage_IntervalRequest{
				IntervalRequest: &schedulespb.IntervalRequest{
					ScheduleName: scheduleName,
				},
			},
		})
		if err != nil {
			logger.Errorf("could not handle trigger for schedule %s: %s", scheduleName, err.Error())
			ctx.Error("could not handle trigger", 500)
			return
		}

		ctx.SuccessString("text/plain", "success")
	}
}

// Converts the GCP event type to our abstract event type
func notificationEventToEventType(eventType string) (*storagepb.BlobEventType, error) {
	switch eventType {
	case "OBJECT_FINALIZE":
		return storagepb.BlobEventType_Created.Enum(), nil
	case "OBJECT_DELETE":
		return storagepb.BlobEventType_Deleted.Enum(), nil
	default:
		return nil, fmt.Errorf("unsupported bucket notification event type %s", eventType)
	}
}

func (g *gcpMiddleware) handleBucketNotification(opts *gateway.GatewayStartOpts) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !eventAuthorised(ctx) {
			ctx.Error("Unauthorized", 401)
			return
		}

		bodyBytes := ctx.Request.Body()

		// Check if the payload contains a pubsub event
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil && pubsubEvent.Subscription != "" {
			bucketName := ctx.UserValue("name").(string)

			key := pubsubEvent.Message.Attributes["objectId"]
			eventType, err := notificationEventToEventType(pubsubEvent.Message.Attributes["eventType"])
			if err != nil {
				ctx.Error(err.Error(), 400)
				return
			}

			resp, err := opts.StorageListenerPlugin.HandleRequest(&storagepb.ServerMessage{
				Content: &storagepb.ServerMessage_BlobEventRequest{
					BlobEventRequest: &storagepb.BlobEventRequest{
						BucketName: bucketName,
						Event: &storagepb.BlobEventRequest_BlobEvent{
							BlobEvent: &storagepb.BlobEvent{
								Key:  key,
								Type: *eventType,
							},
						},
					},
				},
			})
			if err != nil {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
				return
			}

			if !resp.GetBlobEventResponse().Success {
				ctx.Error("Error handling event", 500)
				return
			}

			ctx.SuccessString("text/plain", "success")
		}
	}
}

func (g *gcpMiddleware) router(r *router.Router, opts *gateway.GatewayStartOpts) {
	r.ANY(base_http.DefaultTopicRoute, g.handleSubscription(opts))
	r.ANY(base_http.DefaultScheduleRoute, g.handleSchedule(opts))
	r.ANY(base_http.DefaultBucketNotificationRoute, g.handleBucketNotification(opts))
}

// New - Create a New cloudrun gateway plugin
func New(provider resource.GcpResourceResolver) (gateway.GatewayService, error) {
	mw := &gcpMiddleware{
		provider: provider,
	}

	return base_http.NewHttpGateway(&base_http.HttpGatewayOptions{
		RouteRegistrationHook: mw.router,
	})
}
