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
	"github.com/nitric-dev/membrane/sdk"
	"log"

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/cognito"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/sns"
	httpGateway "github.com/nitric-dev/membrane/plugins/gateway/ecs"
	lambdaGateway "github.com/nitric-dev/membrane/plugins/gateway/lambda"
	documents "github.com/nitric-dev/membrane/plugins/kv/dynamodb"
	queue "github.com/nitric-dev/membrane/plugins/queue/sqs"
	storage "github.com/nitric-dev/membrane/plugins/storage/s3"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	gatewayEnv := utils.GetEnv("GATEWAY_ENVIRONMENT", "lambda")

	// Load the appropriate gateway based on the environment.
	var gatewayPlugin sdk.GatewayService
	switch gatewayEnv {
	case "lambda":
		gatewayPlugin, _ = lambdaGateway.New()
	default:
		gatewayPlugin, _ = httpGateway.New()
	}

	eventingPlugin, _ := eventing.New()
	keyValuePlugin, _ := documents.New()
	storagePlugin, _ := storage.New()
	queuePlugin, _ := queue.New()
	authPlugin, _ := auth.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		EventingPlugin:          eventingPlugin,
		KvPlugin:                keyValuePlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		GatewayPlugin:           gatewayPlugin,
		AuthPlugin:              authPlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	// Start the Membrane server
	m.Start()
}
