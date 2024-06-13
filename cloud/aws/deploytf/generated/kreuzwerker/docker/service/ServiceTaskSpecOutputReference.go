package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecOutputReference interface {
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
	ContainerSpec() ServiceTaskSpecContainerSpecOutputReference
	ContainerSpecInput() *ServiceTaskSpecContainerSpec
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	ForceUpdate() *float64
	SetForceUpdate(val *float64)
	ForceUpdateInput() *float64
	// Experimental.
	Fqn() *string
	InternalValue() *ServiceTaskSpec
	SetInternalValue(val *ServiceTaskSpec)
	LogDriver() ServiceTaskSpecLogDriverOutputReference
	LogDriverInput() *ServiceTaskSpecLogDriver
	NetworksAdvanced() ServiceTaskSpecNetworksAdvancedList
	NetworksAdvancedInput() interface{}
	Placement() ServiceTaskSpecPlacementOutputReference
	PlacementInput() *ServiceTaskSpecPlacement
	Resources() ServiceTaskSpecResourcesOutputReference
	ResourcesInput() *ServiceTaskSpecResources
	RestartPolicy() ServiceTaskSpecRestartPolicyOutputReference
	RestartPolicyInput() *ServiceTaskSpecRestartPolicy
	Runtime() *string
	SetRuntime(val *string)
	RuntimeInput() *string
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
	PutContainerSpec(value *ServiceTaskSpecContainerSpec)
	PutLogDriver(value *ServiceTaskSpecLogDriver)
	PutNetworksAdvanced(value interface{})
	PutPlacement(value *ServiceTaskSpecPlacement)
	PutResources(value *ServiceTaskSpecResources)
	PutRestartPolicy(value *ServiceTaskSpecRestartPolicy)
	ResetForceUpdate()
	ResetLogDriver()
	ResetNetworksAdvanced()
	ResetPlacement()
	ResetResources()
	ResetRestartPolicy()
	ResetRuntime()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecOutputReference
type jsiiProxy_ServiceTaskSpecOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ContainerSpec() ServiceTaskSpecContainerSpecOutputReference {
	var returns ServiceTaskSpecContainerSpecOutputReference
	_jsii_.Get(
		j,
		"containerSpec",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ContainerSpecInput() *ServiceTaskSpecContainerSpec {
	var returns *ServiceTaskSpecContainerSpec
	_jsii_.Get(
		j,
		"containerSpecInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ForceUpdate() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"forceUpdate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ForceUpdateInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"forceUpdateInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) InternalValue() *ServiceTaskSpec {
	var returns *ServiceTaskSpec
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) LogDriver() ServiceTaskSpecLogDriverOutputReference {
	var returns ServiceTaskSpecLogDriverOutputReference
	_jsii_.Get(
		j,
		"logDriver",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) LogDriverInput() *ServiceTaskSpecLogDriver {
	var returns *ServiceTaskSpecLogDriver
	_jsii_.Get(
		j,
		"logDriverInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) NetworksAdvanced() ServiceTaskSpecNetworksAdvancedList {
	var returns ServiceTaskSpecNetworksAdvancedList
	_jsii_.Get(
		j,
		"networksAdvanced",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) NetworksAdvancedInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"networksAdvancedInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) Placement() ServiceTaskSpecPlacementOutputReference {
	var returns ServiceTaskSpecPlacementOutputReference
	_jsii_.Get(
		j,
		"placement",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) PlacementInput() *ServiceTaskSpecPlacement {
	var returns *ServiceTaskSpecPlacement
	_jsii_.Get(
		j,
		"placementInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) Resources() ServiceTaskSpecResourcesOutputReference {
	var returns ServiceTaskSpecResourcesOutputReference
	_jsii_.Get(
		j,
		"resources",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) ResourcesInput() *ServiceTaskSpecResources {
	var returns *ServiceTaskSpecResources
	_jsii_.Get(
		j,
		"resourcesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) RestartPolicy() ServiceTaskSpecRestartPolicyOutputReference {
	var returns ServiceTaskSpecRestartPolicyOutputReference
	_jsii_.Get(
		j,
		"restartPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) RestartPolicyInput() *ServiceTaskSpecRestartPolicy {
	var returns *ServiceTaskSpecRestartPolicy
	_jsii_.Get(
		j,
		"restartPolicyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) Runtime() *string {
	var returns *string
	_jsii_.Get(
		j,
		"runtime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) RuntimeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"runtimeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceTaskSpecOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceTaskSpecOutputReference_Override(s ServiceTaskSpecOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetForceUpdate(val *float64) {
	if err := j.validateSetForceUpdateParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"forceUpdate",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetInternalValue(val *ServiceTaskSpec) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetRuntime(val *string) {
	if err := j.validateSetRuntimeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"runtime",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutContainerSpec(value *ServiceTaskSpecContainerSpec) {
	if err := s.validatePutContainerSpecParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putContainerSpec",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutLogDriver(value *ServiceTaskSpecLogDriver) {
	if err := s.validatePutLogDriverParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putLogDriver",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutNetworksAdvanced(value interface{}) {
	if err := s.validatePutNetworksAdvancedParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putNetworksAdvanced",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutPlacement(value *ServiceTaskSpecPlacement) {
	if err := s.validatePutPlacementParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putPlacement",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutResources(value *ServiceTaskSpecResources) {
	if err := s.validatePutResourcesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putResources",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) PutRestartPolicy(value *ServiceTaskSpecRestartPolicy) {
	if err := s.validatePutRestartPolicyParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putRestartPolicy",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetForceUpdate() {
	_jsii_.InvokeVoid(
		s,
		"resetForceUpdate",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetLogDriver() {
	_jsii_.InvokeVoid(
		s,
		"resetLogDriver",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetNetworksAdvanced() {
	_jsii_.InvokeVoid(
		s,
		"resetNetworksAdvanced",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetPlacement() {
	_jsii_.InvokeVoid(
		s,
		"resetPlacement",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetResources() {
	_jsii_.InvokeVoid(
		s,
		"resetResources",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetRestartPolicy() {
	_jsii_.InvokeVoid(
		s,
		"resetRestartPolicy",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ResetRuntime() {
	_jsii_.InvokeVoid(
		s,
		"resetRuntime",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

