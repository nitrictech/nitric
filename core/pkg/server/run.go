// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nitrictech/nitric/core/pkg/logger"
)

// Run - Run a runtime server until a signal is received or an error occurs
func Run(m *NitricServer) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	logger.SetLogLevel(logger.INFO)

	errChan := make(chan error)

	// Start the runtime server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case serverErr := <-errChan:
		logger.Errorf("Nitric Server Error: %v, exiting\n", serverErr)
	case sigTerm := <-term:
		logger.Infof("Received %v, exiting\n", sigTerm)
	}

	m.Stop()
}
