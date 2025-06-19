package awslambda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nitrictech/nitric/server/runtime/service"
)

type awslambdaService struct {
	proxy service.Proxy
}

func (a *awslambdaService) httpNotFound() *events.LambdaFunctionURLStreamingResponse {
	return &events.LambdaFunctionURLStreamingResponse{
		StatusCode: http.StatusNotFound,
		Body:       strings.NewReader("Not Found"),
	}
}

func (a *awslambdaService) Start(proxy service.Proxy) error {
	a.proxy = proxy

	lambda.Start(func(ctx context.Context, evt json.RawMessage) (interface{}, error) {
		// Try to parse as API Gateway v2 HTTP event first
		fmt.Println("handling event", string(evt))

		var httpEvent events.LambdaFunctionURLRequest
		var scheduleEvent events.EventBridgeEvent
		if err := json.Unmarshal(evt, &httpEvent); err == nil && httpEvent.RequestContext.HTTP.Method != "" {
			return a.handleHTTPEvent(ctx, &httpEvent)
		} else if err := json.Unmarshal(evt, &scheduleEvent); err == nil {
			return a.handleScheduleEvent(ctx, &scheduleEvent)
		}

		// Handle other event types here if needed
		return nil, nil
	})

	return nil
}

type ScheduleEventDetail struct {
	Path string `json:"path"`
}

func (a *awslambdaService) handleScheduleEvent(ctx context.Context, evt *events.EventBridgeEvent) (interface{}, error) {
	// Get the path from the event
	var detail ScheduleEventDetail
	err := json.Unmarshal(evt.Detail, &detail)
	if err != nil {
		return nil, err
	}

	path := detail.Path

	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	resp, err := a.proxy.Forward(ctx, req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return nil, nil
}

func (a *awslambdaService) handleHTTPEvent(ctx context.Context, evt *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
	// Convert the event to an HTTP request and process it through the proxy
	// TODO: Implement full HTTP request handling logic with proxy

	// translate the event to a golang net/http request
	req, err := http.NewRequest(evt.RequestContext.HTTP.Method, evt.RawPath, strings.NewReader(evt.Body))
	if err != nil {
		return nil, err
	}

	// make sure headers are added to the request
	for k, v := range evt.Headers {
		req.Header.Add(k, v)
	}

	req.URL.RawQuery = evt.RawQueryString

	resp, err := a.proxy.Forward(ctx, req)
	if err != nil {
		return nil, err
	}

	// Translate response headers to a map
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	return &events.LambdaFunctionURLStreamingResponse{
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Headers:    headers,
	}, nil
}

func Plugin() (service.Service, error) {
	return &awslambdaService{}, nil
}
