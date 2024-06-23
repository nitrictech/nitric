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

//go:build no_runtime_type_checking

package keyvalue

// Building without runtime type checking enabled, so all the below just return nil

func (k *jsiiProxy_Keyvalue) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (k *jsiiProxy_Keyvalue) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (k *jsiiProxy_Keyvalue) validateGetStringParameters(output *string) error {
	return nil
}

func (k *jsiiProxy_Keyvalue) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (k *jsiiProxy_Keyvalue) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateKeyvalue_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKeyvalue_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Keyvalue) validateSetKvstoreNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Keyvalue) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewKeyvalueParameters(scope constructs.Construct, id *string, config *KeyvalueConfig) error {
	return nil
}

