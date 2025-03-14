// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploytf

import (
	"slices"

	"github.com/aws/jsii-runtime-go"
	acmcertificate "github.com/cdktf/cdktf-provider-aws-go/aws/v19/acmcertificate"
	acmcertificatevalidation "github.com/cdktf/cdktf-provider-aws-go/aws/v19/acmcertificatevalidation"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v19/provider"
	route53record "github.com/cdktf/cdktf-provider-aws-go/aws/v19/route53record"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/cdn"
	"github.com/samber/lo"
)

type ApiGateway struct {
	GatewayURL *string `json:"gateway_url"`
}

type WebsiteConfig struct {
	BucketDomainName *string    `json:"bucket_domain_name"`
	BucketArn        *string    `json:"bucket_arn"`
	BucketId         *string    `json:"bucket_id"`
	BasePath         *string    `json:"base_path"`
	ChangedFiles     *[]*string `json:"changed_files"`
}

type RootWebsite struct {
	Name          *string `json:"name"`
	IndexDocument *string `json:"index_document"`
	ErrorDocument *string `json:"error_document"`
}

// function to create a new cdn
func (a *NitricAwsTerraformProvider) NewCdn(tfstack cdktf.TerraformStack) cdn.Cdn {
	apiGateways := make(map[string]ApiGateway)

	sortedApiKeys := lo.Keys(a.Apis)
	slices.Sort(sortedApiKeys)

	for _, apiName := range sortedApiKeys {
		api := a.Apis[apiName]
		apiGateways[apiName] = ApiGateway{
			GatewayURL: api.EndpointOutput(),
		}
	}

	websites := make(map[string]WebsiteConfig)
	for websiteName, website := range a.Websites {
		websiteFiles := cdktf.Token_AsList(website.ChangedFilesOutput(), nil)

		websites[websiteName] = WebsiteConfig{
			BucketDomainName: website.WebsiteBucketDomainOutput(),
			BucketArn:        website.WebsiteArnOutput(),
			BucketId:         website.WebsiteIdOutput(),
			BasePath:         website.BasePath(),
			ChangedFiles:     websiteFiles,
		}
	}

	domainName := a.AwsConfig.Cdn.Domain
	var certificateArn *string
	var zoneId *string
	if domainName != "" {
		zoneId = getZoneIds([]string{domainName})[domainName]

		// ACM Provider in us-east-1
		acmProvider := awsprovider.NewAwsProvider(tfstack, jsii.String("AWSUsEast1"), &awsprovider.AwsProviderConfig{
			Region: jsii.String("us-east-1"),
			Alias:  jsii.String("us-east-1"),
		})

		// ACM Certificate (must be in us-east-1)
		cert := acmcertificate.NewAcmCertificate(tfstack, jsii.String("CdnCert"), &acmcertificate.AcmCertificateConfig{
			DomainName:       jsii.String(domainName),
			ValidationMethod: jsii.String("DNS"),
			Provider:         acmProvider, // Ensure ACM is deployed in us-east-1
		})

		// Route 53 Record for DNS validation (remains in the main region)
		validationRecord := route53record.NewRoute53Record(tfstack, jsii.String("CdnCertValidation"), &route53record.Route53RecordConfig{
			ZoneId: zoneId,
			Name:   cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordName(),
			Type:   cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordType(),
			Records: &[]*string{
				cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordValue(),
			},
			Ttl: jsii.Number(600),
		})

		// ACM Certificate Validation (must be in us-east-1)
		validation := acmcertificatevalidation.NewAcmCertificateValidation(tfstack, jsii.String("CertValidation"), &acmcertificatevalidation.AcmCertificateValidationConfig{
			CertificateArn: cert.Arn(),
			ValidationRecordFqdns: &[]*string{
				validationRecord.Fqdn(),
			},
			Provider: acmProvider, // Use us-east-1 provider
		})

		certificateArn = validation.CertificateArn()
	}

	return cdn.NewCdn(tfstack, jsii.String("cdn"), &cdn.CdnConfig{
		StackName:             a.Stack.StackIdOutput(),
		Websites:              websites,
		Apis:                  apiGateways,
		RootWebsite:           &a.RootWebsite,
		CertificateArn:        certificateArn,
		DomainName:            jsii.String(domainName),
		SkipCacheInvalidation: jsii.Bool(a.AwsConfig.Cdn.SkipCacheInvalidation),
		ZoneId:                zoneId,
		DependsOn:             &[]cdktf.ITerraformDependable{a.Stack},
	})
}
