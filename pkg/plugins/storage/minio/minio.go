// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package minio_storage_service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/nitrictech/nitric/pkg/plugins/storage"
	s3_service "github.com/nitrictech/nitric/pkg/plugins/storage/s3"
	"github.com/nitrictech/nitric/pkg/utils"
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

func nameSelector(nitricName string) (*string, error) {
	return &nitricName, nil
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
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("error creating new session")
	}

	s3Client := s3.New(newSession)

	return s3_service.NewWithClient(nil, s3Client, s3_service.WithSelector(nameSelector))
}
