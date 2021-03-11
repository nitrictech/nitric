package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/sources"
)

// HttpHandler - The http handler for the membrane when operating in HTTP_PROXY mode
type HttpHandler struct {
	// Get the host we're sending to
	host string
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *HttpHandler) HandleEvent(source *sources.Event) error {
	address := fmt.Sprint("http://%s/subscriptions/%s", h.host, source.Topic)
	httpRequest, _ := http.NewRequest("POST", address, ioutil.NopCloser(bytes.NewReader(source.Payload)))
	httpRequest.Header.Add("x-nitric-request-id", source.ID)
	httpRequest.Header.Add("x-nitric-source-type", "SUBSCRIPTION")
	httpRequest.Header.Add("x-nitric-source", source.Topic)

	// TODO: Handle response or error and response appropriately
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
func (h *HttpHandler) HandleHttpRequest(source *sources.HttpRequest) *http.Response {
	address := fmt.Sprint("http://%s%s", h.host, source.Path)
	httpRequest, err := http.NewRequest(source.Method, address, source.Body)
	httpRequest.Header = source.Header
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
