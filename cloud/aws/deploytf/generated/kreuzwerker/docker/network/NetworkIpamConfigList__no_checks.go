//go:build no_runtime_type_checking

package network

// Building without runtime type checking enabled, so all the below just return nil

func (n *jsiiProxy_NetworkIpamConfigList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (n *jsiiProxy_NetworkIpamConfigList) validateGetParameters(index *float64) error {
	return nil
}

func (n *jsiiProxy_NetworkIpamConfigList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_NetworkIpamConfigList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_NetworkIpamConfigList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_NetworkIpamConfigList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_NetworkIpamConfigList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewNetworkIpamConfigListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

