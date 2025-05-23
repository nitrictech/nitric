package runtime

import (
	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
)

type GrpcServer struct {
	pubsubpb.UnimplementedPubsubServer
	storagepb.UnimplementedStorageServer
}
