//go:build !no_runtime_type_checking

package stack

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func (s *jsiiProxy_Stack) validateAddOverrideParameters(path *string, value interface{}) error {
	if path == nil {
		return fmt.Errorf("parameter path is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (s *jsiiProxy_Stack) validateAddProviderParameters(provider interface{}) error {
	if provider == nil {
		return fmt.Errorf("parameter provider is required, but nil was provided")
	}
	switch provider.(type) {
	case cdktf.TerraformProvider:
		// ok
	case *cdktf.TerraformModuleProvider:
		provider := provider.(*cdktf.TerraformModuleProvider)
		if err := _jsii_.ValidateStruct(provider, func() string { return "parameter provider" }); err != nil {
			return err
		}
	case cdktf.TerraformModuleProvider:
		provider_ := provider.(cdktf.TerraformModuleProvider)
		provider := &provider_
		if err := _jsii_.ValidateStruct(provider, func() string { return "parameter provider" }); err != nil {
			return err
		}
	default:
		if !_jsii_.IsAnonymousProxy(provider) {
			return fmt.Errorf("parameter provider must be one of the allowed types: cdktf.TerraformProvider, *cdktf.TerraformModuleProvider; received %#v (a %T)", provider, provider)
		}
	}

	return nil
}

func (s *jsiiProxy_Stack) validateGetStringParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (s *jsiiProxy_Stack) validateInterpolationForOutputParameters(moduleOutput *string) error {
	if moduleOutput == nil {
		return fmt.Errorf("parameter moduleOutput is required, but nil was provided")
	}

	return nil
}

func (s *jsiiProxy_Stack) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	if newLogicalId == nil {
		return fmt.Errorf("parameter newLogicalId is required, but nil was provided")
	}

	return nil
}

func validateStack_IsConstructParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validateStack_IsTerraformElementParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func (j *jsiiProxy_Stack) validateSetLocationParameters(val *string) error {
	if val == nil {
		return fmt.Errorf("parameter val is required, but nil was provided")
	}

	return nil
}

func (j *jsiiProxy_Stack) validateSetStackNameParameters(val *string) error {
	if val == nil {
		return fmt.Errorf("parameter val is required, but nil was provided")
	}

	return nil
}

func validateNewStackParameters(scope constructs.Construct, id *string, config *StackConfig) error {
	if scope == nil {
		return fmt.Errorf("parameter scope is required, but nil was provided")
	}

	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if config == nil {
		return fmt.Errorf("parameter config is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(config, func() string { return "parameter config" }); err != nil {
		return err
	}

	return nil
}

