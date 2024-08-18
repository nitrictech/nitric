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

package job

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"google.golang.org/grpc"
)

// NitricJobServer is a membrane for job
// Created as a separate membrane type to avoid overloading the service membrane
type NitricJobServer struct {
	// The Command that will be executed to run the job
	cmd string

	// Runtime plugins (for reading/writing to cloud services)
	topicServer     topicspb.TopicsServer
	storageServer   storagepb.StorageServer
	queueServer     queuespb.QueuesServer
	secretServer    secretspb.SecretManagerServer
	sqlServer       sqlpb.SqlServer
	kvStoreServer   kvstorepb.KvStoreServer
	websocketServer websocketspb.WebsocketServer
	batchServer     batchpb.BatchServer
}

func (j *NitricJobServer) Run() error {
	// Start the gRPC server for runtime services
	grpcServer := grpc.NewServer()

	topicspb.RegisterTopicsServer(grpcServer, j.topicServer)
	storagepb.RegisterStorageServer(grpcServer, j.storageServer)
	queuespb.RegisterQueuesServer(grpcServer, j.queueServer)
	secretspb.RegisterSecretManagerServer(grpcServer, j.secretServer)
	sqlpb.RegisterSqlServer(grpcServer, j.sqlServer)
	kvstorepb.RegisterKvStoreServer(grpcServer, j.kvStoreServer)
	websocketspb.RegisterWebsocketServer(grpcServer, j.websocketServer)
	batchpb.RegisterBatchServer(grpcServer, j.batchServer)

	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return err
	}

	// Start the grpc services
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("an error occurred with the nitric runtime server: %v", err)
		}
	}()

	defer grpcServer.GracefulStop()

	// Run the command and wait for it to exit
	cmdParts := strings.Split(j.cmd, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...) //#nosec G204 -- This is by design inputs are determined at compile time

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdEnv := append([]string{}, os.Environ()...)

	cmdEnv = append(cmdEnv, fmt.Sprintf("SERVICE_ADDRESS=%s", lis.Addr().String()))
	// copy the current environment variables
	cmd.Env = cmdEnv

	return cmd.Run()
}

func NewJobServer(cmd string, options ...nitricJobServerOption) *NitricJobServer {
	membrane := &NitricJobServer{
		cmd:             cmd,
		topicServer:     topicspb.UnimplementedTopicsServer{},
		storageServer:   storagepb.UnimplementedStorageServer{},
		queueServer:     queuespb.UnimplementedQueuesServer{},
		secretServer:    secretspb.UnimplementedSecretManagerServer{},
		sqlServer:       sqlpb.UnimplementedSqlServer{},
		batchServer:     batchpb.UnimplementedBatchServer{},
		kvStoreServer:   kvstorepb.UnimplementedKvStoreServer{},
		websocketServer: websocketspb.UnimplementedWebsocketServer{},
	}

	for _, option := range options {
		// Apply the option to the membrane
		option(membrane)
	}

	return membrane
}
