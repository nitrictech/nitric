//go:build no_runtime_type_checking

package roles

// Building without runtime type checking enabled, so all the below just return nil

func (r *jsiiProxy_Roles) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (r *jsiiProxy_Roles) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (r *jsiiProxy_Roles) validateGetStringParameters(output *string) error {
	return nil
}

func (r *jsiiProxy_Roles) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (r *jsiiProxy_Roles) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateRoles_IsConstructParameters(x interface{}) error {
	return nil
}

func validateRoles_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Roles) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Roles) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewRolesParameters(scope constructs.Construct, id *string, config *RolesConfig) error {
	return nil
}

