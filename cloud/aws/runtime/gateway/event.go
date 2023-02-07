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
	"github.com/aws/aws-lambda-go/events"
)

type eventType int

const (
	unknown eventType = iota
	sns
	httpEvent
	healthcheck
	cloudwatch
	xforwardHeader string = "x-forwarded-for"
)

type healthCheckEvent struct {
	Check bool `json:"x-nitric-healthcheck,omitempty"`
}

// An event struct that embeds the AWS event types that we handle
type Event struct {
	events.APIGatewayV2HTTPRequest `json:",omitempty"`
	// The base serializable type here matches other embeddable repeated constructs
	events.SNSEvent        `json:",omitempty"`
	events.CloudWatchEvent `json:",omitempty"`
	healthCheckEvent
}

func (e *Event) Type() eventType {
	// check if this event type contains valid data
	if e.APIGatewayV2HTTPRequest.RouteKey != "" {
		return httpEvent
	} else if e.Check {
		return healthcheck
	} else if len(e.Records) > 0 && e.Records[0].EventSource == "aws:sns" {
		return sns
	} else if e.Source != "" {
		return cloudwatch
	}

	return unknown
}
