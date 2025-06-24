package awss3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	"github.com/iancoleman/strcase"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/nitrictech/nitric/runtime/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type s3Storage struct {
	storagepb.UnimplementedStorageServer
	nitricStackId string
	s3Client      *s3.Client
	preSignClient *s3.PresignClient
}

func (s *s3Storage) getS3BucketName(bucket string) string {
	normalizedBucketName := strcase.ToKebab(bucket)

	// We want to build the bucket name from convention
	return fmt.Sprintf("%s-%s", s.nitricStackId, normalizedBucketName)
}

func isS3AccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "S3" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

func (s *s3Storage) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)

	resp, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to read file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error reading file: %v", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &storagepb.StorageReadResponse{
		Body: bodyBytes,
	}, nil
}

func detectContentType(filename string, content []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}

	return http.DetectContentType(content)
}

func (s *s3Storage) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)
	contentType := detectContentType(req.Key, req.Body)

	if _, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Body:        bytes.NewReader(req.Body),
		ContentType: &contentType,
		Key:         aws.String(req.Key),
	}); err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to write file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error writing file: %v", err)
	}

	return &storagepb.StorageWriteResponse{}, nil
}

func (s *s3Storage) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)

	if _, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(req.Key),
	}); err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to delete file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error deleting file: %v", err)
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

func (s *s3Storage) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)

	switch req.Operation {
	case storagepb.StoragePreSignUrlRequest_READ:
		response, err := s.preSignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(req.Key),
		}, s3.WithPresignExpires(req.Expiry.AsDuration()))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate signed READ URL: %v", err)
		}
		return &storagepb.StoragePreSignUrlResponse{
			Url: response.URL,
		}, err
	case storagepb.StoragePreSignUrlRequest_WRITE:
		req, err := s.preSignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(req.Key),
		}, s3.WithPresignExpires(req.Expiry.AsDuration()))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate signed WRITE URL: %v", err)
		}
		return &storagepb.StoragePreSignUrlResponse{
			Url: req.URL,
		}, err
	default:
		return nil, status.Errorf(codes.Unimplemented, "requested operation not supported for pre-signed AWS S3 URLs")
	}
}

func (s *s3Storage) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)

	objects, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(req.Prefix),
	})
	if err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to list files, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error listing files: %v", err)
	}

	files := make([]*storagepb.Blob, 0, len(objects.Contents))
	for _, o := range objects.Contents {
		files = append(files, &storagepb.Blob{
			Key: *o.Key,
		})
	}

	return &storagepb.StorageListBlobsResponse{
		Blobs: files,
	}, nil
}

func (s *s3Storage) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	bucketName := s.getS3BucketName(req.BucketName)

	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to check if file exists, this may be due to a missing permissions request in your code.")
		}

		return &storagepb.StorageExistsResponse{
			Exists: false,
		}, nil
	}

	return &storagepb.StorageExistsResponse{
		Exists: true,
	}, nil
}

func Plugin() (storage.Storage, error) {
	nitricStackId := os.Getenv("NITRIC_STACK_ID")
	if nitricStackId == "" {
		return nil, fmt.Errorf("NITRIC_STACK_ID is not set")
	}

	// Create a new aws s3 client
	// TODO: Set the region
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	preSignClient := s3.NewPresignClient(s3Client)

	return &s3Storage{
		s3Client:      s3Client,
		preSignClient: preSignClient,
		nitricStackId: nitricStackId,
	}, nil
}
