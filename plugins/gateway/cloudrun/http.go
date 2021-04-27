// The GCP HTTP gateway plugin for CloudRun
package cloudrun_plugin

import (
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
	"github.com/valyala/fasthttp"
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

func httpHandler(handler handler.TriggerHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		bodyBytes := ctx.Request.Body()

		// Check if the payload contains a pubsub event
		// TODO: We probably want to use a simpler method than this
		// like reading off the request origin to ensure it is from pubsub
		var pubsubEvent PubSubMessage
		if err := json.Unmarshal(bodyBytes, &pubsubEvent); err == nil {
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
				ctx.SuccessString("text/plain", "success")
			} else {
				ctx.Error(fmt.Sprintf("Error handling event %v", err), 500)
			}

			return
		}

		httpTrigger := triggers.FromHttpRequest(ctx)
		response, err := handler.HandleHttpRequest(httpTrigger)

		if err != nil {
			ctx.Error(fmt.Sprintf("Error handling HTTP Request: %v", err), 500)
		}
		// responseBody, _ := ioutil.ReadAll(response.Body)
		if response.Header != nil {
			// Set headers...
			response.Header.VisitAll(func(key []byte, val []byte) {
				ctx.Response.Header.AddBytesKV(key, val)
			})
		}

		// Avoid content length header duplication
		ctx.Response.Header.Del("Content-Length")
		ctx.Response.SetStatusCode(response.StatusCode)
		ctx.Response.SetBody(response.Body)
	}
}

func (s *HttpProxyGateway) Start(handler handler.TriggerHandler) error {
	httpError := fasthttp.ListenAndServe(s.address, httpHandler(handler))

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
