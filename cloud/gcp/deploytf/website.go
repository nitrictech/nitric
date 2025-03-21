package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/website"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricGcpTerraformProvider) deployEntrypoint(stack cdktf.TerraformStack) error {
	// Return gRPC unimplemented error
	// pathRules := compute.URLMapPathMatcherPathRuleArray{}

	// // Add deployed API gateways to the URLMap
	// for name, api := range a.ApiGateways {
	// 	neg, err := compute.NewRegionNetworkEndpointGroup(ctx, fmt.Sprintf("%s-apigw-neg", name), &compute.RegionNetworkEndpointGroupArgs{
	// 		NetworkEndpointType: pulumi.String("SERVERLESS"),
	// 		Region:              api.Region,
	// 		ServerlessDeployment: compute.RegionNetworkEndpointGroupServerlessDeploymentArgs{
	// 			Platform: pulumi.String("apigateway.googleapis.com"),
	// 			Resource: api.GatewayId,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}

	// 	bs, err := compute.NewBackendService(ctx, fmt.Sprintf("%s-apigw-bs", name), &compute.BackendServiceArgs{
	// 		Backends: compute.BackendServiceBackendArray{
	// 			compute.BackendServiceBackendArgs{
	// 				Group: neg.SelfLink,
	// 			},
	// 		},
	// 		// EnableCdn: pulumi.Bool(true),
	// 		Protocol: pulumi.String("HTTPS"),
	// 	}, nil)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	pr := compute.URLMapPathMatcherPathRuleArgs{
	// 		Service: bs.ID(),
	// 		Paths:   pulumi.StringArray{pulumi.Sprintf("/apis/%s/*", name)},
	// 		RouteAction: compute.URLMapPathMatcherPathRuleRouteActionArgs{
	// 			UrlRewrite: compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs{
	// 				PathPrefixRewrite: pulumi.String("/"),
	// 				HostRewrite:       api.DefaultHostname,
	// 			},
	// 		},
	// 	}

	// 	pathRules = append(pathRules, pr)
	// }
	// var defaultService pulumi.StringOutput
	// for name, website := range a.WebsiteBuckets {
	// 	normalizedName := strings.Replace(name, "/", "", -1)
	// 	if normalizedName == "" {
	// 		normalizedName = "default"
	// 	}

	// 	backend, err := compute.NewBackendBucket(ctx, fmt.Sprintf("%s-site-bucket", normalizedName), &compute.BackendBucketArgs{
	// 		BucketName: website.Name,
	// 		EnableCdn:  pulumi.Bool(true),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if name == "/" {
	// 		defaultService = backend.SelfLink
	// 	} else {
	// 		pr := compute.URLMapPathMatcherPathRuleArgs{
	// 			Service: backend.ID(),
	// 			Paths:   pulumi.StringArray{pulumi.String(filepath.Join("/", name, "./*"))},
	// 			RouteAction: compute.URLMapPathMatcherPathRuleRouteActionArgs{
	// 				UrlRewrite: compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs{
	// 					PathPrefixRewrite: pulumi.String("/"),
	// 				},
	// 			},
	// 		}

	// 		pathRules = append(pathRules, pr)
	// 	}
	// }

	// // Provision a global IP address for the CDN.
	// ip, err := compute.NewGlobalAddress(ctx, "ip", nil)
	// if err != nil {
	// 	return err
	// }

	// // Create a URLMap to route requests to the storage bucket.
	// httpsUrlMap, err := compute.NewURLMap(ctx, "https-site-url-map", &compute.URLMapArgs{
	// 	DefaultService: defaultService,
	// 	HostRules: compute.URLMapHostRuleArray{
	// 		compute.URLMapHostRuleArgs{
	// 			Hosts:       pulumi.StringArray{pulumi.String("*")},
	// 			PathMatcher: pulumi.String("all-paths"),
	// 		},
	// 	},
	// 	PathMatchers: compute.URLMapPathMatcherArray{
	// 		compute.URLMapPathMatcherArgs{
	// 			Name:           pulumi.String("all-paths"),
	// 			DefaultService: defaultService,
	// 			PathRules:      pathRules,
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	// // If a domain is specified in the config, then lookup to see if there is a GCP managed zone for it
	// managedZone, err := dns.LookupManagedZone(ctx, &dns.LookupManagedZoneArgs{
	// 	Name: a.GcpConfig.CdnDomain.ZoneName,
	// })
	// if err != nil {
	// 	return err
	// }

	// // Add root zone, to ensure reliable comparisons (i.e. trailing dot, e.g. example.com.)
	// var subDomain = strings.ToLower(a.GcpConfig.CdnDomain.DomainName)
	// if !strings.HasSuffix(subDomain, ".") {
	// 	subDomain = subDomain + "."
	// }

	// if !strings.HasSuffix(subDomain, managedZone.DnsName) {
	// 	return fmt.Errorf("CDN domain %s is not a subdomain of managed zone %s", subDomain, managedZone.DnsName)
	// }

	// // Create root DNS record for the IP address
	// _, err = dns.NewRecordSet(ctx, "cdn-dns-record", &dns.RecordSetArgs{
	// 	Name:        pulumi.String(subDomain),
	// 	ManagedZone: pulumi.String(managedZone.Name),
	// 	Type:        pulumi.String("A"),
	// 	Rrdatas:     pulumi.StringArray{ip.Address},
	// 	Ttl:         pulumi.IntPtr(300),
	// })
	// if err != nil {
	// 	return err
	// }

	// _, err = dns.NewRecordSet(ctx, "www-cdn-dns-record", &dns.RecordSetArgs{
	// 	Name:        pulumi.String(fmt.Sprintf("www.%s", subDomain)),
	// 	ManagedZone: pulumi.String(managedZone.Name),
	// 	Type:        pulumi.String("A"),
	// 	Rrdatas:     pulumi.StringArray{ip.Address},
	// 	Ttl:         pulumi.IntPtr(300),
	// })
	// if err != nil {
	// 	return err
	// }

	// // The certificate will use Load Balancer authorization (as opposed to DNS auth).
	// sslCert, err := certificatemanager.NewCertificate(ctx, "cdn-cert", &certificatemanager.CertificateArgs{
	// 	Scope: pulumi.String("DEFAULT"),
	// 	Managed: certificatemanager.CertificateManagedArgs{
	// 		Domains: pulumi.StringArray{
	// 			// Removing trailing dot (root zone), it's unsupported by certificate manager
	// 			pulumi.String(subDomain[:len(subDomain)-1]),
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	// certMap, err := certificatemanager.NewCertificateMapResource(ctx, "cert-map", &certificatemanager.CertificateMapResourceArgs{})
	// if err != nil {
	// 	return err
	// }

	// _, err = certificatemanager.NewCertificateMapEntry(ctx, "default", &certificatemanager.CertificateMapEntryArgs{
	// 	Name:        pulumi.String("cdn-cert-map-entry"),
	// 	Description: pulumi.String("CDN Certificate Map Entry"),
	// 	Map:         certMap.Name,
	// 	Certificates: pulumi.StringArray{
	// 		sslCert.ID(),
	// 	},
	// 	Matcher: pulumi.String("PRIMARY"),
	// })
	// if err != nil {
	// 	return err
	// }

	// // Create an HTTP proxy to route requests to the URLMap.
	// // https://www.pulumi.com/registry/packages/gcp/api-docs/compute/targethttpsproxy/#target-https-proxy-certificate-manager-certificate
	// httpsProxy, err := compute.NewTargetHttpsProxy(ctx, "http-proxy", &compute.TargetHttpsProxyArgs{
	// 	// CertificateManagerCertificates: pulumi.StringArray{pulumi.Sprintf("//certificatemanager.googleapis.com/%v", sslCert.ID())},
	// 	CertificateMap: pulumi.Sprintf("//certificatemanager.googleapis.com/%v", certMap.ID()),
	// 	UrlMap:         httpsUrlMap.SelfLink,
	// })
	// if err != nil {
	// 	return err
	// }

	// // Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
	// _, err = compute.NewGlobalForwardingRule(ctx, "http-forwarding-rule", &compute.GlobalForwardingRuleArgs{
	// 	IpAddress:  ip.Address,
	// 	IpProtocol: pulumi.String("TCP"),
	// 	PortRange:  pulumi.String("443"),
	// 	Target:     httpsProxy.SelfLink,
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a *NitricGcpTerraformProvider) Website(stack cdktf.TerraformStack, name string, config *deploymentspb.Website) error {
	// Deploy a website
	a.Websites[name] = website.NewWebsite(stack, jsii.Sprintf("website_%s", name), &website.WebsiteConfig{
		WebsiteName:    jsii.String(name),
		StackId:        a.Stack.StackIdOutput(),
		BasePath:       jsii.String(config.BasePath),
		LocalDirectory: jsii.String(config.GetLocalDirectory()),
		Region:         jsii.String(a.Region),
		ErrorDocument:  jsii.String(config.GetErrorDocument()),
		IndexDocument:  jsii.String(config.GetIndexDocument()),
	})

	return nil
}
