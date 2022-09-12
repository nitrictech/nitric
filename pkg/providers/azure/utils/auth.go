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

package utils

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// GetServicePrincipalToken - Retrieves the service principal token from env
func GetServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error) {
	config, err := auth.GetSettingsFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve azure auth settings: %v", err)
	}

	if fileCred, fileErr := auth.GetSettingsFromFile(); fileErr == nil {
		return fileCred.ServicePrincipalTokenFromClientCredentialsWithResource(resource)
	} else if clientCred, clientErr := config.GetClientCredentials(); clientErr == nil {
		clientCred.Resource = resource
		return clientCred.ServicePrincipalToken()
	} else if clientCert, certErr := config.GetClientCertificate(); certErr == nil {
		clientCert.Resource = resource
		return clientCert.ServicePrincipalToken()
	} else if userPass, userErr := config.GetUsernamePassword(); userErr == nil {
		userPass.Resource = resource
		return userPass.ServicePrincipalToken()
	} else {
		fmt.Printf("error retrieving credentials:\n -> %v\n -> %v\n -> %v\n -> %v\n", fileErr, clientErr, certErr, userErr)
	}

	msiConf := config.GetMSI()
	msiConf.Resource = resource
	return msiConf.ServicePrincipalToken()
}
