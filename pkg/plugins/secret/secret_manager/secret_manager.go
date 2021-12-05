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

package secret_manager_secret_service

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/secret"
	"golang.org/x/oauth2/google"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	pbcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SecretManagerClient - iface that exposes utilized subset of generated SecretManagerServiceClient
// Used with gomock to assert create client -> service interaction in unit tests
type SecretManagerClient interface {
	AccessSecretVersion(context.Context, *secretmanagerpb.AccessSecretVersionRequest, ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	AddSecretVersion(context.Context, *secretmanagerpb.AddSecretVersionRequest, ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	CreateSecret(context.Context, *secretmanagerpb.CreateSecretRequest, ...gax.CallOption) (*secretmanagerpb.Secret, error)
	GetSecret(context.Context, *secretmanagerpb.GetSecretRequest, ...gax.CallOption) (*secretmanagerpb.Secret, error)
	UpdateSecret(context.Context, *secretmanagerpb.UpdateSecretRequest, ...gax.CallOption) (*secretmanagerpb.Secret, error)
}

type secretManagerSecretService struct {
	secret.UnimplementedSecretPlugin
	client    SecretManagerClient
	projectId string
}

func validateNewSecret(sec *secret.Secret, val []byte) error {
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

func (s *secretManagerSecretService) getParentName() string {
	return fmt.Sprintf("projects/%s", s.projectId)
}

func (s *secretManagerSecretService) buildSecretName(sec *secret.Secret) (string, error) {
	if len(sec.Name) == 0 {
		return "", fmt.Errorf("provide non-blank name")
	}

	return fmt.Sprintf("%s/secrets/%s", s.getParentName(), sec.Name), nil
}

func (s *secretManagerSecretService) buildSecretVersionName(sv *secret.SecretVersion) (string, error) {
	parent, err := s.buildSecretName(sv.Secret)

	if err != nil {
		return "", err
	}

	if len(sv.Version) == 0 {
		return "", fmt.Errorf("provide non-blank version")
	}

	return fmt.Sprintf("%s/versions/%s", parent, sv.Version), nil
}

// ensure a secret container exists for storing secret versions
func (s *secretManagerSecretService) ensureSecret(sec *secret.Secret) (*secretmanagerpb.Secret, error) {
	secName, err := s.buildSecretName(sec)

	if err != nil {
		return nil, err
	}

	getReq := &secretmanagerpb.GetSecretRequest{
		Name: secName,
	}

	result, err := s.client.GetSecret(context.TODO(), getReq)

	if err != nil {
		// check error status, if it was an RPC NOT_FOUND error then continue
		if s, ok := status.FromError(err); ok && s.Code() != pbcodes.NotFound {
			return nil, err
		} else if !ok {
			// It wasn't an RPC error so return
			return nil, err
		}
	}

	if result == nil {
		// Creates the secret container
		secReq := &secretmanagerpb.CreateSecretRequest{
			Parent:   s.getParentName(),
			SecretId: sec.Name,
			Secret: &secretmanagerpb.Secret{
				Replication: &secretmanagerpb.Replication{
					Replication: &secretmanagerpb.Replication_Automatic_{
						Automatic: &secretmanagerpb.Replication_Automatic{},
					},
				},
			},
		}

		result, err = s.client.CreateSecret(context.TODO(), secReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create new secret: %v", err)
		}
	}

	return result, nil
}

// Put - Creates a new secret if one doesn't exist, or just adds a new secret version
func (s *secretManagerSecretService) Put(sec *secret.Secret, val []byte) (*secret.SecretPutResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SecretManagerSecretService.Put",
		map[string]interface{}{
			"secret": sec,
		},
	)

	if err := validateNewSecret(sec, val); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret",
			err,
		)
	}

	// ensure the secret container exists...
	parentSec, err := s.ensureSecret(sec)

	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error ensuring secret container exists",
			err,
		)
	}

	verResult, err := s.client.AddSecretVersion(context.TODO(), &secretmanagerpb.AddSecretVersionRequest{
		Parent: parentSec.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: val,
		},
	})

	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to add new secret version",
			err,
		)
	}

	versionStringParts := strings.Split(verResult.Name, "/")
	version := versionStringParts[len(versionStringParts)-1]

	return &secret.SecretPutResponse{
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sec.Name,
			},
			Version: version,
		},
	}, nil
}

// Get - Retrieves a secret given a name and a version
func (s *secretManagerSecretService) Access(sv *secret.SecretVersion) (*secret.SecretAccessResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SecretManagerSecretService.Access",
		map[string]interface{}{
			"version": sv,
		},
	)

	fullName, err := s.buildSecretVersionName(sv)

	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret version",
			err,
		)
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fullName,
	}

	result, err := s.client.AccessSecretVersion(context.TODO(), req)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to access secret version",
			err,
		)
	}

	return &secret.SecretAccessResponse{
		// Return the original secret version payload
		SecretVersion: sv,
		Value:         result.Payload.GetData(),
	}, nil
}

// New - Creates a new Nitric secret service with GCP Secret Manager provider
func New() (secret.SecretService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, secretmanager.DefaultAuthScopes()...)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := secretmanager.NewClient(ctx)
	if clientError != nil {
		return nil, fmt.Errorf("secret manager client error: %v", clientError)
	}

	return &secretManagerSecretService{
		client:    client,
		projectId: credentials.ProjectID,
	}, nil
}
