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

