// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/membrane/handler"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/plugins/sdk/sources"
	"github.com/nitric-dev/membrane/utils"
)

type HttpGateway struct {
	address string
	sdk.UnimplementedGatewayPlugin
}

func (s *HttpGateway) Start(handler handler.SourceHandler) error {
	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle the HTTP response...
		headers := req.Header

		var sourceTypeString = headers.Get("x-nitric-source-type")

		// Handle Event/Subscription Request Types
		if strings.ToLower(sourceTypeString) == "subscription" {
			sourceType = sdk.Subscription
			source = headers.Get("x-nitric-source")
			requestId = headers.Get("x-nitric-request-id")
			payload, _ := ioutil.ReadAll(req.Body)

			handler.HandleEvent(&sources.Event{
				ID: requestId,
				Topic: source,
				Payload: payload
			})

			resp.WriteHeader(response.Status)
			resp.Write(response.Body)
			// return here...
			return
		}

		// Handle HTTP Request Types
		response := handler.HandleHttpRequest(&sources.FromHttpRequest(req))

		responsePayload, _ := ioutil.ReadAll(response.Body)

		for key, val := range response.Header {
			resp.Header().Add(key, val)
		}
		
		resp.WriteHeader(response.StatusCode)
		resp.Write(responsePayload)
	})

	// Start a HTTP Proxy server here...
	httpError := http.ListenAndServe(fmt.Sprintf("%s", s.address), nil)

	return httpError
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpGateway{
		address: address,
	}, nil
}
