package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference interface {
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
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	DriverName() *string
	SetDriverName(val *string)
	DriverNameInput() *string
	DriverOptions() *map[string]*string
	SetDriverOptions(val *map[string]*string)
	DriverOptionsInput() *map[string]*string
	// Experimental.
	Fqn() *string
	InternalValue() *ServiceTaskSpecContainerSpecMountsVolumeOptions
	SetInternalValue(val *ServiceTaskSpecContainerSpecMountsVolumeOptions)
	Labels() ServiceTaskSpecContainerSpecMountsVolumeOptionsLabelsList
	LabelsInput() interface{}
	NoCopy() interface{}
	SetNoCopy(val interface{})
	NoCopyInput() interface{}
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
	PutLabels(value interface{})
	ResetDriverName()
	ResetDriverOptions()
	ResetLabels()
	ResetNoCopy()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference
type jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) DriverName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driverName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) DriverNameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driverNameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) DriverOptions() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"driverOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) DriverOptionsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"driverOptionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) InternalValue() *ServiceTaskSpecContainerSpecMountsVolumeOptions {
	var returns *ServiceTaskSpecContainerSpecMountsVolumeOptions
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) Labels() ServiceTaskSpecContainerSpecMountsVolumeOptionsLabelsList {
	var returns ServiceTaskSpecContainerSpecMountsVolumeOptionsLabelsList
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) LabelsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) NoCopy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCopy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) NoCopyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCopyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference_Override(s ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetDriverName(val *string) {
	if err := j.validateSetDriverNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"driverName",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetDriverOptions(val *map[string]*string) {
	if err := j.validateSetDriverOptionsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"driverOptions",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetInternalValue(val *ServiceTaskSpecContainerSpecMountsVolumeOptions) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetNoCopy(val interface{}) {
	if err := j.validateSetNoCopyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"noCopy",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) PutLabels(value interface{}) {
	if err := s.validatePutLabelsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putLabels",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ResetDriverName() {
	_jsii_.InvokeVoid(
		s,
		"resetDriverName",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ResetDriverOptions() {
	_jsii_.InvokeVoid(
		s,
		"resetDriverOptions",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ResetLabels() {
	_jsii_.InvokeVoid(
		s,
		"resetLabels",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ResetNoCopy() {
	_jsii_.InvokeVoid(
		s,
		"resetNoCopy",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecMountsVolumeOptionsOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

