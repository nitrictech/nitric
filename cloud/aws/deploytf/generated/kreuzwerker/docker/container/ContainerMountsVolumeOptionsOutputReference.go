package container

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/container/internal"
)

type ContainerMountsVolumeOptionsOutputReference interface {
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
	InternalValue() *ContainerMountsVolumeOptions
	SetInternalValue(val *ContainerMountsVolumeOptions)
	Labels() ContainerMountsVolumeOptionsLabelsList
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

// The jsii proxy struct for ContainerMountsVolumeOptionsOutputReference
type jsiiProxy_ContainerMountsVolumeOptionsOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) DriverName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driverName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) DriverNameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driverNameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) DriverOptions() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"driverOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) DriverOptionsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"driverOptionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) InternalValue() *ContainerMountsVolumeOptions {
	var returns *ContainerMountsVolumeOptions
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) Labels() ContainerMountsVolumeOptionsLabelsList {
	var returns ContainerMountsVolumeOptionsLabelsList
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) LabelsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) NoCopy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCopy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) NoCopyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCopyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}


func NewContainerMountsVolumeOptionsOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ContainerMountsVolumeOptionsOutputReference {
	_init_.Initialize()

	if err := validateNewContainerMountsVolumeOptionsOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ContainerMountsVolumeOptionsOutputReference{}

	_jsii_.Create(
		"docker.container.ContainerMountsVolumeOptionsOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewContainerMountsVolumeOptionsOutputReference_Override(c ContainerMountsVolumeOptionsOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.container.ContainerMountsVolumeOptionsOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		c,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetDriverName(val *string) {
	if err := j.validateSetDriverNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"driverName",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetDriverOptions(val *map[string]*string) {
	if err := j.validateSetDriverOptionsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"driverOptions",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetInternalValue(val *ContainerMountsVolumeOptions) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetNoCopy(val interface{}) {
	if err := j.validateSetNoCopyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"noCopy",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ContainerMountsVolumeOptionsOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
	if err := c.validateGetAnyMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"getAnyMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := c.validateGetBooleanAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		c,
		"getBooleanAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
	if err := c.validateGetBooleanMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*bool

	_jsii_.Invoke(
		c,
		"getBooleanMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
	if err := c.validateGetListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"getListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
	if err := c.validateGetNumberAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		c,
		"getNumberAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
	if err := c.validateGetNumberListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*float64

	_jsii_.Invoke(
		c,
		"getNumberListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
	if err := c.validateGetNumberMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*float64

	_jsii_.Invoke(
		c,
		"getNumberMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetStringAttribute(terraformAttribute *string) *string {
	if err := c.validateGetStringAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		c,
		"getStringAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
	if err := c.validateGetStringMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*string

	_jsii_.Invoke(
		c,
		"getStringMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		c,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
	if err := c.validateInterpolationForAttributeParameters(property); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		c,
		"interpolationForAttribute",
		[]interface{}{property},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) PutLabels(value interface{}) {
	if err := c.validatePutLabelsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putLabels",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ResetDriverName() {
	_jsii_.InvokeVoid(
		c,
		"resetDriverName",
		nil, // no parameters
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ResetDriverOptions() {
	_jsii_.InvokeVoid(
		c,
		"resetDriverOptions",
		nil, // no parameters
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ResetLabels() {
	_jsii_.InvokeVoid(
		c,
		"resetLabels",
		nil, // no parameters
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ResetNoCopy() {
	_jsii_.InvokeVoid(
		c,
		"resetNoCopy",
		nil, // no parameters
	)
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
	if err := c.validateResolveParameters(_context); err != nil {
		panic(err)
	}
	var returns interface{}

	_jsii_.Invoke(
		c,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_ContainerMountsVolumeOptionsOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

