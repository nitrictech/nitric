package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecPlacementOutputReference interface {
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
	Constraints() *[]*string
	SetConstraints(val *[]*string)
	ConstraintsInput() *[]*string
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	// Experimental.
	Fqn() *string
	InternalValue() *ServiceTaskSpecPlacement
	SetInternalValue(val *ServiceTaskSpecPlacement)
	MaxReplicas() *float64
	SetMaxReplicas(val *float64)
	MaxReplicasInput() *float64
	Platforms() ServiceTaskSpecPlacementPlatformsList
	PlatformsInput() interface{}
	Prefs() *[]*string
	SetPrefs(val *[]*string)
	PrefsInput() *[]*string
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
	PutPlatforms(value interface{})
	ResetConstraints()
	ResetMaxReplicas()
	ResetPlatforms()
	ResetPrefs()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecPlacementOutputReference
type jsiiProxy_ServiceTaskSpecPlacementOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) Constraints() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"constraints",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ConstraintsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"constraintsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) InternalValue() *ServiceTaskSpecPlacement {
	var returns *ServiceTaskSpecPlacement
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) MaxReplicas() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxReplicas",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) MaxReplicasInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxReplicasInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) Platforms() ServiceTaskSpecPlacementPlatformsList {
	var returns ServiceTaskSpecPlacementPlatformsList
	_jsii_.Get(
		j,
		"platforms",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) PlatformsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"platformsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) Prefs() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"prefs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) PrefsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"prefsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecPlacementOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceTaskSpecPlacementOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecPlacementOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecPlacementOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecPlacementOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceTaskSpecPlacementOutputReference_Override(s ServiceTaskSpecPlacementOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecPlacementOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetConstraints(val *[]*string) {
	if err := j.validateSetConstraintsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"constraints",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetInternalValue(val *ServiceTaskSpecPlacement) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetMaxReplicas(val *float64) {
	if err := j.validateSetMaxReplicasParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"maxReplicas",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetPrefs(val *[]*string) {
	if err := j.validateSetPrefsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"prefs",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecPlacementOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) PutPlatforms(value interface{}) {
	if err := s.validatePutPlatformsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putPlatforms",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ResetConstraints() {
	_jsii_.InvokeVoid(
		s,
		"resetConstraints",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ResetMaxReplicas() {
	_jsii_.InvokeVoid(
		s,
		"resetMaxReplicas",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ResetPlatforms() {
	_jsii_.InvokeVoid(
		s,
		"resetPlatforms",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ResetPrefs() {
	_jsii_.InvokeVoid(
		s,
		"resetPrefs",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecPlacementOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

