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

package secret

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/secretsmanageriface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	secretpb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
)

type SecretsManagerSecretService struct {
	client   secretsmanageriface.SecretsManagerAPI
	provider resource.AwsResourceProvider
}

var _ secretpb.SecretManagerServer = &SecretsManagerSecretService{}

func (s *SecretsManagerSecretService) validateNewSecret(sec *secretpb.Secret, val []byte) error {
	if sec == nil {
		return fmt.Errorf("secret cannot be empty")
	}
	if len(sec.Name) == 0 {
		return fmt.Errorf("secret name cannot be empty")
	}
	if len(val) == 0 {
		return fmt.Errorf("secret value cannot be empty")
	}

	return nil
}

// getSecretArn - Retrieve the ARN for a given secret name
func (s *SecretsManagerSecretService) getSecretArn(ctx context.Context, sec string) (string, error) {
	secrets, err := s.provider.GetResources(ctx, resource.AwsResource_Secret)
	if err != nil {
		return "", fmt.Errorf("error retrieving secrets list: %w", err)
	}

	if secret, ok := secrets[sec]; ok {
		return secret.ARN, nil
	}

	return "", fmt.Errorf("secret %s does not exist", sec)
}

func isSecretsManagerAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "Secrets Manager" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

// Put - Store a new secret value
func (s *SecretsManagerSecretService) Put(ctx context.Context, req *secretpb.SecretPutRequest) (*secretpb.SecretPutResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SecretManagerSecretService.Put")

	if err := s.validateNewSecret(req.Secret, req.Value); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret",
			err,
		)
	}

	secretId, err := s.getSecretArn(ctx, req.Secret.Name)
	if err != nil {
		return nil, newErr(codes.NotFound, "secret not found", err)
	}

	result, err := s.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretId),
		SecretBinary: req.Value,
	})
	if err != nil {
		if isSecretsManagerAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to put secret value, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(codes.Unknown, "unable to put secret", err)
	}

	return &secretpb.SecretPutResponse{
		SecretVersion: &secretpb.SecretVersion{
			Secret: &secretpb.Secret{
				Name: req.Secret.Name,
			},
			Version: *result.VersionId,
		},
	}, nil
}

// Access - Retrieve a secret value
func (s *SecretsManagerSecretService) Access(ctx context.Context, req *secretpb.SecretAccessRequest) (*secretpb.SecretAccessResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SecretManagerSecretService.Access")

	if len(req.SecretVersion.GetSecret().Name) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"secret name cannot be blank or empty",
			nil,
		)
	}

	if len(req.SecretVersion.Version) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"secret version cannot be blank or empty",
			nil,
		)
	}

	secretArn, err := s.getSecretArn(ctx, req.SecretVersion.GetSecret().Name)
	if err != nil {
		return nil, newErr(codes.NotFound, "secret not found", err)
	}

	// Build the request to get the secret
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	}

	// If the requested version is latest then we want
	// to exclude the version from input
	if strings.ToLower(req.SecretVersion.Version) != "latest" {
		input.VersionId = aws.String(req.SecretVersion.Version)
	}

	result, err := s.client.GetSecretValue(ctx, input)
	if err != nil {
		if isSecretsManagerAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to access secret value, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(
			codes.Unknown,
			"failed to retrieve secret value",
			err,
		)
	}

	returnValue := result.SecretBinary

	if returnValue == nil && result.SecretString != nil {
		returnValue = []byte(*result.SecretString)
	}

	return &secretpb.SecretAccessResponse{
		SecretVersion: &secretpb.SecretVersion{
			Secret:  req.SecretVersion.Secret,
			Version: *result.VersionId,
		},
		Value: returnValue,
	}, nil
}

// Gets a new Secrets Manager Client
func New(provider resource.AwsResourceProvider) (*SecretsManagerSecretService, error) {
	awsRegion := env.AWS_REGION.String()

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := secretsmanager.NewFromConfig(cfg)

	return &SecretsManagerSecretService{
		client:   client,
		provider: provider,
	}, nil
}
