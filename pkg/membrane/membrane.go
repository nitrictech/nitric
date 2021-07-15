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

package membrane

import (
	"fmt"
	grpc3 "github.com/nitric-dev/membrane/pkg/adapters/grpc"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	worker2 "github.com/nitric-dev/membrane/pkg/worker"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	v1 "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/pkg/sdk"
	"google.golang.org/grpc"
)

type MembraneOptions struct {
	ServiceAddress string
	// The address the child will be listening on
	ChildAddress string
	// The command that will be used to invoke the child process
	ChildCommand []string
	// The total time to wait for the child process to be available in seconds
	ChildTimeoutSeconds int

	DocumentPlugin sdk.DocumentService
	EventingPlugin sdk.EventService
	KvPlugin       sdk.KeyValueService
	StoragePlugin  sdk.StorageService
	QueuePlugin    sdk.QueueService
	GatewayPlugin  sdk.GatewayService

	SuppressLogs            bool
	TolerateMissingServices bool

	// The operating mode of the membrane
	Mode *Mode

	// Supply your own worker pool
	Pool worker2.WorkerPool
}

type Membrane struct {
	// Address & port to bind the membrane i/o proxy to
	// This will still be bound even in pass through mode
	// proxyAddress string
	// Address & port to bind the membrane service interfaces to
	serviceAddress string
	// The address the child will be listening on
	childAddress string

	// The URL (including protocol, the child process can be reached on)
	childUrl string

	// The command that will be used to invoke the child process
	childCommand []string

	childTimeoutSeconds int

	// Configured plugins
	documentPlugin sdk.DocumentService
	eventPlugin    sdk.EventService
	kvPlugin       sdk.KeyValueService
	storagePlugin  sdk.StorageService
	gatewayPlugin  sdk.GatewayService
	queuePlugin    sdk.QueueService

	// Tolerate if provider specific plugins aren't available for some services.
	// Not this does not include the gateway service
	tolerateMissingServices bool

	// Suppress println statements in the membrane server
	suppressLogs bool

	// Handler operating mode, e.g. FaaS or HTTP Proxy. Governs how incoming triggers are translated.
	mode Mode

	grpcServer *grpc.Server

	// Worker pool
	pool worker2.WorkerPool
}

func (s *Membrane) log(log string) {
	if !s.suppressLogs {
		fmt.Println(log)
	}
}

// Create a new Nitric Document Server
func (s *Membrane) createDocumentServer() v1.DocumentServiceServer {
	return grpc3.NewDocumentServer(s.documentPlugin)
}

// Create a new Nitric Eventing Server
func (s *Membrane) createEventingServer() v1.EventServer {
	return grpc3.NewEventServer(s.eventPlugin)
}

// Create a new Nitric Topic Server
func (s *Membrane) createTopicServer() v1.TopicServer {
	return grpc3.NewTopicServer(s.eventPlugin)
}

// Create a new Nitric Storage Server
func (s *Membrane) createStorageServer() v1.StorageServer {
	return grpc3.NewStorageServer(s.storagePlugin)
}

func (s *Membrane) createKeyValueServer() v1.KeyValueServer {
	return grpc3.NewKeyValueServer(s.kvPlugin)
}

func (s *Membrane) createQueueServer() v1.QueueServer {
	return grpc3.NewQueueServer(s.queuePlugin)
}

func (s *Membrane) startChildProcess() error {
	// TODO: This is a detached process
	// so it will continue to run until even after the membrane dies

	fmt.Println(fmt.Sprintf("Starting Child Process"))
	childProcess := exec.Command(s.childCommand[0], s.childCommand[1:]...)
	childProcess.Stdout = os.Stdout
	childProcess.Stderr = os.Stderr
	applicationError := childProcess.Start()

	// Actual panic here, we don't want to start if our userland code cannot successfully start
	if applicationError != nil {
		return fmt.Errorf("There was an error starting the child process: %v", applicationError)
	}

	return nil
}

func (s *Membrane) nitricResponseFromError(err error) *sdk.NitricResponse {
	return &sdk.NitricResponse{
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    []byte(err.Error()),
		Status:  503,
	}
}

// Start the membrane
func (s *Membrane) Start() error {
	// Search for known plugins

	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)

	// Load & Register the GRPC service plugins
	documentServer := s.createDocumentServer()
	v1.RegisterDocumentServiceServer(s.grpcServer, documentServer)

	eventingServer := s.createEventingServer()
	v1.RegisterEventServer(s.grpcServer, eventingServer)

	topicServer := s.createTopicServer()
	v1.RegisterTopicServer(s.grpcServer, topicServer)

	kvServer := s.createKeyValueServer()
	v1.RegisterKeyValueServer(s.grpcServer, kvServer)

	storageServer := s.createStorageServer()
	v1.RegisterStorageServer(s.grpcServer, storageServer)

	queueServer := s.createQueueServer()
	v1.RegisterQueueServer(s.grpcServer, queueServer)

	// FaaS server MUST start before the child process
	if s.mode == Mode_Faas {
		faasServer := grpc3.NewFaasServer(s.pool)
		v1.RegisterFaasServer(s.grpcServer, faasServer)
	}
	lis, err := net.Listen("tcp", s.serviceAddress)
	if err != nil {
		return fmt.Errorf("Could not listen on configured service address: %v", err)
	}

	s.log("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		s.log(fmt.Sprintf("Services listening on: %s", s.serviceAddress))
		s.grpcServer.Serve(lis)
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if len(s.childCommand) > 0 {
		if err := s.startChildProcess(); err != nil {
			// Return the error
			return err
		}
	} else {
		s.log("No Child Command Specified, Skipping...")
	}

	// If we aren't in FaaS mode
	// We need to manually register our worker for now
	if s.mode != Mode_Faas {
		var wrkr worker2.Worker
		var workerErr error
		if s.mode == Mode_HttpProxy {
			wrkr, workerErr = worker2.NewHttpWorker(s.childAddress)
		} else if s.mode == Mode_HttpFaas {
			wrkr, workerErr = worker2.NewFaasHttpWorker(s.childAddress)
		}

		if workerErr == nil {
			if err := s.pool.AddWorker(wrkr); err != nil {
				return err
			}
		} else {
			return workerErr
		}
	}

	// Wait for active workers to be available before begining the gateway
	// This will ensure requests can be facilitated as soon as the gateway is ready
	s.log("Waiting for active workers")
	err = s.pool.WaitForActiveWorkers(s.childTimeoutSeconds)

	if err != nil {
		return err
	}

	gatewayErrchan := make(chan error)
	poolErrchan := make(chan error)

	// Start the gateway
	go func(errch chan error) {
		s.log("Starting Gateway")
		errch <- s.gatewayPlugin.Start(s.pool)
	}(gatewayErrchan)

	// Start the worker pool monitor
	go func(errch chan error) {
		s.log("Starting Worker Supervisor")
		errch <- s.pool.Monitor()
	}(poolErrchan)

	var exitErr error

	// Wait and fail on either
	select {
	case gatewayErr := <-gatewayErrchan:
		if err == nil {
			// Normal Gateway shutdown
			// Allowing the membrane to exit
			return nil
		}
		exitErr = fmt.Errorf(fmt.Sprintf("Gateway Error: %v, exiting", gatewayErr))
	case poolErr := <-poolErrchan:
		exitErr = fmt.Errorf(fmt.Sprintf("Supervisor error: %v, exiting", poolErr))
	}

	return exitErr
}

func (s *Membrane) Stop() {
	s.gatewayPlugin.Stop()
	s.grpcServer.Stop()
}

// Create a new Membrane server
func New(options *MembraneOptions) (*Membrane, error) {

	// Get unset options from env or defaults
	if options.ServiceAddress == "" {
		options.ServiceAddress = utils2.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	}

	if options.ChildAddress == "" {
		options.ChildAddress = utils2.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	}

	// Pull child command from command line args or environment variable if not provided.
	if len(options.ChildCommand) < 1 {
		// Get the command line arguments, minus the program name in index 0.
		if len(os.Args) > 1 && len(os.Args[1:]) > 0 {
			options.ChildCommand = os.Args[1:]
		} else {
			options.ChildCommand = strings.Fields(utils2.GetEnv("INVOKE", ""))
			if len(options.ChildCommand) > 0 {
				fmt.Println("Warning: use of INVOKE environment variable is deprecated and may be removed in a future version")
			}
		}
	}

	if !options.TolerateMissingServices {
		tolerateMissing, err := strconv.ParseBool(utils2.GetEnv("TOLERATE_MISSING_SERVICES", "false"))
		if err != nil {
			return nil, err
		}
		options.TolerateMissingServices = tolerateMissing
	}

	if options.Mode == nil {
		mode, err := ModeFromString(utils2.GetEnv("MEMBRANE_MODE", "FAAS"))
		if err != nil {
			return nil, err
		}
		options.Mode = &mode
	}

	if options.ChildTimeoutSeconds < 1 {
		options.ChildTimeoutSeconds = 10
	}

	if options.GatewayPlugin == nil {
		return nil, fmt.Errorf("Missing gateway plugin, Gateway plugin must not be nil")
	}

	if !options.TolerateMissingServices {
		if options.EventingPlugin == nil || options.StoragePlugin == nil || options.KvPlugin == nil || options.QueuePlugin == nil {
			return nil, fmt.Errorf("Missing membrane plugins, if you meant to load with missing plugins set options.TolerateMissingServices to true")
		}
	}

	if options.Pool == nil {
		// Create new pool with defaults
		options.Pool = worker2.NewProcessPool(&worker2.ProcessPoolOptions{
			MaxWorkers: 1,
		})
	}

	return &Membrane{
		serviceAddress:          options.ServiceAddress,
		childAddress:            options.ChildAddress,
		childUrl:                fmt.Sprintf("http://%s", options.ChildAddress),
		childCommand:            options.ChildCommand,
		childTimeoutSeconds:     options.ChildTimeoutSeconds,
		documentPlugin:          options.DocumentPlugin,
		eventPlugin:             options.EventingPlugin,
		storagePlugin:           options.StoragePlugin,
		kvPlugin:                options.KvPlugin,
		queuePlugin:             options.QueuePlugin,
		gatewayPlugin:           options.GatewayPlugin,
		suppressLogs:            options.SuppressLogs,
		tolerateMissingServices: options.TolerateMissingServices,
		mode:                    *options.Mode,
		pool:                    options.Pool,
	}, nil
}
