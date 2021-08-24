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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault/keyvaultapi"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	"github.com/nitric-dev/membrane/pkg/utils"
)

type keyVaultSecretService struct {
	secret.UnimplementedSecretPlugin
	client      keyvaultapi.BaseClientAPI
	accessToken AzureAccessToken
	vaultName   string
}

type AzureAccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

func GetToken(tenantId string, clientId string, clientSecret string) (AzureAccessToken, error) {
	requestAccessTokenUri := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantId)
	requestBody := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"resource":      {"https://management.azure.com/"},
	}
	resp, err := http.PostForm(requestAccessTokenUri, requestBody)
	if err != nil {
		return AzureAccessToken{}, err
	}

	var result AzureAccessToken

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
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
	ctx.Value(map[string]string{
		"Authorization": s.accessToken.TokenType + " " + s.accessToken.AccessToken,
	})
	stringVal := string(val[:])

	result, err := s.client.SetSecret(
		ctx,
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
	ctx.Value(map[string]string{
		"Authorization": fmt.Sprintf("%s %s", s.accessToken.TokenType, s.accessToken.AccessToken),
	})

	//Key vault will default to latest if an empty string is provided
	version := sv.Version
	if version == "latest" {
		version = ""
	}
	result, err := s.client.GetSecret(
		ctx,
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
	newErr := errors.ErrorsWithScope("KeyVaultSecretService.New")

	subscriptionId := utils.GetEnv("AZURE_SUBSCRIPTION_ID", "")
	vaultName := utils.GetEnv("AZURE_VAULT_NAME", "")
	tenantId := utils.GetEnv("AZURE_TENANT_ID", "")
	clientId := utils.GetEnv("AZURE_CLIENT_ID", "")
	clientSecret := utils.GetEnv("AZURE_CLIENT_SECRET", "")

	if len(tenantId) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"AZURE_TENANT_ID not configured",
			fmt.Errorf(""),
		)
	}

	if len(clientId) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"AZURE_CLIENT_ID not configured",
			fmt.Errorf(""),
		)
	}

	if len(clientSecret) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"AZURE_CLIENT_SECRET not configured",
			fmt.Errorf(""),
		)
	}

	if len(subscriptionId) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"AZURE_SUBSCRIPTION_ID not configured",
			fmt.Errorf(""),
		)
	}
	if len(vaultName) == 0 {
		return nil, newErr(
			codes.InvalidArgument,
			"AZURE_VAULT_NAME not configured",
			fmt.Errorf(""),
		)
	}

	client := keyvault.New()
	token, err := GetToken(tenantId, clientId, clientSecret)
	if err != nil {
		return nil, newErr(
			codes.Unauthenticated,
			"Error authenticating key vault",
			err,
		)
	}
	return &keyVaultSecretService{
		client:      client,
		vaultName:   vaultName,
		accessToken: token,
	}, nil
}

func NewWithClient(client keyvaultapi.BaseClientAPI) secret.SecretService {
	return &keyVaultSecretService{
		client:    client,
		vaultName: "localvault",
	}
}
