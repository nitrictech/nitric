// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
	"github.com/valyala/fasthttp"
)

type HttpGateway struct {
	address string
	sdk.UnimplementedGatewayPlugin
}

// TODO: Lets bind this to a struct...
func httpHandler(handler handler.TriggerHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var triggerTypeString = string(ctx.Request.Header.Peek("x-nitric-source-type"))

		// Handle Event/Subscription Request Types
		if strings.ToUpper(triggerTypeString) == triggers.TriggerType_Subscription.String() {
			trigger := string(ctx.Request.Header.Peek("x-nitric-source"))
			requestId := string(ctx.Request.Header.Peek("x-nitric-request-id"))
			payload := ctx.PostBody()

			err := handler.HandleEvent(&triggers.Event{
				ID:      requestId,
				Topic:   trigger,
				Payload: payload,
			})

			if err != nil {
				// TODO: Make this more informative
				ctx.Error("There was an error processing the event", 500)
			} else {
				ctx.SuccessString("text/plain", "Successfully Handled the Event")
			}

			// return here...
			return
		}

		httpReq := triggers.FromHttpRequest(ctx)
		// Handle HTTP Request Types
		response, err := handler.HandleHttpRequest(httpReq)

		ctx.Response = fasthttp.Response{}

		response.Header.VisitAll(func(key []byte, val []byte) {
			if len(ctx.Response.Header.PeekBytes(key)) < 1 {
				ctx.Response.Header.AddBytesKV(key, val)
			}
		})
		ctx.Response.Header.Del("Content-Length")
		ctx.Response.Header.Del("Connection")

		ctx.Response.SetBody(response.Body)
		ctx.Response.SetStatusCode(response.StatusCode)

		if err != nil {
			// TODO: Redact message in production
			ctx.Error(err.Error(), 500)
		}
	}
}

func (s *HttpGateway) Start(handler handler.TriggerHandler) error {

	// Start the fasthttp server
	fasthttp.ListenAndServe(s.address, httpHandler(handler))

	// Setup the function handler for the default (catch all route)
	//http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
	//	// Handle the HTTP response...
	//	headers := req.Header

	//	var triggerTypeString = headers.Get("x-nitric-source-type")

	//	// Handle Event/Subscription Request Types
	//	if strings.ToUpper(triggerTypeString) == triggers.TriggerType_Subscription.String() {
	//		trigger := headers.Get("x-nitric-source")
	//		requestId := headers.Get("x-nitric-request-id")
	//		payload, _ := ioutil.ReadAll(req.Body)

	//		err := handler.HandleEvent(&triggers.Event{
	//			ID:      requestId,
	//			Topic:   trigger,
	//			Payload: payload,
	//		})

	//		if err != nil {
	//			// TODO: Make this more informative
	//			resp.WriteHeader(500)
	//			resp.Write([]byte("There was an error processing the event"))
	//		} else {
	//			resp.WriteHeader(200)
	//			resp.Write([]byte("Successfully Handled the Event"))
	//		}

	//		// return here...
	//		return
	//	}

	//	httpReq := triggers.FromHttpRequest(req)
	//	// Handle HTTP Request Types
	//	response := handler.HandleHttpRequest(httpReq)
	//	responsePayload, _ := ioutil.ReadAll(response.Body)

	//	for key := range response.Header {
	//		resp.Header().Add(key, response.Header.Get(key))
	//	}

	//	resp.WriteHeader(response.StatusCode)
	//	resp.Write(responsePayload)
	//})

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
