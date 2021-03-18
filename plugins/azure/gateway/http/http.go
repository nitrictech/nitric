package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/mitchellh/mapstructure"
	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/sources"
	"github.com/nitric-dev/membrane/utils"
)

// HttpService - The HTTP gateway plugin for Azure
type HttpService struct {
	address string
}

func (s *HttpService) handleSubscriptionValidation(w http.ResponseWriter, events []eventgrid.Event) {
	subPayload := events[0]
	var validateData eventgrid.SubscriptionValidationEventData
	if err := mapstructure.Decode(subPayload.Data, &validateData); err != nil {
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

func (s *HttpService) handleNotifications(w http.ResponseWriter, events []eventgrid.Event, handler handler.SourceHandler) {
	// FIXME: As we are batch handling events in azure
	// how do we notify of failed event handling?
	for _, event := range events {
		// XXX: Assume we have a nitric event for now
		// We have a valid nitric event
		// Decode and pass to our function
		var payloadBytes []byte
		if stringData, ok := event.Data.(string); ok {
			payloadBytes = []byte(stringData)
		} else if byteData, ok := event.Data.([]byte); ok {
			payloadBytes = byteData
		} else {
			// Assume a json serializable struct for now...
			payloadBytes, _ = json.Marshal(event.Data)
		}

		// FIXME: Handle error...

		// FIXME: Handle error
		handler.HandleEvent(&sources.Event{
			// FIXME: Check if ID is nil
			ID:      *event.ID,
			Topic:   *event.Topic,
			Payload: payloadBytes,
		})
	}

	// Return 200 OK (TODO: Determine how we could mark individual events for failure)
	// Or potentially requeue them here internally...
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func (s *HttpService) handleRequest(w http.ResponseWriter, r *http.Request, handler handler.SourceHandler) {
	response := handler.HandleHttpRequest(sources.FromHttpRequest(r))

	for name := range response.Header {
		w.Header().Add(name, response.Header.Get(name))
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	// Pass through the function response
	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}

func (s *HttpService) Start(handler handler.SourceHandler) error {

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
