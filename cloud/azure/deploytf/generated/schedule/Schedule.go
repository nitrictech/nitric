package schedule

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/schedule/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/schedule/internal"
)

// Defines an Schedule based on a Terraform module.
//
// Source at ./.nitric/modules/schedule
type Schedule interface {
	cdktf.TerraformModule
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	ContainerAppEnvironmentId() *string
	SetContainerAppEnvironmentId(val *string)
	CronExpression() *string
	SetCronExpression(val *string)
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
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	TargetAppId() *string
	SetTargetAppId(val *string)
	TargetEventToken() *string
	SetTargetEventToken(val *string)
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

// The jsii proxy struct for Schedule
type jsiiProxy_Schedule struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Schedule) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) ContainerAppEnvironmentId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerAppEnvironmentId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) CronExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cronExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) TargetAppId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetAppId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) TargetEventToken() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetEventToken",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Schedule) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewSchedule(scope constructs.Construct, id *string, config *ScheduleConfig) Schedule {
	_init_.Initialize()

	if err := validateNewScheduleParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Schedule{}

	_jsii_.Create(
		"schedule.Schedule",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewSchedule_Override(s Schedule, scope constructs.Construct, id *string, config *ScheduleConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"schedule.Schedule",
		[]interface{}{scope, id, config},
		s,
	)
}

func (j *jsiiProxy_Schedule)SetContainerAppEnvironmentId(val *string) {
	if err := j.validateSetContainerAppEnvironmentIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"containerAppEnvironmentId",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetCronExpression(val *string) {
	if err := j.validateSetCronExpressionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cronExpression",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetTargetAppId(val *string) {
	if err := j.validateSetTargetAppIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"targetAppId",
		val,
	)
}

func (j *jsiiProxy_Schedule)SetTargetEventToken(val *string) {
	if err := j.validateSetTargetEventTokenParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"targetEventToken",
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
func Schedule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSchedule_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"schedule.Schedule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Schedule_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSchedule_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"schedule.Schedule",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) AddOverride(path *string, value interface{}) {
	if err := s.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (s *jsiiProxy_Schedule) AddProvider(provider interface{}) {
	if err := s.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addProvider",
		[]interface{}{provider},
	)
}

func (s *jsiiProxy_Schedule) GetString(output *string) *string {
	if err := s.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		s,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := s.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) OverrideLogicalId(newLogicalId *string) {
	if err := s.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (s *jsiiProxy_Schedule) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		s,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (s *jsiiProxy_Schedule) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Schedule) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

