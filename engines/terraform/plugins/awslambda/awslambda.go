package awslambda

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nitrictech/nitric/server/runtime/service"
)

type awslambdaService struct {
	proxy service.Proxy
}

func (a *awslambdaService) Start(proxy service.Proxy) error {
	a.proxy = proxy

	lambda.Start(func(ctx context.Context, evt json.RawMessage) (interface{}, error) {
		// Try to parse as API Gateway v2 HTTP event first
		var httpEvent events.APIGatewayV2HTTPRequest
		if err := json.Unmarshal(evt, &httpEvent); err == nil && httpEvent.RequestContext.HTTP.Method != "" {
			return a.handleHTTPEvent(ctx, &httpEvent)
		}

		// Handle other event types here if needed
		return nil, nil
	})

	return nil
}

func (a *awslambdaService) handleHTTPEvent(ctx context.Context, evt *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
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

	for k, v := range evt.QueryStringParameters {
		req.URL.Query().Add(k, v)
	}

	resp, err := a.proxy.Forward(ctx, req)
	if err != nil {
		return nil, err
	}

	// translate the response to an API Gateway v2 HTTP response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Translate response headers to a map
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Headers:    headers,
	}, nil
}

func Plugin() (service.Service, error) {
	return &awslambdaService{}, nil
}
