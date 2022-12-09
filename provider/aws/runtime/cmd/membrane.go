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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nitrictech/nitric/core/pkg/membrane"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/provider/aws/runtime/core"
	dynamodb_service "github.com/nitrictech/nitric/provider/aws/runtime/documents"
	sns_service "github.com/nitrictech/nitric/provider/aws/runtime/events"
	lambda_service "github.com/nitrictech/nitric/provider/aws/runtime/gateway"
	sqs_service "github.com/nitrictech/nitric/provider/aws/runtime/queue"
	secrets_manager_secret_service "github.com/nitrictech/nitric/provider/aws/runtime/secret"
	s3_service "github.com/nitrictech/nitric/provider/aws/runtime/storage"
	base_http "github.com/nitrictech/nitric/provider/common/runtime/gateway"
)

func main() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	gatewayEnv := utils.GetEnv("GATEWAY_ENVIRONMENT", "lambda")

	membraneOpts := membrane.DefaultMembraneOptions()

	provider, err := core.New()
	if err != nil {
		log.Fatalf("could not create aws provider: %v", err)
		return
	}

	// Load the appropriate gateway based on the environment.
	switch gatewayEnv {
	case "lambda":
		membraneOpts.GatewayPlugin, _ = lambda_service.New(provider)
	default:
		membraneOpts.GatewayPlugin, _ = base_http.New(nil)
	}

	membraneOpts.SecretPlugin, _ = secrets_manager_secret_service.New(provider)
	membraneOpts.DocumentPlugin, _ = dynamodb_service.New(provider)
	membraneOpts.EventsPlugin, _ = sns_service.New(provider)
	membraneOpts.QueuePlugin, _ = sqs_service.New(provider)
	membraneOpts.StoragePlugin, _ = s3_service.New(provider)
	membraneOpts.ResourcesPlugin = provider
	membraneOpts.CreateTracerProvider = newTracerProvider

	m, err := membrane.New(membraneOpts)
	if err != nil {
		log.Default().Fatalf("There was an error initialising the membrane server: %v", err)
	}

	errChan := make(chan error)
	// Start the Membrane server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case membraneError := <-errChan:
		log.Default().Printf("Membrane Error: %v, exiting\n", membraneError)
	case sigTerm := <-term:
		log.Default().Printf("Received %v, exiting\n", sigTerm)
	}

	m.Stop()
}
