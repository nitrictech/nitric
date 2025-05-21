package runtime

import (
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	pubsubpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type GrpcServer struct {
	pubsubpb.UnimplementedTopicsServer
	storagepb.UnimplementedStorageServer
}
