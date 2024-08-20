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
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/core/pkg/gateway"
)

type LambdaRuntimeHandler func(interface{})

type LambdaGateway struct {
	resolver   resource.AwsResourceResolver
	runtime    LambdaRuntimeHandler
	routeEvent LambdaEventRouter
	gateway.UnimplementedGatewayPlugin
	finished chan int
}

var _ gateway.GatewayService = &LambdaGateway{}

// Start - poll the lambda runtime for events and route the to handlers for processing
func (s *LambdaGateway) Start(opts *gateway.GatewayStartOpts) error {
	handlers := &Handlers{
		Apis:               opts.ApiPlugin,
		Https:              opts.HttpPlugin,
		Schedules:          opts.SchedulesPlugin,
		Subscriptions:      opts.TopicsListenerPlugin,
		StorageListeners:   opts.StorageListenerPlugin,
		WebsocketListeners: opts.WebsocketListenerPlugin,
	}

	// Begin polling lambda for incoming requests...
	s.runtime(func(ctx context.Context, evt json.RawMessage) (interface{}, error) {
		return s.routeEvent(ctx, s.resolver, handlers, evt)
	})
	// Unblock the 'Stop' function if it's waiting.
	go func() { s.finished <- 1 }()
	return nil
}

// Stop - block until the lambda runtime is finished
func (s *LambdaGateway) Stop() error {
	// This is a NO_OP Process, as this is a pull based system
	// We don't need to stop listening to anything
	log.Default().Println("gateway 'Stop' called, waiting for lambda runtime to finish")
	// IT CANNOT BE STOPPED!!! Lambda is done when it wants to be and you won't change its mind.
	// But seriously we set the with SIGTERM option in Start for automatic graceful shutdown
	<-s.finished
	return nil
}

// New - Create a new LambdaGateway
func New(resolver resource.AwsResourceResolver, opts ...lambdaGatewayOption) *LambdaGateway {
	lg := &LambdaGateway{
		resolver:   resolver,
		finished:   make(chan int),
		runtime:    lambda.Start,
		routeEvent: StandardEventRouter,
	}

	for _, opt := range opts {
		opt(lg)
	}

	return lg
}
