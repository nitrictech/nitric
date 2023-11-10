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
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"google.golang.org/grpc/codes"

	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	azureutils "github.com/nitrictech/nitric/cloud/azure/runtime/utils"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	secretpb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
)

type KeyVaultClient interface {
	SetSecret(ctx context.Context, vaultBaseURL string, secretName string, parameters keyvault.SecretSetParameters) (result keyvault.SecretBundle, err error)
	GetSecret(ctx context.Context, vaultBaseURL string, secretName string, secretVersion string) (result keyvault.SecretBundle, err error)
}

// KeyVaultSecretService - Nitric Secret Service implementation for Azure Key Vault
type KeyVaultSecretService struct {
	client    KeyVaultClient
	vaultName string
}

var _ secretpb.SecretManagerServer = &KeyVaultSecretService{}

// versionIdFromUrl - Extracts a secret version ID from a full secret version URL
// the expected versionUrl format is https://{VAULT_NAME}.vault.azure.net/secrets/{SECRET_NAME}/{SECRET_VERSION}
func versionIdFromUrl(versionUrl string) string {
	urlParts := strings.Split(versionUrl, "/")
	return urlParts[len(urlParts)-1]
}

func (s *KeyVaultSecretService) Put(ctx context.Context, req *secretpb.SecretPutRequest) (*secretpb.SecretPutResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("KeyVaultSecretService.Put")
	stringVal := string(req.Value[:])

	result, err := s.client.SetSecret(
		ctx,
		fmt.Sprintf("https://%s.vault.azure.net", s.vaultName),
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

	return &secretpb.SecretPutResponse{
		SecretVersion: &secretpb.SecretVersion{
			Secret: &secretpb.Secret{
				Name: req.Secret.Name,
			},
			Version: versionIdFromUrl(*result.ID),
		},
	}, nil
}

func (s *KeyVaultSecretService) Access(ctx context.Context, req *secretpb.SecretAccessRequest) (*secretpb.SecretAccessResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("KeyVaultSecretService.Access")

	// Key vault will default to latest if an empty string is provided
	version := req.SecretVersion.Version
	if version == "latest" {
		version = ""
	}
	result, err := s.client.GetSecret(
		ctx,
		fmt.Sprintf("https://%s.vault.azure.net", s.vaultName), // https://myvault.vault.azure.net.
		req.SecretVersion.Secret.Name,
		version,
	)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to access secret",
			err,
		)
	}
	// Returned Secret ID: https://myvault.vault.azure.net/secrets/mysecret/11a536561da34d6b8b452d880df58f3a
	// Split to get the version
	return &secretpb.SecretAccessResponse{
		// Return the original secret version payload
		SecretVersion: &secretpb.SecretVersion{
			Secret: &secretpb.Secret{
				Name: req.SecretVersion.Secret.Name,
			},
			Version: versionIdFromUrl(*result.ID),
		},
		Value: []byte(*result.Value),
	}, nil
}

// New - Creates a new Nitric secret service with Azure Key Vault Provider
func New() (*KeyVaultSecretService, error) {
	vaultName := env.KVAULT_NAME.String()
	if len(vaultName) == 0 {
		return nil, fmt.Errorf("KVAULT_NAME not configured")
	}

	// Auth requires:
	// AZURE_TENANT_ID: Your Azure tenant ID
	// AZURE_CLIENT_ID: Your Azure client ID. This will be an app ID from your AAD.
	spt, err := azureutils.GetServicePrincipalToken(azure.PublicCloud.ResourceIdentifiers.KeyVault)
	if err != nil {
		return nil, err
	}

	client := keyvault.New()
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	return &KeyVaultSecretService{
		client:    client,
		vaultName: vaultName,
	}, nil
}

func NewWithClient(client KeyVaultClient) *KeyVaultSecretService {
	return &KeyVaultSecretService{
		client:    client,
		vaultName: "localvault",
	}
}
