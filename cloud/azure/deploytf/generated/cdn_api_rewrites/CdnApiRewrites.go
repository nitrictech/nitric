package cdn_api_rewrites

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn_api_rewrites/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn_api_rewrites/internal"
)

// Defines an CdnApiRewrites based on a Terraform module.
//
// Source at ./.nitric/modules/cdn_api_rewrites
type CdnApiRewrites interface {
	cdktf.TerraformModule
	ApiHostName() *string
	SetApiHostName(val *string)
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CdnFrontdoorProfileId() *string
	SetCdnFrontdoorProfileId(val *string)
	CdnFrontdoorRuleSetId() *string
	SetCdnFrontdoorRuleSetId(val *string)
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
	Name() *string
	SetName(val *string)
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	RuleOrder() *float64
	SetRuleOrder(val *float64)
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

// The jsii proxy struct for CdnApiRewrites
type jsiiProxy_CdnApiRewrites struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_CdnApiRewrites) ApiHostName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiHostName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) CdnFrontdoorProfileId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnFrontdoorProfileId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) CdnFrontdoorRuleSetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnFrontdoorRuleSetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) RuleOrder() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"ruleOrder",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CdnApiRewrites) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewCdnApiRewrites(scope constructs.Construct, id *string, config *CdnApiRewritesConfig) CdnApiRewrites {
	_init_.Initialize()

	if err := validateNewCdnApiRewritesParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_CdnApiRewrites{}

	_jsii_.Create(
		"cdn_api_rewrites.CdnApiRewrites",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewCdnApiRewrites_Override(c CdnApiRewrites, scope constructs.Construct, id *string, config *CdnApiRewritesConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"cdn_api_rewrites.CdnApiRewrites",
		[]interface{}{scope, id, config},
		c,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetApiHostName(val *string) {
	if err := j.validateSetApiHostNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"apiHostName",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetCdnFrontdoorProfileId(val *string) {
	if err := j.validateSetCdnFrontdoorProfileIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cdnFrontdoorProfileId",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetCdnFrontdoorRuleSetId(val *string) {
	if err := j.validateSetCdnFrontdoorRuleSetIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cdnFrontdoorRuleSetId",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CdnApiRewrites)SetRuleOrder(val *float64) {
	_jsii_.Set(
		j,
		"ruleOrder",
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
func CdnApiRewrites_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateCdnApiRewrites_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdn_api_rewrites.CdnApiRewrites",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func CdnApiRewrites_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateCdnApiRewrites_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdn_api_rewrites.CdnApiRewrites",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) AddOverride(path *string, value interface{}) {
	if err := c.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (c *jsiiProxy_CdnApiRewrites) AddProvider(provider interface{}) {
	if err := c.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addProvider",
		[]interface{}{provider},
	)
}

func (c *jsiiProxy_CdnApiRewrites) GetString(output *string) *string {
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

func (c *jsiiProxy_CdnApiRewrites) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
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

func (c *jsiiProxy_CdnApiRewrites) OverrideLogicalId(newLogicalId *string) {
	if err := c.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (c *jsiiProxy_CdnApiRewrites) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		c,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CdnApiRewrites) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_CdnApiRewrites) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

