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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nitric-dev/membrane/pkg/membrane"
	dynamodb_service "github.com/nitric-dev/membrane/pkg/plugins/document/dynamodb"
	sns_service "github.com/nitric-dev/membrane/pkg/plugins/events/sns"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	ecs_service "github.com/nitric-dev/membrane/pkg/plugins/gateway/ecs"
	lambda_service "github.com/nitric-dev/membrane/pkg/plugins/gateway/lambda"
	sqs_service "github.com/nitric-dev/membrane/pkg/plugins/queue/sqs"
	s3_service "github.com/nitric-dev/membrane/pkg/plugins/storage/s3"
	"github.com/nitric-dev/membrane/pkg/utils"
)

func main() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	gatewayEnv := utils.GetEnv("GATEWAY_ENVIRONMENT", "lambda")

	// Load the appropriate gateway based on the environment.
	var gatewayPlugin gateway.GatewayService
	switch gatewayEnv {
	case "lambda":
		gatewayPlugin, _ = lambda_service.New()
	default:
		gatewayPlugin, _ = ecs_service.New()
	}
	documentPlugin, _ := dynamodb_service.New()
	eventsPlugin, _ := sns_service.New()
	queuePlugin, _ := sqs_service.New()
	storagePlugin, _ := s3_service.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		DocumentPlugin: documentPlugin,
		EventsPlugin:   eventsPlugin,
		GatewayPlugin:  gatewayPlugin,
		QueuePlugin:    queuePlugin,
		StoragePlugin:  storagePlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	errChan := make(chan error)
	// Start the Membrane server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case membraneError := <-errChan:
		fmt.Println(fmt.Sprintf("Membrane Error: %v, exiting", membraneError))
	case sigTerm := <-term:
		fmt.Println(fmt.Sprintf("Received %v, exiting", sigTerm))
	}

	m.Stop()
}
