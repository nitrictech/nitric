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

package http_proxy

// Building without runtime type checking enabled, so all the below just return nil

func (h *jsiiProxy_HttpProxy) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (h *jsiiProxy_HttpProxy) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (h *jsiiProxy_HttpProxy) validateGetStringParameters(output *string) error {
	return nil
}

func (h *jsiiProxy_HttpProxy) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (h *jsiiProxy_HttpProxy) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateHttpProxy_IsConstructParameters(x interface{}) error {
	return nil
}

func validateHttpProxy_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetInvokerEmailParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetTargetServiceUrlParameters(val *string) error {
	return nil
}

func validateNewHttpProxyParameters(scope constructs.Construct, id *string, config *HttpProxyConfig) error {
	return nil
}

