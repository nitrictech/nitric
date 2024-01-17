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

package resource

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	resourcepb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type AzProvider interface {
	GetResources(context.Context, AzResource) (map[string]AzGenericResource, error)
	SubscriptionId() string
	ResourceGroupName() string
	ServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error)
}

var _ resourcepb.ResourcesServer = &AzureResourceService{}

type AzResource = string

const (
	AzResource_Topic AzResource = "Microsoft.Eventgrid/topics"
	// Collections are handled by mongodb in azure
	// AzResource_Collection AzResource = "TODO"
	AzResource_Api    AzResource = "Microsoft.ApiManagement/service"
	AzResource_Queue  AzResource = "Microsoft.Storage/storageAccounts/queueServices"
	AzResource_Bucket AzResource = "Microsoft.Storage/storageAccounts/blobServices"
	AzResource_Secret AzResource = "Microsoft.KeyVault/vaults/secrets"
)

type AzGenericResource struct {
	Name       string
	Type       string
	Location   string
	Properties interface{}
}

type azResourceCache = map[AzResource]map[string]AzGenericResource

type AzureResourceService struct {
	env       auth.EnvironmentSettings
	rclient   resources.Client
	srvClient apimanagement.ServiceClient
	subId     string
	rgName    string
	stackId   string
	cache     azResourceCache
}

func (p *AzureResourceService) getApiDetails(ctx context.Context, name string) (*resourcepb.ResourceDetailsResponse, error) {
	res, err := p.srvClient.ListByResourceGroupComplete(ctx, p.rgName)
	if err != nil {
		return nil, err
	}

	for res.NotDone() {
		service := res.Value()

		if t, ok := service.Tags[fmt.Sprintf("x-nitric-stackId-%s", p.stackId)]; ok && t != nil && *t == name {
			return &resourcepb.ResourceDetailsResponse{
				Id:       *service.ID,
				Provider: "azure",
				Service:  "ApiManagement",
				Details: &resourcepb.ResourceDetailsResponse_Api{
					Api: &resourcepb.ApiResourceDetails{
						Url: *service.GatewayURL,
					},
				},
			}, nil
		}

		err := res.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("api resource %s not found", name)
}

func (p *AzureResourceService) Declare(ctx context.Context, req *resourcepb.ResourceDeclareRequest) (*resourcepb.ResourceDeclareResponse, error) {
	return &resourcepb.ResourceDeclareResponse{}, nil
}

func (p *AzureResourceService) Details(ctx context.Context, req *resourcepb.ResourceDetailsRequest) (*resourcepb.ResourceDetailsResponse, error) {
	switch req.Id.Type {
	case resourcepb.ResourceType_Api:
		return p.getApiDetails(ctx, req.Id.Name)
	default:
		return nil, fmt.Errorf("unsupported resource type %s", req.Id.Type)
	}
}

func (p *AzureResourceService) GetResources(ctx context.Context, r AzResource) (map[string]AzGenericResource, error) {
	filter := fmt.Sprintf("resourceType eq '%s'", r)
	if _, ok := p.cache[r]; !ok {
		// populate the cache
		results, err := p.rclient.ListByResourceGroupComplete(ctx, p.rgName, filter, "", nil)
		if err != nil {
			return nil, err
		}

		p.cache[r] = map[string]AzGenericResource{}

		for results.NotDone() {
			resource := results.Value()

			resourceNameKey := tags.GetResourceNameKey(p.stackId)
			if tagV, ok := resource.Tags[resourceNameKey]; ok && tagV != nil {
				// Add it to the cache
				p.cache[r][*tagV] = AzGenericResource{
					Name:       *resource.Name,
					Type:       *resource.Type,
					Location:   *resource.Location,
					Properties: resource.Properties,
				}
			}

			err := results.NextWithContext(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	// otherwise return the results
	return p.cache[r], nil
}

func (p *AzureResourceService) ServicePrincipalToken(resource string) (*adal.ServicePrincipalToken, error) {
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
		msiConf := p.env.GetMSI()
		msiConf.Resource = resource
		return msiConf.ServicePrincipalToken()
	}
}

func (p *AzureResourceService) SubscriptionId() string {
	return p.subId
}

func (p *AzureResourceService) ResourceGroupName() string {
	return p.rgName
}

var _ AzProvider = &AzureResourceService{}

func New() (*AzureResourceService, error) {
	rgName := os.Getenv(AZURE_RESOURCE_GROUP)
	if rgName == "" {
		return nil, fmt.Errorf("envvar %s is not set", AZURE_RESOURCE_GROUP)
	}

	subId := os.Getenv(AZURE_SUBSCRIPTION_ID)
	if subId == "" {
		return nil, fmt.Errorf("envvar %s is not set", AZURE_SUBSCRIPTION_ID)
	}

	stackId := os.Getenv(NITRIC_STACK_ID)
	if stackId == "" {
		return nil, fmt.Errorf("envvar %s is not set", NITRIC_STACK_ID)
	}

	config, err := auth.GetSettingsFromEnvironment()
	if err != nil {
		return nil, err
	}

	prov := &AzureResourceService{
		rgName:  rgName,
		subId:   subId,
		env:     config,
		cache:   make(map[string]map[string]AzGenericResource),
		stackId: stackId,
	}

	spt, err := prov.ServicePrincipalToken("https://management.azure.com")
	if err != nil {
		return nil, err
	}

	sClient := apimanagement.NewServiceClient(subId)
	sClient.Authorizer = autorest.NewBearerAuthorizer(spt)

	rclient := resources.NewClient(subId)
	rclient.Authorizer = autorest.NewBearerAuthorizer(spt)
	prov.rclient = rclient
	prov.srvClient = sClient

	return prov, nil
}
