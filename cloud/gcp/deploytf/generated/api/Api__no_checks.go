//go:build no_runtime_type_checking

package api

// Building without runtime type checking enabled, so all the below just return nil

func (a *jsiiProxy_Api) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (a *jsiiProxy_Api) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (a *jsiiProxy_Api) validateGetStringParameters(output *string) error {
	return nil
}

func (a *jsiiProxy_Api) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (a *jsiiProxy_Api) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateApi_IsConstructParameters(x interface{}) error {
	return nil
}

func validateApi_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Api) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Api) validateSetOpenapiSpecParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Api) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Api) validateSetTargetServicesParameters(val *map[string]*string) error {
	return nil
}

func validateNewApiParameters(scope constructs.Construct, id *string, config *ApiConfig) error {
	return nil
}

