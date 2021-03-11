// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/sources"
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
		if strings.ToUpper(sourceTypeString) == sources.SourceType_Subscription.String() {
			source := headers.Get("x-nitric-source")
			requestId := headers.Get("x-nitric-request-id")
			payload, _ := ioutil.ReadAll(req.Body)

			err := handler.HandleEvent(&sources.Event{
				ID:      requestId,
				Topic:   source,
				Payload: payload,
			})

			if err != nil {
				// TODO: Make this more informative
				resp.WriteHeader(500)
				resp.Write([]byte("There was an error processing the event"))
			} else {
				resp.WriteHeader(200)
				resp.Write([]byte("Successfully Handled the Event"))
			}

			// return here...
			return
		}

		// Handle HTTP Request Types
		response := handler.HandleHttpRequest(sources.FromHttpRequest(req))

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
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpGateway{
		address: address,
	}, nil
}
