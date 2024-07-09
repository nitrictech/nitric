package job

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"google.golang.org/grpc"
)

// JobMembrane is a membrane for job
// Created as a sepearate membrane type to avoid overloading the service membrane
type JobMembrane struct {
	// The Command that will be executed to run the job
	cmd []string

	// Runtime plugins (for reading/writing to cloud services)
	topicServer     topicspb.TopicsServer
	storageServer   storagepb.StorageServer
	queueServer     queuespb.QueuesServer
	secretServer    secretspb.SecretManagerServer
	sqlServer       sqlpb.SqlServer
	resourcesServer resourcespb.ResourcesServer
}

func (j *JobMembrane) Run() error {
	// Start the gRPC server for runtime services
	grpcServer := grpc.NewServer()

	topicspb.RegisterTopicsServer(grpcServer, j.topicServer)
	storagepb.RegisterStorageServer(grpcServer, j.storageServer)
	queuespb.RegisterQueuesServer(grpcServer, j.queueServer)
	secretspb.RegisterSecretManagerServer(grpcServer, j.secretServer)
	sqlpb.RegisterSqlServer(grpcServer, j.sqlServer)
	resourcespb.RegisterResourcesServer(grpcServer, j.resourcesServer)

	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return err
	}

	// Start the grpc services
	go func() {
		fmt.Printf("Starting nitric gRPC server on %s\n", lis.Addr().String())
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
		lis.Close()
	}()

	// Give the grpc server time to start up
	// TODO: Replace this with logic for waiting until the port is active
	fmt.Printf("Waiting for gRPC server to start\n")
	time.Sleep(1 * time.Second)

	fmt.Printf("Running command: %+v\n", j.cmd)

	// Run the command and wait for it to exit
	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdEnv := append([]string{}, os.Environ()...)

	cmdEnv = append(cmdEnv, fmt.Sprintf("SERVICE_ADDRESS=%s", lis.Addr().String()))
	// copy the current environment variables
	cmd.Env = cmdEnv

	fmt.Printf("Running with env variables: %+v\n", cmd.Env)

	return cmd.Run()
}

func NewJobMembrane(cmd []string, options ...JobMembraneOption) *JobMembrane {
	membrane := &JobMembrane{
		cmd:             cmd,
		topicServer:     topicspb.UnimplementedTopicsServer{},
		storageServer:   storagepb.UnimplementedStorageServer{},
		queueServer:     queuespb.UnimplementedQueuesServer{},
		secretServer:    secretspb.UnimplementedSecretManagerServer{},
		sqlServer:       sqlpb.UnimplementedSqlServer{},
		resourcesServer: resourcespb.UnimplementedResourcesServer{},
	}

	for _, option := range options {
		// Apply the option to the membrane
		option(membrane)
	}

	return membrane
}
