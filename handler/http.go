package handler

import (
	"fmt"

	"github.com/nitric-dev/membrane/triggers"
	"github.com/valyala/fasthttp"
)

// HttpHandler - The http handler for the membrane when operating in HTTP_PROXY mode
type HttpHandler struct {
	// Get the host we're sending to
	address string
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *HttpHandler) HandleEvent(trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s/subscriptions/%s", h.address, trigger.Topic)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)
	httpRequest.Header.Add("x-nitric-request-id", trigger.ID)
	httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Subscription.String())
	httpRequest.Header.Add("x-nitric-source", trigger.Topic)

	var resp fasthttp.Response

	httpRequest.SetBody(trigger.Payload)
	httpRequest.Header.SetContentLength(len(trigger.Payload))

	// TODO: Handle response or error and respond appropriately
	err := fasthttp.Do(httpRequest, &resp)

	if &resp != nil && resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
		return nil
	} else if &resp != nil {
		return fmt.Errorf("Error processing event (%d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return fmt.Errorf("Error processing event: %s", err.Error())
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *HttpHandler) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	address := fmt.Sprintf("http://%s%s", h.address, trigger.Path)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)

	for key, val := range trigger.Query {
		httpRequest.URI().QueryArgs().Add(key, val)
	}

	for key, val := range trigger.Header {
		httpRequest.Header.Add(key, val)
	}

	httpRequest.Header.Del("Content-Length")
	httpRequest.SetBody(trigger.Body)
	httpRequest.Header.SetContentLength(len(trigger.Body))

	var resp fasthttp.Response
	err := fasthttp.Do(httpRequest, &resp)

	if err != nil {
		return nil, err
	}

	return triggers.FromHttpResponse(&resp), nil
}

func NewHttpHandler(host string) *HttpHandler {
	return &HttpHandler{
		address: host,
	}
}
