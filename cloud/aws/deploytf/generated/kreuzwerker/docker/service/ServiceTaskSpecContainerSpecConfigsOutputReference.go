package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecContainerSpecConfigsOutputReference interface {
	cdktf.ComplexObject
	// the index of the complex object in a list.
	// Experimental.
	ComplexObjectIndex() interface{}
	// Experimental.
	SetComplexObjectIndex(val interface{})
	// set to true if this item is from inside a set and needs tolist() for accessing it set to "0" for single list items.
	// Experimental.
	ComplexObjectIsFromSet() *bool
	// Experimental.
	SetComplexObjectIsFromSet(val *bool)
	ConfigId() *string
	SetConfigId(val *string)
	ConfigIdInput() *string
	ConfigName() *string
	SetConfigName(val *string)
	ConfigNameInput() *string
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	FileGid() *string
	SetFileGid(val *string)
	FileGidInput() *string
	FileMode() *float64
	SetFileMode(val *float64)
	FileModeInput() *float64
	FileName() *string
	SetFileName(val *string)
	FileNameInput() *string
	FileUid() *string
	SetFileUid(val *string)
	FileUidInput() *string
	// Experimental.
	Fqn() *string
	InternalValue() interface{}
	SetInternalValue(val interface{})
	// Experimental.
	TerraformAttribute() *string
	// Experimental.
	SetTerraformAttribute(val *string)
	// Experimental.
	TerraformResource() cdktf.IInterpolatingParent
	// Experimental.
	SetTerraformResource(val cdktf.IInterpolatingParent)
	// Experimental.
	ComputeFqn() *string
	// Experimental.
	GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{}
	// Experimental.
	GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable
	// Experimental.
	GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool
	// Experimental.
	GetListAttribute(terraformAttribute *string) *[]*string
	// Experimental.
	GetNumberAttribute(terraformAttribute *string) *float64
	// Experimental.
	GetNumberListAttribute(terraformAttribute *string) *[]*float64
	// Experimental.
	GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64
	// Experimental.
	GetStringAttribute(terraformAttribute *string) *string
	// Experimental.
	GetStringMapAttribute(terraformAttribute *string) *map[string]*string
	// Experimental.
	InterpolationAsList() cdktf.IResolvable
	// Experimental.
	InterpolationForAttribute(property *string) cdktf.IResolvable
	ResetConfigName()
	ResetFileGid()
	ResetFileMode()
	ResetFileUid()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecContainerSpecConfigsOutputReference
type jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ConfigId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"configId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ConfigIdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"configIdInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ConfigName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"configName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ConfigNameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"configNameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileGid() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileGid",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileGidInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileGidInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileMode() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"fileMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileModeInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"fileModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileNameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileNameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileUid() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileUid",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) FileUidInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileUidInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) InternalValue() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecContainerSpecConfigsOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, complexObjectIndex *float64, complexObjectIsFromSet *bool) ServiceTaskSpecContainerSpecConfigsOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecContainerSpecConfigsOutputReferenceParameters(terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecConfigsOutputReference",
		[]interface{}{terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet},
		&j,
	)

	return &j
}

func NewServiceTaskSpecContainerSpecConfigsOutputReference_Override(s ServiceTaskSpecContainerSpecConfigsOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, complexObjectIndex *float64, complexObjectIsFromSet *bool) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecConfigsOutputReference",
		[]interface{}{terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetConfigId(val *string) {
	if err := j.validateSetConfigIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"configId",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetConfigName(val *string) {
	if err := j.validateSetConfigNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"configName",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetFileGid(val *string) {
	if err := j.validateSetFileGidParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"fileGid",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetFileMode(val *float64) {
	if err := j.validateSetFileModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"fileMode",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetFileName(val *string) {
	if err := j.validateSetFileNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"fileName",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetFileUid(val *string) {
	if err := j.validateSetFileUidParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"fileUid",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetInternalValue(val interface{}) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
	if err := s.validateGetAnyMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"getAnyMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := s.validateGetBooleanAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"getBooleanAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
	if err := s.validateGetBooleanMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*bool

	_jsii_.Invoke(
		s,
		"getBooleanMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
	if err := s.validateGetListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"getListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
	if err := s.validateGetNumberAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"getNumberAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
	if err := s.validateGetNumberListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*float64

	_jsii_.Invoke(
		s,
		"getNumberListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
	if err := s.validateGetNumberMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*float64

	_jsii_.Invoke(
		s,
		"getNumberMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetStringAttribute(terraformAttribute *string) *string {
	if err := s.validateGetStringAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		s,
		"getStringAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
	if err := s.validateGetStringMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*string

	_jsii_.Invoke(
		s,
		"getStringMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
	if err := s.validateInterpolationForAttributeParameters(property); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationForAttribute",
		[]interface{}{property},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ResetConfigName() {
	_jsii_.InvokeVoid(
		s,
		"resetConfigName",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ResetFileGid() {
	_jsii_.InvokeVoid(
		s,
		"resetFileGid",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ResetFileMode() {
	_jsii_.InvokeVoid(
		s,
		"resetFileMode",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ResetFileUid() {
	_jsii_.InvokeVoid(
		s,
		"resetFileUid",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
	if err := s.validateResolveParameters(_context); err != nil {
		panic(err)
	}
	var returns interface{}

	_jsii_.Invoke(
		s,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

