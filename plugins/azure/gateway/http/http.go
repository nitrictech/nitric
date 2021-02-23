package http_service

import (
	"fmt"

	"github.com/nitric-dev/membrane/plugins/sdk"
)

// HttpService - The HTTP gateway plugin for Azure
type HttpService struct {
	address string
}
 
func (s *HttpService) Start(handler sdk.GatewayHandler) error {

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		headers := req.Header

		var sourceType = sdk.Request

		var source = headers.Get("User-Agent")
		var contentType = headers.Get("Content-Type")
		requestId := headers.Get("x-nitric-request-id")
		payloadType := headers.Get("x-nitric-payload-type")
		var payload = bytes

		// TODO: We need to acknowledge event grid messages sent to us as being valid

		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// Return a http error here...
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(500)
			// TODO: Remove this unless in debug mode...
			resp.Write([]byte(err.Error()))
		}

		// Example eventgrid handshake
		//[
		//	{
		//		"id": "2d1781af-3a4c-4d7c-bd0c-e34b19da4e66",
		//		"topic": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		//		"subject": "",
		//		"data": {
		//			"validationCode": "512d38b6-c7b8-40c8-89fe-f46f9e9622b6",
		//			"validationUrl": "https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=512d38b6-c7b8-40c8-89fe-f46f9e9622b6&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1A1A1A1A"
		//		},
		//		"eventType": "Microsoft.EventGrid.SubscriptionValidationEvent",
		//		"eventTime": "2018-01-25T22:12:19.4556811Z",
		//		"metadataVersion": "1",
		//		"dataVersion": "1"
		//	}
		//]

		// Validate subscription path
		if headers.Get("aeg-event-type") == "SubscriptionValidation" {
			var payload = bytes
			jsonBody := make([]map[string]interface{}, 0)
			// TODO: verify topic for validity
			if err = json.Unmarshal(bytes, &jsonBody); err == nil {
				subPayload := jsonBody[0]
				// We just need to get the data and echo it
				if data, ok := subPayload["data"]; ok {
					validatationData := data.(map[string]string)
					validationCode := validatationData["validationCode"]
					resp.Header().Add("Content-Type", "application/json")
					resp.WriteHeader(200)
					// TODO: Remove this unless in debug mode...
					resp.Write([]byte(fmt.Sprintf("{\"validationResponse\":\"%s\"}", validationCode))
					return
				}
			}

			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(200)
			// TODO: Remove this unless in debug mode...
			resp.Write([]byte("There was an error validating eventgrid subscription")
			return
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