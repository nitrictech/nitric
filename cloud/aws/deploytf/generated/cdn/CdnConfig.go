package cdn

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CdnConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// The ARN for the website bucket.
	WebsiteBucketArn *string `field:"required" json:"websiteBucketArn" yaml:"websiteBucketArn"`
	// The domain name for the website bucket.
	WebsiteBucketDomainName *string `field:"required" json:"websiteBucketDomainName" yaml:"websiteBucketDomainName"`
	// The ID for the website bucket.
	WebsiteBucketId *string `field:"required" json:"websiteBucketId" yaml:"websiteBucketId"`
	// Map of APIs and their gateway information.
	Apis interface{} `field:"optional" json:"apis" yaml:"apis"`
	// The website error document 404.html.
	WebsiteErrorDocument *string `field:"optional" json:"websiteErrorDocument" yaml:"websiteErrorDocument"`
	// The website index document index.html.
	WebsiteIndexDocument *string `field:"optional" json:"websiteIndexDocument" yaml:"websiteIndexDocument"`
}

