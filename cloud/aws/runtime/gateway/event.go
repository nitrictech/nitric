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

package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type eventType int

const (
	unknown eventType = iota
	sns
	s3
	httpEvent
	websocketEvent
	healthcheck
	// cloudwatch
	schedule
	xforwardHeader string = "x-forwarded-for"
)

type healthCheckEvent struct {
	Check bool `json:"x-nitric-healthcheck,omitempty"`
}

// An event struct that embeds the AWS event types that we handle
type Record struct {
	EventSource      string
	EventSourceArn   string
	EventName        string
	ResponseElements map[string]string
	S3               events.S3Entity
	SNS              events.SNSEntity
}

type nitricScheduleEvent struct {
	Schedule string `json:"x-nitric-schedule,omitempty"`
}

// An event struct that embeds the AWS event types that we handle
type Event struct {
	events.APIGatewayV2HTTPRequest
	events.APIGatewayWebsocketProxyRequest
	healthCheckEvent
	Records []Record
	nitricScheduleEvent
}

func (e *Event) Type() eventType {
	// check if this event type contains valid data
	if e.APIGatewayWebsocketProxyRequest.RequestContext.ConnectionID != "" {
		return websocketEvent
	} else if e.APIGatewayV2HTTPRequest.RouteKey != "" || e.APIGatewayV2HTTPRequest.RequestContext.APIID != "" {
		return httpEvent
	} else if e.Check {
		return healthcheck
	} else if len(e.Records) > 0 && e.Records[0].EventSource == "aws:sns" {
		return sns
	} else if len(e.Records) > 0 && e.Records[0].EventSource == "aws:s3" {
		return s3
	} else if e.Schedule != "" {
		return schedule
	}

	return unknown
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var err error

	switch e.getEventType(data) {
	case s3:
		s3Event := &events.S3Event{}
		err = json.Unmarshal(data, s3Event)
		if err != nil {
			return err
		}

		e.Records = make([]Record, 0)

		for _, s3Record := range s3Event.Records {
			e.Records = append(e.Records, Record{
				EventSource:      s3Record.EventSource,
				EventSourceArn:   s3Record.S3.Bucket.Arn,
				EventName:        s3Record.EventName,
				ResponseElements: s3Record.ResponseElements,
				S3:               s3Record.S3,
			})
		}

		return nil

	case sns:
		snsEvent := &events.SNSEvent{}
		err = json.Unmarshal(data, snsEvent)
		if err != nil {
			return err
		}

		e.Records = make([]Record, 0)

		for _, snsRecord := range snsEvent.Records {
			e.Records = append(e.Records, Record{
				EventSource:      snsRecord.EventSource,
				EventSourceArn:   snsRecord.SNS.TopicArn,
				EventName:        snsRecord.SNS.Type,
				ResponseElements: map[string]string{},
				SNS:              snsRecord.SNS,
			})
		}

	case httpEvent:
		apiEvent := events.APIGatewayV2HTTPRequest{}
		err = json.Unmarshal(data, &apiEvent)
		if err != nil {
			return err
		}

		e.APIGatewayV2HTTPRequest = apiEvent
	case websocketEvent:
		websocketEvent := events.APIGatewayWebsocketProxyRequest{}
		err = json.Unmarshal(data, &websocketEvent)
		if err != nil {
			return err
		}

		e.APIGatewayWebsocketProxyRequest = websocketEvent
	case schedule:
		nitricSchedule := nitricScheduleEvent{}
		err = json.Unmarshal(data, &nitricSchedule)
		if err != nil {
			return err
		}

		e.nitricScheduleEvent = nitricSchedule
	case healthcheck:
		checkEvent := healthCheckEvent{}
		err = json.Unmarshal(data, &checkEvent)
		if err != nil {
			return err
		}

		e.healthCheckEvent = checkEvent
	default:
		return fmt.Errorf("unhandled lambda event type")
	}

	return nil
}

func (e *Event) getEventType(data []byte) eventType {
	temp := make(map[string]interface{})
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return unknown
	}

	requestContext, isRequest := temp["requestContext"].(map[string]interface{})

	// Handle non-record events
	if isRequest {
		if _, ok := requestContext["connectionId"]; ok {
			return websocketEvent
		}
		return httpEvent
	} else if _, ok := temp["x-nitric-healthcheck"]; ok {
		return healthcheck
	} else if _, ok := temp["x-nitric-schedule"]; ok {
		return schedule
	}

	// Handle Events
	recordsList, _ := temp["Records"].([]interface{})
	if len(recordsList) == 0 {
		return unknown
	}

	record, _ := recordsList[0].(map[string]interface{})

	var eventSource string

	if es, ok := record["EventSource"]; ok {
		eventSource = es.(string)
	} else if es, ok := record["eventSource"]; ok {
		eventSource = es.(string)
	}

	switch eventSource {
	case "aws:s3":
		return s3
	case "aws:sns":
		return sns
	}

	return unknown
}
