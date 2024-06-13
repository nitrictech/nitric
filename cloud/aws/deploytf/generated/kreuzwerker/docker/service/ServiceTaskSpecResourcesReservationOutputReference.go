package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecResourcesReservationOutputReference interface {
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
	// Experimental.
	Fqn() *string
	GenericResources() ServiceTaskSpecResourcesReservationGenericResourcesOutputReference
	GenericResourcesInput() *ServiceTaskSpecResourcesReservationGenericResources
	InternalValue() *ServiceTaskSpecResourcesReservation
	SetInternalValue(val *ServiceTaskSpecResourcesReservation)
	MemoryBytes() *float64
	SetMemoryBytes(val *float64)
	MemoryBytesInput() *float64
	NanoCpus() *float64
	SetNanoCpus(val *float64)
	NanoCpusInput() *float64
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
	PutGenericResources(value *ServiceTaskSpecResourcesReservationGenericResources)
	ResetGenericResources()
	ResetMemoryBytes()
	ResetNanoCpus()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecResourcesReservationOutputReference
type jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GenericResources() ServiceTaskSpecResourcesReservationGenericResourcesOutputReference {
	var returns ServiceTaskSpecResourcesReservationGenericResourcesOutputReference
	_jsii_.Get(
		j,
		"genericResources",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GenericResourcesInput() *ServiceTaskSpecResourcesReservationGenericResources {
	var returns *ServiceTaskSpecResourcesReservationGenericResources
	_jsii_.Get(
		j,
		"genericResourcesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) InternalValue() *ServiceTaskSpecResourcesReservation {
	var returns *ServiceTaskSpecResourcesReservation
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) MemoryBytes() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memoryBytes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) MemoryBytesInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memoryBytesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) NanoCpus() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"nanoCpus",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) NanoCpusInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"nanoCpusInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecResourcesReservationOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceTaskSpecResourcesReservationOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecResourcesReservationOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecResourcesReservationOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceTaskSpecResourcesReservationOutputReference_Override(s ServiceTaskSpecResourcesReservationOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecResourcesReservationOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetInternalValue(val *ServiceTaskSpecResourcesReservation) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetMemoryBytes(val *float64) {
	if err := j.validateSetMemoryBytesParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"memoryBytes",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetNanoCpus(val *float64) {
	if err := j.validateSetNanoCpusParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"nanoCpus",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) PutGenericResources(value *ServiceTaskSpecResourcesReservationGenericResources) {
	if err := s.validatePutGenericResourcesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putGenericResources",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ResetGenericResources() {
	_jsii_.InvokeVoid(
		s,
		"resetGenericResources",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ResetMemoryBytes() {
	_jsii_.InvokeVoid(
		s,
		"resetMemoryBytes",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ResetNanoCpus() {
	_jsii_.InvokeVoid(
		s,
		"resetNanoCpus",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecResourcesReservationOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

