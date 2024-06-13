package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceUpdateConfigOutputReference interface {
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
	Delay() *string
	SetDelay(val *string)
	DelayInput() *string
	FailureAction() *string
	SetFailureAction(val *string)
	FailureActionInput() *string
	// Experimental.
	Fqn() *string
	InternalValue() *ServiceUpdateConfig
	SetInternalValue(val *ServiceUpdateConfig)
	MaxFailureRatio() *string
	SetMaxFailureRatio(val *string)
	MaxFailureRatioInput() *string
	Monitor() *string
	SetMonitor(val *string)
	MonitorInput() *string
	Order() *string
	SetOrder(val *string)
	OrderInput() *string
	Parallelism() *float64
	SetParallelism(val *float64)
	ParallelismInput() *float64
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
	ResetDelay()
	ResetFailureAction()
	ResetMaxFailureRatio()
	ResetMonitor()
	ResetOrder()
	ResetParallelism()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceUpdateConfigOutputReference
type jsiiProxy_ServiceUpdateConfigOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) Delay() *string {
	var returns *string
	_jsii_.Get(
		j,
		"delay",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) DelayInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"delayInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) FailureAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"failureAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) FailureActionInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"failureActionInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) InternalValue() *ServiceUpdateConfig {
	var returns *ServiceUpdateConfig
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) MaxFailureRatio() *string {
	var returns *string
	_jsii_.Get(
		j,
		"maxFailureRatio",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) MaxFailureRatioInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"maxFailureRatioInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) Monitor() *string {
	var returns *string
	_jsii_.Get(
		j,
		"monitor",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) MonitorInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"monitorInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) Order() *string {
	var returns *string
	_jsii_.Get(
		j,
		"order",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) OrderInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"orderInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) Parallelism() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"parallelism",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) ParallelismInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"parallelismInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceUpdateConfigOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceUpdateConfigOutputReference {
	_init_.Initialize()

	if err := validateNewServiceUpdateConfigOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceUpdateConfigOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceUpdateConfigOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceUpdateConfigOutputReference_Override(s ServiceUpdateConfigOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceUpdateConfigOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetDelay(val *string) {
	if err := j.validateSetDelayParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"delay",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetFailureAction(val *string) {
	if err := j.validateSetFailureActionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"failureAction",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetInternalValue(val *ServiceUpdateConfig) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetMaxFailureRatio(val *string) {
	if err := j.validateSetMaxFailureRatioParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"maxFailureRatio",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetMonitor(val *string) {
	if err := j.validateSetMonitorParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"monitor",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetOrder(val *string) {
	if err := j.validateSetOrderParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"order",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetParallelism(val *float64) {
	if err := j.validateSetParallelismParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"parallelism",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceUpdateConfigOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetDelay() {
	_jsii_.InvokeVoid(
		s,
		"resetDelay",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetFailureAction() {
	_jsii_.InvokeVoid(
		s,
		"resetFailureAction",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetMaxFailureRatio() {
	_jsii_.InvokeVoid(
		s,
		"resetMaxFailureRatio",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetMonitor() {
	_jsii_.InvokeVoid(
		s,
		"resetMonitor",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetOrder() {
	_jsii_.InvokeVoid(
		s,
		"resetOrder",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ResetParallelism() {
	_jsii_.InvokeVoid(
		s,
		"resetParallelism",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceUpdateConfigOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

