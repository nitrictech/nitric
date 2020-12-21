package membrane

import (
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/eventing"
)

// GRPC Interface for registered Nitric Eventing Plugins
type EventingServer struct {
	pb.UnimplementedEventingServer
}
