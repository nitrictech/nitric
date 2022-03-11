package core

import (
	"fmt"
	"os"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzProvider interface {
	SubscriptionId() string
	ResourceGroupName() string
	ServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error)
}

type azProviderImpl struct {
	env    auth.EnvironmentSettings
	subId  string
	rgName string
}

func (p *azProviderImpl) ServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error) {
	if fileCred, fileErr := auth.GetSettingsFromFile(); fileErr == nil {
		return fileCred.ServicePrincipalTokenFromClientCredentialsWithResource(resource)
	} else if clientCred, clientErr := p.env.GetClientCredentials(); clientErr == nil {
		clientCred.Resource = resource
		return clientCred.ServicePrincipalToken()
	} else if clientCert, certErr := p.env.GetClientCertificate(); certErr == nil {
		clientCert.Resource = resource
		return clientCert.ServicePrincipalToken()
	} else if userPass, userErr := p.env.GetUsernamePassword(); userErr == nil {
		userPass.Resource = resource
		return userPass.ServicePrincipalToken()
	} else {
		fmt.Printf("error retrieving credentials:\n -> %v\n -> %v\n -> %v\n -> %v\n", fileErr, clientErr, certErr, userErr)
	}

	msiConf := p.env.GetMSI()
	msiConf.Resource = resource
	return msiConf.ServicePrincipalToken()
}

func (p *azProviderImpl) SubscriptionId() string {
	return p.subId
}

func (p *azProviderImpl) ResourceGroupName() string {
	return p.rgName
}

var _ AzProvider = &azProviderImpl{}

func New() (*azProviderImpl, error) {
	rgName := os.Getenv(AZURE_RESOUCE_GROUP)
	if rgName == "" {
		return nil, fmt.Errorf("envvar %s is not set", AZURE_RESOUCE_GROUP)
	}

	subId := os.Getenv(AZURE_SUBSCRIPTION_ID)
	if subId == "" {
		return nil, fmt.Errorf("envvar %s is not set", AZURE_SUBSCRIPTION_ID)
	}

	config, err := auth.GetSettingsFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &azProviderImpl{
		rgName: rgName,
		subId:  subId,
		env:    config,
	}, nil
}
