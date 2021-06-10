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
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nitric-dev/membrane/utils"
	"github.com/nitric-dev/membrane/worker"

	grpc2 "github.com/nitric-dev/membrane/adapters/grpc"

	v1 "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/sdk"
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

	EventingPlugin sdk.EventService
	KvPlugin       sdk.KeyValueService
	StoragePlugin  sdk.StorageService
	QueuePlugin    sdk.QueueService
	GatewayPlugin  sdk.GatewayService
	AuthPlugin     sdk.UserService

	SuppressLogs            bool
	TolerateMissingServices bool

	// The operating mode of the membrane
	Mode *Mode
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
	eventPlugin   sdk.EventService
	kvPlugin      sdk.KeyValueService
	storagePlugin sdk.StorageService
	gatewayPlugin sdk.GatewayService
	queuePlugin   sdk.QueueService
	authPlugin    sdk.UserService

	// Tolerate if provider specific plugins aren't available for some services.
	// Not this does not include the gateway service
	tolerateMissingServices bool

	// Suppress println statements in the membrane server
	suppressLogs bool

	// Handler operating mode, e.g. FaaS or HTTP Proxy. Governs how incoming triggers are translated.
	mode Mode

	grpcServer *grpc.Server
}

func (s *Membrane) log(log string) {
	if !s.suppressLogs {
		fmt.Println(log)
	}
}

// Create a new Nitric Eventing Server
func (s *Membrane) createEventingServer() v1.EventServer {
	return grpc2.NewEventServer(s.eventPlugin)
}

func (s *Membrane) createTopicServer() v1.TopicServer {
	return grpc2.NewTopicServer(s.eventPlugin)
}

// Create a new Nitric Storage Server
func (s *Membrane) createStorageServer() v1.StorageServer {
	return grpc2.NewStorageServer(s.storagePlugin)
}

func (s *Membrane) createKeyValueServer() v1.KeyValueServer {
	return grpc2.NewKeyValueServer(s.kvPlugin)
}

func (s *Membrane) createQueueServer() v1.QueueServer {
	return grpc2.NewQueueServer(s.queuePlugin)
}

func (s *Membrane) createUserServer() v1.UserServer {
	return grpc2.NewUserServer(s.authPlugin)
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

	authServer := s.createUserServer()
	v1.RegisterUserServer(s.grpcServer, authServer)

	// Start with a maximum of a single worker
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{
		MaxWorkers: 1,
	})
	// FaaS server MUST start before the child process
	if s.mode == Mode_Faas {
		faasServer := grpc2.NewFaasServer(pool)
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

	if s.mode == Mode_HttpProxy {
		if httpWorker, err := worker.NewHttpWorker(s.childAddress); err != nil {
			return err
		} else {
			if err := pool.AddWorker(httpWorker); err != nil {
				return err
			}
		}
	}

	// FIXME: Only do this in Gateway mode...
	// Otherwise always pass through to the provided child address
	// Start the Gateway Server

	// Start the gateway, this will provide us an entrypoint for
	// data ingress/egress to our userland code
	// The gateway should block the main thread but will
	// use this callback as a control mechanism
	s.log("Waiting for active workers")
	err = pool.WaitForActiveWorkers(5)

	if err != nil {
		return err
	}
	s.log("Starting Gateway")
	return s.gatewayPlugin.Start(pool)
}

func (s *Membrane) Stop() {
	_ = s.gatewayPlugin.Stop()
	s.grpcServer.GracefulStop()
}

// Create a new Membrane server
func New(options *MembraneOptions) (*Membrane, error) {

	// Get unset options from env or defaults
	if options.ServiceAddress == "" {
		options.ServiceAddress = utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	}

	if options.ChildAddress == "" {
		options.ChildAddress = utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	}

	// Pull child command from command line args or environment variable if not provided.
	if len(options.ChildCommand) < 1 {
		// Get the command line arguments, minus the program name in index 0.
		if len(os.Args) > 1 && len(os.Args[1:]) > 0 {
			options.ChildCommand = os.Args[1:]
		} else {
			options.ChildCommand = strings.Fields(utils.GetEnv("INVOKE", ""))
			if len(options.ChildCommand) > 0 {
				fmt.Println("Warning: use of INVOKE environment variable is deprecated and may be removed in a future version")
			}
		}
	}

	if !options.TolerateMissingServices {
		tolerateMissing, err := strconv.ParseBool(utils.GetEnv("TOLERATE_MISSING_SERVICES", "false"))
		if err != nil {
			return nil, err
		}
		options.TolerateMissingServices = tolerateMissing
	}

	if options.Mode == nil {
		mode, err := ModeFromString(utils.GetEnv("MEMBRANE_MODE", "FAAS"))
		if err != nil {
			return nil, err
		}
		options.Mode = &mode
	}

	if options.ChildTimeoutSeconds < 1 {
		options.ChildTimeoutSeconds = 5
	}

	if options.GatewayPlugin == nil {
		return nil, fmt.Errorf("Missing gateway plugin, Gateway plugin must not be nil")
	}

	if !options.TolerateMissingServices {
		if options.EventingPlugin == nil || options.StoragePlugin == nil || options.KvPlugin == nil || options.QueuePlugin == nil || options.AuthPlugin == nil {
			return nil, fmt.Errorf("Missing membrane plugins, if you meant to load with missing plugins set options.TolerateMissingServices to true")
		}
	}

	return &Membrane{
		serviceAddress:          options.ServiceAddress,
		childAddress:            options.ChildAddress,
		childUrl:                fmt.Sprintf("http://%s", options.ChildAddress),
		childCommand:            options.ChildCommand,
		childTimeoutSeconds:     options.ChildTimeoutSeconds,
		authPlugin:              options.AuthPlugin,
		eventPlugin:             options.EventingPlugin,
		storagePlugin:           options.StoragePlugin,
		kvPlugin:                options.KvPlugin,
		queuePlugin:             options.QueuePlugin,
		gatewayPlugin:           options.GatewayPlugin,
		suppressLogs:            options.SuppressLogs,
		tolerateMissingServices: options.TolerateMissingServices,
		mode:                    *options.Mode,
	}, nil
}
