package storage_service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
	s3_service "github.com/nitric-dev/membrane/pkg/plugins/storage/s3"
	"github.com/nitric-dev/membrane/pkg/utils"
)

const (
	MINIO_ENDPOINT_ENV   = "MINIO_ENDPOINT"
	MINIO_ACCESS_KEY_ENV = "MINIO_ACCESS_KEY"
	MINIO_SECRET_KEY_ENV = "MINIO_SECRET_KEY"
)

type minioConfig struct {
	endpoint  string
	accessKey string
	secretKey string
}

func configFromEnv() (*minioConfig, error) {
	endpoint := utils.GetEnv(MINIO_ENDPOINT_ENV, "")
	accKey := utils.GetEnv(MINIO_ACCESS_KEY_ENV, "")
	secKey := utils.GetEnv(MINIO_SECRET_KEY_ENV, "")

	configErrors := make([]error, 0)

	if endpoint == "" {
		configErrors = append(configErrors, fmt.Errorf("%s not configured", MINIO_ENDPOINT_ENV))
	}

	if accKey == "" {
		configErrors = append(configErrors, fmt.Errorf("%s not configured", MINIO_ACCESS_KEY_ENV))
	}

	if secKey == "" {
		configErrors = append(configErrors, fmt.Errorf("%s not configured", MINIO_SECRET_KEY_ENV))
	}

	if len(configErrors) > 0 {
		return nil, fmt.Errorf("configuration errors: %v", configErrors)
	}

	return &minioConfig{
		endpoint:  endpoint,
		accessKey: accKey,
		secretKey: secKey,
	}, nil
}

func New() (storage.StorageService, error) {

	conf, err := configFromEnv()

	if err != nil {
		return nil, err
	}

	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.accessKey, conf.secretKey, ""),
		Endpoint:         aws.String(conf.endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	return s3_service.NewWithClient(s3Client)
}
