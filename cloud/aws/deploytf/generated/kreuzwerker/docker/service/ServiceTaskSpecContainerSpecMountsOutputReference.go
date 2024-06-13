package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecContainerSpecMountsOutputReference interface {
	cdktf.ComplexObject
	BindOptions() ServiceTaskSpecContainerSpecMountsBindOptionsOutputReference
	BindOptionsInput() *ServiceTaskSpecContainerSpecMountsBindOptions
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
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	// Experimental.
	Fqn() *string
	InternalValue() interface{}
	SetInternalValue(val interface{})
	ReadOnly() interface{}
	SetReadOnly(val interface{})
	ReadOnlyInput() interface{}
	Source() *string
	SetSource(val *string)
	SourceInput() *string
	Target() *string
	SetTarget(val *string)
	TargetInput() *string
	// Experimental.
	TerraformAttribute() *string
	// Experimental.
	SetTerraformAttribute(val *string)
	// Experimental.
	TerraformResource() cdktf.IInterpolatingParent
	// Experimental.
	SetTerraformResource(val cdktf.IInterpolatingParent)
	TmpfsOptions() ServiceTaskSpecContainerSpecMountsTmpfsOptionsOutputReference
	TmpfsOptionsInput() *ServiceTaskSpecContainerSpecMountsTmpfsOptions
	Type() *string
	SetType(val *string)
	TypeInput() *string
	VolumeOptions() ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference
	VolumeOptionsInput() *ServiceTaskSpecContainerSpecMountsVolumeOptions
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
	PutBindOptions(value *ServiceTaskSpecContainerSpecMountsBindOptions)
	PutTmpfsOptions(value *ServiceTaskSpecContainerSpecMountsTmpfsOptions)
	PutVolumeOptions(value *ServiceTaskSpecContainerSpecMountsVolumeOptions)
	ResetBindOptions()
	ResetReadOnly()
	ResetSource()
	ResetTmpfsOptions()
	ResetVolumeOptions()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecContainerSpecMountsOutputReference
type jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) BindOptions() ServiceTaskSpecContainerSpecMountsBindOptionsOutputReference {
	var returns ServiceTaskSpecContainerSpecMountsBindOptionsOutputReference
	_jsii_.Get(
		j,
		"bindOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) BindOptionsInput() *ServiceTaskSpecContainerSpecMountsBindOptions {
	var returns *ServiceTaskSpecContainerSpecMountsBindOptions
	_jsii_.Get(
		j,
		"bindOptionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) InternalValue() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ReadOnly() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnly",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ReadOnlyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnlyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) SourceInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) Target() *string {
	var returns *string
	_jsii_.Get(
		j,
		"target",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TargetInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TmpfsOptions() ServiceTaskSpecContainerSpecMountsTmpfsOptionsOutputReference {
	var returns ServiceTaskSpecContainerSpecMountsTmpfsOptionsOutputReference
	_jsii_.Get(
		j,
		"tmpfsOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TmpfsOptionsInput() *ServiceTaskSpecContainerSpecMountsTmpfsOptions {
	var returns *ServiceTaskSpecContainerSpecMountsTmpfsOptions
	_jsii_.Get(
		j,
		"tmpfsOptionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) TypeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"typeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) VolumeOptions() ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference {
	var returns ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference
	_jsii_.Get(
		j,
		"volumeOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) VolumeOptionsInput() *ServiceTaskSpecContainerSpecMountsVolumeOptions {
	var returns *ServiceTaskSpecContainerSpecMountsVolumeOptions
	_jsii_.Get(
		j,
		"volumeOptionsInput",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecContainerSpecMountsOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, complexObjectIndex *float64, complexObjectIsFromSet *bool) ServiceTaskSpecContainerSpecMountsOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecContainerSpecMountsOutputReferenceParameters(terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecMountsOutputReference",
		[]interface{}{terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet},
		&j,
	)

	return &j
}

func NewServiceTaskSpecContainerSpecMountsOutputReference_Override(s ServiceTaskSpecContainerSpecMountsOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, complexObjectIndex *float64, complexObjectIsFromSet *bool) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecMountsOutputReference",
		[]interface{}{terraformResource, terraformAttribute, complexObjectIndex, complexObjectIsFromSet},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetInternalValue(val interface{}) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetReadOnly(val interface{}) {
	if err := j.validateSetReadOnlyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"readOnly",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetSource(val *string) {
	if err := j.validateSetSourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"source",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetTarget(val *string) {
	if err := j.validateSetTargetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"target",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference)SetType(val *string) {
	if err := j.validateSetTypeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"type",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) PutBindOptions(value *ServiceTaskSpecContainerSpecMountsBindOptions) {
	if err := s.validatePutBindOptionsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putBindOptions",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) PutTmpfsOptions(value *ServiceTaskSpecContainerSpecMountsTmpfsOptions) {
	if err := s.validatePutTmpfsOptionsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putTmpfsOptions",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) PutVolumeOptions(value *ServiceTaskSpecContainerSpecMountsVolumeOptions) {
	if err := s.validatePutVolumeOptionsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putVolumeOptions",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ResetBindOptions() {
	_jsii_.InvokeVoid(
		s,
		"resetBindOptions",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ResetReadOnly() {
	_jsii_.InvokeVoid(
		s,
		"resetReadOnly",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ResetSource() {
	_jsii_.InvokeVoid(
		s,
		"resetSource",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ResetTmpfsOptions() {
	_jsii_.InvokeVoid(
		s,
		"resetTmpfsOptions",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ResetVolumeOptions() {
	_jsii_.InvokeVoid(
		s,
		"resetVolumeOptions",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

