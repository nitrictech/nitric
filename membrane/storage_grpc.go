package membrane

import (
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/storage"
)

// GRPC Interface for registered Nitric Storage Plugins
type StorageServer struct {
	pb.UnimplementedStorageServer
}
