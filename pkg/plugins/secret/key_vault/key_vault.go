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

package key_vault_secret_service

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/armkeyvault"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	"github.com/nitric-dev/membrane/pkg/utils"
)

const DEFAULT_SUBSCRIPTION_ID = "subscription-id"
const DEFAULT_RESOURCE_GROUP = "resource-group"
const DEFAULT_VAULT_NAME = "vault-name"

// KeyVaultClient - iface that exposes utilized subset of generated KeyVaultSecretClient
// Used with gomock to assert create client -> service interaction in unit tests
type KeyVaultClient interface {
	Get(ctx context.Context, resourceGroupName string, vaultName string, secretName string, options *armkeyvault.SecretsGetOptions) (armkeyvault.SecretResponse, error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, secretName string, parameters armkeyvault.SecretCreateOrUpdateParameters, options *armkeyvault.SecretsCreateOrUpdateOptions) (armkeyvault.SecretResponse, error)
}
type keyVaultSecretService struct {
	secret.UnimplementedSecretPlugin
	client            KeyVaultClient
	resourceGroupName string
	vaultName         string
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

func validateSecretVersion(sec *secret.SecretVersion) error {
	if sec == nil {
		return fmt.Errorf("provide non-nil versioned secret")
	}
	if sec.Secret == nil {
		return fmt.Errorf("provide non-nil secret")
	}
	if len(sec.Secret.Name) == 0 {
		return fmt.Errorf("provide non-blank secret name")
	}
	if len(sec.Version) == 0 {
		return fmt.Errorf("provide non-blank secret version")
	}
	return nil
}

func (s *keyVaultSecretService) Put(sec *secret.Secret, val []byte) (*secret.SecretPutResponse, error) {
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.Put")

	if err := validateNewSecret(sec, val); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret",
			err,
		)
	}
	ctx := context.Background()
	stringVal := string(val[:])
	result, err := s.client.CreateOrUpdate(
		ctx,
		s.resourceGroupName,
		s.vaultName,
		sec.Name,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Value: &stringVal,
			},
		},
		nil,
	)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error putting secret",
			err,
		)
	}
	return &secret.SecretPutResponse{
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sec.Name,
			},
			Version: *result.Secret.Properties.SecretURIWithVersion,
		},
	}, nil
}

//GET https://{vaultBaseUrl}/secrets/{secret-name}/{secret-version}?api-version={api-version}
func (s *keyVaultSecretService) Access(sv *secret.SecretVersion) (*secret.SecretAccessResponse, error) {
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.Access")

	if err := validateSecretVersion(sv); err != nil {
		return nil, newErr(
			codes.Internal,
			"invalid secret version",
			err,
		)
	}
	ctx := context.Background()
	result, err := s.client.Get(
		ctx,
		s.resourceGroupName,
		s.vaultName,
		sv.Version,
		&armkeyvault.SecretsGetOptions{},
	)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to access secret",
			err,
		)
	}
	return &secret.SecretAccessResponse{
		// Return the original secret version payload
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sv.Secret.Name,
			},
			Version: *result.Secret.Properties.SecretURIWithVersion,
		},
		Value: []byte(*result.Secret.Properties.Value),
	}, nil
}

// New - Creates a new Nitric secret service with Azure Key Vault Provider
func New() (secret.SecretService, error) {
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.New")

	subscriptionId := utils.GetEnv("AZURE_SUBSCRIPTION_ID", DEFAULT_SUBSCRIPTION_ID)
	resouceGroup := utils.GetEnv("AZURE_RESOURCE_GROUP", DEFAULT_RESOURCE_GROUP)
	vaultName := utils.GetEnv("AZURE_VAULT_NAME", DEFAULT_VAULT_NAME)

	credentials, credentialsError := azidentity.NewDefaultAzureCredential(nil)
	if credentialsError != nil {
		return nil, newErr(
			codes.Internal,
			"azure credentials error",
			credentialsError,
		)
	}
	conn := armcore.NewDefaultConnection(credentials, nil)
	client := armkeyvault.NewSecretsClient(conn, subscriptionId)

	return &keyVaultSecretService{
		client:            client,
		resourceGroupName: resouceGroup,
		vaultName:         vaultName,
	}, nil
}
