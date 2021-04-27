// The Digital Ocean App Platform HTTP gateway plugin
package appplatform_service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
)

type HttpGateway struct {
	address string
	sdk.UnimplementedGatewayPlugin
}

func (s *HttpGateway) Start(handler handler.TriggerHandler) error {
	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle Event/Subscription Request Types
		// TODO: Determine how we will handle nitric events for digital ocean

		httpReq := triggers.FromHttpRequest(req)
		// Handle HTTP Request Types
		response := handler.HandleHttpRequest(httpReq)
		responsePayload, _ := ioutil.ReadAll(response.Body)

		for key := range response.Header {
			resp.Header().Add(key, response.Header.Get(key))
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
	address := utils.GetEnv("GATEWAY_ADDRESS", ":9001")

	return &HttpGateway{
		address: address,
	}, nil
}
