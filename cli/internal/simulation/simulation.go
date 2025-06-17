package simulation

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/netx"
	"github.com/nitrictech/nitric/cli/internal/simulation/service"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

type SimulationServer struct {
	fs      afero.Fs
	appDir  string
	appSpec *schema.Application
	storagepb.UnimplementedStorageServer
	pubsubpb.UnimplementedPubsubServer

	services map[string]*service.ServiceSimulation
}

const DEFAULT_SERVER_PORT = "50051"

var nitric = style.Purple(icons.Lightning + " Nitric")

func nitricIntro(addr string, dashUrl string, appSpec *schema.Application) string {
	version := version.GetShortVersion()

	intro := fmt.Sprintf("%s %s\n- App: %s\n- Addr: %s\n- Dashboard: %s\n", nitric, style.Gray(version), appSpec.Name, addr, dashUrl)

	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(intro)
}

func (s *SimulationServer) startNitricApis() error {
	srv := grpc.NewServer()

	storagepb.RegisterStorageServer(srv, s)
	pubsubpb.RegisterPubsubServer(srv, s)

	host := os.Getenv("NITRIC_HOST")
	port := os.Getenv("NITRIC_PORT")
	if port == "" {
		port = DEFAULT_SERVER_PORT
	}

	addr := net.JoinHostPort(host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Println(nitricIntro(addr, "https://app.nitric.io/dashboard", s.appSpec))
	go srv.Serve(lis)

	return nil
}

const (
	ENTRYPOINT_MIN_PORT = 3000
	ENTRYPOINT_MAX_PORT = 3999
)

var greenCheck = style.Green(icons.Check)

func (s *SimulationServer) startEntrypoints(services map[string]*service.ServiceSimulation) error {
	// TODO: Possibly handle on multiple ports
	serviceProxies := map[string]*httputil.ReverseProxy{}
	for serviceName, service := range services {
		url := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", service.GetPort()),
		}

		serviceProxies[serviceName] = httputil.NewSingleHostReverseProxy(url)
	}

	for entrypointName, entrypoint := range s.appSpec.GetEntrypointIntents() {
		// Reserve a port
		reservedPort, err := netx.GetNextPort(netx.MinPort(ENTRYPOINT_MIN_PORT), netx.MaxPort(ENTRYPOINT_MAX_PORT))
		if err != nil {
			return err
		}

		router := http.NewServeMux()

		for route, target := range entrypoint.Routes {
			// TODO: Handle other target types
			targetProxy := serviceProxies[target.TargetName]
			router.Handle(route, http.StripPrefix(strings.TrimSuffix(route, "/"), targetProxy))
		}

		go http.ListenAndServe(fmt.Sprintf(":%d", reservedPort), router)

		fmt.Printf("%s Starting %s http://localhost:%d\n", greenCheck, styledName(entrypointName, style.Orange), reservedPort)
	}

	return nil
}

func (s *SimulationServer) startServices(output io.Writer) (<-chan service.ServiceEvent, error) {
	serviceIntents := s.appSpec.GetServiceIntents()

	eventChans := []<-chan service.ServiceEvent{}

	for serviceName, serviceIntent := range serviceIntents {
		port, err := netx.GetNextPort()
		if err != nil {
			return nil, err
		}

		simulatedService, eventChan, err := service.NewServiceSimulation(serviceName, *serviceIntent, port)
		if err != nil {
			return nil, err
		}

		eventChans = append(eventChans, eventChan)
		s.services[serviceName] = simulatedService

		fmt.Fprintf(output, "%s Starting %s\n", greenCheck, styledName(serviceName, style.Teal))

	}

	for _, simulatedService := range s.services {
		go func() {
			err := simulatedService.Start(true)
			if err != nil {
				// TODO: Handle the start error
			}
		}()
	}

	// Combine all of the events
	combinedEventsChan := lo.FanIn(100, eventChans...)

	return combinedEventsChan, nil
}

func (s *SimulationServer) handleServiceOutputs(output io.Writer, events <-chan service.ServiceEvent) {

	serviceWriters := make(map[string]io.Writer, len(s.appSpec.GetServiceIntents()))
	for serviceName := range s.appSpec.GetServiceIntents() {
		serviceWriters[serviceName] = NewPrefixWriter(styledName(serviceName, style.Teal)+" ", output)
	}

	for {
		event := <-events

		if event.Output != nil {
			// write some kind of output for that service
			if writer, ok := serviceWriters[event.GetName()]; ok {
				writer.Write(event.Content)
			} else {
				log.Fatalf("failed to retrieve output writer for service %s", event.GetName())
			}
		}

		if event.PreviousStatus != event.GetStatus() {
			// If the status has changed write about it
		}

		// Handle output logging on the channels
	}
}

var styledNames = map[string]string{}

func styledName(name string, styleFunc func(...string) string) string {
	_, ok := styledNames[name]
	if !ok {
		styledNames[name] = styleFunc(fmt.Sprintf("[%s]", name))
	}

	return styledNames[name]
}

func (s *SimulationServer) Start(output io.Writer) error {
	err := s.startNitricApis()
	if err != nil {
		return err
	}

	var svcEvents <-chan service.ServiceEvent

	if len(s.appSpec.GetServiceIntents()) > 0 {
		fmt.Fprintf(output, "%s\n\n", style.Teal("Services"))
		svcEvents, err = s.startServices(output)
		if err != nil {
			return err
		}
		fmt.Fprint(output, "\n")
	}

	if len(s.appSpec.GetEntrypointIntents()) > 0 {
		fmt.Fprintf(output, "%s\n\n", style.Orange("Entrypoints"))
		err = s.startEntrypoints(s.services)
		if err != nil {
			return err
		}
		fmt.Fprint(output, "\n")
	}

	// block on handling service outputs for now
	s.handleServiceOutputs(output, svcEvents)

	return nil
}

type SimulationServerOption func(*SimulationServer)

func WithAppDirectory(appDir string) SimulationServerOption {
	return func(s *SimulationServer) {
		s.appDir = appDir
	}
}

func NewSimulationServer(fs afero.Fs, appSpec *schema.Application, opts ...SimulationServerOption) *SimulationServer {
	simServer := &SimulationServer{
		fs:       fs,
		appSpec:  appSpec,
		appDir:   ".",
		services: make(map[string]*service.ServiceSimulation),
	}

	for _, o := range opts {
		o(simServer)
	}

	return simServer
}
