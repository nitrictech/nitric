//go:build no_runtime_type_checking

package topic

// Building without runtime type checking enabled, so all the below just return nil

func (t *jsiiProxy_Topic) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (t *jsiiProxy_Topic) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (t *jsiiProxy_Topic) validateGetStringParameters(output *string) error {
	return nil
}

func (t *jsiiProxy_Topic) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (t *jsiiProxy_Topic) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateTopic_IsConstructParameters(x interface{}) error {
	return nil
}

func validateTopic_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Topic) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Topic) validateSetSubscriberServicesParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Topic) validateSetTopicNameParameters(val *string) error {
	return nil
}

func validateNewTopicParameters(scope constructs.Construct, id *string, config *TopicConfig) error {
	return nil
}

