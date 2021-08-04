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

// The AWS HTTP gateway plugin for ECS
package ecs_service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/utils"
	"github.com/nitric-dev/membrane/pkg/worker"

	"github.com/valyala/fasthttp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
)

const (
	AMZ_MESSAGE_ID   = "x-amz-sns-message-id"
	AMZ_MESSAGE_TYPE = "x-amz-sns-message-type"
	AMZ_TOPIC_ARN    = "x-amz-sns-topic-arn"
)

type HttpProxyGateway struct {
	client  *sns.SNS
	address string
	server  *fasthttp.Server
}

func (s *HttpProxyGateway) httpHandler(pool worker.WorkerPool) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		wrkr, err := pool.GetWorker()
		if err != nil {
			ctx.Error("Unable to get work to handle this event", 500)
		}

		var trigger = ctx.UserAgent()

		if string(trigger) == "Amazon Simple Notification Service Agent" {
			// If its a subscribe or unsubscribe notification then we need to handle it
			amzMessageType := string(ctx.Request.Header.Peek(AMZ_MESSAGE_TYPE))
			topicArn := string(ctx.Request.Header.Peek(AMZ_TOPIC_ARN))
			id := string(ctx.Request.Header.Peek(AMZ_MESSAGE_ID))

			payload := ctx.Request.Body()

			// SNS bodies are always JSON
			var jsonBody map[string]interface{}
			unmarshalError := json.Unmarshal(payload, &jsonBody)
			if unmarshalError != nil {
				ctx.Error("There was an error unmarshalling an SNS message", 403)
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
				ctx.SuccessString("text/plain", "success")
				return
			} else if amzMessageType == "UnsubscribeConfirmation" {
				// FIXME: Decide how we need to handle this
				ctx.SuccessString("text/plain", "success")
				return
			}

			if err := wrkr.HandleEvent(&triggers.Event{
				ID: id,
				// FIXME: Split this to retrive the nitric topic name
				Topic:   topicArn,
				Payload: payload,
			}); err == nil {
				ctx.SuccessString("text/plain", "success")
			} else {
				ctx.Error("Internal Server Error", 500)
			}

			return
		}

		// Otherwise treat as a normal http request
		response, err := wrkr.HandleHttpRequest(triggers.FromHttpRequest(ctx))

		if err != nil {
			ctx.Error(err.Error(), 500)
			return
		}

		response.Header.CopyTo(&ctx.Response.Header)
		ctx.Response.SetStatusCode(response.Header.StatusCode())
		ctx.Response.SetBody(response.Body)
	}
}

func (s *HttpProxyGateway) Start(pool worker.WorkerPool) error {
	// Start the fasthttp server
	s.server = &fasthttp.Server{
		IdleTimeout:     time.Second * 1,
		CloseOnShutdown: true,
		Handler:         s.httpHandler(pool),
	}

	return s.server.ListenAndServe(s.address)
}

func (s *HttpProxyGateway) Stop() error {
	return s.server.Shutdown()
}

// Create new AWS HTTP Gateway service
func New() (gateway.GatewayService, error) {
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
