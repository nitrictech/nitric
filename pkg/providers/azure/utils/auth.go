package utils

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func GetServicePrincipalToken() (*adal.ServicePrincipalToken, error) {
	config, err := auth.GetSettingsFromEnvironment()

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve Azure auth settings: %v", err)
	}

	if clientCred, err := config.GetClientCredentials(); err == nil {
		return clientCred.ServicePrincipalToken()
	} else if clientCert, err := config.GetClientCertificate(); err == nil {
		return clientCert.ServicePrincipalToken()
	} else if userPass, err := config.GetUsernamePassword(); err == nil {
		return userPass.ServicePrincipalToken()
	}

	msiConf := config.GetMSI()
	return msiConf.ServicePrincipalToken()
}
