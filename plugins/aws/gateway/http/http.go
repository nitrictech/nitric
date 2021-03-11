// The AWS HTTP gateway plugin
package http_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	gw "github.com/nitric-dev/membrane/sdk"
)

type HttpProxyGateway struct {
	client  *sns.SNS
	address string
}

func (s *HttpProxyGateway) Start(handler gw.GatewayHandler) error {

	// Setup the function handler for the default (catch all route)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// Handle the HTTP response...
		headers := req.Header

		var sourceType = gw.Request
		var source = strings.Join(headers["User-Agent"], "")
		var requestId = strings.Join(headers["x-nitric-request-id"], "")
		var payloadType = strings.Join(headers["x-nitric-payload-type"], "")
		// var contentType = strings.Join(headers["Content-Type"], "")
		// var timestamp = &timestamp.Timestamp{}
		var payload, _ = ioutil.ReadAll(req.Body)

		if source == "Amazon Simple Notification Service Agent" {
			// If its a subscribe or unsubscribe notification then we need to handle it
			amzMessageType := strings.Join(headers["x-amz-sns-message-type"], "")
			topicArn := strings.Join(headers["x-amz-sns-topic-arn"], "")

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

			// We know that the source is now a topic
			sourceType = gw.Subscription
			topicParts := strings.Split(topicArn, ":")
			source = topicParts[len(topicParts)-1] // get the topic name from the full ARN.
		}

		nitricContext := &gw.NitricContext{
			RequestId:   requestId,
			PayloadType: payloadType,
			Source:      source,
			SourceType:  sourceType,
		}

		// Call the membrane function handler
		response := handler(&gw.NitricRequest{
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

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
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
