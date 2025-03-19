// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/samber/lo"

	cdn "github.com/pulumi/pulumi-azure-native-sdk/cdn/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricAzurePulumiProvider) createStaticWebsite(ctx *pulumi.Context, resources []*pulumix.NitricPulumiResource[any]) error {
	var rootWebsite *deploymentspb.Website
	var err error

	for _, resource := range resources {
		config, ok := resource.Config.(*deploymentspb.Resource_Website)

		if ok && config.Website.BasePath == "/" {
			rootWebsite = config.Website
			break
		}
	}

	if rootWebsite == nil {
		return fmt.Errorf("no root website configuration found")
	}

	p.staticWebsite, err = storage.NewStorageAccountStaticWebsite(ctx, "website", &storage.StorageAccountStaticWebsiteArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		AccountName:       p.StorageAccount.Name,
		IndexDocument:     pulumi.String(rootWebsite.IndexDocument),
		Error404Document:  pulumi.String(rootWebsite.ErrorDocument),
	}, pulumi.DependsOn([]pulumi.Resource{p.StorageAccount}))
	if err != nil {
		return err
	}

	return nil
}

func purgeCDNEndpoint(subscriptionID, resourceGroup, profileName, endpointName string, contentPaths []*string) error {
	// Authenticate with Azure
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}

	// Create a CDN endpoint client
	client, err := armcdn.NewAFDEndpointsClient(subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create CDN endpoint client: %w", err)
	}

	// Call PurgeContent
	_, err = client.BeginPurgeContent(
		context.Background(),
		resourceGroup,
		profileName,
		endpointName,
		armcdn.AfdPurgeParameters{
			ContentPaths: contentPaths,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to purge CDN endpoint: %w", err)
	}

	return nil
}

var (
	blobClientInstance *azblob.ServiceURL
	once               sync.Once
)

func getBlobServiceClient(subscriptionID, resourceGroupName, accountName string) (*azblob.ServiceURL, error) {
	var err error

	once.Do(func() {
		// Authenticate using Default Azure Credentials
		cred, e := azidentity.NewDefaultAzureCredential(nil)
		if e != nil {
			err = fmt.Errorf("failed to create Azure credential: %w", e)
			return
		}

		clientFactory, e := armstorage.NewClientFactory(subscriptionID, cred, nil)
		if e != nil {
			err = fmt.Errorf("failed to create storage client factory: %w", e)
			return
		}

		accountsClient := clientFactory.NewAccountsClient()

		// Get the account keys for the Storage Account
		result, e := accountsClient.ListKeys(context.Background(), resourceGroupName, accountName, nil)
		if e != nil {
			log.Fatalf("failed to get account keys: %v", e)
			return
		}

		accountKey := *result.Keys[0].Value

		// Create shared key credentials
		sharedKeyCred, e := azblob.NewSharedKeyCredential(accountName, accountKey)
		if e != nil {
			err = fmt.Errorf("failed to create shared key credential: %w", e)
			return
		}

		// Create a pipeline with the credentials
		pipeline := azblob.NewPipeline(sharedKeyCred, azblob.PipelineOptions{})

		// Build the storage account URL
		urlStr := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)

		parsedUrl, e := url.Parse(urlStr)
		if e != nil {
			err = fmt.Errorf("failed to parse URL: %w", e)
			return
		}

		// Create the service client with the URL and the pipeline
		blobClient := azblob.NewServiceURL(*parsedUrl, pipeline)
		blobClientInstance = &blobClient
	})

	return blobClientInstance, err
}

// Get MD5 hash of a blob (if it exists)
func getBlobMD5(serviceURL *azblob.ServiceURL, containerName, blobName string) (string, error) {
	containerURL := serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlobURL(blobName)

	props, err := blobURL.GetProperties(context.TODO(), azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		var storageErr azblob.StorageError

		if errors.As(err, &storageErr) && storageErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
			return "", nil // Return empty string if blob does not exist
		}
		return "", fmt.Errorf("failed to get blob properties: %w", err)
	}

	md5 := props.ContentMD5()
	if md5 == nil {
		return "", nil
	}

	// Convert the byte slice (MD5) to Base64
	base64MD5 := base64.StdEncoding.EncodeToString(md5)

	return base64MD5, nil
}

// Website - Implements the Website deployment method for the Azure provider
func (p *NitricAzurePulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{p.staticWebsite})}

	localDir, ok := config.AssetSource.(*deploymentspb.Website_LocalDirectory)
	if !ok {
		return fmt.Errorf("unsupported asset source type for website: %s", name)
	}

	cleanedPath := filepath.ToSlash(filepath.Clean(localDir.LocalDirectory))

	if p.staticWebsite == nil {
		return fmt.Errorf("website storage account not found")
	}

	// Walk the directory and upload each file to the storage account
	err := filepath.WalkDir(cleanedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get file info to check for special types
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Skip non-regular files (e.g., symlinks, sockets, devices)
		if info.Mode()&fs.ModeType != 0 {
			return nil
		}

		// Determine the content type based on the file extension
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Generate the object key to include the folder structure
		var objectKey string

		filePath := path[len(cleanedPath):]

		// If the base path is not the root, include it in the object key
		if config.BasePath == "/" {
			objectKey = filepath.ToSlash(filePath)
		} else {
			objectKey = filepath.ToSlash(filepath.Join(config.BasePath, filePath))
		}

		name := strings.TrimPrefix(objectKey, "/")

		existingMd5Output := pulumi.All(p.ResourceGroup.Name, p.StorageAccount.Name, p.staticWebsite.ContainerName).ApplyT(func(args []interface{}) (string, error) {
			resourceGroupName := args[0].(string)
			accountName := args[1].(string)
			containerName := args[2].(string)

			serviceClient, err := getBlobServiceClient(p.ClientConfig.SubscriptionId, resourceGroupName, accountName)
			if err != nil {
				return "", err
			}

			existingMd5, err := getBlobMD5(serviceClient, containerName, name)
			if err != nil {
				return "", err
			}

			return existingMd5, nil
		}).(pulumi.StringOutput)

		blob, err := storage.NewBlob(ctx, name, &storage.BlobArgs{
			ResourceGroupName: p.ResourceGroup.Name,
			AccountName:       p.StorageAccount.Name,
			ContainerName:     p.staticWebsite.ContainerName,
			Source:            pulumi.NewFileAsset(path),
			ContentType:       pulumi.String(contentType),
		}, opts...)
		if err != nil {
			return err
		}

		keyToInvalidate := pulumi.All(blob.ContentMd5, existingMd5Output).ApplyT(func(args []any) (string, error) {
			newMd5, newMd5Ok := args[0].(*string)
			existingMd5, existingMd5Ok := args[1].(string)

			if !newMd5Ok || !existingMd5Ok {
				return "", fmt.Errorf("failed to assert md5 types")
			}

			// if an existing ETag is present and it is different from the new ETag, return the key to invalidate
			if existingMd5 != "" && newMd5 != nil && *newMd5 != existingMd5 {
				return objectKey, nil
			}

			return "", nil
		}).(pulumi.StringOutput)

		p.staticWebsiteChangedFileOutputs = append(p.staticWebsiteChangedFileOutputs, keyToInvalidate)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Deploy CDN
func (p *NitricAzurePulumiProvider) deployCDN(ctx *pulumi.Context) error {
	profile, err := cdn.NewProfile(ctx, "website-profile", &cdn.ProfileArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		Location:          pulumi.String("Global"),
		Sku: &cdn.SkuArgs{
			// TODO could make this a config option, standard or premium
			Name: cdn.SkuName_Standard_AzureFrontDoor,
		},
	})
	if err != nil {
		return err
	}

	// Pull the hostname out of the storage-account endpoint.
	originHostname := p.StorageAccount.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
		parsed, err := url.Parse(endpoints.Web)
		if err != nil {
			return "", err
		}
		return parsed.Hostname(), nil
	}).(pulumi.StringOutput)

	endpointName := fmt.Sprintf("%s-cdn", p.StackId)

	defaultOriginGroupName := fmt.Sprintf("%s-default-origin-group", p.StackId)

	defaultOriginName := fmt.Sprintf("%s-default-origin", p.StackId)

	endpoint, err := cdn.NewAFDEndpoint(ctx, endpointName, &cdn.AFDEndpointArgs{
		EndpointName:      pulumi.String(endpointName),
		ResourceGroupName: p.ResourceGroup.Name,
		ProfileName:       profile.Name,
		Location:          profile.Location,
		EnabledState:      cdn.EnabledStateEnabled,
	}, pulumi.DependsOn([]pulumi.Resource{profile}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor endpoint: %w", err)
	}

	defaultOriginGroup, err := cdn.NewAFDOriginGroup(ctx, defaultOriginGroupName, &cdn.AFDOriginGroupArgs{
		OriginGroupName:   pulumi.String(defaultOriginGroupName),
		ResourceGroupName: p.ResourceGroup.Name,
		ProfileName:       profile.Name,
		LoadBalancingSettings: &cdn.LoadBalancingSettingsParametersArgs{
			AdditionalLatencyInMilliseconds: pulumi.Int(200), // Lower latency tolerance for faster failover
			SampleSize:                      pulumi.Int(5),   // More samples for better decision-making
			SuccessfulSamplesRequired:       pulumi.Int(3),   // Keep at 3 to maintain reliability
		},
	}, pulumi.DependsOn([]pulumi.Resource{profile}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor origin group: %w", err)
	}

	_, err = cdn.NewAFDOrigin(ctx, defaultOriginName, &cdn.AFDOriginArgs{
		OriginName:        pulumi.String(defaultOriginName),
		OriginGroupName:   defaultOriginGroup.Name,
		ResourceGroupName: p.ResourceGroup.Name,
		ProfileName:       profile.Name,
		HostName:          originHostname,
		OriginHostHeader:  originHostname,
		HttpPort:          pulumi.Int(80),
		HttpsPort:         pulumi.Int(443),
		EnabledState:      cdn.EnabledStateEnabled,
	}, pulumi.DependsOn([]pulumi.Resource{defaultOriginGroup}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor origin: %w", err)
	}

	// Create a rule set for the CDN endpoint if there are APIs
	ruleSets := cdn.ResourceReferenceArray{}

	if len(p.Apis) > 0 {
		ruleSetName := "apiruleset"
		ruleOrder := 1

		ruleSet, err := cdn.NewRuleSet(ctx, ruleSetName, &cdn.RuleSetArgs{
			RuleSetName:       pulumi.String(ruleSetName),
			ResourceGroupName: p.ResourceGroup.Name,
			ProfileName:       profile.Name,
		}, pulumi.DependsOn([]pulumi.Resource{endpoint}))
		if err != nil {
			return fmt.Errorf("failed to create Frontdoor rule set: %w", err)
		}

		ruleSets = append(ruleSets, &cdn.ResourceReferenceArgs{
			Id: ruleSet.ID(),
		})

		// Sort the APIs by name
		sortedApiKeys := lo.Keys(p.Apis)
		slices.Sort(sortedApiKeys)

		// Create a delivery rule for each API
		for _, apiName := range sortedApiKeys {
			api := p.Apis[apiName]

			apiHostName := api.ApiManagementService.GatewayUrl.ApplyT(func(gatewayUrl string) (string, error) {
				parsed, err := url.Parse(gatewayUrl)
				if err != nil {
					return "", err
				}

				return parsed.Hostname(), nil
			}).(pulumi.StringOutput)

			apiOriginGroupName := fmt.Sprintf("api-origin-group-%s", apiName)

			apiOriginGroup, err := cdn.NewAFDOriginGroup(ctx, apiOriginGroupName, &cdn.AFDOriginGroupArgs{
				OriginGroupName:   pulumi.String(apiOriginGroupName),
				ResourceGroupName: p.ResourceGroup.Name,
				ProfileName:       profile.Name,
				LoadBalancingSettings: &cdn.LoadBalancingSettingsParametersArgs{
					AdditionalLatencyInMilliseconds: pulumi.Int(100), // Quick failover for API requests
					SampleSize:                      pulumi.Int(5),   // More accurate health assessment
					SuccessfulSamplesRequired:       pulumi.Int(2),   // Faster recovery for healthy endpoints
				},
			}, pulumi.DependsOn([]pulumi.Resource{ruleSet}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor origin group: %w", err)
			}

			apiOriginName := fmt.Sprintf("api-origin-%s", apiName)

			origin, err := cdn.NewAFDOrigin(ctx, apiOriginName, &cdn.AFDOriginArgs{
				OriginName:        pulumi.String(apiOriginName),
				OriginGroupName:   apiOriginGroup.Name,
				ResourceGroupName: p.ResourceGroup.Name,
				ProfileName:       profile.Name,
				HostName:          apiHostName,
				OriginHostHeader:  apiHostName,
				HttpPort:          pulumi.Int(80),
				HttpsPort:         pulumi.Int(443),
			}, pulumi.DependsOn([]pulumi.Resource{apiOriginGroup}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor origin: %w", err)
			}

			ruleName := fmt.Sprintf("apirule%s", apiName)

			_, err = cdn.NewRule(ctx, ruleName, &cdn.RuleArgs{
				Order:             pulumi.Int(ruleOrder),
				RuleName:          pulumi.String(ruleName),
				RuleSetName:       ruleSet.Name,
				ProfileName:       profile.Name,
				ResourceGroupName: p.ResourceGroup.Name,
				Conditions: pulumi.ToArray(
					[]interface{}{
						cdn.DeliveryRuleUrlPathConditionArgs{
							Name: pulumi.String(cdn.MatchVariableUrlPath),
							Parameters: cdn.UrlPathMatchConditionParametersArgs{
								MatchValues: pulumi.StringArray{
									pulumi.Sprintf("/api/%s", apiName),
								},
								Transforms: pulumi.StringArray{
									pulumi.String(cdn.TransformLowercase),
								},
								TypeName: pulumi.String("DeliveryRuleUrlPathMatchConditionParameters"),
								Operator: pulumi.String(cdn.OperatorBeginsWith),
							},
						},
					}),
				Actions: pulumi.ToArray([]interface{}{
					cdn.DeliveryRuleRouteConfigurationOverrideActionArgs{
						Name: pulumi.String(cdn.DeliveryRuleActionRouteConfigurationOverride),
						Parameters: cdn.RouteConfigurationOverrideActionParametersArgs{
							OriginGroupOverride: cdn.OriginGroupOverrideArgs{
								ForwardingProtocol: pulumi.String(cdn.ForwardingProtocolHttpsOnly),
								OriginGroup: &cdn.ResourceReferenceArgs{
									Id: apiOriginGroup.ID(),
								},
							},
							CacheConfiguration: cdn.CacheConfigurationArgs{
								CacheBehavior:              pulumi.String(cdn.RuleCacheBehaviorHonorOrigin),
								QueryStringCachingBehavior: pulumi.String(cdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
							},
							TypeName: pulumi.String("DeliveryRuleRouteConfigurationOverrideActionParameters"),
						},
					},
					cdn.UrlRewriteActionArgs{
						Name: pulumi.String(cdn.DeliveryRuleActionUrlRewrite),
						Parameters: cdn.UrlRewriteActionParametersArgs{
							Destination:   pulumi.String("/"),
							SourcePattern: pulumi.String(fmt.Sprintf("/api/%s/", apiName)),
							TypeName:      pulumi.String("DeliveryRuleUrlRewriteActionParameters"),
						},
					},
				}),
			}, pulumi.DependsOn([]pulumi.Resource{ruleSet, origin, apiOriginGroup}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor rule: %w", err)
			}

			ruleOrder++
		}
	}

	routeName := fmt.Sprintf("%s-main-route", p.StackId)

	_, err = cdn.NewRoute(ctx, routeName, &cdn.RouteArgs{
		RouteName:         pulumi.String(routeName),
		ResourceGroupName: p.ResourceGroup.Name,
		// TODO make this a config option for custom domains
		LinkToDefaultDomain: pulumi.String(cdn.LinkToDefaultDomainEnabled),
		ProfileName:         profile.Name,
		EndpointName:        endpoint.Name,
		SupportedProtocols: pulumi.StringArray{
			pulumi.String(cdn.AFDEndpointProtocolsHttps),
		},
		ForwardingProtocol: pulumi.String(cdn.ForwardingProtocolHttpsOnly),
		HttpsRedirect:      pulumi.String(cdn.HttpsRedirectDisabled),
		PatternsToMatch:    pulumi.ToStringArray([]string{"/*"}),
		EnabledState:       cdn.EnabledStateEnabled,
		OriginGroup: &cdn.ResourceReferenceArgs{
			Id: defaultOriginGroup.ID(),
		},
		CacheConfiguration: &cdn.AfdRouteCacheConfigurationArgs{
			CompressionSettings: &cdn.CompressionSettingsArgs{
				ContentTypesToCompress: pulumi.StringArray{
					pulumi.String("text/html"),
					pulumi.String("text/css"),
					pulumi.String("application/javascript"),
					pulumi.String("application/json"),
					pulumi.String("image/svg+xml"),
					pulumi.String("font/woff"),
					pulumi.String("font/woff2"),
				},
				IsCompressionEnabled: pulumi.Bool(true),
			},
			QueryStringCachingBehavior: pulumi.String(cdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
		},
		RuleSets: ruleSets,
	}, pulumi.DependsOn([]pulumi.Resource{defaultOriginGroup}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor route: %w", err)
	}

	// Purge the CDN endpoint if content has changed
	pulumi.All(p.ResourceGroup.Name, profile.Name, p.staticWebsiteChangedFileOutputs.ToStringArrayOutput()).ApplyT(func(args []interface{}) error {
		resourceGroupName := args[0].(string)
		profileName := args[1].(string)
		websiteChangedFileKeys := []*string{}

		for _, key := range args[2].([]string) {
			if key != "" {
				// require to purge the index.html file served at root of cdn
				if strings.HasSuffix(key, "/index.html") {
					key = strings.TrimSuffix(key, "index.html")
				}

				websiteChangedFileKeys = append(websiteChangedFileKeys, &key)
			}
		}

		if len(websiteChangedFileKeys) > 0 {
			err := purgeCDNEndpoint(p.ClientConfig.SubscriptionId, resourceGroupName, profileName, endpointName, websiteChangedFileKeys)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// Export the CDN endpoint hostname.
	ctx.Export("cdn", pulumi.Sprintf("https://%s", endpoint.HostName))

	p.Endpoint = endpoint

	return nil
}
