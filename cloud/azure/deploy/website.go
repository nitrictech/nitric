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
	"fmt"
	"io/fs"
	"mime"
	"net/url"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/samber/lo"

	cdn "github.com/pulumi/pulumi-azure-native-sdk/cdn/v2"
	network "github.com/pulumi/pulumi-azure-native-sdk/network/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Website - Implements the Website deployment method for the Azure provider
func (p *NitricAzurePulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	var err error
	normalizedName := strings.ReplaceAll(config.BasePath, "/", "")

	if normalizedName == "" {
		normalizedName = "root"
	}

	// create website storage account
	p.WebsiteStorageAccounts[config.BasePath], err = storage.NewStorageAccount(ctx, ResourceName(ctx, normalizedName, StorageAccountRTW), &storage.StorageAccountArgs{
		AccessTier:        storage.AccessTierHot,
		ResourceGroupName: p.ResourceGroup.Name,
		Kind:              pulumi.String("StorageV2"),
		Sku: storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
	})
	if err != nil {
		return err
	}

	// double check default documents are set
	if config.IndexDocument == "" {
		config.IndexDocument = "index.html"
	}

	if config.ErrorDocument == "" {
		config.ErrorDocument = "404.html"
	}

	p.WebsiteContainers[config.BasePath], err = storage.NewStorageAccountStaticWebsite(ctx, ResourceName(ctx, normalizedName, StorageContainerRTW), &storage.StorageAccountStaticWebsiteArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		AccountName:       p.WebsiteStorageAccounts[config.BasePath].Name,
		IndexDocument:     pulumi.String(config.IndexDocument),
		Error404Document:  pulumi.String(config.ErrorDocument),
	}, pulumi.DependsOn([]pulumi.Resource{p.WebsiteStorageAccounts[config.BasePath]}))
	if err != nil {
		return err
	}

	localDir, ok := config.AssetSource.(*deploymentspb.Website_LocalDirectory)
	if !ok {
		return fmt.Errorf("unsupported asset source type for website: %s", name)
	}

	cleanedPath := filepath.ToSlash(filepath.Clean(localDir.LocalDirectory))

	// add the website contain as a dependency for uploads
	opts = append(opts, pulumi.DependsOn([]pulumi.Resource{p.WebsiteContainers[config.BasePath]}))

	// Walk the directory and upload each file to the storage account
	err = filepath.WalkDir(cleanedPath, func(path string, d fs.DirEntry, err error) error {
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
		objectKey = filepath.ToSlash(filePath)

		name := strings.TrimPrefix(objectKey, "/")

		uniqueName := fmt.Sprintf("%s-%s", normalizedName, name)

		blob, err := storage.NewBlob(ctx, uniqueName, &storage.BlobArgs{
			BlobName:          pulumi.String(name),
			ResourceGroupName: p.ResourceGroup.Name,
			AccountName:       p.WebsiteStorageAccounts[config.BasePath].Name,
			ContainerName:     p.WebsiteContainers[config.BasePath].ContainerName,
			Source:            pulumi.NewFileAsset(path),
			ContentType:       pulumi.String(contentType),
		}, opts...)
		if err != nil {
			return err
		}

		// Get the MD5 hash of the blob and store it in the outputs
		p.websiteFileMd5Outputs = append(p.websiteFileMd5Outputs, blob.ContentMd5)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ensureValidSubdomain(domain string, subdomain string) error {
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	subdomain = strings.ToLower(strings.TrimSuffix(subdomain, "."))

	if subdomain == domain || strings.HasSuffix(subdomain, "."+domain) {
		return nil
	}

	return fmt.Errorf("%s is not a valid subdomain of %s", subdomain, domain)
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

	endpointName := fmt.Sprintf("%s-cdn", p.StackId)

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

	customDomains := cdn.ActivatedResourceReferenceArray{}

	if p.AzureConfig.CdnDomain.DomainName != "" && p.AzureConfig.CdnDomain.ZoneName != "" {
		// both are required if one is set
		if p.AzureConfig.CdnDomain.ZoneName == "" {
			return fmt.Errorf("zone-name is required for custom domain")
		}

		if p.AzureConfig.CdnDomain.DomainName == "" {
			return fmt.Errorf("domain-name is required for custom domain")
		}

		if p.AzureConfig.CdnDomain.ZoneResourceGroup == "" {
			return fmt.Errorf("zone-resource-group is required for custom domain")
		}

		dnsZone, err := network.LookupZone(ctx, &network.LookupZoneArgs{
			ResourceGroupName: p.AzureConfig.CdnDomain.ZoneResourceGroup,
			ZoneName:          p.AzureConfig.CdnDomain.ZoneName,
		})
		if err != nil {
			return err
		}

		subDomain := strings.ToLower(p.AzureConfig.CdnDomain.DomainName)

		// check if the domain name is a subdomain of the zone name
		err = ensureValidSubdomain(p.AzureConfig.CdnDomain.ZoneName, subDomain)
		if err != nil {
			return err
		}

		// if it is a subdomain, remove the zone name from the subdomain
		subDomain = strings.ReplaceAll(strings.TrimSuffix(subDomain, p.AzureConfig.CdnDomain.ZoneName), ".", "")

		isApexDomain := subDomain == ""

		// if domain is an apex domain, create a unique subdomain for naming purposes
		if isApexDomain {
			subDomain = fmt.Sprintf("%s-%s", p.StackId, strings.ReplaceAll(p.AzureConfig.CdnDomain.ZoneName, ".", "-"))
		}

		domain, err := cdn.NewAFDCustomDomain(ctx, p.AzureConfig.CdnDomain.DomainName, &cdn.AFDCustomDomainArgs{
			ResourceGroupName: p.ResourceGroup.Name,
			ProfileName:       profile.Name,
			CustomDomainName:  pulumi.String(subDomain),
			HostName:          pulumi.String(p.AzureConfig.CdnDomain.DomainName),
			AzureDnsZone: &cdn.ResourceReferenceArgs{
				Id: pulumi.String(dnsZone.Id),
			},
			TlsSettings: &cdn.AFDDomainHttpsParametersArgs{
				CertificateType:   cdn.AfdCertificateTypeManagedCertificate,
				MinimumTlsVersion: cdn.AfdMinimumTlsVersionTLS12,
			},
		}, pulumi.DependsOn([]pulumi.Resource{profile, endpoint}))
		if err != nil {
			return err
		}

		relativeRecordSetName := fmt.Sprintf("_dnsauth.%s", subDomain)

		if isApexDomain {
			relativeRecordSetName = "_dnsauth"
		}

		// Create a TXT record for domain validation
		_, err = network.NewRecordSet(ctx, fmt.Sprintf("validate-%s", subDomain), &network.RecordSetArgs{
			RecordType:            pulumi.String("TXT"),
			RelativeRecordSetName: pulumi.String(relativeRecordSetName),
			ResourceGroupName:     pulumi.String(p.AzureConfig.CdnDomain.ZoneResourceGroup),
			Ttl:                   pulumi.Float64(3600), // Set TTL to one hour
			TxtRecords: network.TxtRecordArray{
				&network.TxtRecordArgs{
					Value: pulumi.StringArray{
						domain.ValidationProperties.ValidationToken(),
					},
				},
			},
			ZoneName: pulumi.String(dnsZone.Name),
		}, pulumi.DependsOn([]pulumi.Resource{domain}))
		if err != nil {
			return err
		}

		if isApexDomain {
			// Create an APEX record for the custom domain
			_, err = network.NewRecordSet(ctx, fmt.Sprintf("cdn-%s", subDomain), &network.RecordSetArgs{
				RecordType:            pulumi.String("A"),
				RelativeRecordSetName: pulumi.String("@"),
				TargetResource: &network.SubResourceArgs{
					Id: endpoint.ID(),
				},
				ResourceGroupName: pulumi.String(p.AzureConfig.CdnDomain.ZoneResourceGroup),
				Ttl:               pulumi.Float64(3600), // Set TTL to one hour
				ZoneName:          pulumi.String(dnsZone.Name),
			}, pulumi.DependsOn([]pulumi.Resource{domain, endpoint}))
			if err != nil {
				return err
			}
		} else {
			// Create a CNAME record for the custom domain
			_, err = network.NewRecordSet(ctx, fmt.Sprintf("cdn-%s", subDomain), &network.RecordSetArgs{
				RecordType:            pulumi.String("CNAME"),
				RelativeRecordSetName: pulumi.String(subDomain),
				ResourceGroupName:     pulumi.String(p.AzureConfig.CdnDomain.ZoneResourceGroup),
				Ttl:                   pulumi.Float64(3600), // Set TTL to one hour
				CnameRecord: &network.CnameRecordArgs{
					Cname: endpoint.HostName,
				},
				ZoneName: pulumi.String(dnsZone.Name),
			}, pulumi.DependsOn([]pulumi.Resource{domain}))
			if err != nil {
				return err
			}
		}

		customDomains = append(customDomains, &cdn.ActivatedResourceReferenceArgs{
			Id: domain.ID(),
		})
	}

	// Pull the hostname out of the root storage-account endpoint.
	originHostname := getWebEndpointHostName(p.WebsiteStorageAccounts["/"])
	defaultOriginGroupName := fmt.Sprintf("%s-default-origin-group", p.StackId)

	defaultOriginName := fmt.Sprintf("%s-default-origin", p.StackId)

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

	ruleSets := cdn.ResourceReferenceArray{}

	// Create a default rule set for the CDN endpoint
	ruleSetName := "default"
	defaultRuleSet, err := cdn.NewRuleSet(ctx, ResourceName(ctx, ruleSetName, FrontDoorRuleSetRT), &cdn.RuleSetArgs{
		RuleSetName:       pulumi.String(ruleSetName),
		ResourceGroupName: p.ResourceGroup.Name,
		ProfileName:       profile.Name,
	}, pulumi.DependsOn([]pulumi.Resource{endpoint}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor rule set: %w", err)
	}

	ruleSets = append(ruleSets, &cdn.ResourceReferenceArgs{
		Id: defaultRuleSet.ID(),
	})

	// Create a rule for redirecting paths that end in a slash
	_, err = cdn.NewRule(ctx, ResourceName(ctx, "redirectslash", FrontDoorRuleRT), &cdn.RuleArgs{
		Order:             pulumi.Int(1),
		RuleName:          pulumi.String("redirectslash"),
		RuleSetName:       defaultRuleSet.Name,
		ProfileName:       profile.Name,
		ResourceGroupName: p.ResourceGroup.Name,
		Conditions: pulumi.ToArray(
			[]interface{}{
				cdn.DeliveryRuleUrlPathConditionArgs{
					Name: pulumi.String(cdn.MatchVariableUrlPath),
					Parameters: cdn.UrlPathMatchConditionParametersArgs{
						MatchValues: pulumi.StringArray{
							pulumi.String(".*\\/$"),
						},
						TypeName: pulumi.String("DeliveryRuleUrlPathMatchConditionParameters"),
						Operator: pulumi.String(cdn.OperatorRegEx),
					},
				},
			}),
		Actions: pulumi.ToArray([]interface{}{
			cdn.UrlRedirectActionArgs{
				Name: pulumi.String(cdn.DeliveryRuleActionUrlRedirect),
				Parameters: cdn.UrlRedirectActionParametersArgs{
					RedirectType:        pulumi.String(cdn.RedirectTypeFound),
					DestinationProtocol: pulumi.String(cdn.DestinationProtocolMatchRequest),
					CustomPath:          pulumi.String("/{url_path:0:-1}"),
					TypeName:            pulumi.String("DeliveryRuleUrlRedirectActionParameters"),
				},
			},
		}),
	}, pulumi.DependsOn([]pulumi.Resource{defaultRuleSet}))
	if err != nil {
		return fmt.Errorf("failed to create Frontdoor rule for redirecting paths: %w", err)
	}

	// Add origins, origin groups and rule for sub websites
	if len(p.WebsiteStorageAccounts) > 1 {
		// Sort the storage accounts by name
		sortedStorageKeys := lo.Keys(p.WebsiteStorageAccounts)
		slices.Sort(sortedStorageKeys)

		// Create an origin group for each storage account
		for _, basePath := range sortedStorageKeys {
			// Skip the root storage account
			if basePath == "/" {
				continue
			}

			storageAccount := p.WebsiteStorageAccounts[basePath]

			normalizedName := strings.ReplaceAll(basePath, "/", "")

			originGroupName := fmt.Sprintf("%s-%s-origin-group", p.StackId, normalizedName)

			subsiteOriginGroup, err := cdn.NewAFDOriginGroup(ctx, originGroupName, &cdn.AFDOriginGroupArgs{
				OriginGroupName:   pulumi.String(originGroupName),
				ResourceGroupName: p.ResourceGroup.Name,
				ProfileName:       profile.Name,
				LoadBalancingSettings: &cdn.LoadBalancingSettingsParametersArgs{
					AdditionalLatencyInMilliseconds: pulumi.Int(200), // Lower latency tolerance for faster failover
					SampleSize:                      pulumi.Int(5),   // More samples for better decision-making
					SuccessfulSamplesRequired:       pulumi.Int(3),   // Keep at 3 to maintain reliability
				},
			}, pulumi.DependsOn([]pulumi.Resource{profile}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor origin group for subsite: %w", err)
			}

			originName := fmt.Sprintf("%s-%s-origin", p.StackId, normalizedName)

			originHostname := getWebEndpointHostName(storageAccount)

			subsiteOrigin, err := cdn.NewAFDOrigin(ctx, originName, &cdn.AFDOriginArgs{
				OriginName:        pulumi.String(originName),
				ResourceGroupName: p.ResourceGroup.Name,
				ProfileName:       profile.Name,
				HostName:          originHostname,
				OriginHostHeader:  originHostname,
				HttpPort:          pulumi.Int(80),
				HttpsPort:         pulumi.Int(443),
				EnabledState:      cdn.EnabledStateEnabled,
				OriginGroupName:   subsiteOriginGroup.Name,
			}, pulumi.DependsOn([]pulumi.Resource{subsiteOriginGroup}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor origin for subsite: %w", err)
			}

			ruleName := ResourceName(ctx, normalizedName, FrontDoorRuleRT)
			// create a override rule for the subsite
			_, err = cdn.NewRule(ctx, ruleName, &cdn.RuleArgs{
				Order:             pulumi.Int(nameToUniqueNumber(normalizedName)),
				RuleName:          pulumi.String(ruleName),
				RuleSetName:       defaultRuleSet.Name,
				ProfileName:       profile.Name,
				ResourceGroupName: p.ResourceGroup.Name,
				Conditions: pulumi.ToArray(
					[]interface{}{
						cdn.DeliveryRuleUrlPathConditionArgs{
							Name: pulumi.String(cdn.MatchVariableUrlPath),
							Parameters: cdn.UrlPathMatchConditionParametersArgs{
								MatchValues: pulumi.StringArray{
									// Match the base path and any sub-paths, azure requires no leading slash
									pulumi.Sprintf("%s(/.*)?$", strings.TrimPrefix(basePath, "/")),
								},
								Transforms: pulumi.StringArray{
									pulumi.String(cdn.TransformLowercase),
								},
								TypeName: pulumi.String("DeliveryRuleUrlPathMatchConditionParameters"),
								Operator: pulumi.String(cdn.OperatorRegEx),
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
									Id: subsiteOriginGroup.ID(),
								},
							},
							CacheConfiguration: cdn.CacheConfigurationArgs{
								CacheBehavior:              pulumi.String(cdn.RuleCacheBehaviorHonorOrigin),
								IsCompressionEnabled:       cdn.RuleIsCompressionEnabledEnabled,
								QueryStringCachingBehavior: pulumi.String(cdn.AfdQueryStringCachingBehaviorUseQueryString),
							},
							TypeName: pulumi.String("DeliveryRuleRouteConfigurationOverrideActionParameters"),
						},
					},
					cdn.UrlRewriteActionArgs{
						Name: pulumi.String(cdn.DeliveryRuleActionUrlRewrite),
						Parameters: cdn.UrlRewriteActionParametersArgs{
							Destination:           pulumi.String("/"),
							SourcePattern:         pulumi.String(strings.TrimSuffix(basePath, "/")),
							PreserveUnmatchedPath: pulumi.Bool(true),
							TypeName:              pulumi.String("DeliveryRuleUrlRewriteActionParameters"),
						},
					},
				}),
			}, pulumi.DependsOn([]pulumi.Resource{defaultRuleSet, subsiteOrigin, subsiteOriginGroup}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor rule for subsite: %w", err)
			}
		}
	}

	// Create a rule set for the CDN endpoint if there are APIs
	if len(p.Apis) > 0 {
		ruleSetName := ResourceName(ctx, "apis", FrontDoorRuleSetRT)

		ruleSet, err := cdn.NewRuleSet(ctx, ResourceName(ctx, ruleSetName, FrontDoorRuleSetRT), &cdn.RuleSetArgs{
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
			ruleOrder := nameToUniqueNumber(apiName)

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
				EnabledState:      cdn.EnabledStateEnabled,
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

			ruleName := ResourceName(ctx, fmt.Sprintf("api%s", apiName), FrontDoorRuleRT)

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
								IsCompressionEnabled:       cdn.RuleIsCompressionEnabledEnabled,
								QueryStringCachingBehavior: pulumi.String(cdn.AfdQueryStringCachingBehaviorUseQueryString),
							},
							TypeName: pulumi.String("DeliveryRuleRouteConfigurationOverrideActionParameters"),
						},
					},
					cdn.UrlRewriteActionArgs{
						Name: pulumi.String(cdn.DeliveryRuleActionUrlRewrite),
						Parameters: cdn.UrlRewriteActionParametersArgs{
							Destination:           pulumi.String("/"),
							SourcePattern:         pulumi.String(fmt.Sprintf("/api/%s/", apiName)),
							PreserveUnmatchedPath: pulumi.Bool(true),
							TypeName:              pulumi.String("DeliveryRuleUrlRewriteActionParameters"),
						},
					},
				}),
			}, pulumi.DependsOn([]pulumi.Resource{ruleSet, origin, apiOriginGroup}))
			if err != nil {
				return fmt.Errorf("failed to create Frontdoor rule: %w", err)
			}
		}
	}

	routeName := fmt.Sprintf("%s-main-route", p.StackId)

	_, err = cdn.NewRoute(ctx, routeName, &cdn.RouteArgs{
		RouteName:         pulumi.String(routeName),
		CustomDomains:     customDomains,
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

	// Export the CDN endpoint hostname.
	ctx.Export("cdn", pulumi.Sprintf("https://%s", endpoint.HostName))

	p.Endpoint = endpoint

	if p.AzureConfig.CdnDomain.SkipCacheInvalidation != nil && *p.AzureConfig.CdnDomain.SkipCacheInvalidation {
		return nil
	}

	// Apply a function to sort the array
	sortedMd5Result := p.websiteFileMd5Outputs.ToArrayOutput().ApplyT(func(arr []interface{}) string {
		// Convert each element to string
		md5Strings := []string{}
		for _, md5 := range arr {
			if md5Str, ok := md5.(string); ok {
				if md5Str != "" {
					md5Strings = append(md5Strings, md5Str)
				}
			}
		}

		sort.Strings(md5Strings)

		return strings.Join(md5Strings, "")
	}).(pulumi.StringOutput)

	// Purge the CDN endpoint if content has changed
	_, err = local.NewCommand(ctx, "invalidate-cache", &local.CommandArgs{
		Create: pulumi.Sprintf("MSYS_NO_PATHCONV=1 az afd endpoint purge -g %s --profile-name %s --endpoint-name %s --subscription %s --content-paths '/*' --no-wait",
			p.ResourceGroup.Name, profile.Name, endpointName, p.ClientConfig.SubscriptionId),
		Triggers: pulumi.Array{
			sortedMd5Result,
		},
		Logging: local.LoggingStdoutAndStderr,
		Interpreter: pulumi.StringArray{
			pulumi.String("bash"),
			pulumi.String("-c"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{endpoint}))
	if err != nil {
		return fmt.Errorf("failed to create command to purge CDN endpoint: %w", err)
	}

	return nil
}

func getWebEndpointHostName(storageAccount *storage.StorageAccount) pulumi.StringOutput {
	return storageAccount.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
		parsed, err := url.Parse(endpoints.Web)
		if err != nil {
			return "", err
		}
		return parsed.Hostname(), nil
	}).(pulumi.StringOutput)
}

// Convert a name to a unique number
// This creates a number that's guaranteed unique for different strings
// Required due to conflicts with Azure's rule order during updates/replacements
// https://learn.microsoft.com/en-us/answers/questions/2103790/why-does-a-azure-front-door-rule-set-rule-has-an-o
func nameToUniqueNumber(name string) int {
	// Start at a high base to avoid conflicts with other rules
	base := 10000

	// Use character position and value to guarantee uniqueness
	// This is essentially creating a custom numeric representation
	for i, char := range name {
		// Multiply by position+1 to weight characters differently
		// Use prime number multiplier to reduce collision risk
		base += int(char) * (i + 1) * 31
	}

	return base
}
