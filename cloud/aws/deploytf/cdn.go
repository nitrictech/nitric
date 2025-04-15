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

	return cdn.NewCdn(tfstack, jsii.String("cdn"), &cdn.CdnConfig{
		StackName:             a.Stack.StackIdOutput(),
		Websites:              websites,
		Apis:                  apiGateways,
		RootWebsite:           &a.RootWebsite,
		SkipCacheInvalidation: jsii.Bool(a.AwsConfig.Cdn.SkipCacheInvalidation),
		DependsOn:             &[]cdktf.ITerraformDependable{a.Stack},
	})
}
