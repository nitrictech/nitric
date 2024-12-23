package api

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/api/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/api/internal"
)

// Defines an Api based on a Terraform module.
//
// Source at ./.nitric/modules/api
type Api interface {
	cdktf.TerraformModule
	AppIdentity() *string
	SetAppIdentity(val *string)
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	Description() *string
	SetDescription(val *string)
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	Location() *string
	SetLocation(val *string)
	Name() *string
	SetName(val *string)
	// The tree node.
	Node() constructs.Node
	OpenapiSpec() *string
	SetOpenapiSpec(val *string)
	OperationPolicyTemplates() *map[string]*string
	SetOperationPolicyTemplates(val *map[string]*string)
	// Experimental.
	Providers() *[]interface{}
	PublisherEmail() *string
	SetPublisherEmail(val *string)
	PublisherName() *string
	SetPublisherName(val *string)
	// Experimental.
	RawOverrides() interface{}
	ResourceGroupName() *string
	SetResourceGroupName(val *string)
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	// Experimental.
	Version() *string
	// Experimental.
	AddOverride(path *string, value interface{})
	// Experimental.
	AddProvider(provider interface{})
	// Experimental.
	GetString(output *string) *string
	// Experimental.
	InterpolationForOutput(moduleOutput *string) cdktf.IResolvable
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
	// Resets a previously passed logical Id to use the auto-generated logical id again.
	// Experimental.
	ResetOverrideLogicalId()
	SynthesizeAttributes() *map[string]interface{}
	SynthesizeHclAttributes() *map[string]interface{}
	// Experimental.
	ToHclTerraform() interface{}
	// Experimental.
	ToMetadata() interface{}
	// Returns a string representation of this construct.
	ToString() *string
	// Experimental.
	ToTerraform() interface{}
}

// The jsii proxy struct for Api
type jsiiProxy_Api struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Api) AppIdentity() *string {
	var returns *string
	_jsii_.Get(
		j,
		"appIdentity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Location() *string {
	var returns *string
	_jsii_.Get(
		j,
		"location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) OpenapiSpec() *string {
	var returns *string
	_jsii_.Get(
		j,
		"openapiSpec",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) OperationPolicyTemplates() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"operationPolicyTemplates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) PublisherEmail() *string {
	var returns *string
	_jsii_.Get(
		j,
		"publisherEmail",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) PublisherName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"publisherName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) ResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Api) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewApi(scope constructs.Construct, id *string, config *ApiConfig) Api {
	_init_.Initialize()

	if err := validateNewApiParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Api{}

	_jsii_.Create(
		"api.Api",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewApi_Override(a Api, scope constructs.Construct, id *string, config *ApiConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"api.Api",
		[]interface{}{scope, id, config},
		a,
	)
}

func (j *jsiiProxy_Api)SetAppIdentity(val *string) {
	if err := j.validateSetAppIdentityParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"appIdentity",
		val,
	)
}

func (j *jsiiProxy_Api)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Api)SetDescription(val *string) {
	if err := j.validateSetDescriptionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_Api)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Api)SetLocation(val *string) {
	if err := j.validateSetLocationParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"location",
		val,
	)
}

func (j *jsiiProxy_Api)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_Api)SetOpenapiSpec(val *string) {
	if err := j.validateSetOpenapiSpecParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"openapiSpec",
		val,
	)
}

func (j *jsiiProxy_Api)SetOperationPolicyTemplates(val *map[string]*string) {
	if err := j.validateSetOperationPolicyTemplatesParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"operationPolicyTemplates",
		val,
	)
}

func (j *jsiiProxy_Api)SetPublisherEmail(val *string) {
	if err := j.validateSetPublisherEmailParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"publisherEmail",
		val,
	)
}

func (j *jsiiProxy_Api)SetPublisherName(val *string) {
	if err := j.validateSetPublisherNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"publisherName",
		val,
	)
}

func (j *jsiiProxy_Api)SetResourceGroupName(val *string) {
	if err := j.validateSetResourceGroupNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"resourceGroupName",
		val,
	)
}

// Checks if `x` is a construct.
//
// Use this method instead of `instanceof` to properly detect `Construct`
// instances, even when the construct library is symlinked.
//
// Explanation: in JavaScript, multiple copies of the `constructs` library on
// disk are seen as independent, completely different libraries. As a
// consequence, the class `Construct` in each copy of the `constructs` library
// is seen as a different class, and an instance of one class will not test as
// `instanceof` the other class. `npm install` will not create installations
// like this, but users may manually symlink construct libraries together or
// use a monorepo tool: in those cases, multiple copies of the `constructs`
// library can be accidentally installed, and `instanceof` will behave
// unpredictably. It is safest to avoid using `instanceof`, and using
// this type-testing method instead.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
func Api_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateApi_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"api.Api",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Api_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateApi_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"api.Api",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) AddOverride(path *string, value interface{}) {
	if err := a.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		a,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (a *jsiiProxy_Api) AddProvider(provider interface{}) {
	if err := a.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		a,
		"addProvider",
		[]interface{}{provider},
	)
}

func (a *jsiiProxy_Api) GetString(output *string) *string {
	if err := a.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		a,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := a.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		a,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) OverrideLogicalId(newLogicalId *string) {
	if err := a.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		a,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (a *jsiiProxy_Api) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		a,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (a *jsiiProxy_Api) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		a,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		a,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		a,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		a,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (a *jsiiProxy_Api) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		a,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

