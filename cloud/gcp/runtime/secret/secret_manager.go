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

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"

	ifaces_gcloud_secret "github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_secret"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	secretpb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
)

type SecretManagerSecretService struct {
	client    ifaces_gcloud_secret.SecretManagerClient
	projectId string
	stackName string
	cache     map[string]string
}

var _ secretpb.SecretManagerServer = &SecretManagerSecretService{}

func validateNewSecret(sec *secretpb.Secret, val []byte) error {
	if sec == nil {
		return fmt.Errorf("provide non-nil secret")
	}
	if len(sec.Name) == 0 {
		return fmt.Errorf("provide non-blank secret name")
	}
	if len(val) == 0 {
		return fmt.Errorf("provide non-blank secret value")
	}

	return nil
}

func (s *SecretManagerSecretService) getParentName() string {
	return fmt.Sprintf("projects/%s", s.projectId)
}

func (s *SecretManagerSecretService) buildSecretVersionName(ctx context.Context, sv *secretpb.SecretVersion) (string, error) {
	if sv == nil {
		return "", fmt.Errorf("provide non-nil secret version")
	}

	if sv.Secret == nil {
		return "", fmt.Errorf("provide non-nil secret")
	}

	if len(sv.Secret.Name) == 0 {
		return "", fmt.Errorf("provide non-blank name")
	}

	if len(sv.Version) == 0 {
		return "", fmt.Errorf("provide non-blank version")
	}

	parent, inCache := s.cache[sv.Secret.Name]
	if !inCache {
		realSec, err := s.getSecret(ctx, sv.Secret)
		if err != nil {
			return "", err
		}

		parent = realSec.Name
	}

	return fmt.Sprintf("%s/versions/%s", parent, sv.Version), nil
}

// ensure a secret container exists for storing secret versions
func (s *SecretManagerSecretService) getSecret(ctx context.Context, sec *secretpb.Secret) (*secretmanagerpb.Secret, error) {
	iter := s.client.ListSecrets(ctx, &secretmanagerpb.ListSecretsRequest{
		Parent: s.getParentName(),
		Filter: fmt.Sprintf("labels.%s=%s", tags.GetResourceNameKey(env.GetNitricStackID()), sec.Name),
	})

	result, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil, status.Error(grpccodes.NotFound, "secret not found")
	}

	if err != nil {
		return nil, err
	}

	s.cache[sec.Name] = result.Name

	return result, nil
}

// Put - Creates a new secret if one doesn't exist, or just adds a new secret version
func (s *SecretManagerSecretService) Put(ctx context.Context, req *secretpb.SecretPutRequest) (*secretpb.SecretPutResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SecretManagerSecretService.Put")

	if err := validateNewSecret(req.Secret, req.Value); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret",
			err,
		)
	}

	// ensure the secret container exists...
	parentSec, err := s.getSecret(ctx, req.Secret)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"error ensuring secret container exists",
			err,
		)
	}

	verResult, err := s.client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: parentSec.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: req.Value,
		},
	})
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this secret?", err)
		}

		return nil, newErr(
			codes.Internal,
			"failed to add new secret version",
			err,
		)
	}

	versionStringParts := strings.Split(verResult.Name, "/")
	version := versionStringParts[len(versionStringParts)-1]

	return &secretpb.SecretPutResponse{
		SecretVersion: &secretpb.SecretVersion{
			Secret: &secretpb.Secret{
				Name: req.Secret.Name,
			},
			Version: version,
		},
	}, nil
}

// Get - Retrieves a secret given a name and a version
func (s *SecretManagerSecretService) Access(ctx context.Context, req *secretpb.SecretAccessRequest) (*secretpb.SecretAccessResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SecretManagerSecretService.Access")

	fullName, err := s.buildSecretVersionName(ctx, req.SecretVersion)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret version",
			err,
		)
	}

	secretRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fullName,
	}

	result, err := s.client.AccessSecretVersion(ctx, secretRequest)
	if err != nil {
		errStatus, _ := status.FromError(err)

		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this secret?", err)
		}

		return nil, newErr(
			codes.Internal,
			"failed to access secret version",
			err,
		)
	}

	return &secretpb.SecretAccessResponse{
		// Return the original secret version payload
		SecretVersion: req.SecretVersion,
		Value:         result.Payload.GetData(),
	}, nil
}

// New - Creates a new Nitric secret service with GCP Secret Manager provider
func New() (*SecretManagerSecretService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, secretmanager.DefaultAuthScopes()...)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}

	client, clientError := ifaces_gcloud_secret.NewClient(ctx)
	if clientError != nil {
		return nil, fmt.Errorf("secret manager client error: %w", clientError)
	}

	return &SecretManagerSecretService{
		client:    client,
		projectId: credentials.ProjectID,
		stackName: commonenv.NITRIC_STACK_ID.String(),
		cache:     make(map[string]string),
	}, nil
}
