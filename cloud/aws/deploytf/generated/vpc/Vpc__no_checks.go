//go:build no_runtime_type_checking

package vpc

// Building without runtime type checking enabled, so all the below just return nil

func (v *jsiiProxy_Vpc) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (v *jsiiProxy_Vpc) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (v *jsiiProxy_Vpc) validateGetStringParameters(output *string) error {
	return nil
}

func (v *jsiiProxy_Vpc) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (v *jsiiProxy_Vpc) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateVpc_IsConstructParameters(x interface{}) error {
	return nil
}

func validateVpc_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateNewVpcParameters(scope constructs.Construct, id *string, config *VpcConfig) error {
	return nil
}

