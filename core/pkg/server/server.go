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

package server

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/nitrictech/nitric/core/pkg/decorators"
	"github.com/nitrictech/nitric/core/pkg/env"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	"github.com/nitrictech/nitric/core/pkg/logger"
	pm "github.com/nitrictech/nitric/core/pkg/process"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
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
	"github.com/nitrictech/nitric/core/pkg/server/runtime"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/jobs"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type NitricServer struct {
	processManager pm.ProcessManager
	grpcServer     *grpc.Server

	// Options
	ServiceAddress string
	// The command that will be used to invoke the child process
	ChildCommand []string
	// Commands that will be started before all others
	PreCommands [][]string

	// The total time to wait for the child process to be available in seconds
	ChildTimeoutSeconds int

	// The minimum number of workers that need to be available
	MinWorkers int

	// The provider adapter gateway
	GatewayPlugin gateway.GatewayService

	ResourcesPlugin resourcespb.ResourcesServer

	// Resource access plugins
	KeyValuePlugin      kvstorepb.KvStoreServer
	TopicsPlugin        topicspb.TopicsServer
	StoragePlugin       storagepb.StorageServer
	SecretManagerPlugin secretspb.SecretManagerServer
	WebsocketPlugin     websocketspb.WebsocketServer
	QueuesPlugin        queuespb.QueuesServer
	SqlPlugin           sqlpb.SqlServer
	BatchPlugin         batchpb.BatchServer

	// Worker plugins
	ApiPlugin               apis.ApiRequestHandler
	HttpPlugin              http.HttpRequestHandler
	SchedulesPlugin         schedules.ScheduleRequestHandler
	TopicsListenerPlugin    topics.SubscriptionRequestHandler
	StorageListenerPlugin   storage.BucketRequestHandler
	WebsocketListenerPlugin websockets.WebsocketRequestHandler
	JobHandlerPlugin        jobs.JobRequestHandler
}

func (s *NitricServer) WorkerCount() int {
	return s.ApiPlugin.WorkerCount() +
		s.HttpPlugin.WorkerCount() +
		s.SchedulesPlugin.WorkerCount() +
		s.TopicsListenerPlugin.WorkerCount() +
		s.StorageListenerPlugin.WorkerCount() +
		s.WebsocketListenerPlugin.WorkerCount() +
		s.JobHandlerPlugin.WorkerCount()
}

func (s *NitricServer) waitForMinimumWorkers(timeout int) error {
	waitUntil := time.Now().Add(time.Duration(timeout) * time.Second)
	ticker := time.NewTicker(time.Duration(5) * time.Millisecond)

	// stop the ticker on exit
	defer ticker.Stop()

	for {
		if s.WorkerCount() >= s.MinWorkers {
			break
		}

		// wait for the next tick
		time := <-ticker.C

		if time.After(waitUntil) {
			return fmt.Errorf("available workers below required minimum of %d, %d available, timed out waiting for more workers", s.MinWorkers, s.WorkerCount())
		}
	}

	return nil
}

type ServerStartOptions func(m *NitricServer)

func WithGrpcServer(s *grpc.Server) ServerStartOptions {
	return func(m *NitricServer) {
		m.grpcServer = s
	}
}

// Start the server
func (s *NitricServer) Start(startOpts ...ServerStartOptions) error {
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

	if maxWorkers < 0 || maxWorkers > math.MaxUint32 {
		return fmt.Errorf("MAX_WORKERS not in range 0-%d: %d", math.MaxUint32, maxWorkers)
	}

	if s.grpcServer == nil {
		opts := []grpc.ServerOption{
			grpc.MaxConcurrentStreams(uint32(maxWorkers)), //#nosec G115 -- max workers checked for potential out of range or overflow errors
		}

		s.grpcServer = grpc.NewServer(opts...)
	}

	// Register the listener servers
	if s.ApiPlugin == nil {
		s.ApiPlugin = apis.New()
	}
	apispb.RegisterApiServer(s.grpcServer, s.ApiPlugin)

	if s.TopicsListenerPlugin == nil {
		s.TopicsListenerPlugin = topics.New()
	}
	topicspb.RegisterSubscriberServer(s.grpcServer, s.TopicsListenerPlugin)

	if s.StorageListenerPlugin == nil {
		s.StorageListenerPlugin = storage.New()
	}
	storagepb.RegisterStorageListenerServer(s.grpcServer, s.StorageListenerPlugin)

	if s.SchedulesPlugin == nil {
		s.SchedulesPlugin = schedules.New()
	}
	schedulespb.RegisterSchedulesServer(s.grpcServer, s.SchedulesPlugin)

	if s.WebsocketListenerPlugin == nil {
		s.WebsocketListenerPlugin = websockets.NewWebsocketManager()
	}
	websocketspb.RegisterWebsocketHandlerServer(s.grpcServer, s.WebsocketListenerPlugin)

	if s.HttpPlugin == nil {
		s.HttpPlugin = http.New()
	}
	httppb.RegisterHttpServer(s.grpcServer, s.HttpPlugin)

	if s.JobHandlerPlugin == nil {
		s.JobHandlerPlugin = jobs.New()
	}
	batchpb.RegisterJobServer(s.grpcServer, s.JobHandlerPlugin)

	// Load & Register the service plugins
	secretsServerWithValidation := decorators.SecretsServerWithValidation(s.SecretManagerPlugin)
	keyvalueServerWithCompat := decorators.KeyValueServerWithCompat(s.KeyValuePlugin)

	kvstorepb.RegisterKvStoreServer(s.grpcServer, keyvalueServerWithCompat)
	keyvaluepb.RegisterKeyValueServer(s.grpcServer, keyvalueServerWithCompat)
	topicspb.RegisterTopicsServer(s.grpcServer, s.TopicsPlugin)
	storagepb.RegisterStorageServer(s.grpcServer, s.StoragePlugin)
	secretspb.RegisterSecretManagerServer(s.grpcServer, secretsServerWithValidation)
	resourcespb.RegisterResourcesServer(s.grpcServer, s.ResourcesPlugin)
	websocketspb.RegisterWebsocketServer(s.grpcServer, s.WebsocketPlugin)
	queuespb.RegisterQueuesServer(s.grpcServer, s.QueuesPlugin)
	sqlpb.RegisterSqlServer(s.grpcServer, s.SqlPlugin)
	batchpb.RegisterBatchServer(s.grpcServer, s.BatchPlugin)

	lis, err := net.Listen("tcp", s.ServiceAddress)
	if err != nil {
		return fmt.Errorf("could not listen on configured service address: %w", err)
	}

	logger.Debug("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		logger.Debugf("Services listening on: %s", s.ServiceAddress)
		err := s.grpcServer.Serve(lis)
		if err != nil {
			logger.Errorf("grpc serve %v", err)
		}
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if err := s.processManager.StartUserProcess(fmt.Sprintf("SERVICE_ADDRESS=%s", lis.Addr().String())); err != nil {
		return err
	}

	// Wait for the minimum number of active workers to be available before beginning the gateway
	// This ensures workers have registered and can handle triggers as soon the gateway is ready, if a minimum > 1 has been set
	logger.Debug("Waiting for active workers")
	err = s.waitForMinimumWorkers(s.ChildTimeoutSeconds)
	if err != nil {
		return err
	}

	gatewayErrchan := make(chan error)

	// Start the gateway
	go func(errch chan error) {
		logger.Debugf("Starting Gateway, %d workers currently available", s.WorkerCount())

		errch <- s.GatewayPlugin.Start(&gateway.GatewayStartOpts{
			ApiPlugin:               s.ApiPlugin,
			HttpPlugin:              s.HttpPlugin,
			SchedulesPlugin:         s.SchedulesPlugin,
			TopicsListenerPlugin:    s.TopicsListenerPlugin,
			StorageListenerPlugin:   s.StorageListenerPlugin,
			WebsocketListenerPlugin: s.WebsocketListenerPlugin,
			JobHandlerPlugin:        s.JobHandlerPlugin,
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
			// Allowing the server to exit
			return nil
		}
		exitErr = fmt.Errorf("Gateway Error: %w, exiting", gatewayErr)
	case processErr := <-processErrchan:
		exitErr = fmt.Errorf("Process error: %w, exiting", processErr)
	}

	return exitErr
}

func (s *NitricServer) Stop() {
	_ = s.GatewayPlugin.Stop()
	s.grpcServer.Stop()
	s.processManager.StopAll()
}

// New - Create a new nitric server
func New(opts ...ServerOption) (*NitricServer, error) {
	m := &NitricServer{
		// The resource service is defaulted, because it typically isn't required to be implemented for runtime servers.
		ResourcesPlugin: &runtime.RuntimeResourceService{},
		MinWorkers:      -1,
	}

	for _, opt := range opts {
		opt(m)
	}

	// Get unset options from env or defaults
	if m.ServiceAddress == "" {
		m.ServiceAddress = env.SERVICE_ADDRESS.String()
	}

	minWorkersEnv, err := env.MIN_WORKERS.Int()
	if err == nil && m.MinWorkers < 0 {
		logger.Debugf("MIN_WORKERS environment variable set to %d", minWorkersEnv)
		m.MinWorkers = minWorkersEnv
	}

	workerTimeout, err := env.WORKER_TIMEOUT.Int()
	if m.ChildTimeoutSeconds < 1 {
		m.ChildTimeoutSeconds = workerTimeout
	} else if err != nil {
		return nil, fmt.Errorf("invalid WORKER_TIMEOUT: %w", err)
	}

	if m.ChildCommand == nil {
		if len(os.Args) > 1 {
			m.ChildCommand = os.Args[1:]
		}
	}

	if m.GatewayPlugin == nil {
		return nil, errors.New("missing gateway plugin, Gateway plugin must not be nil")
	}

	m.processManager = pm.NewProcessManager(m.ChildCommand, m.PreCommands)

	return m, nil
}
