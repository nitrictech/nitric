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
