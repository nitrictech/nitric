package service

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ServiceConfig struct {
	// Experimental.
	Connection interface{} `field:"optional" json:"connection" yaml:"connection"`
	// Experimental.
	Count interface{} `field:"optional" json:"count" yaml:"count"`
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Lifecycle *cdktf.TerraformResourceLifecycle `field:"optional" json:"lifecycle" yaml:"lifecycle"`
	// Experimental.
	Provider cdktf.TerraformProvider `field:"optional" json:"provider" yaml:"provider"`
	// Experimental.
	Provisioners *[]interface{} `field:"optional" json:"provisioners" yaml:"provisioners"`
	// Name of the service.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#name Service#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// task_spec block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#task_spec Service#task_spec}
	TaskSpec *ServiceTaskSpec `field:"required" json:"taskSpec" yaml:"taskSpec"`
	// auth block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#auth Service#auth}
	Auth *ServiceAuth `field:"optional" json:"auth" yaml:"auth"`
	// converge_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#converge_config Service#converge_config}
	ConvergeConfig *ServiceConvergeConfig `field:"optional" json:"convergeConfig" yaml:"convergeConfig"`
	// endpoint_spec block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#endpoint_spec Service#endpoint_spec}
	EndpointSpec *ServiceEndpointSpec `field:"optional" json:"endpointSpec" yaml:"endpointSpec"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#id Service#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
	// labels block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#labels Service#labels}
	Labels interface{} `field:"optional" json:"labels" yaml:"labels"`
	// mode block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#mode Service#mode}
	Mode *ServiceMode `field:"optional" json:"mode" yaml:"mode"`
	// rollback_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#rollback_config Service#rollback_config}
	RollbackConfig *ServiceRollbackConfig `field:"optional" json:"rollbackConfig" yaml:"rollbackConfig"`
	// update_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#update_config Service#update_config}
	UpdateConfig *ServiceUpdateConfig `field:"optional" json:"updateConfig" yaml:"updateConfig"`
}

