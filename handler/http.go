package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/triggers"
)

// HttpHandler - The http handler for the membrane when operating in HTTP_PROXY mode
type HttpHandler struct {
	// Get the host we're sending to
	host string
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *HttpHandler) HandleEvent(trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s/subscriptions/%s", h.host, trigger.Topic)
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
func (h *HttpHandler) HandleHttpRequest(trigger *triggers.HttpRequest) *http.Response {
	address := fmt.Sprintf("http://%s%s", h.host, trigger.Path)
	httpRequest, err := http.NewRequest(trigger.Method, address, trigger.Body)
	httpRequest.Header = trigger.Header
	defaultErr := &http.Response{
		Status:     "Internal Server Error",
		StatusCode: 500,
	}

	if err != nil {
		return defaultErr
	}

	resp, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return defaultErr
	}

	return resp
}

func NewHttpHandler(host string) *HttpHandler {
	return &HttpHandler{
		host: host,
	}
}
