package simulation

import (
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/lipgloss"
	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

type SimulationServer struct {
	fs     afero.Fs
	appDir string
	storagepb.UnimplementedStorageServer
	pubsubpb.UnimplementedPubsubServer
}

const DEFAULT_SERVER_PORT = "50051"

var (
	styledNitric = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(13)).Render("nitric")
)

func (s *SimulationServer) Start() error {
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

	fmt.Printf("\n%s starting on: %s\n\n", styledNitric, addr)
	return srv.Serve(lis)
}

type SimulationServerOption func(*SimulationServer)

func WithAppDirectory(appDir string) SimulationServerOption {
	return func(s *SimulationServer) {
		s.appDir = appDir
	}
}

func NewSimulationServer(fs afero.Fs, opts ...SimulationServerOption) *SimulationServer {
	simServer := &SimulationServer{
		fs:     fs,
		appDir: ".",
	}

	for _, o := range opts {
		o(simServer)
	}

	return simServer
}
