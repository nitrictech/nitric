package runtime

import (
	"context"

	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/nitrictech/nitric/server/runtime/storage"
)

func (s *GrpcServer) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.Read(ctx, req)
}

func (s *GrpcServer) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.Write(ctx, req)
}

func (s *GrpcServer) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.Delete(ctx, req)
}

func (s *GrpcServer) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.PreSignUrl(ctx, req)
}

func (s *GrpcServer) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.ListBlobs(ctx, req)
}

func (s *GrpcServer) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	plugin := storage.GetPlugin("default")

	return plugin.Exists(ctx, req)
}
