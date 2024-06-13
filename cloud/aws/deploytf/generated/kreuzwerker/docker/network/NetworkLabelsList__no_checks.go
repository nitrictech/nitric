//go:build no_runtime_type_checking

package network

// Building without runtime type checking enabled, so all the below just return nil

func (n *jsiiProxy_NetworkLabelsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (n *jsiiProxy_NetworkLabelsList) validateGetParameters(index *float64) error {
	return nil
}

func (n *jsiiProxy_NetworkLabelsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_NetworkLabelsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_NetworkLabelsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_NetworkLabelsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_NetworkLabelsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewNetworkLabelsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

