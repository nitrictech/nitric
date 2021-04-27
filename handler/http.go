package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

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
	httpRequest, _ := http.NewRequest("POST", address, ioutil.NopCloser(bytes.NewReader(trigger.Payload)))
	httpRequest.Header.Add("x-nitric-request-id", trigger.ID)
	httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Subscription.String())
	httpRequest.Header.Add("x-nitric-source", trigger.Topic)

	// TODO: Handle response or error and respond appropriately
	resp, err := http.DefaultClient.Do(httpRequest)

	if resp != nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	} else if resp != nil {
		respMessage, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("Error processing event (%d): %s", resp.StatusCode, string(respMessage))
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
