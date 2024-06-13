//go:build no_runtime_type_checking

package secret

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_SecretLabelsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (s *jsiiProxy_SecretLabelsList) validateGetParameters(index *float64) error {
	return nil
}

func (s *jsiiProxy_SecretLabelsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_SecretLabelsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_SecretLabelsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_SecretLabelsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_SecretLabelsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewSecretLabelsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

