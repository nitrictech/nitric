//go:build no_runtime_type_checking

package http

// Building without runtime type checking enabled, so all the below just return nil

func (h *jsiiProxy_Http) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (h *jsiiProxy_Http) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (h *jsiiProxy_Http) validateGetStringParameters(output *string) error {
	return nil
}

func (h *jsiiProxy_Http) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (h *jsiiProxy_Http) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateHttp_IsConstructParameters(x interface{}) error {
	return nil
}

func validateHttp_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Http) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Http) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Http) validateSetTargetLambdaFunctionParameters(val *string) error {
	return nil
}

func validateNewHttpParameters(scope constructs.Construct, id *string, config *HttpConfig) error {
	return nil
}

