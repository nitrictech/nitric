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
	// The id of the container app environment.
	ContainerAppEnvironmentId *string `field:"required" json:"containerAppEnvironmentId" yaml:"containerAppEnvironmentId"`
	// The cpu limit for the container.
	Cpu *float64 `field:"required" json:"cpu" yaml:"cpu"`
	// The environment variables to set The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	Env *map[string]*string `field:"required" json:"env" yaml:"env"`
	// The image uri for the container.
	ImageUri *string `field:"required" json:"imageUri" yaml:"imageUri"`
	// Maximum number of replicas for the service.
	MaxReplicas *float64 `field:"required" json:"maxReplicas" yaml:"maxReplicas"`
	// The memory limit for the container.
	Memory *string `field:"required" json:"memory" yaml:"memory"`
	// Minimum number of replicas for the service.
	MinReplicas *float64 `field:"required" json:"minReplicas" yaml:"minReplicas"`
	// The name of the service.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The login server for the container registry.
	RegistryLoginServer *string `field:"required" json:"registryLoginServer" yaml:"registryLoginServer"`
	// The password for the container registry.
	RegistryPassword *string `field:"required" json:"registryPassword" yaml:"registryPassword"`
	// The username for the container registry.
	RegistryUsername *string `field:"required" json:"registryUsername" yaml:"registryUsername"`
	// The name of the resource group.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// The tags to apply to the service The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	Tags *map[string]*string `field:"required" json:"tags" yaml:"tags"`
}

