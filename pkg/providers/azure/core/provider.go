package core

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzProvider interface {
	GetResources(AzResource) (map[string]AzGenericResource, error)
	SubscriptionId() string
	ResourceGroupName() string
	ServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error)
}

type AzResource = string

const (
	AzResource_Topic AzResource = "Microsoft.Eventgrid/topics"
	// Collections are handled by mongodb in azure
	// AzResource_Collection AzResource = "TODO"
	AzResource_Queue  AzResource = "Microsoft.Storage/storageAccounts/queueServices"
	AzResource_Bucket AzResource = "Microsoft.Storage/storageAccounts/blobServices"
	AzResource_Secret AzResource = "Microsoft.KeyVault/vaults/secrets"
)

type AzGenericResource struct {
	Name     string
	Type     string
	Location string
}

type azResourceCache = map[AzResource]map[string]AzGenericResource

type azProviderImpl struct {
	env     auth.EnvironmentSettings
	rclient resources.Client
	subId   string
	rgName  string
	cache   azResourceCache
}

func (p *azProviderImpl) GetResources(r AzResource) (map[string]AzGenericResource, error) {
	filter := fmt.Sprintf("resourceType eq %s", r)

	if _, ok := p.cache[r]; !ok {
		// populate the cache
		results, err := p.rclient.ListByResourceGroupComplete(context.TODO(), p.rgName, filter, "", nil)

		if err != nil {
			return nil, err
		}

		p.cache[r] = map[string]AzGenericResource{}

		for results.NotDone() {
			// TODO: Determine if iterator needs to be advanced first or last
			err := results.NextWithContext(context.TODO())

			if err != nil {
				return nil, err
			}

			resource := results.Value()
			if tagV, ok := resource.Tags["x-nitric-name"]; ok && tagV != nil {
				// Add it to the cache
				p.cache[r][*tagV] = AzGenericResource{
					Name:     *resource.Name,
					Type:     *resource.Type,
					Location: *resource.Location,
				}
			}
		}
	}

	// otherwise return the results
	return p.cache[r], nil
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
		cache:  make(map[string]map[string]AzGenericResource),
	}, nil
}
