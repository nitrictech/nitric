package membrane

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	documentsPb "github.com/nitric-dev/membrane/interfaces/nitric/v1/documents"
	eventingPb "github.com/nitric-dev/membrane/interfaces/nitric/v1/eventing"
	storagePb "github.com/nitric-dev/membrane/interfaces/nitric/v1/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/services"
	"google.golang.org/grpc"
)

type MembraneOptions struct {
	ServiceAddress string
	// The address the child will be listening on
	ChildAddress string
	// The command that will be used to invoke the child process
	ChildCommand string
	// The total time to wait for the child process to be available in seconds
	ChildTimeoutSeconds int

	EventingPlugin  sdk.EventingPlugin
	DocumentsPlugin sdk.DocumentsPlugin
	StoragePlugin   sdk.StoragePlugin
	GatewayPlugin   sdk.GatewayPlugin

	SuppressLogs            bool
	TolerateMissingServices bool
}

type Membrane struct {
	// Address & port to bind the membrane i/o proxy to
	// This will still be bound even in pass through mode
	// proxyAddress string
	// Address & port to bind the membrane service interfaces to
	serviceAddress string
	// The address the child will be listening on
	childAddress string
	// The command that will be used to invoke the child process
	childCommand string

	childTimeoutSeconds int

	// Configured plugins
	eventingPlugin  sdk.EventingPlugin
	documentsPlugin sdk.DocumentsPlugin
	storagePlugin   sdk.StoragePlugin
	gatewayPlugin   sdk.GatewayPlugin

	// Tolerate if services are not available
	// Not this does not include the gateway service
	tolerateMissingServices bool

	// Suppress println statements in the membrane server
	supressLogs bool
}

func (s *Membrane) log(log string) {
	if !s.supressLogs {
		fmt.Println(log)
	}
}

// Create a new Nitric Eventing Server
func (s *Membrane) createEventingServer() (eventingPb.EventingServer, error) {
	// Cast to the new eventing server function
	if s.eventingPlugin != nil {
		return services.NewEventingServer(s.eventingPlugin), nil
	} else {
		return nil, fmt.Errorf("Eventing plugin not configured")
	}
}

// Create a new Nitric Storage Server
func (s *Membrane) createStorageServer() (storagePb.StorageServer, error) {
	// Cast to the new storage server function
	if s.storagePlugin != nil {
		return services.NewStorageServer(s.storagePlugin), nil
	} else {
		return nil, fmt.Errorf("Storage plugin not configured")
	}
}

func (s *Membrane) createDocumentsServer() (documentsPb.DocumentsServer, error) {
	// Cast to the new documents server function
	if s.documentsPlugin != nil {
		return services.NewDocumentsServer(s.documentsPlugin), nil
	} else {
		return nil, fmt.Errorf("Documents plugin not configured")
	}
}

// Provides a means for the nitric membrane to accept and normalize input/output for a given interface
// TODO: Create a entrypoint plugin for different styles of membrane
// data ingress/egress, this could include
// - HTTP Proxy plugin (for providing a HTTP proxy down to user land code/applications)
// - AWS Lambda plugin (for querying the AWS lambda service and directing/normalizing input to user land code)
// - Kafka Plugin (for providing a streaming server)
func (s *Membrane) loadGatewayPlugin() (sdk.GatewayPlugin, error) {
	if s.gatewayPlugin != nil {
		// There was an error loading the eventing plugin
		return s.gatewayPlugin, nil
	}

	return nil, fmt.Errorf("Gateway plugin not configured")
}

func (s *Membrane) startChildProcess() {
	// TODO: This is a detached process
	// so it will continue to run until even after the director dies
	commandArgs := strings.Fields(s.childCommand)

	fmt.Println(fmt.Sprintf("Starting Function"))
	childProcess := exec.Command(commandArgs[0], commandArgs[1:len(commandArgs)]...)
	childProcess.Stdout = os.Stdout
	childProcess.Stderr = os.Stderr
	applicationError := childProcess.Start()

	// Actual panic here, we don't want to start if our userland code cannot successfully start
	if applicationError != nil {
		panic(fmt.Errorf("Function failed to start in time: %v", applicationError))
	}

	// Dial the child port to see if it's open and ready...
	// Only wait for 10s, if we timeout that will be it
	// TODO: make app startup time configurable
	maxWaitTime := time.Duration(s.childTimeoutSeconds) * time.Second
	pollInterval := time.Duration(200) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		conn, _ := net.Dial("tcp", s.childAddress)
		if conn != nil {
			conn.Close()
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				panic(fmt.Errorf("Function failed to start in time"))
			}
		}
	}
}

// Start the membrane
func (s *Membrane) Start() {
	// Search for known plugins

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Load & Register the GRPC service plugins
	eventingServer, err := s.createEventingServer()

	// There was a failure loading the eventing server
	// XXX: For now we will gracefully continue
	// However we will likely want to use env variables to determine required services
	// for a given function
	if err == nil {
		// Register the service
		eventingPb.RegisterEventingServer(grpcServer, eventingServer)
		s.log("Registered Eventing Plugin")
	} else if s.tolerateMissingServices {
		s.log(fmt.Sprintf("Failed to load eventing plugin %v", err))
	} else {
		panic(fmt.Errorf("Fatal error loading eventing plugin %v", err))
	}

	documentsServer, err := s.createDocumentsServer()
	if err == nil {
		// Register the service
		documentsPb.RegisterDocumentsServer(grpcServer, documentsServer)
		s.log("Registered Documents Plugin")
	} else if s.tolerateMissingServices {
		s.log(fmt.Sprintf("Failed to load documents plugin %v", err))
	} else {
		panic(fmt.Errorf("Fatal error loading documents plugin %v", err))
	}

	storageServer, err := s.createStorageServer()
	if err == nil {
		// Register the service
		storagePb.RegisterStorageServer(grpcServer, storageServer)
		s.log("Registered Storage Plugin")
	} else if s.tolerateMissingServices {
		s.log(fmt.Sprintf("Failed to load storage plugin %v", err))
	} else {
		panic(fmt.Errorf("Fatal error loading storage plugin %v", err))
	}

	lis, err := net.Listen("tcp", s.serviceAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gateway, err := s.loadGatewayPlugin()
	if err != nil {
		panic(err)
	}

	s.log("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		s.log(fmt.Sprintf("Services listening on: %s", s.serviceAddress))
		grpcServer.Serve(lis)
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if s.childCommand != "" {
		s.startChildProcess()
	} else {
		s.log("No Child Configured Specified, Skipping...")
	}

	// FIXME: Only do this in Gateway mode...
	// Otherwise always pass through to the provided child address
	// Start the Gateway Server

	// Start the gateway, this will provide us an entrypoint for
	// data ingress/egress to our userland code
	// The gateway should block the main thread but will
	// use this callback as a control mechanism
	s.log("Starting Gateway")
	err = gateway.Start(func(request *sdk.NitricRequest) *sdk.NitricResponse {
		childUrl := fmt.Sprintf("http://%s", s.childAddress)

		httpRequest, err := http.NewRequest("POST", childUrl, bytes.NewReader(request.Payload))

		// There was an error creating the HTTP request
		if err != nil {
			// return an error to the Gateway
			return &sdk.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(err.Error()),
				Status:  503,
			}
		}

		// Encode NitricContext into HTTP headers
		httpRequest.Header.Add("Content-Type", http.DetectContentType(request.Payload))
		httpRequest.Header.Add("x-nitric-payload-type", request.Context.PayloadType)
		httpRequest.Header.Add("x-nitric-request-id", request.Context.RequestId)
		httpRequest.Header.Add("x-nitric-source-type", request.Context.SourceType.String())
		httpRequest.Header.Add("x-nitric-source", request.Context.Source)

		// Send the request down to our function
		// Here we'll be making a normal http request to the function server
		// From here we will return the response from the server
		// Always do a post request to the local function???
		response, err := http.DefaultClient.Do(httpRequest)

		if err != nil {
			// there was an error calling the HTTP service
			return &sdk.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(err.Error()),
				Status:  503,
			}
		}

		responseBody, err := ioutil.ReadAll(response.Body)

		if err != nil {
			// There was an error reading the http response
			return &sdk.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(err.Error()),
				Status:  503,
			}
		}

		headers := map[string]string{}
		for name, value := range response.Header {
			headers[name] = value[0]
		}

		// Pass the response back to the gateway
		return &sdk.NitricResponse{
			Headers: headers,
			Body:    responseBody,
			Status:  response.StatusCode,
		}
	})
	// The gateway process has exited

	panic(err)
}

// Create a new Membrane server
func New(options *MembraneOptions) (*Membrane, error) {
	var childTimeout = 5
	if options.ChildTimeoutSeconds > 0 {
		childTimeout = options.ChildTimeoutSeconds
	}

	return &Membrane{
		serviceAddress:          options.ServiceAddress,
		childAddress:            options.ChildAddress,
		childCommand:            options.ChildCommand,
		childTimeoutSeconds:     childTimeout,
		eventingPlugin:          options.EventingPlugin,
		storagePlugin:           options.StoragePlugin,
		documentsPlugin:         options.DocumentsPlugin,
		gatewayPlugin:           options.GatewayPlugin,
		supressLogs:             options.SuppressLogs,
		tolerateMissingServices: options.TolerateMissingServices,
	}, nil
}
