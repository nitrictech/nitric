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
	"plugin"
	"strings"
	"time"

	gw "github.com/nitric-dev/membrane-plugin-sdk"
	documentsPb "github.com/nitric-dev/membrane-plugin-sdk/nitric/v1/documents"
	eventingPb "github.com/nitric-dev/membrane-plugin-sdk/nitric/v1/eventing"
	storagePb "github.com/nitric-dev/membrane-plugin-sdk/nitric/v1/storage"
	"google.golang.org/grpc"
)

// TODO: return eventing server
type NewEventingServer func() (eventingPb.EventingServer, error)
type NewStorageServer func() (storagePb.StorageServer, error)
type NewDocumentsServer func() (documentsPb.DocumentsServer, error)
type NewGateway func() (gw.Gateway, error)

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
	// The plugin directory for loading membrane plugins
	pluginDir string

	// Plugin file names
	eventingPluginFile  string
	documentsPluginFile string
	storagePluginFile   string
	gatewayPluginFile   string
}

// Create a new Nitric Eventing Server
func (s *Membrane) createEventingServer() (eventingPb.EventingServer, error) {
	pluginLocation := fmt.Sprintf("%s/%s", s.pluginDir, s.eventingPluginFile)
	eventingPlugin, err := plugin.Open(pluginLocation)
	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("There was an issue loading the Nitric eventing plugin from %s: %v", pluginLocation, err)
	}

	// Lookup the New method for the eventing server
	newEventingServer, err := eventingPlugin.Lookup("New")

	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("Interface for Eventing Server was incorrect: %v", err)
	}

	// Cast to the new eventing server function
	newServerFunc := newEventingServer.(func() (eventingPb.EventingServer, error))

	// Return the new eventing server
	return newServerFunc()
}

// Create a new Nitric Storage Server
func (s *Membrane) createStorageServer() (storagePb.StorageServer, error) {
	pluginLocation := fmt.Sprintf("%s/%s", s.pluginDir, s.storagePluginFile)
	storagePlugin, err := plugin.Open(pluginLocation)
	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("There was an issue loading the Nitric storage plugin from %s: %v", pluginLocation, err)
	}

	// Lookup the New method for the eventing server
	newEventingServer, err := storagePlugin.Lookup("New")

	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("Interface for Storage Server was incorrect: %v", err)
	}

	// Cast to the new storage server function
	newServerFunc := newEventingServer.(func() (storagePb.StorageServer, error))

	// Return the new storage server
	return newServerFunc()
}

func (s *Membrane) createDocumentsServer() (documentsPb.DocumentsServer, error) {
	pluginLocation := fmt.Sprintf("%s/%s", s.pluginDir, s.documentsPluginFile)
	documentsPlugin, err := plugin.Open(pluginLocation)
	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("There was an issue loading the Nitric documents plugin from %s: %v", pluginLocation, err)
	}

	// Lookup the New method for the eventing server
	newDocumentsServer, err := documentsPlugin.Lookup("New")

	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("Interface for Documents Server was incorrect: %v", err)
	}

	// Cast to the new documents server function
	newServerFunc := newDocumentsServer.(func() (documentsPb.DocumentsServer, error))

	// Return the new documents server
	documentsServer, error := newServerFunc()

	return documentsServer.(documentsPb.DocumentsServer), error
}

// Provides a means for the nitric membrane to accept and normalize input/output for a given interface
// TODO: Create a entrypoint plugin for different styles of membrane
// data ingress/egress, this could include
// - HTTP Proxy plugin (for providing a HTTP proxy down to user land code/applications)
// - AWS Lambda plugin (for querying the AWS lambda service and directing/normalizing input to user land code)
// - Kafka Plugin (for providing a streaming server)
func (s *Membrane) loadGatewayPlugin() (gw.Gateway, error) {
	pluginLocation := fmt.Sprintf("%s/%s", s.pluginDir, s.gatewayPluginFile)
	// We expect that the gateway plugin will block the primary thread while it is processing
	// userland input
	gatewayPlugin, err := plugin.Open(pluginLocation)
	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("There was an issue loading the Nitric gateway plugin from %s: %v", pluginLocation, err)
	}

	// Lookup the New method for the eventing server
	newGatewayPlugin, err := gatewayPlugin.Lookup("New")

	if err != nil {
		// There was an error loading the eventing plugin
		return nil, fmt.Errorf("Interface for Gateway Server was incorrect: %v", err)
	}

	// Cast to the new documents server function
	newServerFunc := newGatewayPlugin.(func() (gw.Gateway, error))

	// Return the new documents server
	return newServerFunc()
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
		log.Fatalf("Function failed to start in time: %v", applicationError)
	}

	// Dial the child port to see if it's open and ready...
	// Only wait for 10s, if we timeout that will be it
	// TODO: make app startup time configurable
	maxWaitTime := time.Duration(5) * time.Second
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
				log.Fatalf("Function failed to start in time")
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
	eventingServer, error := s.createEventingServer()

	// There was a failure loading the eventing server
	// XXX: For now we will gracefully continue
	// However we will likely want to use env variables to determine required services
	// for a given function
	if error == nil {
		// Register the service
		eventingPb.RegisterEventingServer(grpcServer, eventingServer)
		fmt.Println("Registered Eventing Plugin")
	} else {
		fmt.Println("Failed to load eventing plugin", error)
	}

	documentsServer, error := s.createDocumentsServer()
	if error == nil {
		// Register the service
		documentsPb.RegisterDocumentsServer(grpcServer, documentsServer)
		fmt.Println("Registered Documents Plugin")
	} else {
		fmt.Println("Failed to load documents plugin", error)
	}

	storageServer, error := s.createStorageServer()
	if error == nil {
		// Register the service
		storagePb.RegisterStorageServer(grpcServer, storageServer)
		fmt.Println("Registered Storage Plugin")
	} else {
		fmt.Println("Failed to load storage plugin", error)
	}

	lis, error := net.Listen("tcp", s.serviceAddress)
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	gateway, error := s.loadGatewayPlugin()
	if error != nil {
		panic(error)
	}

	fmt.Println("Registered Gateway Plugin")

	// Start the gRPC server
	go (func() {
		fmt.Println(fmt.Sprintf("Services listening on: %s", s.serviceAddress))
		grpcServer.Serve(lis)
	})()

	// Start our child process
	// This will block until our child process is ready to accept incoming connections
	if s.childCommand != "" {
		s.startChildProcess()
	} else {
		fmt.Println("No Child Configured Specified, Skipping...")
	}

	// FIXME: Only do this in Gateway mode...
	// Otherwise always pass through to the provided child address
	// Start the Gateway Server

	// Start the gateway, this will provide us an entrypoint for
	// data ingress/egress to our userland code
	// The gateway should block the main thread but will
	// use this callback as a control mechanism
	fmt.Println("Starting Gateway")
	err := gateway.Start(func(request *gw.NitricRequest) *gw.NitricResponse {
		childUrl := fmt.Sprintf("http://%s", s.childAddress)

		httpRequest, error := http.NewRequest("POST", childUrl, bytes.NewReader(request.Payload))

		// There was an error creating the HTTP request
		if error != nil {
			// return an error to the Gateway
			return &gw.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(error.Error()),
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
		response, error := http.DefaultClient.Do(httpRequest)

		if error != nil {
			// there was an error calling the HTTP service
			return &gw.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(error.Error()),
				Status:  503,
			}
		}

		responseBody, error := ioutil.ReadAll(response.Body)

		if error != nil {
			// There was an error reading the http response
			return &gw.NitricResponse{
				Headers: map[string]string{"Content-Type": "text/plain"},
				Body:    []byte(error.Error()),
				Status:  503,
			}
		}

		headers := map[string]string{}
		for name, value := range response.Header {
			headers[name] = value[0]
		}

		// Pass the response back to the gateway
		return &gw.NitricResponse{
			Headers: headers,
			Body:    responseBody,
			Status:  response.StatusCode,
		}
	})
	// The gateway process has exited

	panic(err)
}

// Create a new Membrane server
func New(serviceAddress string, childAddress string, childCommand string, pluginDir string, eventingPlugin string, documentsPlugin string, storagePlugin string, gatewayPlugin string) (*Membrane, error) {
	return &Membrane{
		serviceAddress:      serviceAddress,
		childAddress:        childAddress,
		childCommand:        childCommand,
		pluginDir:           pluginDir,
		eventingPluginFile:  eventingPlugin,
		storagePluginFile:   storagePlugin,
		documentsPluginFile: documentsPlugin,
		gatewayPluginFile:   gatewayPlugin,
	}, nil
}
