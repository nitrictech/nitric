//go:build no_runtime_type_checking

package container

// Building without runtime type checking enabled, so all the below just return nil

func (c *jsiiProxy_Container) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateImportFromParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validateMoveToIdParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Container) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutCapabilitiesParameters(value *ContainerCapabilities) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutDevicesParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutHealthcheckParameters(value *ContainerHealthcheck) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutHostParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutLabelsParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutMountsParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutNetworksAdvancedParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutPortsParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutUlimitParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutUploadParameters(value interface{}) error {
	return nil
}

func (c *jsiiProxy_Container) validatePutVolumesParameters(value interface{}) error {
	return nil
}

func validateContainer_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateContainer_IsConstructParameters(x interface{}) error {
	return nil
}

func validateContainer_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateContainer_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetAttachParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetCgroupnsModeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetCommandParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetContainerReadRefreshTimeoutMillisecondsParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetCpuSetParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetCpuSharesParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetDestroyGraceSecondsParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetDnsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetDnsOptsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetDnsSearchParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetDomainnameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetEntrypointParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetEnvParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetGpusParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetGroupAddParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetHostnameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetImageParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetInitParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetIpcModeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetLogDriverParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetLogOptsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetLogsParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetMaxRetryCountParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetMemoryParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetMemorySwapParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetMustRunParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetNetworkModeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetPidModeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetPrivilegedParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetPublishAllPortsParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetReadOnlyParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetRemoveVolumesParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetRestartParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetRmParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetRuntimeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetSecurityOptsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetShmSizeParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetStartParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetStdinOpenParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetStopSignalParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetStopTimeoutParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetStorageOptsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetSysctlsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetTmpfsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetTtyParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetUserParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetUsernsModeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetWaitParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetWaitTimeoutParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Container) validateSetWorkingDirParameters(val *string) error {
	return nil
}

func validateNewContainerParameters(scope constructs.Construct, id *string, config *ContainerConfig) error {
	return nil
}

