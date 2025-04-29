//go:build no_runtime_type_checking

package bucket

// Building without runtime type checking enabled, so all the below just return nil

func (b *jsiiProxy_Bucket) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (b *jsiiProxy_Bucket) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (b *jsiiProxy_Bucket) validateGetStringParameters(output *string) error {
	return nil
}

func (b *jsiiProxy_Bucket) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (b *jsiiProxy_Bucket) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateBucket_IsConstructParameters(x interface{}) error {
	return nil
}

func validateBucket_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Bucket) validateSetListenersParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Bucket) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Bucket) validateSetStorageAccountIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Bucket) validateSetTagsParameters(val *map[string]*string) error {
	return nil
}

func validateNewBucketParameters(scope constructs.Construct, id *string, config *BucketConfig) error {
	return nil
}

