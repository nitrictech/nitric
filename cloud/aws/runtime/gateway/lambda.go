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

package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type LambdaRuntimeHandler func(interface{})

func (s *LambdaGateway) getTopicNameForArn(ctx context.Context, topicArn string) (string, error) {
	topics, err := s.provider.GetResources(ctx, resource.AwsResource_Topic)
	if err != nil {
		return "", fmt.Errorf("error retrieving topics: %w", err)
	}

	for name, topic := range topics {
		if topic.ARN == topicArn {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find topic for arn %s", topicArn)
}

func (s *LambdaGateway) getBucketNameForArn(ctx context.Context, bucketArn string) (string, error) {
	buckets, err := s.provider.GetResources(ctx, resource.AwsResource_Bucket)
	if err != nil {
		return "", fmt.Errorf("error retrieving topics: %w", err)
	}

	for name, bucket := range buckets {
		if bucket.ARN == bucketArn {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find topic for arn %s", bucketArn)
}

type LambdaGateway struct {
	provider resource.AwsResourceProvider
	runtime  LambdaRuntimeHandler
	gateway.UnimplementedGatewayPlugin
	finished chan int
}

var _ gateway.GatewayService = &LambdaGateway{}

// isRejectedConnection returns true if the client message was a rejection response to a connection request.
func isRejectedConnection(resp *websocketspb.ClientMessage) bool {
	eventResponse := resp.GetWebsocketEventResponse()
	if eventResponse == nil {
		return false
	}
	connectionResponse := resp.GetWebsocketEventResponse().GetConnectionResponse()
	if connectionResponse == nil {
		return false
	}

	return connectionResponse.GetReject()
}

// handleWebsocketEvent translates AWS Websocket API events to Nitric Websocket events and forwards them to be handled by registered workers.
func (s *LambdaGateway) handleWebsocketEvent(ctx context.Context, websockets websockets.WebsocketRequestHandler, evt events.APIGatewayWebsocketProxyRequest) (interface{}, error) {
	api, err := s.provider.GetApiGatewayById(ctx, evt.RequestContext.APIID)
	if err != nil {
		return nil, err
	}

	stackID := commonenv.NITRIC_STACK_ID.String()
	nitricName, ok := api.Tags[tags.GetResourceNameKey(stackID)]
	if !ok {
		return nil, fmt.Errorf("received websocket trigger from non-nitric API gateway")
	}

	// Use the routekey to get the event type
	wsEvent := &websocketspb.ServerMessage_WebsocketEventRequest{
		WebsocketEventRequest: &websocketspb.WebsocketEventRequest{
			ConnectionId: evt.RequestContext.ConnectionID,
			SocketName:   nitricName,
			WebsocketEvent: &websocketspb.WebsocketEventRequest_Message{
				Message: &websocketspb.WebsocketMessageEvent{
					Body: []byte(evt.Body),
				},
			},
		},
	}
	switch evt.RequestContext.RouteKey {
	case "$connect":
		queryParams := map[string]*websocketspb.QueryValue{}
		for k, v := range evt.QueryStringParameters {
			queryParams[k] = &websocketspb.QueryValue{
				Value: []string{v},
			}
		}
		wsEvent = &websocketspb.ServerMessage_WebsocketEventRequest{
			WebsocketEventRequest: &websocketspb.WebsocketEventRequest{
				ConnectionId: evt.RequestContext.ConnectionID,
				SocketName:   nitricName,
				WebsocketEvent: &websocketspb.WebsocketEventRequest_Connection{
					Connection: &websocketspb.WebsocketConnectionEvent{
						QueryParams: queryParams,
					},
				},
			},
		}
	case "$disconnect":
		wsEvent = &websocketspb.ServerMessage_WebsocketEventRequest{
			WebsocketEventRequest: &websocketspb.WebsocketEventRequest{
				ConnectionId: evt.RequestContext.ConnectionID,
				SocketName:   nitricName,
				WebsocketEvent: &websocketspb.WebsocketEventRequest_Disconnection{
					Disconnection: &websocketspb.WebsocketDisconnectionEvent{},
				},
			},
		}
	}

	req := &websocketspb.ServerMessage{
		Content: wsEvent,
	}

	resp, err := websockets.HandleRequest(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error processing lambda request",
			IsBase64Encoded: false,
		}, nil
	}

	if isRejectedConnection(resp) {
		return events.APIGatewayProxyResponse{
			StatusCode:      401,
			Body:            "not authorized",
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func handleApiGatewayRequest(ctx context.Context, nitricName string, apismanager apis.ApiRequestHandler, evt events.APIGatewayV2HTTPRequest) (interface{}, error) {
	// Copy the headers and re-write for the proxy
	headerCopy := map[string]*apispb.HeaderValue{}

	for key, val := range evt.Headers {
		if strings.ToLower(key) == "host" {
			headerCopy[xforwardHeader] = &apispb.HeaderValue{
				Value: []string{val},
			}
		} else {
			if headerCopy[key] == nil {
				headerCopy[key] = &apispb.HeaderValue{}
			}
			headerCopy[key].Value = append(headerCopy[key].Value, val)
		}
	}

	// Copy the cookies over
	headerCopy["Cookie"] = &apispb.HeaderValue{
		Value: evt.Cookies,
	}

	// Parse the raw query string
	qVals, err := url.ParseQuery(evt.RawQueryString)
	if err != nil {
		return nil, fmt.Errorf("error parsing query for httpEvent: %w", err)
	}
	query := map[string]*apispb.QueryValue{}
	for k, v := range qVals {
		query[k] = &apispb.QueryValue{
			Value: v,
		}
	}

	data := []byte(evt.Body)
	if evt.IsBase64Encoded {
		data, err = base64.StdEncoding.DecodeString(evt.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      400,
				Body:            "Error processing lambda request",
				IsBase64Encoded: false,
			}, nil
		}
	}

	req := &apispb.ServerMessage{
		Content: &apispb.ServerMessage_HttpRequest{
			HttpRequest: &apispb.HttpRequest{
				Method:      evt.RequestContext.HTTP.Method,
				Path:        evt.RawPath,
				QueryParams: query,
				Headers:     headerCopy,
				Body:        data,
				PathParams:  evt.PathParameters,
			},
		},
	}

	resp, err := apismanager.HandleRequest(nitricName, req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "Internal Server Error",
			IsBase64Encoded: false,
		}, nil
	}

	lambdaHTTPHeaders := make(map[string]string)
	if resp.GetHttpResponse().Headers != nil {
		for k, v := range resp.GetHttpResponse().Headers {
			lambdaHTTPHeaders[k] = v.Value[0]
		}
	}

	responseString := base64.StdEncoding.EncodeToString(resp.GetHttpResponse().Body)

	return events.APIGatewayProxyResponse{
		StatusCode:      int(resp.GetHttpResponse().Status),
		Headers:         lambdaHTTPHeaders,
		Body:            responseString,
		IsBase64Encoded: true,
	}, nil
}

func handleHttpProxyRequest(ctx context.Context, httpmanager http.HttpRequestHandler, evt events.APIGatewayV2HTTPRequest) (interface{}, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.Header.SetMethod(evt.RequestContext.HTTP.Method)
	request.SetRequestURI(evt.RawPath)
	request.SetBody([]byte(evt.Body))

	// Copy the headers and re-write for the proxy
	for key, val := range evt.Headers {
		request.Header.Add(key, val)
	}

	// Copy the cookies over
	for _, cookie := range evt.Cookies {
		request.Header.Add("Cookie", cookie)
	}

	resp, err := httpmanager.HandleRequest(request)
	if err != nil {
		return nil, err
	}

	lambdaHTTPHeaders := make(map[string]string)
	resp.Header.VisitAll(func(key, value []byte) {
		lambdaHTTPHeaders[string(key)] = string(value)
	})

	responseString := base64.StdEncoding.EncodeToString(resp.Body())

	return events.APIGatewayProxyResponse{
		StatusCode:      resp.StatusCode(),
		Headers:         lambdaHTTPHeaders,
		Body:            responseString,
		IsBase64Encoded: true,
	}, nil
}

// handleApiEvent translates AWS API events to Nitric API events and forwards them to be handled by registered workers.
func (s *LambdaGateway) handleApiEvent(ctx context.Context, apismanager apis.ApiRequestHandler, httpmanager http.HttpRequestHandler, evt events.APIGatewayV2HTTPRequest) (interface{}, error) {
	api, err := s.provider.GetApiGatewayById(ctx, evt.RequestContext.APIID)
	if err != nil {
		return nil, err
	}

	stackID := commonenv.NITRIC_STACK_ID.String()
	nitricName, ok := api.Tags[tags.GetResourceNameKey(stackID)]
	if !ok {
		return nil, fmt.Errorf("received request from non-nitric API gateway")
	}

	nitricType, ok := api.Tags[tags.GetResourceTypeKey(stackID)]
	if !ok {
		return nil, fmt.Errorf("received request from non-nitric API gateway")
	}

	if nitricType == "http-proxy" {
		return handleHttpProxyRequest(ctx, httpmanager, evt)
	} else {
		return handleApiGatewayRequest(ctx, nitricName, apismanager, evt)
	}
}

type ScheduleMessage struct {
	Schedule string
}

// handleScheduleEvent translates AWS schedule events to Nitric schedule intervals and forwards them to be handled by registered workers.
func (s *LambdaGateway) handleScheduleEvent(ctx context.Context, schedules schedules.ScheduleRequestHandler, evt nitricScheduleEvent) (interface{}, error) {
	if evt.Schedule == "" {
		return nil, fmt.Errorf("unable to identify source nitric schedule")
	}

	request := &schedulespb.ServerMessage{
		// Send empty data for now (no reason to send data for schedules at the moment)
		Content: &schedulespb.ServerMessage_IntervalRequest{
			IntervalRequest: &schedulespb.IntervalRequest{
				ScheduleName: evt.Schedule,
			},
		},
	}

	_, err := schedules.HandleRequest(request)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleSnsEvents translates AWS SNS events to Nitric topic events and forwards them to be handled by registered workers.
func (s *LambdaGateway) handleSnsEvents(ctx context.Context, subscriptions topics.SubscriptionRequestHandler, records []Record) (interface{}, error) {
	for _, snsRecord := range records {
		messageString := snsRecord.SNS.Message
		attrs := map[string]string{}

		for k, v := range snsRecord.SNS.MessageAttributes {
			sv, ok := v.(string)
			if ok {
				attrs[k] = sv
			}
		}

		tName, err := s.getTopicNameForArn(ctx, snsRecord.SNS.TopicArn)
		if err != nil {
			logger.Errorf("unable to find nitric topic: %v", err)
			continue
		}

		messageBytes, err := base64.StdEncoding.DecodeString(messageString)
		if err != nil {
			logger.Errorf("unable decode SNS payload: %v", err)
			continue
		}

		var message topicspb.Message

		if err := proto.Unmarshal(messageBytes, &message); err != nil {
			logger.Errorf("unable to unmarshal nitric message from SNS trigger: %v", err)
			continue
		}

		request := &topicspb.ServerMessage{
			Content: &topicspb.ServerMessage_MessageRequest{
				MessageRequest: &topicspb.MessageRequest{
					TopicName: tName,
					Message:   &message,
				},
			},
		}

		resp, err := subscriptions.HandleRequest(request)
		if err != nil {
			return nil, err
		}

		if !resp.GetMessageResponse().Success {
			return nil, fmt.Errorf("event processing failed")
		}
	}

	return nil, nil
}

// handleHealthCheck responds to AWS Lambda service health checks with a 'healthy' response.
func (s *LambdaGateway) handleHealthCheck(ctx context.Context, evt healthCheckEvent) (interface{}, error) {
	return map[string]interface{}{
		"healthy": true,
	}, nil
}

// Converts an AWS Lambda S3 event type to the corresponding nitric blob event type
func s3EventTypeToNitricBlobEventType(eventType string) (*storagepb.BlobEventType, error) {
	if ok := strings.Contains(eventType, "ObjectCreated:"); ok {
		return storagepb.BlobEventType_Created.Enum(), nil
	} else if ok := strings.Contains(eventType, "ObjectRemoved:"); ok {
		return storagepb.BlobEventType_Deleted.Enum(), nil
	}
	return nil, fmt.Errorf("unsupported blob event type, expected ObjectCreated or ObjectRemoved, got %s", eventType)
}

func (s *LambdaGateway) processS3Event(ctx context.Context, storageListeners storage.BucketRequestHandler, records []Record) (interface{}, error) {
	for _, s3Record := range records {
		bucketName, err := s.getBucketNameForArn(ctx, s3Record.EventSourceArn)
		if err != nil {
			logger.Errorf("unable to find nitric bucket: %s", err.Error())
			return nil, fmt.Errorf("unable to find nitric bucket: %w", err)
		}

		eventType, err := s3EventTypeToNitricBlobEventType(s3Record.EventName)
		if err != nil {
			return nil, err
		}

		msg := &storagepb.ServerMessage{
			Content: &storagepb.ServerMessage_BlobEventRequest{
				BlobEventRequest: &storagepb.BlobEventRequest{
					BucketName: bucketName,
					Event: &storagepb.BlobEventRequest_BlobEvent{
						BlobEvent: &storagepb.BlobEvent{
							Key:  s3Record.S3.Object.Key,
							Type: *eventType,
						},
					},
				},
			},
		}

		resp, err := storageListeners.HandleRequest(msg)
		if err != nil {
			return nil, err
		}

		if !resp.GetBlobEventResponse().Success {
			return nil, fmt.Errorf("failed to process blob event")
		}
	}

	return nil, nil
}

func (s *LambdaGateway) routeEvent(ctx context.Context, opts *gateway.GatewayStartOpts, evt Event) (interface{}, error) {
	switch evt.Type() {
	case websocketEvent:
		return s.handleWebsocketEvent(ctx, opts.WebsocketListenerPlugin, evt.APIGatewayWebsocketProxyRequest)
	case httpEvent:
		return s.handleApiEvent(ctx, opts.ApiPlugin, opts.HttpPlugin, evt.APIGatewayV2HTTPRequest)
	case healthcheck:
		return s.handleHealthCheck(ctx, evt.healthCheckEvent)
	case sns:
		return s.handleSnsEvents(ctx, opts.TopicsListenerPlugin, evt.Records)
	case s3:
		return s.processS3Event(ctx, opts.StorageListenerPlugin, evt.Records)
	case schedule:
		return s.handleScheduleEvent(ctx, opts.SchedulesPlugin, evt.nitricScheduleEvent)
	default:
		return nil, fmt.Errorf("unhandled lambda event type: %+v", evt)
	}
}

// Start polling the lambda runtime for events and route the to workers for processing
func (s *LambdaGateway) Start(opts *gateway.GatewayStartOpts) error {
	// Begin polling lambda for incoming requests...
	s.runtime(func(ctx context.Context, evt Event) (interface{}, error) {
		a, err := s.routeEvent(ctx, opts, evt)

		tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
		if ok {
			_ = tp.ForceFlush(ctx)
		}

		return a, err
	})
	// Unblock the 'Stop' function if it's waiting.
	go func() { s.finished <- 1 }()
	return nil
}

// Stop will block until the lambda runtime is finished
func (s *LambdaGateway) Stop() error {
	// This is a NO_OP Process, as this is a pull based system
	// We don't need to stop listening to anything
	log.Default().Println("gateway 'Stop' called, waiting for lambda runtime to finish")
	// IT CANNOT BE STOPPED!!! Lambda is done when it wants to be and you won't change its mind.
	// But seriously we set the with SIGTERM option in Start for automatic graceful shutdown
	<-s.finished
	return nil
}

func New(provider *resource.AwsResourceService) (gateway.GatewayService, error) {
	return NewWithRuntime(provider, lambda.Start)
}

func NewWithRuntime(provider resource.AwsResourceProvider, runtime LambdaRuntimeHandler) (gateway.GatewayService, error) {
	return &LambdaGateway{
		provider: provider,
		runtime:  runtime,
		finished: make(chan int),
	}, nil
}
