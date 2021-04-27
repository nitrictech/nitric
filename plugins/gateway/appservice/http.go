package http_service

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/mitchellh/mapstructure"
	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"
	"github.com/valyala/fasthttp"
)

// HttpService - The HTTP gateway plugin for Azure
type HttpService struct {
	address string
}

func handleSubscriptionValidation(ctx *fasthttp.RequestCtx, events []eventgrid.Event) {
	subPayload := events[0]
	var validateData eventgrid.SubscriptionValidationEventData
	if err := mapstructure.Decode(subPayload.Data, &validateData); err != nil {
		ctx.Error("Invalid subscription event data", 400)
		return
	}

	response := eventgrid.SubscriptionValidationResponse{
		ValidationResponse: validateData.ValidationCode,
	}

	responseBody, _ := json.Marshal(response)

	ctx.Success("application/json", responseBody)
}

func handleNotifications(ctx *fasthttp.RequestCtx, events []eventgrid.Event, handler handler.TriggerHandler) {
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

		// FIXME: Handle error
		handler.HandleEvent(&triggers.Event{
			// FIXME: Check if ID is nil
			ID:      *event.ID,
			Topic:   *event.Topic,
			Payload: payloadBytes,
		})
	}

	// Return 200 OK (TODO: Determine how we could mark individual events for failure)
	// Or potentially requeue them here internally...
	ctx.SuccessString("text/plain", "success")
}

func handleRequest(ctx *fasthttp.RequestCtx, handler handler.TriggerHandler) {
	response, err := handler.HandleHttpRequest(triggers.FromHttpRequest(ctx))

	if err != nil {
		ctx.Error(fmt.Sprintf("Error Handling Request: %v", err), 500)
		return
	}
	if response.Header != nil {
		response.Header.VisitAll(func(key []byte, val []byte) {
			ctx.Response.Header.AddBytesKV(key, val)
		})
	}

	// Avoid content length header duplication
	ctx.Response.Header.Del("Content-Length")
	ctx.Response.SetStatusCode(response.StatusCode)
	ctx.Response.SetBody(response.Body)
}

func httpHandler(handler handler.TriggerHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// Handle Event/Subscription Request Types
		eventType := string(ctx.Request.Header.Peek("aeg-event-type"))

		// Handle an eventgrid webhook event
		if eventType != "" {
			var eventgridEvents []eventgrid.Event
			bytes := ctx.Request.Body()
			// TODO: verify topic for validity
			if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
				ctx.Error("Invalid event grid types", 400)
				return
			}

			// Handle Eventgrid event
			if eventType == "SubscriptionValidation" {
				// Validate a subscription
				handleSubscriptionValidation(ctx, eventgridEvents)
				return
			} else if eventType == "Notification" {
				// Handle notifications
				handleNotifications(ctx, eventgridEvents, handler)
				return
			}
		}

		// Handle a standard HTTP request
		handleRequest(ctx, handler)
	}
}

func (s *HttpService) Start(handler handler.TriggerHandler) error {
	httpError := fasthttp.ListenAndServe(s.address, httpHandler(handler))

	return httpError
}

// Create a new HTTP Gateway plugin
func New() (sdk.GatewayService, error) {
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	return &HttpService{
		address: address,
	}, nil
}
