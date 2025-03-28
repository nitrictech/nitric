package cdn

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/cdn/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/cdn/internal"
)

// Defines an Cdn based on a Terraform module.
//
// Source at ./.nitric/modules/cdn
type Cdn interface {
	cdktf.TerraformModule
	ApiGateways() interface{}
	SetApiGateways(val interface{})
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CdnDomain() interface{}
	SetCdnDomain(val interface{})
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	// The tree node.
	Node() constructs.Node
	ProjectId() *string
	SetProjectId(val *string)
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	Region() *string
	SetRegion(val *string)
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	StackId() *string
	SetStackId(val *string)
	// Experimental.
	Version() *string
	WebsiteBuckets() interface{}
	SetWebsiteBuckets(val interface{})
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

// The jsii proxy struct for Cdn
type jsiiProxy_Cdn struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Cdn) ApiGateways() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"apiGateways",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdnDomain() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cdnDomain",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) ProjectId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) WebsiteBuckets() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"websiteBuckets",
		&returns,
	)
	return returns
}


func NewCdn(scope constructs.Construct, id *string, config *CdnConfig) Cdn {
	_init_.Initialize()

	if err := validateNewCdnParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Cdn{}

	_jsii_.Create(
		"cdn.Cdn",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewCdn_Override(c Cdn, scope constructs.Construct, id *string, config *CdnConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"cdn.Cdn",
		[]interface{}{scope, id, config},
		c,
	)
}

func (j *jsiiProxy_Cdn)SetApiGateways(val interface{}) {
	if err := j.validateSetApiGatewaysParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"apiGateways",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetCdnDomain(val interface{}) {
	if err := j.validateSetCdnDomainParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cdnDomain",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetProjectId(val *string) {
	if err := j.validateSetProjectIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"projectId",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetRegion(val *string) {
	if err := j.validateSetRegionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"region",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetStackId(val *string) {
	if err := j.validateSetStackIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackId",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetWebsiteBuckets(val interface{}) {
	if err := j.validateSetWebsiteBucketsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"websiteBuckets",
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
func Cdn_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateCdn_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdn.Cdn",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Cdn_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateCdn_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdn.Cdn",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) AddOverride(path *string, value interface{}) {
	if err := c.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (c *jsiiProxy_Cdn) AddProvider(provider interface{}) {
	if err := c.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addProvider",
		[]interface{}{provider},
	)
}

func (c *jsiiProxy_Cdn) GetString(output *string) *string {
	if err := c.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		c,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := c.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		c,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) OverrideLogicalId(newLogicalId *string) {
	if err := c.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (c *jsiiProxy_Cdn) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		c,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Cdn) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Cdn) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

