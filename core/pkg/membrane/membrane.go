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
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"

	grpc2 "github.com/nitrictech/nitric/core/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/plugins/queue"
	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
	"github.com/nitrictech/nitric/core/pkg/plugins/secret"
	"github.com/nitrictech/nitric/core/pkg/plugins/storage"
	"github.com/nitrictech/nitric/core/pkg/pm"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/core/pkg/worker"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

type MembraneOptions struct {
	ServiceAddress string
	// The address the child will be listening on
	ChildAddress string
	// The command that will be used to invoke the child process
	ChildCommand []string
	// Commands that will be started before all others
	PreCommands [][]string

	// The total time to wait for the child process to be available in seconds
	ChildTimeoutSeconds int

	DocumentPlugin  document.DocumentService
	EventsPlugin    events.EventService
	StoragePlugin   storage.StorageService
	QueuePlugin     queue.QueueService
	GatewayPlugin   gateway.GatewayService
	SecretPlugin    secret.SecretService
	ResourcesPlugin resource.ResourceService

	CreateTracerProvider func(ctx context.Context) (*sdktrace.TracerProvider, error)

	SuppressLogs            bool
	TolerateMissingServices bool

	// Supply your own worker pool
	Pool pool.WorkerPool
}

type Membrane struct {
	// Address & port to bind the membrane service interfaces to
	serviceAddress string
	// The address the child will be listening on
	childAddress string

	// The URL (including protocol, the child process can be reached on)
	childUrl string

	processManager       pm.ProcessManager
	tracerProvider       *sdktrace.TracerProvider
	createTracerProvider func(ctx context.Context) (*sdktrace.TracerProvider, error)

	childTimeoutSeconds int

	// Configured plugins
	documentPlugin document.DocumentService
	eventsPlugin   events.EventService
	storagePlugin  storage.StorageService
	gatewayPlugin  gateway.GatewayService
	queuePlugin    queue.QueueService
	secretPlugin   secret.SecretService
	resourcePlugin resource.ResourceService

	// Tolerate if provider specific plugins aren't available for some services.
	// Not this does not include the gateway service
	tolerateMissingServices bool

	// Suppress println statements in the membrane server
	suppressLogs bool

	grpcServer *grpc.Server

	// Worker pool
	pool pool.WorkerPool
}

func (s *Membrane) log(msg string) {
	if !s.suppressLogs {
		log.Default().Println(msg)
	}
}

func (s *Membrane) createSecretServer() v1.SecretServiceServer {
	return grpc2.NewSecretServer(s.secretPlugin)
}

// Create a new Nitric Document Server
func (s *Membrane) createDocumentServer() v1.DocumentServiceServer {
	return grpc2.NewDocumentServer(s.documentPlugin)
}

// Create a new Nitric events Server
func (s *Membrane) createEventsServer() v1.EventServiceServer {
	return grpc2.NewEventServiceServer(s.eventsPlugin)
}

// Create a new Nitric Topic Server
func (s *Membrane) createTopicServer() v1.TopicServiceServer {
	return grpc2.NewTopicServiceServer(s.eventsPlugin)
}

// Create a new Nitric Storage Server
func (s *Membrane) createStorageServer() v1.StorageServiceServer {
	return grpc2.NewStorageServiceServer(s.storagePlugin)
}

func (s *Membrane) createQueueServer() v1.QueueServiceServer {
	return grpc2.NewQueueServiceServer(s.queuePlugin)
}

// Start the membrane
func (s *Membrane) Start() error {
	if err := s.processManager.StartPreProcesses(); err != nil {
		return err
	}

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(uint32(s.pool.GetMaxWorkers())),
	}

	if s.createTracerProvider != nil {
		tp, err := s.createTracerProvider(context.Background())
		if err != nil {
			s.log(fmt.Sprintf("traceProvider %v", err))
			return err
		}

		if tp != nil {
			s.log(fmt.Sprintf("traceProvider connected"))
			otel.SetTracerProvider(tp)
		}

		interceptorOpts := []otelgrpc.Option{
			otelgrpc.WithPropagators(propagation.TraceContext{}),
		}

		opts = append(opts,
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(interceptorOpts...)),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(interceptorOpts...)),
		)
	}

	s.grpcServer = grpc.NewServer(opts...)

	// Load & Register the GRPC service plugins
	documentServer := s.createDocumentServer()
	v1.RegisterDocumentServiceServer(s.grpcServer, documentServer)

	eventsServer := s.createEventsServer()
	v1.RegisterEventServiceServer(s.grpcServer, eventsServer)

	topicServer := s.createTopicServer()
	v1.RegisterTopicServiceServer(s.grpcServer, topicServer)

	storageServer := s.createStorageServer()
	v1.RegisterStorageServiceServer(s.grpcServer, storageServer)

	queueServer := s.createQueueServer()
	v1.RegisterQueueServiceServer(s.grpcServer, queueServer)

	secretServer := s.createSecretServer()
	v1.RegisterSecretServiceServer(s.grpcServer, secretServer)

	resourceServer := grpc2.NewResourcesServiceServer(grpc2.WithResourcePlugin(s.resourcePlugin))
	v1.RegisterResourceServiceServer(s.grpcServer, resourceServer)

	// FaaS server MUST start before the child process
	faasServer := grpc2.NewFaasServer(s.pool)
	v1.RegisterFaasServiceServer(s.grpcServer, faasServer)

	lis, err := net.Listen("tcp", s.serviceAddress)
	if err != nil {
		return fmt.Errorf("could not listen on configured service address: %w", err)
	}

	s.log("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		s.log(fmt.Sprintf("Services listening on: %s", s.serviceAddress))
		err := s.grpcServer.Serve(lis)
		if err != nil {
			s.log(fmt.Sprintf("grpc serve %v", err))
		}
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if err := s.processManager.StartUserProcess(); err != nil {
		return err
	}

	// Wait for the minimum number of active workers to be available before beginning the gateway
	// This ensures workers have registered and can handle triggers as soon the gateway is ready, if a minimum > 1 has been set
	s.log("Waiting for active workers")
	err = s.pool.WaitForMinimumWorkers(s.childTimeoutSeconds)
	if err != nil {
		return err
	}

	gatewayErrchan := make(chan error)
	poolErrchan := make(chan error)

	// Start the gateway
	go func(errch chan error) {
		s.log(fmt.Sprintf("Starting Gateway, %d workers currently available", s.pool.GetWorkerCount()))
		errch <- s.gatewayPlugin.Start(s.pool)
	}(gatewayErrchan)

	// Start the worker pool monitor
	go func(errch chan error) {
		s.log("Starting Worker Supervisor")
		errch <- s.pool.Monitor()
	}(poolErrchan)

	processErrchan := make(chan error)
	go func(errch chan error) {
		errch <- s.processManager.Monitor()
	}(processErrchan)

	var exitErr error

	// Wait and fail on either
	select {
	case gatewayErr := <-gatewayErrchan:
		if gatewayErr == nil {
			// Normal Gateway shutdown
			// Allowing the membrane to exit
			return nil
		}
		exitErr = fmt.Errorf(fmt.Sprintf("Gateway Error: %v, exiting", gatewayErr))
	case poolErr := <-poolErrchan:
		exitErr = fmt.Errorf(fmt.Sprintf("Supervisor error: %v, exiting", poolErr))
	case processErr := <-processErrchan:
		exitErr = fmt.Errorf(fmt.Sprintf("Process error: %v, exiting", processErr))
	}

	return exitErr
}

func (s *Membrane) Stop() {
	if s.tracerProvider != nil {
		_ = s.tracerProvider.Shutdown(context.Background())
	}
	_ = s.gatewayPlugin.Stop()
	s.grpcServer.Stop()
	s.processManager.StopAll()
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

	if !options.TolerateMissingServices {
		tolerateMissing, err := strconv.ParseBool(utils.GetEnv("TOLERATE_MISSING_SERVICES", "false"))
		if err != nil {
			return nil, err
		}
		options.TolerateMissingServices = tolerateMissing
	}

	if options.ChildTimeoutSeconds < 1 {
		options.ChildTimeoutSeconds = 10
	}

	if options.GatewayPlugin == nil {
		return nil, errors.New("missing gateway plugin, Gateway plugin must not be nil")
	}

	if !options.TolerateMissingServices {
		if options.EventsPlugin == nil || options.StoragePlugin == nil || options.DocumentPlugin == nil || options.QueuePlugin == nil {
			return nil, errors.New("missing membrane plugins, if you meant to load with missing plugins set options.TolerateMissingServices to true")
		}
	}

	if options.Pool == nil {
		// Create new pool with defaults
		minWorkersEnv := utils.GetEnv("MIN_WORKERS", "1")
		minWorkers, err := strconv.Atoi(minWorkersEnv)
		if err != nil || minWorkers < 0 {
			return nil, fmt.Errorf("invalid MIN_WORKERS env var, expected non-negative integer value, got %v", minWorkersEnv)
		}

		maxWorkersEnv := utils.GetEnv("MAX_WORKERS", "300")
		maxWorkers, err := strconv.Atoi(maxWorkersEnv)
		if err != nil || minWorkers < 0 {
			return nil, fmt.Errorf("invalid MAX_WORKERS env var, expected non-negative integer value, got %v", maxWorkersEnv)
		}

		if minWorkers > maxWorkers {
			maxWorkers = minWorkers
		}

		options.Pool = pool.NewProcessPool(&pool.ProcessPoolOptions{
			MinWorkers: minWorkers,
			MaxWorkers: maxWorkers,
		})
	}

	bin := utils.GetEnv("OTELCOL_BIN", "/usr/bin/otelcol-contrib")
	config := utils.GetEnv("OTELCOL_CONFIG", "/etc/otelcol/config.yaml")
	createTracerProvider := options.CreateTracerProvider

	if createTracerProvider != nil && fileExists(bin) && fileExists(config) {
		log.Default().Println("Tracing is enabled")

		options.PreCommands = [][]string{
			{
				bin, "--config", config,
			},
		}

		options.Pool = &pool.InstrumentedWorkerPool{
			WorkerPool: options.Pool,
			Wrapper:    worker.InstrumentedWorkerFn,
		}
	} else {
		log.Default().Printf("Tracing is disabled %v %v %v", createTracerProvider != nil, fileExists(bin), fileExists(config))
		createTracerProvider = nil
	}

	return &Membrane{
		serviceAddress:          options.ServiceAddress,
		childAddress:            options.ChildAddress,
		childUrl:                fmt.Sprintf("http://%s", options.ChildAddress),
		processManager:          pm.NewProcessManager(options.ChildCommand, options.PreCommands),
		createTracerProvider:    createTracerProvider,
		childTimeoutSeconds:     options.ChildTimeoutSeconds,
		documentPlugin:          options.DocumentPlugin,
		eventsPlugin:            options.EventsPlugin,
		storagePlugin:           options.StoragePlugin,
		queuePlugin:             options.QueuePlugin,
		gatewayPlugin:           options.GatewayPlugin,
		secretPlugin:            options.SecretPlugin,
		resourcePlugin:          options.ResourcesPlugin,
		suppressLogs:            options.SuppressLogs,
		tolerateMissingServices: options.TolerateMissingServices,
		pool:                    options.Pool,
	}, nil
}
