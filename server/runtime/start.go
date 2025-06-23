package runtime

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/nitrictech/nitric/server/runtime/plugin"
	"github.com/nitrictech/nitric/server/runtime/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterPlugins[T any](register plugin.Register[T], plugins map[string]plugin.Constructor[T]) error {
	// Register the plugins
	for name, constructor := range plugins {
		err := register(name, constructor)
		if err != nil {
			return err
		}
	}

	return nil
}

// waitForPort attempts to connect to the given port until it succeeds or times out
func waitForPort(host string, port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for port %s to be available", port)
}

func Start(cmd string) {
	fmt.Println("Starting server with command:", cmd)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	// Register plugin router
	storagepb.RegisterStorageServer(srv, &GrpcServer{})
	pubsubpb.RegisterPubsubServer(srv, &GrpcServer{})

	// Register reflection service on gRPC server
	reflection.Register(srv)

	log.Printf("Starting server on %s", lis.Addr().String())

	// Start the runtime services in a goroutine
	go func() {
		srv.Serve(lis)
	}()

	// Start the actual nitric service
	cmdParts := strings.Split(cmd, " ")
	runCmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	runCmd.Env = os.Environ()
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	servicePort := os.Getenv("NITRIC_GUEST_PORT")
	if servicePort == "" {
		servicePort = os.Getenv("PORT")
	} else {
		runCmd.Env = append(runCmd.Env, fmt.Sprintf("PORT=%s", servicePort))
	}

	if err := runCmd.Start(); err != nil {
		log.Fatalf("failed to start service: %v", err)
	}

	// Wait for the service to be ready (up to 30 seconds)
	// DO NOT CHANGE THE ORDER OF THIS
	// The Guest application port MUST be ready before the service gateway is started
	// Otherwise CPU may be throttled in serverless environments like AWS Lambda if the service gateway is started too early
	// Meaning the guest application may never get a chance to start
	log.Printf("Waiting for service to be ready on port %s...", servicePort)
	if err := waitForPort("localhost", servicePort, 10*time.Second); err != nil {
		log.Fatalf("service failed to start: %v", err)
	}
	log.Printf("Service is ready on port %s", servicePort)

	// Start the service gateway and proxy
	err = service.Start(service.NewHttpServerProxy(fmt.Sprintf("localhost:%s", servicePort)))
	if err != nil {
		log.Fatalf("failed to start ingress: %v", err)
	}
}
