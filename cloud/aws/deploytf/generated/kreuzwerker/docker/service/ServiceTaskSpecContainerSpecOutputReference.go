package service

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/service/internal"
)

type ServiceTaskSpecContainerSpecOutputReference interface {
	cdktf.ComplexObject
	Args() *[]*string
	SetArgs(val *[]*string)
	ArgsInput() *[]*string
	Command() *[]*string
	SetCommand(val *[]*string)
	CommandInput() *[]*string
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
	Configs() ServiceTaskSpecContainerSpecConfigsList
	ConfigsInput() interface{}
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	Dir() *string
	SetDir(val *string)
	DirInput() *string
	DnsConfig() ServiceTaskSpecContainerSpecDnsConfigOutputReference
	DnsConfigInput() *ServiceTaskSpecContainerSpecDnsConfig
	Env() *map[string]*string
	SetEnv(val *map[string]*string)
	EnvInput() *map[string]*string
	// Experimental.
	Fqn() *string
	Groups() *[]*string
	SetGroups(val *[]*string)
	GroupsInput() *[]*string
	Healthcheck() ServiceTaskSpecContainerSpecHealthcheckOutputReference
	HealthcheckInput() *ServiceTaskSpecContainerSpecHealthcheck
	Hostname() *string
	SetHostname(val *string)
	HostnameInput() *string
	Hosts() ServiceTaskSpecContainerSpecHostsList
	HostsInput() interface{}
	Image() *string
	SetImage(val *string)
	ImageInput() *string
	InternalValue() *ServiceTaskSpecContainerSpec
	SetInternalValue(val *ServiceTaskSpecContainerSpec)
	Isolation() *string
	SetIsolation(val *string)
	IsolationInput() *string
	Labels() ServiceTaskSpecContainerSpecLabelsList
	LabelsInput() interface{}
	Mounts() ServiceTaskSpecContainerSpecMountsList
	MountsInput() interface{}
	Privileges() ServiceTaskSpecContainerSpecPrivilegesOutputReference
	PrivilegesInput() *ServiceTaskSpecContainerSpecPrivileges
	ReadOnly() interface{}
	SetReadOnly(val interface{})
	ReadOnlyInput() interface{}
	Secrets() ServiceTaskSpecContainerSpecSecretsList
	SecretsInput() interface{}
	StopGracePeriod() *string
	SetStopGracePeriod(val *string)
	StopGracePeriodInput() *string
	StopSignal() *string
	SetStopSignal(val *string)
	StopSignalInput() *string
	Sysctl() *map[string]*string
	SetSysctl(val *map[string]*string)
	SysctlInput() *map[string]*string
	// Experimental.
	TerraformAttribute() *string
	// Experimental.
	SetTerraformAttribute(val *string)
	// Experimental.
	TerraformResource() cdktf.IInterpolatingParent
	// Experimental.
	SetTerraformResource(val cdktf.IInterpolatingParent)
	User() *string
	SetUser(val *string)
	UserInput() *string
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
	PutConfigs(value interface{})
	PutDnsConfig(value *ServiceTaskSpecContainerSpecDnsConfig)
	PutHealthcheck(value *ServiceTaskSpecContainerSpecHealthcheck)
	PutHosts(value interface{})
	PutLabels(value interface{})
	PutMounts(value interface{})
	PutPrivileges(value *ServiceTaskSpecContainerSpecPrivileges)
	PutSecrets(value interface{})
	ResetArgs()
	ResetCommand()
	ResetConfigs()
	ResetDir()
	ResetDnsConfig()
	ResetEnv()
	ResetGroups()
	ResetHealthcheck()
	ResetHostname()
	ResetHosts()
	ResetIsolation()
	ResetLabels()
	ResetMounts()
	ResetPrivileges()
	ResetReadOnly()
	ResetSecrets()
	ResetStopGracePeriod()
	ResetStopSignal()
	ResetSysctl()
	ResetUser()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ServiceTaskSpecContainerSpecOutputReference
type jsiiProxy_ServiceTaskSpecContainerSpecOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Args() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"args",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ArgsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"argsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Command() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"command",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) CommandInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"commandInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Configs() ServiceTaskSpecContainerSpecConfigsList {
	var returns ServiceTaskSpecContainerSpecConfigsList
	_jsii_.Get(
		j,
		"configs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ConfigsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"configsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Dir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) DirInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dirInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) DnsConfig() ServiceTaskSpecContainerSpecDnsConfigOutputReference {
	var returns ServiceTaskSpecContainerSpecDnsConfigOutputReference
	_jsii_.Get(
		j,
		"dnsConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) DnsConfigInput() *ServiceTaskSpecContainerSpecDnsConfig {
	var returns *ServiceTaskSpecContainerSpecDnsConfig
	_jsii_.Get(
		j,
		"dnsConfigInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Env() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) EnvInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"envInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Groups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GroupsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groupsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Healthcheck() ServiceTaskSpecContainerSpecHealthcheckOutputReference {
	var returns ServiceTaskSpecContainerSpecHealthcheckOutputReference
	_jsii_.Get(
		j,
		"healthcheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) HealthcheckInput() *ServiceTaskSpecContainerSpecHealthcheck {
	var returns *ServiceTaskSpecContainerSpecHealthcheck
	_jsii_.Get(
		j,
		"healthcheckInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Hostname() *string {
	var returns *string
	_jsii_.Get(
		j,
		"hostname",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) HostnameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"hostnameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Hosts() ServiceTaskSpecContainerSpecHostsList {
	var returns ServiceTaskSpecContainerSpecHostsList
	_jsii_.Get(
		j,
		"hosts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) HostsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"hostsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Image() *string {
	var returns *string
	_jsii_.Get(
		j,
		"image",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ImageInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) InternalValue() *ServiceTaskSpecContainerSpec {
	var returns *ServiceTaskSpecContainerSpec
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Isolation() *string {
	var returns *string
	_jsii_.Get(
		j,
		"isolation",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) IsolationInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"isolationInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Labels() ServiceTaskSpecContainerSpecLabelsList {
	var returns ServiceTaskSpecContainerSpecLabelsList
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) LabelsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Mounts() ServiceTaskSpecContainerSpecMountsList {
	var returns ServiceTaskSpecContainerSpecMountsList
	_jsii_.Get(
		j,
		"mounts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) MountsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mountsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Privileges() ServiceTaskSpecContainerSpecPrivilegesOutputReference {
	var returns ServiceTaskSpecContainerSpecPrivilegesOutputReference
	_jsii_.Get(
		j,
		"privileges",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PrivilegesInput() *ServiceTaskSpecContainerSpecPrivileges {
	var returns *ServiceTaskSpecContainerSpecPrivileges
	_jsii_.Get(
		j,
		"privilegesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ReadOnly() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnly",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ReadOnlyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnlyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Secrets() ServiceTaskSpecContainerSpecSecretsList {
	var returns ServiceTaskSpecContainerSpecSecretsList
	_jsii_.Get(
		j,
		"secrets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) SecretsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"secretsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) StopGracePeriod() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopGracePeriod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) StopGracePeriodInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopGracePeriodInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) StopSignal() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopSignal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) StopSignalInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopSignalInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Sysctl() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"sysctl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) SysctlInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"sysctlInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) User() *string {
	var returns *string
	_jsii_.Get(
		j,
		"user",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) UserInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userInput",
		&returns,
	)
	return returns
}


func NewServiceTaskSpecContainerSpecOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ServiceTaskSpecContainerSpecOutputReference {
	_init_.Initialize()

	if err := validateNewServiceTaskSpecContainerSpecOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ServiceTaskSpecContainerSpecOutputReference{}

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewServiceTaskSpecContainerSpecOutputReference_Override(s ServiceTaskSpecContainerSpecOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.service.ServiceTaskSpecContainerSpecOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		s,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetArgs(val *[]*string) {
	if err := j.validateSetArgsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"args",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetCommand(val *[]*string) {
	if err := j.validateSetCommandParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"command",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetDir(val *string) {
	if err := j.validateSetDirParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dir",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetEnv(val *map[string]*string) {
	if err := j.validateSetEnvParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"env",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetGroups(val *[]*string) {
	if err := j.validateSetGroupsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"groups",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetHostname(val *string) {
	if err := j.validateSetHostnameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"hostname",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetImage(val *string) {
	if err := j.validateSetImageParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"image",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetInternalValue(val *ServiceTaskSpecContainerSpec) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetIsolation(val *string) {
	if err := j.validateSetIsolationParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"isolation",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetReadOnly(val interface{}) {
	if err := j.validateSetReadOnlyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"readOnly",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetStopGracePeriod(val *string) {
	if err := j.validateSetStopGracePeriodParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stopGracePeriod",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetStopSignal(val *string) {
	if err := j.validateSetStopSignalParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stopSignal",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetSysctl(val *map[string]*string) {
	if err := j.validateSetSysctlParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"sysctl",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference)SetUser(val *string) {
	if err := j.validateSetUserParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"user",
		val,
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetStringAttribute(terraformAttribute *string) *string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutConfigs(value interface{}) {
	if err := s.validatePutConfigsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putConfigs",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutDnsConfig(value *ServiceTaskSpecContainerSpecDnsConfig) {
	if err := s.validatePutDnsConfigParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putDnsConfig",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutHealthcheck(value *ServiceTaskSpecContainerSpecHealthcheck) {
	if err := s.validatePutHealthcheckParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putHealthcheck",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutHosts(value interface{}) {
	if err := s.validatePutHostsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putHosts",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutLabels(value interface{}) {
	if err := s.validatePutLabelsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putLabels",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutMounts(value interface{}) {
	if err := s.validatePutMountsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putMounts",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutPrivileges(value *ServiceTaskSpecContainerSpecPrivileges) {
	if err := s.validatePutPrivilegesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putPrivileges",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) PutSecrets(value interface{}) {
	if err := s.validatePutSecretsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"putSecrets",
		[]interface{}{value},
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetArgs() {
	_jsii_.InvokeVoid(
		s,
		"resetArgs",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetCommand() {
	_jsii_.InvokeVoid(
		s,
		"resetCommand",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetConfigs() {
	_jsii_.InvokeVoid(
		s,
		"resetConfigs",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetDir() {
	_jsii_.InvokeVoid(
		s,
		"resetDir",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetDnsConfig() {
	_jsii_.InvokeVoid(
		s,
		"resetDnsConfig",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetEnv() {
	_jsii_.InvokeVoid(
		s,
		"resetEnv",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetGroups() {
	_jsii_.InvokeVoid(
		s,
		"resetGroups",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetHealthcheck() {
	_jsii_.InvokeVoid(
		s,
		"resetHealthcheck",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetHostname() {
	_jsii_.InvokeVoid(
		s,
		"resetHostname",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetHosts() {
	_jsii_.InvokeVoid(
		s,
		"resetHosts",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetIsolation() {
	_jsii_.InvokeVoid(
		s,
		"resetIsolation",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetLabels() {
	_jsii_.InvokeVoid(
		s,
		"resetLabels",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetMounts() {
	_jsii_.InvokeVoid(
		s,
		"resetMounts",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetPrivileges() {
	_jsii_.InvokeVoid(
		s,
		"resetPrivileges",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetReadOnly() {
	_jsii_.InvokeVoid(
		s,
		"resetReadOnly",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetSecrets() {
	_jsii_.InvokeVoid(
		s,
		"resetSecrets",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetStopGracePeriod() {
	_jsii_.InvokeVoid(
		s,
		"resetStopGracePeriod",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetStopSignal() {
	_jsii_.InvokeVoid(
		s,
		"resetStopSignal",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetSysctl() {
	_jsii_.InvokeVoid(
		s,
		"resetSysctl",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ResetUser() {
	_jsii_.InvokeVoid(
		s,
		"resetUser",
		nil, // no parameters
	)
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
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

func (s *jsiiProxy_ServiceTaskSpecContainerSpecOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

