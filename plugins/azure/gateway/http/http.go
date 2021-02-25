package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/mitchellh/mapstructure"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// HttpService - The HTTP gateway plugin for Azure
type HttpService struct {
	address string
}

func (s *HttpService) handleSubscriptionValidation(w http.ResponseWriter, events []eventgrid.Event) {
	subPayload := events[0]
	var validateData eventgrid.SubscriptionValidationEventData
	if err := mapstructure.Decode(subPayload.Data, &validateData); err == nil {
		// Some error here...
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte("Invalid subscription event data"))
		return
	}

	response := eventgrid.SubscriptionValidationResponse{
		ValidationResponse: validateData.ValidationCode,
	}

	responseBody, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	// TODO: Remove this unless in debug mode...
	w.Write(responseBody)
}

func (s *HttpService) handleNotifications(w http.ResponseWriter, events []eventgrid.Event, handler sdk.GatewayHandler) {
	for _, event := range events {
		// XXX: Assume we have a nitric event for now
		nitricEvt := sdk.NitricEvent{}
		if err := mapstructure.Decode(event.Data, &nitricEvt); err == nil {
			// We have a valid nitric event
			// Decode and pass to our function
			//requestId = nitricEvt.RequestId
			//payload, _ = json.Marshal(nitricEvt.Payload)
			//payloadType = nitricEvt.PayloadType
			// Carry on if our data isn't formatted in json anyway...
			nitricContext := &sdk.NitricContext{
				RequestId:   nitricEvt.RequestId,
				PayloadType: nitricEvt.PayloadType,
				Source:      *event.Topic,
				SourceType:  sdk.Subscription,
			}

			bytes, _ := json.Marshal(nitricEvt.Payload)
			// Call the membrane function handler
			// TODO: Handle response
			handler(&sdk.NitricRequest{
				Context:     nitricContext,
				Payload:     bytes,
				ContentType: "application/json",
			})
		}
	}

	// Return 200 OK (TODO: Determine how we could mark individual events for failure)
	// Or potentially requeue them here internally...
	w.WriteHeader(200)
	w.Write([]byte(""))
}

func (s *HttpService) handleRequest(w http.ResponseWriter, r *http.Request, handler sdk.GatewayHandler) {
	source := r.Header.Get("User-Agent")
	contentType := r.Header.Get("Content-Type")
	requestId := r.Header.Get("x-nitric-request-id")
	payloadType := r.Header.Get("x-nitric-payload-type")

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Return a http error here...
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(500)
		// TODO: Remove this unless in debug mode...
		w.Write([]byte(err.Error()))

		return
	}

	// Carry on if our data isn't formatted in json anyway...
	nitricContext := &sdk.NitricContext{
		RequestId:   requestId,
		PayloadType: payloadType,
		Source:      source,
		SourceType:  sdk.Request,
	}

	// Call the membrane function handler
	response := handler(&sdk.NitricRequest{
		Context:     nitricContext,
		Payload:     bytes,
		ContentType: contentType,
	})

	for name, value := range response.Headers {
		w.Header().Add(name, value)
	}

	// Pass through the function response
	w.WriteHeader(response.Status)
	w.Write(response.Body)
}

func (s *HttpService) Start(handler sdk.GatewayHandler) error {

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		eventType := req.Header.Get("aeg-event-type")

		// Handle an eventgrid webhook event
		if eventType != "" {
			var eventgridEvents []eventgrid.Event
			bytes, _ := ioutil.ReadAll(req.Body)
			// TODO: verify topic for validity
			if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
				resp.Header().Add("Content-Type", "text/plain")
				resp.WriteHeader(400)
				resp.Write([]byte(fmt.Sprintf("Invalid event grid types")))
				return
			}

			// Handle Eventgrid event
			if eventType == "SubscriptionValidation" {
				// Validate a subscription
				s.handleSubscriptionValidation(resp, eventgridEvents)
				return
			} else if eventType == "Notification" {
				// Handle notifications
				s.handleNotifications(resp, eventgridEvents, handler)
				return
			}
		}

		// Handle a standard HTTP request
		s.handleRequest(resp, req, handler)
	})

	// Start a HTTP server here...
	httpError := http.ListenAndServe(s.address, nil)

	return httpError
}

// Create a new HTTP Gateway plugin
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpService{
		address: address,
	}, nil
}
