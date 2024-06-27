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

//go:build !no_runtime_type_checking

package keyvalue

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func (k *jsiiProxy_Keyvalue) validateAddOverrideParameters(path *string, value interface{}) error {
	if path == nil {
		return fmt.Errorf("parameter path is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (k *jsiiProxy_Keyvalue) validateAddProviderParameters(provider interface{}) error {
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

func (k *jsiiProxy_Keyvalue) validateGetStringParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (k *jsiiProxy_Keyvalue) validateInterpolationForOutputParameters(moduleOutput *string) error {
	if moduleOutput == nil {
		return fmt.Errorf("parameter moduleOutput is required, but nil was provided")
	}

	return nil
}

func (k *jsiiProxy_Keyvalue) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	if newLogicalId == nil {
		return fmt.Errorf("parameter newLogicalId is required, but nil was provided")
	}

	return nil
}

func validateKeyvalue_IsConstructParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validateKeyvalue_IsTerraformElementParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func (j *jsiiProxy_Keyvalue) validateSetKvstoreNameParameters(val *string) error {
	if val == nil {
		return fmt.Errorf("parameter val is required, but nil was provided")
	}

	return nil
}

func (j *jsiiProxy_Keyvalue) validateSetStackIdParameters(val *string) error {
	if val == nil {
		return fmt.Errorf("parameter val is required, but nil was provided")
	}

	return nil
}

func validateNewKeyvalueParameters(scope constructs.Construct, id *string, config *KeyvalueConfig) error {
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
