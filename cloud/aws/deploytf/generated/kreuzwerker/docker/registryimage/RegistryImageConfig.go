package registryimage

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type RegistryImageConfig struct {
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
	// The name of the Docker image.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/registry_image#name RegistryImage#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/registry_image#id RegistryImage#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
	// If `true`, the verification of TLS certificates of the server/registry is disabled. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/registry_image#insecure_skip_verify RegistryImage#insecure_skip_verify}
	InsecureSkipVerify interface{} `field:"optional" json:"insecureSkipVerify" yaml:"insecureSkipVerify"`
	// If true, then the Docker image won't be deleted on destroy operation.
	//
	// If this is false, it will delete the image from the docker registry on destroy operation. Defaults to `false`
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/registry_image#keep_remotely RegistryImage#keep_remotely}
	KeepRemotely interface{} `field:"optional" json:"keepRemotely" yaml:"keepRemotely"`
	// A map of arbitrary strings that, when changed, will force the `docker_registry_image` resource to be replaced.
	//
	// This can be used to repush a local image
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/registry_image#triggers RegistryImage#triggers}
	Triggers *map[string]*string `field:"optional" json:"triggers" yaml:"triggers"`
}

