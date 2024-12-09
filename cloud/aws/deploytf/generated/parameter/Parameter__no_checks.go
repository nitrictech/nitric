//go:build no_runtime_type_checking

package parameter

// Building without runtime type checking enabled, so all the below just return nil

func (p *jsiiProxy_Parameter) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (p *jsiiProxy_Parameter) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (p *jsiiProxy_Parameter) validateGetStringParameters(output *string) error {
	return nil
}

func (p *jsiiProxy_Parameter) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (p *jsiiProxy_Parameter) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateParameter_IsConstructParameters(x interface{}) error {
	return nil
}

func validateParameter_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Parameter) validateSetAccessRoleNamesParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Parameter) validateSetParameterNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Parameter) validateSetParameterValueParameters(val *string) error {
	return nil
}

func validateNewParameterParameters(scope constructs.Construct, id *string, config *ParameterConfig) error {
	return nil
}

