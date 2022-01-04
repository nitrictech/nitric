// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lambda_service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ep "github.com/nitrictech/nitric/pkg/plugins/events"
	"github.com/nitrictech/nitric/pkg/plugins/gateway"
)

type eventType int

const (
	unknown eventType = iota
	sns
	httpEvent
	xforwardHeader string = "x-forwarded-for"
)

type LambdaRuntimeHandler func(handler interface{})

//Event incoming event
type Event struct {
	Requests []triggers.Trigger
}

func (event *Event) getEventType(data []byte) eventType {
	tmp := make(map[string]interface{})
	// Unmarshal so we can get just enough info about the type of event to fully deserialize it
	json.Unmarshal(data, &tmp)

	// If our event is a HTTP request
	if _, ok := tmp["rawPath"]; ok {
		return httpEvent
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

	event.Requests = make([]triggers.Trigger, 0)

	switch event.getEventType(data) {
	case sns:
		snsEvent := &events.SNSEvent{}
		err = json.Unmarshal(data, snsEvent)

		if err == nil {
			// Map over the records and return
			for _, snsRecord := range snsEvent.Records {
				messageString := snsRecord.SNS.Message
				// FIXME: What about non-nitric SNS events???
				messageJson := &ep.NitricEvent{}

				// Populate the JSON
				err = json.Unmarshal([]byte(messageString), messageJson)

				topicArn := snsRecord.SNS.TopicArn
				topicParts := strings.Split(topicArn, ":")
				trigger := topicParts[len(topicParts)-1]
				// get the topic name from the full ARN.
				// Get the topic name from the arn

				if err == nil {
					// Decode the message to see if it's a Nitric message
					payloadMap := messageJson.Payload
					payloadBytes, err := json.Marshal(&payloadMap)

					if err == nil {
						event.Requests = append(event.Requests, &triggers.Event{
							ID:      messageJson.ID,
							Topic:   trigger,
							Payload: payloadBytes,
						})
					}
				}
			}
		}
		break
	case httpEvent:
		evt := &events.APIGatewayV2HTTPRequest{}
		err = json.Unmarshal(data, evt)

		if err == nil {
			// Copy the headers and re-write for the proxy
			headerCopy := make(map[string][]string)

			for key, val := range evt.Headers {
				if strings.ToLower(key) == "host" {
					headerCopy[xforwardHeader] = append(headerCopy[xforwardHeader], string(val))
				} else {
					headerCopy[key] = append(headerCopy[key], string(val))
				}
			}

			// Copy the cookies over
			headerCopy["Cookie"] = evt.Cookies

			// Parse the raw query string
			qVals, err := url.ParseQuery(evt.RawQueryString)

			if err == nil {
				event.Requests = append(event.Requests, &triggers.HttpRequest{
					// FIXME: Translate to http.Header
					Header: headerCopy,
					Body:   []byte(evt.Body),
					Method: evt.RequestContext.HTTP.Method,
					Path:   evt.RawPath,
					Query:  qVals,
				})
			}
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
	pool    worker.WorkerPool
	runtime LambdaRuntimeHandler
	gateway.UnimplementedGatewayPlugin
	finished chan int
}

func (s *LambdaGateway) handle(ctx context.Context, event Event) (interface{}, error) {
	wrkr, err := s.pool.GetWorker()

	if err != nil {
		return nil, fmt.Errorf("Unable to get worker to handle events")
	}

	for _, request := range event.Requests {
		switch request.GetTriggerType() {
		case triggers.TriggerType_Request:
			if httpEvent, ok := request.(*triggers.HttpRequest); ok {
				response, err := wrkr.HandleHttpRequest(httpEvent)

				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: 500,
						Body:       "Error processing lambda request",
						// TODO: Need to determine best case when to use this...
						IsBase64Encoded: true,
					}, nil
				}

				lambdaHTTPHeaders := make(map[string]string)

				if response.Header != nil {
					response.Header.VisitAll(func(key []byte, val []byte) {
						lambdaHTTPHeaders[string(key)] = string(val)
					})
				}

				responseString := base64.StdEncoding.EncodeToString(response.Body)

				// We want to sniff the content type of the body that we have here as lambda cannot gzip it...
				return events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Headers:    lambdaHTTPHeaders,
					Body:       responseString,
					// TODO: Need to determine best case when to use this...
					IsBase64Encoded: true,
				}, nil
			} else {
				return nil, fmt.Errorf("Error!: Found non HttpRequest in event with trigger type: %s", triggers.TriggerType_Request.String())
			}
			break
		case triggers.TriggerType_Subscription:
			if event, ok := request.(*triggers.Event); ok {
				if err := wrkr.HandleEvent(event); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("Error!: Found non Event in event with trigger type: %s", triggers.TriggerType_Subscription.String())
			}
			break
		}
	}
	return nil, nil
}

// Start the lambda gateway handler
func (s *LambdaGateway) Start(pool worker.WorkerPool) error {
	//s.finished = make(chan int)
	s.pool = pool
	// Here we want to begin polling lambda for incoming requests...
	s.runtime(s.handle)
	// Unblock the 'Stop' function if it's waiting.
	go func() { s.finished <- 1 }()
	return nil
}

func (s *LambdaGateway) Stop() error {
	// XXX: This is a NO_OP Process, as this is a pull based system
	// We don't need to stop listening to anything
	fmt.Println("gateway 'Stop' called, waiting for lambda runtime to finish")
	// Lambda can't be stopped, need to wait for it to finish
	<-s.finished
	return nil
}

func New() (gateway.GatewayService, error) {
	return &LambdaGateway{
		runtime:  lambda.Start,
		finished: make(chan int),
	}, nil
}

func NewWithRuntime(runtime LambdaRuntimeHandler) (gateway.GatewayService, error) {
	return &LambdaGateway{
		runtime:  runtime,
		finished: make(chan int),
	}, nil
}
