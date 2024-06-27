//go:build no_runtime_type_checking

package rds

// Building without runtime type checking enabled, so all the below just return nil

func (r *jsiiProxy_Rds) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (r *jsiiProxy_Rds) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (r *jsiiProxy_Rds) validateGetStringParameters(output *string) error {
	return nil
}

func (r *jsiiProxy_Rds) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (r *jsiiProxy_Rds) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateRds_IsConstructParameters(x interface{}) error {
	return nil
}

func validateRds_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Rds) validateSetMaxCapacityParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Rds) validateSetMinCapacityParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Rds) validateSetPrivateSubnetIdsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Rds) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Rds) validateSetVpcIdParameters(val *string) error {
	return nil
}

func validateNewRdsParameters(scope constructs.Construct, id *string, config *RdsConfig) error {
	return nil
}

