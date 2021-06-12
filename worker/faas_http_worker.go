package worker

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/valyala/fasthttp"
)

// A Nitric HTTP worker
type FaasHttpWorker struct {
	address string
}

var METHOD_TYPE = []byte("POST")

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *FaasHttpWorker) HandleEvent(trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s", h.address)
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	// Release resources after finishing
	defer func() {
		request.Reset()
		response.Reset()
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	triggerRequest := &pb.TriggerRequest{
		Data: trigger.Payload,
		Context: &pb.TriggerRequest_Topic{
			Topic: &pb.TopicTriggerContext{
				Topic: trigger.Topic,
			},
		},
	}

	if jsonData, err := json.Marshal(triggerRequest); err != nil {
		request := fasthttp.AcquireRequest()

		request.Header.SetContentType("application/json")
		request.SetBody(jsonData)
		request.SetRequestURI(address)

		err := fasthttp.Do(request, response)

		if err != nil {
			return fmt.Errorf("Function request failed")
		}

		// Response body should contain an instance of triggerResponse
		var triggerResponse pb.TriggerResponse
		err = json.Unmarshal(response.Body(), &triggerResponse)

		if err != nil {
			return err
		}

		topic := triggerResponse.GetTopic()

		if topic != nil {
			if topic.Success {
				return nil
			}

			return fmt.Errorf("Topic context indicated processing was unsuccesful")
		}

		return fmt.Errorf("Response from function did not contain topic context")
	} else {
		return fmt.Errorf("Error marshalling request")
	}
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *FaasHttpWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	address := fmt.Sprintf("http://%s", h.address)
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	// Release resources after finishing
	defer func() {
		request.Reset()
		response.Reset()
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	triggerRequest := &pb.TriggerRequest{
		Data: trigger.Body,
		Context: &pb.TriggerRequest_Http{
			Http: &pb.HttpTriggerContext{
				Method:      trigger.Method,
				QueryParams: trigger.Query,
				PathParams:  make(map[string]string),
			},
		},
	}

	if jsonData, err := json.Marshal(triggerRequest); err != nil {
		request := fasthttp.AcquireRequest()

		request.Header.SetContentType("application/json")
		request.SetBody(jsonData)
		request.SetRequestURI(address)

		err := fasthttp.Do(request, response)

		if err != nil {
			return nil, err
		}

		// Response body should contain an instance of triggerResponse
		var triggerResponse pb.TriggerResponse
		err = json.Unmarshal(response.Body(), &triggerResponse)

		if err != nil {
			return nil, err
		}

		return triggers.FromTriggerResponse(&triggerResponse)
	} else {
		return nil, err
	}
}

// Creates a new FaasHttpWorker
// Will wait to ensure that the provided address is dialable
// before proceeding
func NewFaasHttpWorker(address string) (*FaasHttpWorker, error) {
	// Dial the child port to see if it's open and ready...
	maxWaitTime := time.Duration(5) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		conn, _ := net.Dial("tcp", address)
		if conn != nil {
			conn.Close()
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return nil, fmt.Errorf("Unable to dial http worker, does it expose a http server at: %s?", address)
			}
		}
	}

	// Dial the provided address to ensure its availability
	return &FaasHttpWorker{
		address: address,
	}, nil
}
