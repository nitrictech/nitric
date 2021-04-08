// The AWS HTTP gateway plugin
package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/triggers"
	"github.com/nitric-dev/membrane/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	gw "github.com/nitric-dev/membrane/sdk"
)

const (
	AMZ_MESSAGE_ID   = "x-amz-sns-message-id"
	AMZ_MESSAGE_TYPE = "x-amz-sns-message-type"
	AMZ_TOPIC_ARN    = "x-amz-sns-topic-arn"
)

type HttpProxyGateway struct {
	client  *sns.SNS
	address string
}

func (s *HttpProxyGateway) Start(handler handler.TriggerHandler) error {

	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle the HTTP response...
		headers := req.Header

		var trigger = strings.Join(headers["User-Agent"], "")

		if trigger == "Amazon Simple Notification Service Agent" {
			// If its a subscribe or unsubscribe notification then we need to handle it
			amzMessageType := headers.Get(AMZ_MESSAGE_TYPE)
			topicArn := headers.Get(AMZ_TOPIC_ARN)
			id := headers.Get(AMZ_MESSAGE_ID)

			payload, _ := ioutil.ReadAll(req.Body)

			// SNS bodies are always JSON
			var jsonBody map[string]interface{}
			unmarshalError := json.Unmarshal(payload, &jsonBody)
			if unmarshalError != nil {
				// Return an error response
				resp.Header().Add("Content-Type", "text/plain")
				resp.Write([]byte("There was an error unmarshalling an SNS message"))
				resp.WriteHeader(403)
				return
			}

			if amzMessageType == "SubscriptionConfirmation" {
				token := jsonBody["Token"].(string)
				// call to confirm the subscription and return a 200 OK
				// We don't need to perform any processing on this type of request
				s.client.ConfirmSubscription(&sns.ConfirmSubscriptionInput{
					TopicArn: &topicArn,
					Token:    &token,
				})

				resp.WriteHeader(200)
				return
			} else if amzMessageType == "UnsubscribeConfirmation" {
				// FIXME: Decide how we need to handle this
				resp.WriteHeader(200)
				return
			}

			if err := handler.HandleEvent(&triggers.Event{
				ID: id,
				// FIXME: Split this to retrive the nitric topic name
				Topic:   topicArn,
				Payload: payload,
			}); err == nil {
				// Return a positive (200) response
				resp.WriteHeader(200)
				resp.Write([]byte("Success"))
			} else {
				// return a negative (500) response
				resp.WriteHeader(500)
				// TODO: Debug mode for printing errror here...
				resp.Write([]byte("Internal Server Error"))
			}

			return
		}

		// Otherwise treat as a normal http request
		response := handler.HandleHttpRequest(triggers.FromHttpRequest(req))
		responseBody, _ := ioutil.ReadAll(response.Body)

		for name := range response.Header {
			resp.Header().Add(name, response.Header.Get(name))
		}

		// Pass through the function response
		resp.WriteHeader(response.StatusCode)
		resp.Write(responseBody)
	})

	// Start a HTTP Proxy server here...
	httpError := http.ListenAndServe(fmt.Sprintf("%s", s.address), nil)

	return httpError
}

// Create new AWS HTTP Gateway service
func New() (gw.GatewayService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")
	address := utils.GetEnv("GATEWAY_ADDRESS", "0.0.0.0:9001")

	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use ENV configuration
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	snsClient := sns.New(sess)

	return &HttpProxyGateway{
		client:  snsClient,
		address: address,
	}, nil
}
