package sql

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/sql/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/sql/internal"
)

// Defines an Sql based on a Terraform module.
//
// Source at ./.nitric/modules/sql
type Sql interface {
	cdktf.TerraformModule
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CodebuildRegion() *string
	SetCodebuildRegion(val *string)
	CodebuildRoleArn() *string
	SetCodebuildRoleArn(val *string)
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	CreateDatabaseProjectName() *string
	SetCreateDatabaseProjectName(val *string)
	DbName() *string
	SetDbName(val *string)
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
	ImageUri() *string
	SetImageUri(val *string)
	MigrateCommand() *string
	SetMigrateCommand(val *string)
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	RdsClusterEndpoint() *string
	SetRdsClusterEndpoint(val *string)
	RdsClusterPassword() *string
	SetRdsClusterPassword(val *string)
	RdsClusterUsername() *string
	SetRdsClusterUsername(val *string)
	SecurityGroupIds() *[]*string
	SetSecurityGroupIds(val *[]*string)
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	SubnetIds() *[]*string
	SetSubnetIds(val *[]*string)
	// Experimental.
	Version() *string
	VpcId() *string
	SetVpcId(val *string)
	WorkDir() *string
	SetWorkDir(val *string)
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

// The jsii proxy struct for Sql
type jsiiProxy_Sql struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Sql) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) CodebuildRegion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codebuildRegion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) CodebuildRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codebuildRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) CreateDatabaseProjectName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"createDatabaseProjectName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) DbName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dbName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ImageUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) MigrateCommand() *string {
	var returns *string
	_jsii_.Get(
		j,
		"migrateCommand",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) RdsClusterEndpoint() *string {
	var returns *string
	_jsii_.Get(
		j,
		"rdsClusterEndpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) RdsClusterPassword() *string {
	var returns *string
	_jsii_.Get(
		j,
		"rdsClusterPassword",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) RdsClusterUsername() *string {
	var returns *string
	_jsii_.Get(
		j,
		"rdsClusterUsername",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) SecurityGroupIds() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityGroupIds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) SubnetIds() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"subnetIds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) VpcId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"vpcId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) WorkDir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"workDir",
		&returns,
	)
	return returns
}
func NewSql(scope constructs.Construct, id *string, config *SqlConfig) Sql {
	_init_.Initialize()

	if err := validateNewSqlParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Sql{}

	_jsii_.Create(
		"sql.Sql",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewSql_Override(s Sql, scope constructs.Construct, id *string, config *SqlConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"sql.Sql",
		[]interface{}{scope, id, config},
		s,
	)
}

func (j *jsiiProxy_Sql) SetCodebuildRegion(val *string) {
	if err := j.validateSetCodebuildRegionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"codebuildRegion",
		val,
	)
}

func (j *jsiiProxy_Sql) SetCodebuildRoleArn(val *string) {
	if err := j.validateSetCodebuildRoleArnParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"codebuildRoleArn",
		val,
	)
}

func (j *jsiiProxy_Sql) SetCreateDatabaseProjectName(val *string) {
	if err := j.validateSetCreateDatabaseProjectNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"createDatabaseProjectName",
		val,
	)
}

func (j *jsiiProxy_Sql) SetDbName(val *string) {
	if err := j.validateSetDbNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dbName",
		val,
	)
}

func (j *jsiiProxy_Sql) SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Sql) SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Sql) SetImageUri(val *string) {
	if err := j.validateSetImageUriParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"imageUri",
		val,
	)
}

func (j *jsiiProxy_Sql) SetMigrateCommand(val *string) {
	if err := j.validateSetMigrateCommandParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"migrateCommand",
		val,
	)
}

func (j *jsiiProxy_Sql) SetRdsClusterEndpoint(val *string) {
	if err := j.validateSetRdsClusterEndpointParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"rdsClusterEndpoint",
		val,
	)
}

func (j *jsiiProxy_Sql) SetRdsClusterPassword(val *string) {
	if err := j.validateSetRdsClusterPasswordParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"rdsClusterPassword",
		val,
	)
}

func (j *jsiiProxy_Sql) SetRdsClusterUsername(val *string) {
	if err := j.validateSetRdsClusterUsernameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"rdsClusterUsername",
		val,
	)
}

func (j *jsiiProxy_Sql) SetSecurityGroupIds(val *[]*string) {
	if err := j.validateSetSecurityGroupIdsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"securityGroupIds",
		val,
	)
}

func (j *jsiiProxy_Sql) SetSubnetIds(val *[]*string) {
	if err := j.validateSetSubnetIdsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"subnetIds",
		val,
	)
}

func (j *jsiiProxy_Sql) SetVpcId(val *string) {
	if err := j.validateSetVpcIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"vpcId",
		val,
	)
}

func (j *jsiiProxy_Sql) SetWorkDir(val *string) {
	if err := j.validateSetWorkDirParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"workDir",
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
func Sql_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSql_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"sql.Sql",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Sql_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSql_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"sql.Sql",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) AddOverride(path *string, value interface{}) {
	if err := s.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (s *jsiiProxy_Sql) AddProvider(provider interface{}) {
	if err := s.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addProvider",
		[]interface{}{provider},
	)
}

func (s *jsiiProxy_Sql) GetString(output *string) *string {
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

func (s *jsiiProxy_Sql) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
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

func (s *jsiiProxy_Sql) OverrideLogicalId(newLogicalId *string) {
	if err := s.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (s *jsiiProxy_Sql) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		s,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (s *jsiiProxy_Sql) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}
