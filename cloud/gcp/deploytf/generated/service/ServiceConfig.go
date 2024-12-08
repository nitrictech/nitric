package service

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ServiceConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The base URI for the artifact registry repository the push this services image to.
	ArtifactRegistryRepository *string `field:"required" json:"artifactRegistryRepository" yaml:"artifactRegistryRepository"`
	// The base compute role to use for the service.
	BaseComputeRole *string `field:"required" json:"baseComputeRole" yaml:"baseComputeRole"`
	// Environment variables to set on the lambda function The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	Environment *map[string]*string `field:"required" json:"environment" yaml:"environment"`
	// The docker image to deploy.
	Image *string `field:"required" json:"image" yaml:"image"`
	// The ID of the Google Cloud project where the service is created.
	ProjectId *string `field:"required" json:"projectId" yaml:"projectId"`
	// The region the service is being deployed to.
	Region *string `field:"required" json:"region" yaml:"region"`
	// The name of the service.
	ServiceName *string `field:"required" json:"serviceName" yaml:"serviceName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// cloud run service directvpc config.
	Vpc interface{} `field:"required" json:"vpc" yaml:"vpc"`
	// The number of concurrent requests the CloudRun service can handle 80.
	ContainerConcurrency *float64 `field:"optional" json:"containerConcurrency" yaml:"containerConcurrency"`
	// The amount of cpus to allocate to the CloudRun service 1.
	Cpus *float64 `field:"optional" json:"cpus" yaml:"cpus"`
	// The amount of gpus to allocate to the CloudRun service.
	Gpus *float64 `field:"optional" json:"gpus" yaml:"gpus"`
	// The name of the KMS key to use.
	KmsKey *string `field:"optional" json:"kmsKey" yaml:"kmsKey"`
	// The maximum number of instances to run 10.
	MaxInstances *float64 `field:"optional" json:"maxInstances" yaml:"maxInstances"`
	// The amount of memory to allocate to the CloudRun service in MB 512.
	MemoryMb *float64 `field:"optional" json:"memoryMb" yaml:"memoryMb"`
	// The minimum number of instances to run.
	MinInstances *float64 `field:"optional" json:"minInstances" yaml:"minInstances"`
	// The timeout for the CloudRun service in seconds 10.
	TimeoutSeconds *float64 `field:"optional" json:"timeoutSeconds" yaml:"timeoutSeconds"`
}

