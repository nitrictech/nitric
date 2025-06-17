package simulation

import (
	"context"
	"os"
	"path/filepath"

	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const localNitricBucketDir = "./.nitric/buckets"

func bucketDirFromAppDir(appDir string) string {
	return filepath.Join(appDir, localNitricBucketDir)
}

func ensureBucketDir(appDir, bucketName string) (string, error) {
	bucketPath := filepath.Join(bucketDirFromAppDir(appDir), bucketName)
	err := os.MkdirAll(bucketPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return bucketPath, nil
}

func (s *SimulationServer) bucketFilepath(bucketName, blobName string) (string, error) {
	blobPath := filepath.Join(s.appDir, localNitricBucketDir, bucketName, blobName)

	err := os.MkdirAll(filepath.Dir(blobPath), os.ModePerm)
	if err != nil {
		return "", err
	}

	return blobPath, nil
}

func (s *SimulationServer) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	path, err := s.bucketFilepath(req.BucketName, req.Key)
	if err != nil {

	}

	err = s.fs.Remove(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

func (s *SimulationServer) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	path, err := s.bucketFilepath(req.BucketName, req.Key)
	if err != nil {
		return nil, err
	}

	exists, err := afero.Exists(s.fs, path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageExistsResponse{
		Exists: exists,
	}, nil
}

func (s *SimulationServer) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	path := filepath.Join(s.appDir, localNitricBucketDir, req.BucketName)
	files, err := afero.ReadDir(s.fs, path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	blobs := make([]*storagepb.Blob, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		blobs = append(blobs, &storagepb.Blob{
			Key: file.Name(),
		})
	}

	return &storagepb.StorageListBlobsResponse{
		Blobs: blobs,
	}, nil
}

func (s *SimulationServer) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	// TODO: Host a local HTTP server to serve the files
	return nil, status.Errorf(codes.Unimplemented, "PreSignUrl is not implemented in simulation mode")
}

func (s *SimulationServer) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	path, err := s.bucketFilepath(req.BucketName, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	data, err := afero.ReadFile(s.fs, path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, status.Errorf(codes.NotFound, "Blob %s not found in bucket %s", req.Key, req.BucketName)
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &storagepb.StorageReadResponse{
		Body: data,
	}, nil
}

func (s *SimulationServer) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	path, err := s.bucketFilepath(req.BucketName, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	err = afero.WriteFile(s.fs, path, req.Body, os.ModePerm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageWriteResponse{}, nil
}
