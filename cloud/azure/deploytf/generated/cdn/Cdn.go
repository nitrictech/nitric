package cdn

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn/internal"
)

// Defines an Cdn based on a Terraform module.
//
// Source at ./.nitric/modules/cdn
type Cdn interface {
	cdktf.TerraformModule
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CdnFrontdoorApiRuleSetIdOutput() *string
	CdnFrontdoorDefaultRuleSetIdOutput() *string
	CdnFrontdoorProfileIdOutput() *string
	CdnUrlOutput() *string
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	CustomDomainHostName() *string
	SetCustomDomainHostName(val *string)
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	DomainName() *string
	SetDomainName(val *string)
	EnableApiRewrites() *bool
	SetEnableApiRewrites(val *bool)
	EnableCustomDomain() *bool
	SetEnableCustomDomain(val *bool)
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	IsApexDomain() *bool
	SetIsApexDomain(val *bool)
	// The tree node.
	Node() constructs.Node
	PrimaryWebHost() *string
	SetPrimaryWebHost(val *string)
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	ResourceGroupName() *string
	SetResourceGroupName(val *string)
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	SkipCacheInvalidation() *bool
	SetSkipCacheInvalidation(val *bool)
	// Experimental.
	Source() *string
	StackName() *string
	SetStackName(val *string)
	UploadedFiles() *map[string]*string
	SetUploadedFiles(val *map[string]*string)
	// Experimental.
	Version() *string
	ZoneName() *string
	SetZoneName(val *string)
	ZoneResourceGroupName() *string
	SetZoneResourceGroupName(val *string)
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

func (j *jsiiProxy_Cdn) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdnFrontdoorApiRuleSetIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnFrontdoorApiRuleSetIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdnFrontdoorDefaultRuleSetIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnFrontdoorDefaultRuleSetIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdnFrontdoorProfileIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnFrontdoorProfileIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) CdnUrlOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cdnUrlOutput",
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

func (j *jsiiProxy_Cdn) CustomDomainHostName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"customDomainHostName",
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

func (j *jsiiProxy_Cdn) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) EnableApiRewrites() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableApiRewrites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) EnableCustomDomain() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableCustomDomain",
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

func (j *jsiiProxy_Cdn) IsApexDomain() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isApexDomain",
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

func (j *jsiiProxy_Cdn) PrimaryWebHost() *string {
	var returns *string
	_jsii_.Get(
		j,
		"primaryWebHost",
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

func (j *jsiiProxy_Cdn) ResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupName",
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

func (j *jsiiProxy_Cdn) SkipCacheInvalidation() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipCacheInvalidation",
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

func (j *jsiiProxy_Cdn) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) UploadedFiles() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"uploadedFiles",
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

func (j *jsiiProxy_Cdn) ZoneName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"zoneName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cdn) ZoneResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"zoneResourceGroupName",
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

func (j *jsiiProxy_Cdn)SetCustomDomainHostName(val *string) {
	if err := j.validateSetCustomDomainHostNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"customDomainHostName",
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

func (j *jsiiProxy_Cdn)SetDomainName(val *string) {
	if err := j.validateSetDomainNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"domainName",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetEnableApiRewrites(val *bool) {
	_jsii_.Set(
		j,
		"enableApiRewrites",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetEnableCustomDomain(val *bool) {
	_jsii_.Set(
		j,
		"enableCustomDomain",
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

func (j *jsiiProxy_Cdn)SetIsApexDomain(val *bool) {
	_jsii_.Set(
		j,
		"isApexDomain",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetPrimaryWebHost(val *string) {
	if err := j.validateSetPrimaryWebHostParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"primaryWebHost",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetResourceGroupName(val *string) {
	if err := j.validateSetResourceGroupNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"resourceGroupName",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetSkipCacheInvalidation(val *bool) {
	_jsii_.Set(
		j,
		"skipCacheInvalidation",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetStackName(val *string) {
	if err := j.validateSetStackNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackName",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetUploadedFiles(val *map[string]*string) {
	_jsii_.Set(
		j,
		"uploadedFiles",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetZoneName(val *string) {
	if err := j.validateSetZoneNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"zoneName",
		val,
	)
}

func (j *jsiiProxy_Cdn)SetZoneResourceGroupName(val *string) {
	if err := j.validateSetZoneResourceGroupNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"zoneResourceGroupName",
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

