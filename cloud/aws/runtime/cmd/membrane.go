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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	lambda_service "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	dynamodb_service "github.com/nitrictech/nitric/cloud/aws/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	secrets_manager_secret_service "github.com/nitrictech/nitric/cloud/aws/runtime/secret"
	s3_service "github.com/nitrictech/nitric/cloud/aws/runtime/storage"
	sns_service "github.com/nitrictech/nitric/cloud/aws/runtime/topic"
	"github.com/nitrictech/nitric/cloud/aws/runtime/websocket"
	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	"github.com/nitrictech/nitric/core/pkg/membrane"
)

func main() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	gatewayEnv := env.GATEWAY_ENVIRONMENT.String()

	membraneOpts := membrane.DefaultMembraneOptions()

	resolver, err := resource.New()
	if err != nil {
		logger.Fatalf("could not create aws provider: %v", err)
		return
	}

	// Load the appropriate gateway based on the environment.
	switch gatewayEnv {
	case "lambda":
		membraneOpts.GatewayPlugin, _ = lambda_service.New(resolver)
	default:
		membraneOpts.GatewayPlugin, _ = base_http.NewHttpGateway(nil)
	}

	membraneOpts.SecretManagerPlugin, _ = secrets_manager_secret_service.New(resolver)
	membraneOpts.KeyValuePlugin, _ = dynamodb_service.New(resolver)
	membraneOpts.TopicsPlugin, _ = sns_service.New(resolver)
	membraneOpts.StoragePlugin, _ = s3_service.New(resolver)
	membraneOpts.ResourcesPlugin = resolver
	membraneOpts.WebsocketPlugin, _ = websocket.NewAwsApiGatewayWebsocket(resolver)

	m, err := membrane.New(membraneOpts)
	if err != nil {
		logger.Fatalf("There was an error initializing the membrane server: %v", err)
	}

	errChan := make(chan error)
	// Start the Membrane server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case membraneError := <-errChan:
		logger.Errorf("Membrane Error: %v, exiting\n", membraneError)
	case sigTerm := <-term:
		logger.Debugf("Received %v, exiting\n", sigTerm)
	}

	m.Stop()
}
