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
