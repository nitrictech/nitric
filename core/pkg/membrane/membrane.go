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
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/nitrictech/nitric/core/pkg/decorators"
	"github.com/nitrictech/nitric/core/pkg/env"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	pm "github.com/nitrictech/nitric/core/pkg/process"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	httppb "github.com/nitrictech/nitric/core/pkg/proto/http/v1"
	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type MembraneOptions struct {
	ServiceAddress string
	// The command that will be used to invoke the child process
	ChildCommand []string
	// Commands that will be started before all others
	PreCommands [][]string

	// The total time to wait for the child process to be available in seconds
	ChildTimeoutSeconds int

	// The provider adapter gateway
	GatewayPlugin gateway.GatewayService

	// The minimum number of workers that need to be available
	MinWorkers *int

	// Resource access plugins
	KeyValuePlugin      kvstorepb.KvStoreServer
	TopicsPlugin        topicspb.TopicsServer
	StoragePlugin       storagepb.StorageServer
	SecretManagerPlugin secretspb.SecretManagerServer
	WebsocketPlugin     websocketspb.WebsocketServer
	QueuesPlugin        queuespb.QueuesServer
	SqlPlugin           sqlpb.SqlServer

	// Worker plugins
	ApiPlugin               apis.ApiRequestHandler
	HttpPlugin              http.HttpRequestHandler
	SchedulesPlugin         schedules.ScheduleRequestHandler
	TopicsListenerPlugin    topics.SubscriptionRequestHandler
	StorageListenerPlugin   storage.BucketRequestHandler
	WebsocketListenerPlugin websockets.WebsocketRequestHandler

	// Server listeners

	ResourcesPlugin resourcespb.ResourcesServer

	SuppressLogs            bool
	TolerateMissingServices bool
}

type Membrane struct {
	processManager pm.ProcessManager
	options        MembraneOptions

	// Suppress println statements in the membrane server
	suppressLogs bool

	grpcServer *grpc.Server

	minWorkers int
}

func (s *Membrane) WorkerCount() int {
	return s.options.ApiPlugin.WorkerCount() +
		s.options.HttpPlugin.WorkerCount() +
		s.options.SchedulesPlugin.WorkerCount() +
		s.options.TopicsListenerPlugin.WorkerCount() +
		s.options.StorageListenerPlugin.WorkerCount() +
		s.options.WebsocketListenerPlugin.WorkerCount()
}

func (s *Membrane) waitForMinimumWorkers(timeout int) error {
	waitUntil := time.Now().Add(time.Duration(timeout) * time.Second)
	ticker := time.NewTicker(time.Duration(5) * time.Millisecond)

	// stop the ticker on exit
	defer ticker.Stop()

	for {
		if s.WorkerCount() >= s.minWorkers {
			break
		}

		// wait for the next tick
		time := <-ticker.C

		if time.After(waitUntil) {
			return fmt.Errorf("available workers below required minimum of %d, %d available, timed out waiting for more workers", s.minWorkers, s.WorkerCount())
		}
	}

	return nil
}

type MembraneStartOptions func(m *Membrane)

func WithGrpcServer(s *grpc.Server) MembraneStartOptions {
	return func(m *Membrane) {
		m.grpcServer = s
	}
}

// Start the membrane
func (s *Membrane) Start(startOpts ...MembraneStartOptions) error {
	if err := s.processManager.StartPreProcesses(); err != nil {
		return err
	}

	for _, opt := range startOpts {
		opt(s)
	}

	maxWorkers, err := env.MAX_WORKERS.Int()
	if err != nil {
		return err
	}

	if s.grpcServer == nil {
		opts := []grpc.ServerOption{
			grpc.MaxConcurrentStreams(uint32(maxWorkers)),
		}

		s.grpcServer = grpc.NewServer(opts...)
	}

	// Register the listener servers
	if s.options.ApiPlugin == nil {
		s.options.ApiPlugin = apis.New()
	}
	apispb.RegisterApiServer(s.grpcServer, s.options.ApiPlugin)

	if s.options.TopicsListenerPlugin == nil {
		s.options.TopicsListenerPlugin = topics.New()
	}
	topicspb.RegisterSubscriberServer(s.grpcServer, s.options.TopicsListenerPlugin)

	if s.options.StorageListenerPlugin == nil {
		s.options.StorageListenerPlugin = storage.New()
	}
	storagepb.RegisterStorageListenerServer(s.grpcServer, s.options.StorageListenerPlugin)

	if s.options.SchedulesPlugin == nil {
		s.options.SchedulesPlugin = schedules.New()
	}
	schedulespb.RegisterSchedulesServer(s.grpcServer, s.options.SchedulesPlugin)

	if s.options.WebsocketListenerPlugin == nil {
		s.options.WebsocketListenerPlugin = websockets.NewWebsocketManager()
	}
	websocketspb.RegisterWebsocketHandlerServer(s.grpcServer, s.options.WebsocketListenerPlugin)

	if s.options.HttpPlugin == nil {
		s.options.HttpPlugin = http.New()
	}
	httppb.RegisterHttpServer(s.grpcServer, s.options.HttpPlugin)

	// Load & Register the service plugins
	secretsServerWithValidation := decorators.SecretsServerWithValidation(s.options.SecretManagerPlugin)
	keyvalueServerWithCompat := decorators.KeyValueServerWithCompat(s.options.KeyValuePlugin)

	kvstorepb.RegisterKvStoreServer(s.grpcServer, keyvalueServerWithCompat)
	// TODO: Deprecated, remove in future release
	keyvaluepb.RegisterKeyValueServer(s.grpcServer, keyvalueServerWithCompat)
	topicspb.RegisterTopicsServer(s.grpcServer, s.options.TopicsPlugin)
	storagepb.RegisterStorageServer(s.grpcServer, s.options.StoragePlugin)
	secretspb.RegisterSecretManagerServer(s.grpcServer, secretsServerWithValidation)
	resourcespb.RegisterResourcesServer(s.grpcServer, s.options.ResourcesPlugin)
	websocketspb.RegisterWebsocketServer(s.grpcServer, s.options.WebsocketPlugin)
	queuespb.RegisterQueuesServer(s.grpcServer, s.options.QueuesPlugin)
	sqlpb.RegisterSqlServer(s.grpcServer, s.options.SqlPlugin)

	lis, err := net.Listen("tcp", s.options.ServiceAddress)
	if err != nil {
		return fmt.Errorf("could not listen on configured service address: %w", err)
	}

	logger.Debug("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		logger.Debugf("Services listening on: %s", s.options.ServiceAddress)
		err := s.grpcServer.Serve(lis)
		if err != nil {
			logger.Errorf("grpc serve %v", err)
		}
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if err := s.processManager.StartUserProcess(); err != nil {
		return err
	}

	// Wait for the minimum number of active workers to be available before beginning the gateway
	// This ensures workers have registered and can handle triggers as soon the gateway is ready, if a minimum > 1 has been set
	logger.Debug("Waiting for active workers")
	err = s.waitForMinimumWorkers(s.options.ChildTimeoutSeconds)
	if err != nil {
		return err
	}

	gatewayErrchan := make(chan error)
	// poolErrchan := make(chan error)

	// Start the gateway
	go func(errch chan error) {
		logger.Debugf("Starting Gateway, %d workers currently available", s.WorkerCount())

		errch <- s.options.GatewayPlugin.Start(&gateway.GatewayStartOpts{
			ApiPlugin:               s.options.ApiPlugin,
			HttpPlugin:              s.options.HttpPlugin,
			SchedulesPlugin:         s.options.SchedulesPlugin,
			TopicsListenerPlugin:    s.options.TopicsListenerPlugin,
			StorageListenerPlugin:   s.options.StorageListenerPlugin,
			WebsocketListenerPlugin: s.options.WebsocketListenerPlugin,
		})
	}(gatewayErrchan)

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
	case processErr := <-processErrchan:
		exitErr = fmt.Errorf(fmt.Sprintf("Process error: %v, exiting", processErr))
	}

	return exitErr
}

func (s *Membrane) Stop() {
	_ = s.options.GatewayPlugin.Stop()
	s.grpcServer.Stop()
	s.processManager.StopAll()
}

// Create a new Membrane server
func New(options *MembraneOptions) (*Membrane, error) {
	// Get unset options from env or defaults
	if options.ServiceAddress == "" {
		options.ServiceAddress = env.SERVICE_ADDRESS.String()
	}

	minWorkers, err := env.MIN_WORKERS.Int()
	if options.MinWorkers != nil {
		minWorkers = *options.MinWorkers
	} else if err != nil {
		return nil, err
	}

	workerTimeout, err := env.WORKER_TIMEOUT.Int()
	if options.ChildTimeoutSeconds < 1 {
		options.ChildTimeoutSeconds = workerTimeout
	} else if err != nil {
		return nil, err
	}

	if options.GatewayPlugin == nil {
		return nil, errors.New("missing gateway plugin, Gateway plugin must not be nil")
	}

	return &Membrane{
		processManager: pm.NewProcessManager(options.ChildCommand, options.PreCommands),
		options:        *options,
		minWorkers:     minWorkers,
		suppressLogs:   options.SuppressLogs,
	}, nil
}
