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

func (j *jsiiProxy_HttpProxy) validateSetAppIdentityParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetDescriptionParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetLocationParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetOpenapiSpecParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetOperationPolicyTemplatesParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetPublisherEmailParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetPublisherNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_HttpProxy) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func validateNewHttpProxyParameters(scope constructs.Construct, id *string, config *HttpProxyConfig) error {
	return nil
}

