package http_proxy

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type HttpProxyConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The identity of the app.
	AppIdentity *string `field:"required" json:"appIdentity" yaml:"appIdentity"`
	// The description of the API.
	Description *string `field:"required" json:"description" yaml:"description"`
	// The location of the API.
	Location *string `field:"required" json:"location" yaml:"location"`
	// The name of the API.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The openapi spec to deploy.
	OpenapiSpec *string `field:"required" json:"openapiSpec" yaml:"openapiSpec"`
	// The policy templates to apply The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	OperationPolicyTemplates *map[string]*string `field:"required" json:"operationPolicyTemplates" yaml:"operationPolicyTemplates"`
	// The email of the publisher.
	PublisherEmail *string `field:"required" json:"publisherEmail" yaml:"publisherEmail"`
	// The name of the publisher.
	PublisherName *string `field:"required" json:"publisherName" yaml:"publisherName"`
	// The name of the resource group.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
}

