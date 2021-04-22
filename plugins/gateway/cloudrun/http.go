// The GCP HTTP gateway plugin for CloudRun
package cloudrun_plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
)

type HttpProxyGateway struct {
	address string
}

type PubSubMessage struct {
	Message struct {
		Attributes map[string]string `json:"attributes"`
		Data       []byte            `json:"data,omitempty"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (s *HttpProxyGateway) Start(handler handler.TriggerHandler) error {

	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// Return a http error here...
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(500)
			// TODO: Remove this unless in debug mode...
			resp.Write([]byte(err.Error()))
		}

		// Check if the payload contains a pubsub event
		// TODO: We probably want to use a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err = json.Unmarshal(bodyBytes, &pubsubEvent); err == nil {
			// We have an event from pubsub here...
			event := &triggers.Event{
				ID: pubsubEvent.Message.ID,
				// Set the topic
				Topic: pubsubEvent.Message.Attributes["x-nitric-topic"],
				// Set the payload
				Payload: pubsubEvent.Message.Data,
			}

			if err := handler.HandleEvent(event); err == nil {
				// return a successful response
				resp.WriteHeader(200)
				resp.Write([]byte("Success"))
			} else {
				resp.WriteHeader(500)
				// FIXME: fix operating mode here...
				resp.Write([]byte(fmt.Sprintf("Error handling event %v", err)))
			}

			return
		}
		reader := ioutil.NopCloser(bytes.NewReader(bodyBytes))
		// We don't have an event, so treat as a HTTP request for now
		req.Body = reader
		httpTrigger := triggers.FromHttpRequest(req)
		response := handler.HandleHttpRequest(httpTrigger)
		responseBody, _ := ioutil.ReadAll(response.Body)

		for key, _ := range response.Header {
			resp.Header().Add(key, response.Header.Get(key))
		}

		// Pass through the function response
		resp.WriteHeader(response.StatusCode)
		resp.Write(responseBody)
	})

	// Start a HTTP Proxy server here...
	httpError := http.ListenAndServe(fmt.Sprintf("%s", s.address), nil)

	return httpError
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpProxyGateway{
		address: address,
	}, nil
}
