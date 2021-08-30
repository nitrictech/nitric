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
	"strings"

	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	"github.com/nitric-dev/membrane/pkg/utils"
)

type KeyVaultClient interface {
	SetSecret(ctx context.Context, vaultBaseURL string, secretName string, parameters keyvault.SecretSetParameters) (result keyvault.SecretBundle, err error)
	GetSecret(ctx context.Context, vaultBaseURL string, secretName string, secretVersion string) (result keyvault.SecretBundle, err error)
}

type KeyVaultSecretService struct {
	secret.UnimplementedSecretPlugin
	client    KeyVaultClient
	vaultName string
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

func (s *KeyVaultSecretService) Put(sec *secret.Secret, val []byte) (*secret.SecretPutResponse, error) {
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.Put")

	if err := validateNewSecret(sec, val); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid secret",
			err,
		)
	}
	stringVal := string(val[:])

	result, err := s.client.SetSecret(
		context.Background(),
		fmt.Sprintf("https://%s.vault.azure.net", s.vaultName), //https://myvault.vault.azure.net.
		sec.Name,
		keyvault.SecretSetParameters{
			Value: &stringVal,
		},
	)

	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error putting secret",
			err,
		)
	}
	//Returned Secret ID: https://myvault.vault.azure.net/secrets/{SECRET_NAME}/{SECRET_VERSION}
	//Split to get the version
	versionID := strings.Split(*result.ID, "/")

	return &secret.SecretPutResponse{
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sec.Name,
			},
			Version: versionID[len(versionID)-1],
		},
	}, nil
}

func (s *KeyVaultSecretService) Access(sv *secret.SecretVersion) (*secret.SecretAccessResponse, error) {
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.Access")

	if err := validateSecretVersion(sv); err != nil {
		return nil, newErr(
			codes.Internal,
			"invalid secret version",
			err,
		)
	}

	//Key vault will default to latest if an empty string is provided
	version := sv.Version
	if version == "latest" {
		version = ""
	}
	result, err := s.client.GetSecret(
		context.Background(),
		fmt.Sprintf("https://%s.vault.azure.net", s.vaultName), //https://myvault.vault.azure.net.
		sv.Secret.Name,
		version,
	)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to access secret",
			err,
		)
	}
	//Returned Secret ID: https://myvault.vault.azure.net/secrets/mysecret/11a536561da34d6b8b452d880df58f3a
	//Split to get the version
	versionID := strings.Split(*result.ID, "/")
	return &secret.SecretAccessResponse{
		// Return the original secret version payload
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sv.Secret.Name,
			},
			Version: versionID[len(versionID)-1],
		},
		Value: []byte(*result.Value),
	}, nil
}

// New - Creates a new Nitric secret service with Azure Key Vault Provider
func New() (secret.SecretService, error) {
	vaultName := utils.GetEnv("KVAULT_NAME", "")
	if len(vaultName) == 0 {
		return nil, fmt.Errorf("KVAULT_NAME not configured")
	}

	//Auth requires:
	//AZURE_TENANT_ID: Your Azure tenant ID
	//AZURE_CLIENT_ID: Your Azure client ID. This will be an app ID from your AAD.
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	client := keyvault.New()
	client.Authorizer = authorizer

	return &KeyVaultSecretService{
		client:    client,
		vaultName: vaultName,
	}, nil
}

func NewWithClient(client KeyVaultClient) secret.SecretService {
	return &KeyVaultSecretService{
		client:    client,
		vaultName: "localvault",
	}
}
