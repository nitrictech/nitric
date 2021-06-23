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

	"github.com/nitric-dev/membrane/membrane"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/app_platform"
)

func main() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	gatewayPlugin, _ := gateway.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		GatewayPlugin: gatewayPlugin,
		// FIXME: Hardcode as true as we don't have plugins for other services for digital ocean yet...
		TolerateMissingServices: true,
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
		fmt.Println(fmt.Sprintf("Recieved %v, exiting", sigTerm))
	}

	m.Stop()
}
