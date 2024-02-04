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
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/nitrictech/nitric/cloud/common/runtime/logger"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/api"
	cloudrun_plugin "github.com/nitrictech/nitric/cloud/gcp/runtime/gateway"
	firestore_service "github.com/nitrictech/nitric/cloud/gcp/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	secret_manager_secret_service "github.com/nitrictech/nitric/cloud/gcp/runtime/secret"
	storage_service "github.com/nitrictech/nitric/cloud/gcp/runtime/storage"
	pubsub_service "github.com/nitrictech/nitric/cloud/gcp/runtime/topic"
	"github.com/nitrictech/nitric/core/pkg/membrane"
)

func main() {
	// Setup signal interrupt handling for graceful shutdown
	var err error
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)

	membraneOpts := membrane.DefaultMembraneOptions()
	provider, err := resource.New()
	if err != nil {
		logger.Fatalf("Failed create core provider: %s", err.Error())
	}

	membraneOpts.ApiPlugin = api.NewGcpApiGatewayProvider(provider)

	membraneOpts.SecretManagerPlugin, err = secret_manager_secret_service.New()
	if err != nil {
		logger.Errorf("Failed to load secret plugin: %s", err.Error())
	}

	membraneOpts.KeyValuePlugin, err = firestore_service.New()
	if err != nil {
		logger.Errorf("Failed to load document plugin: %s", err.Error())
	}

	membraneOpts.TopicsPlugin, err = pubsub_service.New(provider)
	if err != nil {
		logger.Errorf("Failed to load events plugin: %s", err.Error())
	}

	membraneOpts.StoragePlugin, err = storage_service.New()
	if err != nil {
		logger.Errorf("Failed to load storage plugin: %s", err.Error())
	}

	membraneOpts.GatewayPlugin, err = cloudrun_plugin.New(provider)
	if err != nil {
		logger.Errorf("Failed to load gateway plugin: %s", err.Error())
	}

	membraneOpts.ResourcesPlugin = provider
	membraneOpts.CreateTracerProvider = newTraceProvider

	m, err := membrane.New(membraneOpts)
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
		log.Errorf(fmt.Sprintf("Membrane Error: %v, exiting", membraneError))
	case sigTerm := <-term:
		log.Errorf(fmt.Sprintf("Received %v, exiting", sigTerm))
	}

	m.Stop()
}
