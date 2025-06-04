package awss3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/nitrictech/nitric/server/runtime/storage"
)

type s3Storage struct {
	storagepb.UnimplementedStorageServer
	s3Client *s3.Client
}

func (s *s3Storage) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *s3Storage) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *s3Storage) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *s3Storage) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *s3Storage) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *s3Storage) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func Plugin() (storage.Storage, error) {
	// Create a new aws s3 client
	// TODO: Set the region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return &s3Storage{
		s3Client: s3Client,
	}, nil
}
