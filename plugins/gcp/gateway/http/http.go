// The GCP HTTP gateway plugin
package http_plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	eventingPb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type HttpProxyGateway struct {
	address string
}

func (s *HttpProxyGateway) Start(handler sdk.GatewayHandler) error {

	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle the HTTP response...
		headers := req.Header

		var sourceType = sdk.Request

		var source = headers.Get("User-Agent")
		var contentType = headers.Get("Content-Type")
		requestId := headers.Get("x-nitric-request-id")
		payloadType := headers.Get("x-nitric-payload-type")

		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// Return a http error here...
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(500)
			// TODO: Remove this unless in debug mode...
			resp.Write([]byte(err.Error()))
		}

		var payload = bytes

		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal(bytes, &jsonBody); err == nil {
			if sub, ok := jsonBody["subscription"]; ok {
				// We have a pubsub request here...
				sourceType = sdk.Subscription
				// TODO: Normalize the topic name from the subscription
				source = sub.(string)
				if message, ok := jsonBody["message"].(map[string]interface{}); ok {
					if bytes, err := base64.StdEncoding.DecodeString(message["data"].(string)); err == nil {
						nitricEvent := eventingPb.NitricEvent{}
						// We'll contine here...
						if err := json.Unmarshal(bytes, &nitricEvent); err == nil {
							// We have an offical NitricEvent payload here...
							requestId = nitricEvent.GetRequestId()
							payloadType = nitricEvent.GetPayloadType()
							if payloadBytes, err := nitricEvent.GetPayload().MarshalJSON(); err == nil {
								payload = payloadBytes
								contentType = http.DetectContentType(payloadBytes)
							}

						} else {
							// We recieved an event with no Nitric related context...
							// Just jam it into the payload and send it
							payload = bytes
							contentType = http.DetectContentType(payload)
						}
					} else {
						// There was a problem capturing the subscription...
						// For now let's log and continue...
					}
				}
			}
		}

		// Carry on if our data isn't formatted in json anyway...
		nitricContext := &sdk.NitricContext{
			RequestId:   requestId,
			PayloadType: payloadType,
			Source:      source,
			SourceType:  sourceType,
		}

		// Call the membrane function handler
		response := handler(&sdk.NitricRequest{
			Context:     nitricContext,
			Payload:     payload,
			ContentType: contentType,
		})

		for name, value := range response.Headers {
			resp.Header().Add(name, value)
		}

		// Pass through the function response
		resp.WriteHeader(response.Status)
		resp.Write(response.Body)
	})

	// Start a HTTP Proxy server here...
	httpError := http.ListenAndServe(fmt.Sprintf("%s", s.address), nil)

	return httpError
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayPlugin, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpProxyGateway{
		address: address,
	}, nil
}
