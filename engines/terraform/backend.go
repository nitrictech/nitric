package terraform

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

// NilTerraformBackend prevents cdktf from automatically defining a backend in the output Terraform configuration.
type NilTerraformBackend struct {
	cdktf.TerraformBackend
}

// ToTerraform returns empty config to ensure an undefined backend in the output.
func (t *NilTerraformBackend) ToTerraform() interface{} {
	return map[string]interface{}{}
}

// NewNilTerraformBackend creates a nil backend that prevents the backend field from being defined in the output terraform.
// We have this here as there is no way to disable the backend in cdktf.
// https://github.com/hashicorp/terraform-cdk/issues/2435
func NewNilTerraformBackend(app constructs.Construct, name *string) *NilTerraformBackend {
	backend := &NilTerraformBackend{}

	cdktf.NewTerraformBackend_Override(backend, app, jsii.String("nil"), name)

	return backend
}
