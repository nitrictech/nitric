// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type HttpGateway struct {
	address string
	sdk.UnimplementedGatewayPlugin
}

func (s *HttpGateway) Start(handler sdk.GatewayHandler) error {
	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle the HTTP response...
		headers := req.Header

		// var source = strings.Join(headers["User-Agent"], "")
		var requestId = strings.Join(headers["x-nitric-request-id"], "")
		var payloadType = strings.Join(headers["x-nitric-payload-type"], "")
		var sourceTypeString = strings.Join(headers["x-nitric-source-type"], "")
		var source = strings.Join(headers["x-nitric-source"], "")
		// var contentType = strings.Join(headers["Content-Type"], "")
		// var timestamp = &timestamp.Timestamp{}
		var payload, _ = ioutil.ReadAll(req.Body)

		// TODO: Create string to enum utility for SourceType
		var sourceType = sdk.Request

		if strings.ToLower(sourceTypeString) == "subscription" {
			sourceType = sdk.Subscription
		}

		nitricContext := &sdk.NitricContext{
			RequestId:   requestId,
			PayloadType: payloadType,
			Source:      source,
			SourceType:  sourceType,
		}

		// Call the membrane function handler
		response := handler(&sdk.NitricRequest{
			Context: nitricContext,
			Payload: payload,
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

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.GatewayPlugin, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpGateway{
		address: address,
	}, nil
}
