package lambda_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	events "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

type eventType int

const (
	unknown eventType = iota
	sns
	http
)

type LambdaRuntimeHandler func(handler interface{})

//Event incoming event
type Event struct {
	Requests []sdk.NitricRequest
}

func (event *Event) getEventType(data []byte) eventType {
	tmp := make(map[string]interface{})
	// Unmarshal so we can get just enough info about the type of event to fully deserialize it
	json.Unmarshal(data, &tmp)

	// If our event is a HTTP request
	if _, ok := tmp["rawPath"]; ok {
		return http
	} else if records, ok := tmp["Records"]; ok {
		recordsList, _ := records.([]interface{})
		record, _ := recordsList[0].(map[string]interface{})
		// We have some kind of event here...
		// we'll assume its an SNS
		var eventSource string
		if es, ok := record["EventSource"]; ok {
			eventSource = es.(string)
		} else if es, ok := record["eventSource"]; ok {
			eventSource = es.(string)
		}

		switch eventSource {
		case "aws:sns":
			return sns
		}
	}

	return unknown
}

// implement the unmarshal interface in order to handle multiple event types
func (event *Event) UnmarshalJSON(data []byte) error {
	var err error

	event.Requests = make([]sdk.NitricRequest, 0)

	switch event.getEventType(data) {
	case sns:
		snsEvent := &events.SNSEvent{}
		err = json.Unmarshal(data, snsEvent)

		if err == nil {
			// Map over the records and return
			for _, snsRecord := range snsEvent.Records {
				messageString := snsRecord.SNS.Message
				// FIXME: What about non-nitric SNS events???
				messageJson := &sdk.NitricEvent{}

				// Populate the JSON
				err = json.Unmarshal([]byte(messageString), messageJson)

				topicArn := snsRecord.SNS.TopicArn
				topicParts := strings.Split(topicArn, ":")
				source := topicParts[len(topicParts)-1]
				// get the topic name from the full ARN.
				// Get the topic name from the arn

				if err == nil {
					// Decode the message to see if it's a Nitric message
					context := sdk.NitricContext{
						SourceType:  sdk.Subscription,
						Source:      source,
						PayloadType: messageJson.PayloadType,
						RequestId:   messageJson.RequestId,
					}

					payloadMap := messageJson.Payload

					payloadBytes, err := json.Marshal(&payloadMap)

					if err == nil {
						event.Requests = append(event.Requests, sdk.NitricRequest{
							Context:     &context,
							Payload:     payloadBytes,
							ContentType: *aws.String("application/json"),
						})
					}
				}
			}
		}
		break
	case http:
		httpEvent := &events.APIGatewayV2HTTPRequest{}

		err = json.Unmarshal(data, httpEvent)

		if err == nil {
			nitricContext := sdk.NitricContext{
				SourceType:  sdk.Request,
				Source:      httpEvent.Headers["User-Agent"],
				PayloadType: httpEvent.Headers["x-nitric-payload-type"],
				RequestId:   httpEvent.Headers["x-nitric-request-id"],
			}

			event.Requests = append(event.Requests, sdk.NitricRequest{
				Context:     &nitricContext,
				ContentType: httpEvent.Headers["Content-Type"],
				Payload:     []byte(httpEvent.Body),
			})
		}

		break
	default:
		jsonEvent := make(map[string]interface{})

		err = json.Unmarshal(data, &jsonEvent)

		if err != nil {
			return err
		}

		err = fmt.Errorf("Unhandled Event Type: %v", data)
	}

	return err
}

type LambdaGateway struct {
	handler sdk.GatewayHandler
	runtime LambdaRuntimeHandler
	sdk.UnimplementedGatewayPlugin
}

func (s *LambdaGateway) handle(ctx context.Context, event Event) (interface{}, error) {
	for _, request := range event.Requests {
		// TODO: Build up an array of responses?
		//in some cases we won't need to send a response as well...
		resp := s.handler(&request)
		// There should only be one request if it was for a http request/response
		// So we'll just return here...
		if request.Context.SourceType == sdk.Request {
			return events.APIGatewayProxyResponse{
				StatusCode:      resp.Status,
				Headers:         resp.Headers,
				Body:            string(resp.Body),
				IsBase64Encoded: false,
			}, nil
		}
	}
	return nil, nil
}

// Start the lambda gateway handler
func (s *LambdaGateway) Start(handler sdk.GatewayHandler) error {
	s.handler = handler
	// Here we want to begin polling lambda for incoming requests...
	// Assuming that this is blocking
	s.runtime(s.handle)

	return fmt.Errorf("Something went wrong causing the lambda runtime to stop")
}

func New() (sdk.GatewayPlugin, error) {
	return &LambdaGateway{
		runtime: lambda.Start,
	}, nil
}

func NewWithRuntime(runtime LambdaRuntimeHandler) (sdk.GatewayPlugin, error) {
	return &LambdaGateway{
		runtime: runtime,
	}, nil
}
