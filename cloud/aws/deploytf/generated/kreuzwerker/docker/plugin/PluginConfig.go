package plugin

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type PluginConfig struct {
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
	// Docker Plugin name.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#name Plugin#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// Docker Plugin alias.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#alias Plugin#alias}
	Alias *string `field:"optional" json:"alias" yaml:"alias"`
	// If `true` the plugin is enabled. Defaults to `true`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#enabled Plugin#enabled}
	Enabled interface{} `field:"optional" json:"enabled" yaml:"enabled"`
	// HTTP client timeout to enable the plugin.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#enable_timeout Plugin#enable_timeout}
	EnableTimeout *float64 `field:"optional" json:"enableTimeout" yaml:"enableTimeout"`
	// The environment variables in the form of `KEY=VALUE`, e.g. `DEBUG=0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#env Plugin#env}
	Env *[]*string `field:"optional" json:"env" yaml:"env"`
	// If true, then the plugin is destroyed forcibly.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#force_destroy Plugin#force_destroy}
	ForceDestroy interface{} `field:"optional" json:"forceDestroy" yaml:"forceDestroy"`
	// If true, then the plugin is disabled forcibly.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#force_disable Plugin#force_disable}
	ForceDisable interface{} `field:"optional" json:"forceDisable" yaml:"forceDisable"`
	// If true, grant all permissions necessary to run the plugin.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#grant_all_permissions Plugin#grant_all_permissions}
	GrantAllPermissions interface{} `field:"optional" json:"grantAllPermissions" yaml:"grantAllPermissions"`
	// grant_permissions block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#grant_permissions Plugin#grant_permissions}
	GrantPermissions interface{} `field:"optional" json:"grantPermissions" yaml:"grantPermissions"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#id Plugin#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
}

