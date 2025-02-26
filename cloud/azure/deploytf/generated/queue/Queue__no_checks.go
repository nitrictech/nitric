//go:build no_runtime_type_checking

package queue

// Building without runtime type checking enabled, so all the below just return nil

func (q *jsiiProxy_Queue) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (q *jsiiProxy_Queue) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (q *jsiiProxy_Queue) validateGetStringParameters(output *string) error {
	return nil
}

func (q *jsiiProxy_Queue) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (q *jsiiProxy_Queue) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateQueue_IsConstructParameters(x interface{}) error {
	return nil
}

func validateQueue_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Queue) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Queue) validateSetStorageAccountNameParameters(val *string) error {
	return nil
}

func validateNewQueueParameters(scope constructs.Construct, id *string, config *QueueConfig) error {
	return nil
}

