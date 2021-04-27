package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/triggers"
	"github.com/valyala/fasthttp"
)

// FaaSHandler - trigger handler for the membrane when operating in FaaS mode
type FaasHandler struct {
	host string
}

// errorToInternalServerError
// Converts a generic golang error to a HTTP(500) response
func errorToInternalServerError(err error) *http.Response {
	return &http.Response{
		Status:     "Internal Server Error",
		StatusCode: 500,
		// TODO: Eat error in non development modes
		// TODO: Log the error to an external log sink
		Body: ioutil.NopCloser(bytes.NewReader([]byte(err.Error()))),
	}
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *FaasHandler) HandleEvent(trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s", h.host)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)
	httpRequest.Header.Add("x-nitric-request-id", trigger.ID)
	httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Subscription.String())
	httpRequest.Header.Add("x-nitric-source", trigger.Topic)

	resp := &fasthttp.Response{}

	err := fasthttp.Do(httpRequest, resp)

	// All good if we got a 2XX error code
	if resp != nil && resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
		return nil
	} else if resp != nil {
		return fmt.Errorf("Error processing event (%d): %s", resp.StatusCode, string(resp.Body()))
	}

	return fmt.Errorf("Error processing event: %s", err.Error())
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *FaasHandler) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	address := fmt.Sprintf("http://%s", h.host)

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetRequestURI(address)

	for key, val := range trigger.Header {
		httpRequest.Header.Add(key, val)
	}

	httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Request.String())
	httpRequest.Header.Add("x-nitric-source", fmt.Sprintf("%s:%s", trigger.Method, trigger.Path))

	var resp fasthttp.Response
	err := fasthttp.Do(httpRequest, &resp)

	if err != nil {
		return nil, err
	}

	return triggers.FromHttpResponse(&resp), nil
}

func NewFaasHandler(host string) *FaasHandler {
	return &FaasHandler{
		host: host,
	}
}
