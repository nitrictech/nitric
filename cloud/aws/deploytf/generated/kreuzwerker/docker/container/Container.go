package container

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/container/internal"
)

// Represents a {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container docker_container}.
type Container interface {
	cdktf.TerraformResource
	Attach() interface{}
	SetAttach(val interface{})
	AttachInput() interface{}
	Bridge() *string
	Capabilities() ContainerCapabilitiesOutputReference
	CapabilitiesInput() *ContainerCapabilities
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CgroupnsMode() *string
	SetCgroupnsMode(val *string)
	CgroupnsModeInput() *string
	Command() *[]*string
	SetCommand(val *[]*string)
	CommandInput() *[]*string
	// Experimental.
	Connection() interface{}
	// Experimental.
	SetConnection(val interface{})
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	ContainerLogs() *string
	ContainerReadRefreshTimeoutMilliseconds() *float64
	SetContainerReadRefreshTimeoutMilliseconds(val *float64)
	ContainerReadRefreshTimeoutMillisecondsInput() *float64
	// Experimental.
	Count() interface{}
	// Experimental.
	SetCount(val interface{})
	CpuSet() *string
	SetCpuSet(val *string)
	CpuSetInput() *string
	CpuShares() *float64
	SetCpuShares(val *float64)
	CpuSharesInput() *float64
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	DestroyGraceSeconds() *float64
	SetDestroyGraceSeconds(val *float64)
	DestroyGraceSecondsInput() *float64
	Devices() ContainerDevicesList
	DevicesInput() interface{}
	Dns() *[]*string
	SetDns(val *[]*string)
	DnsInput() *[]*string
	DnsOpts() *[]*string
	SetDnsOpts(val *[]*string)
	DnsOptsInput() *[]*string
	DnsSearch() *[]*string
	SetDnsSearch(val *[]*string)
	DnsSearchInput() *[]*string
	Domainname() *string
	SetDomainname(val *string)
	DomainnameInput() *string
	Entrypoint() *[]*string
	SetEntrypoint(val *[]*string)
	EntrypointInput() *[]*string
	Env() *[]*string
	SetEnv(val *[]*string)
	EnvInput() *[]*string
	ExitCode() *float64
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	Gpus() *string
	SetGpus(val *string)
	GpusInput() *string
	GroupAdd() *[]*string
	SetGroupAdd(val *[]*string)
	GroupAddInput() *[]*string
	Healthcheck() ContainerHealthcheckOutputReference
	HealthcheckInput() *ContainerHealthcheck
	Host() ContainerHostList
	HostInput() interface{}
	Hostname() *string
	SetHostname(val *string)
	HostnameInput() *string
	Id() *string
	SetId(val *string)
	IdInput() *string
	Image() *string
	SetImage(val *string)
	ImageInput() *string
	Init() interface{}
	SetInit(val interface{})
	InitInput() interface{}
	IpcMode() *string
	SetIpcMode(val *string)
	IpcModeInput() *string
	Labels() ContainerLabelsList
	LabelsInput() interface{}
	// Experimental.
	Lifecycle() *cdktf.TerraformResourceLifecycle
	// Experimental.
	SetLifecycle(val *cdktf.TerraformResourceLifecycle)
	LogDriver() *string
	SetLogDriver(val *string)
	LogDriverInput() *string
	LogOpts() *map[string]*string
	SetLogOpts(val *map[string]*string)
	LogOptsInput() *map[string]*string
	Logs() interface{}
	SetLogs(val interface{})
	LogsInput() interface{}
	MaxRetryCount() *float64
	SetMaxRetryCount(val *float64)
	MaxRetryCountInput() *float64
	Memory() *float64
	SetMemory(val *float64)
	MemoryInput() *float64
	MemorySwap() *float64
	SetMemorySwap(val *float64)
	MemorySwapInput() *float64
	Mounts() ContainerMountsList
	MountsInput() interface{}
	MustRun() interface{}
	SetMustRun(val interface{})
	MustRunInput() interface{}
	Name() *string
	SetName(val *string)
	NameInput() *string
	NetworkData() ContainerNetworkDataList
	NetworkMode() *string
	SetNetworkMode(val *string)
	NetworkModeInput() *string
	NetworksAdvanced() ContainerNetworksAdvancedList
	NetworksAdvancedInput() interface{}
	// The tree node.
	Node() constructs.Node
	PidMode() *string
	SetPidMode(val *string)
	PidModeInput() *string
	Ports() ContainerPortsList
	PortsInput() interface{}
	Privileged() interface{}
	SetPrivileged(val interface{})
	PrivilegedInput() interface{}
	// Experimental.
	Provider() cdktf.TerraformProvider
	// Experimental.
	SetProvider(val cdktf.TerraformProvider)
	// Experimental.
	Provisioners() *[]interface{}
	// Experimental.
	SetProvisioners(val *[]interface{})
	PublishAllPorts() interface{}
	SetPublishAllPorts(val interface{})
	PublishAllPortsInput() interface{}
	// Experimental.
	RawOverrides() interface{}
	ReadOnly() interface{}
	SetReadOnly(val interface{})
	ReadOnlyInput() interface{}
	RemoveVolumes() interface{}
	SetRemoveVolumes(val interface{})
	RemoveVolumesInput() interface{}
	Restart() *string
	SetRestart(val *string)
	RestartInput() *string
	Rm() interface{}
	SetRm(val interface{})
	RmInput() interface{}
	Runtime() *string
	SetRuntime(val *string)
	RuntimeInput() *string
	SecurityOpts() *[]*string
	SetSecurityOpts(val *[]*string)
	SecurityOptsInput() *[]*string
	ShmSize() *float64
	SetShmSize(val *float64)
	ShmSizeInput() *float64
	Start() interface{}
	SetStart(val interface{})
	StartInput() interface{}
	StdinOpen() interface{}
	SetStdinOpen(val interface{})
	StdinOpenInput() interface{}
	StopSignal() *string
	SetStopSignal(val *string)
	StopSignalInput() *string
	StopTimeout() *float64
	SetStopTimeout(val *float64)
	StopTimeoutInput() *float64
	StorageOpts() *map[string]*string
	SetStorageOpts(val *map[string]*string)
	StorageOptsInput() *map[string]*string
	Sysctls() *map[string]*string
	SetSysctls(val *map[string]*string)
	SysctlsInput() *map[string]*string
	// Experimental.
	TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata
	// Experimental.
	TerraformMetaArguments() *map[string]interface{}
	// Experimental.
	TerraformResourceType() *string
	Tmpfs() *map[string]*string
	SetTmpfs(val *map[string]*string)
	TmpfsInput() *map[string]*string
	Tty() interface{}
	SetTty(val interface{})
	TtyInput() interface{}
	Ulimit() ContainerUlimitList
	UlimitInput() interface{}
	Upload() ContainerUploadList
	UploadInput() interface{}
	User() *string
	SetUser(val *string)
	UserInput() *string
	UsernsMode() *string
	SetUsernsMode(val *string)
	UsernsModeInput() *string
	Volumes() ContainerVolumesList
	VolumesInput() interface{}
	Wait() interface{}
	SetWait(val interface{})
	WaitInput() interface{}
	WaitTimeout() *float64
	SetWaitTimeout(val *float64)
	WaitTimeoutInput() *float64
	WorkingDir() *string
	SetWorkingDir(val *string)
	WorkingDirInput() *string
	// Adds a user defined moveTarget string to this resource to be later used in .moveTo(moveTarget) to resolve the location of the move.
	// Experimental.
	AddMoveTarget(moveTarget *string)
	// Experimental.
	AddOverride(path *string, value interface{})
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
	HasResourceMove() interface{}
	// Experimental.
	ImportFrom(id *string, provider cdktf.TerraformProvider)
	// Experimental.
	InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable
	// Move the resource corresponding to "id" to this resource.
	//
	// Note that the resource being moved from must be marked as moved using it's instance function.
	// Experimental.
	MoveFromId(id *string)
	// Moves this resource to the target resource given by moveTarget.
	// Experimental.
	MoveTo(moveTarget *string, index interface{})
	// Moves this resource to the resource corresponding to "id".
	// Experimental.
	MoveToId(id *string)
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
	PutCapabilities(value *ContainerCapabilities)
	PutDevices(value interface{})
	PutHealthcheck(value *ContainerHealthcheck)
	PutHost(value interface{})
	PutLabels(value interface{})
	PutMounts(value interface{})
	PutNetworksAdvanced(value interface{})
	PutPorts(value interface{})
	PutUlimit(value interface{})
	PutUpload(value interface{})
	PutVolumes(value interface{})
	ResetAttach()
	ResetCapabilities()
	ResetCgroupnsMode()
	ResetCommand()
	ResetContainerReadRefreshTimeoutMilliseconds()
	ResetCpuSet()
	ResetCpuShares()
	ResetDestroyGraceSeconds()
	ResetDevices()
	ResetDns()
	ResetDnsOpts()
	ResetDnsSearch()
	ResetDomainname()
	ResetEntrypoint()
	ResetEnv()
	ResetGpus()
	ResetGroupAdd()
	ResetHealthcheck()
	ResetHost()
	ResetHostname()
	ResetId()
	ResetInit()
	ResetIpcMode()
	ResetLabels()
	ResetLogDriver()
	ResetLogOpts()
	ResetLogs()
	ResetMaxRetryCount()
	ResetMemory()
	ResetMemorySwap()
	ResetMounts()
	ResetMustRun()
	ResetNetworkMode()
	ResetNetworksAdvanced()
	// Resets a previously passed logical Id to use the auto-generated logical id again.
	// Experimental.
	ResetOverrideLogicalId()
	ResetPidMode()
	ResetPorts()
	ResetPrivileged()
	ResetPublishAllPorts()
	ResetReadOnly()
	ResetRemoveVolumes()
	ResetRestart()
	ResetRm()
	ResetRuntime()
	ResetSecurityOpts()
	ResetShmSize()
	ResetStart()
	ResetStdinOpen()
	ResetStopSignal()
	ResetStopTimeout()
	ResetStorageOpts()
	ResetSysctls()
	ResetTmpfs()
	ResetTty()
	ResetUlimit()
	ResetUpload()
	ResetUser()
	ResetUsernsMode()
	ResetVolumes()
	ResetWait()
	ResetWaitTimeout()
	ResetWorkingDir()
	SynthesizeAttributes() *map[string]interface{}
	SynthesizeHclAttributes() *map[string]interface{}
	// Experimental.
	ToHclTerraform() interface{}
	// Experimental.
	ToMetadata() interface{}
	// Returns a string representation of this construct.
	ToString() *string
	// Adds this resource to the terraform JSON output.
	// Experimental.
	ToTerraform() interface{}
}

// The jsii proxy struct for Container
type jsiiProxy_Container struct {
	internal.Type__cdktfTerraformResource
}

func (j *jsiiProxy_Container) Attach() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"attach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) AttachInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"attachInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Bridge() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bridge",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Capabilities() ContainerCapabilitiesOutputReference {
	var returns ContainerCapabilitiesOutputReference
	_jsii_.Get(
		j,
		"capabilities",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CapabilitiesInput() *ContainerCapabilities {
	var returns *ContainerCapabilities
	_jsii_.Get(
		j,
		"capabilitiesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CgroupnsMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cgroupnsMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CgroupnsModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cgroupnsModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Command() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"command",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CommandInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"commandInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Connection() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"connection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ContainerLogs() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerLogs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ContainerReadRefreshTimeoutMilliseconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"containerReadRefreshTimeoutMilliseconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ContainerReadRefreshTimeoutMillisecondsInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"containerReadRefreshTimeoutMillisecondsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Count() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"count",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CpuSet() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CpuSetInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSetInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CpuShares() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuShares",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) CpuSharesInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuSharesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DestroyGraceSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"destroyGraceSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DestroyGraceSecondsInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"destroyGraceSecondsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Devices() ContainerDevicesList {
	var returns ContainerDevicesList
	_jsii_.Get(
		j,
		"devices",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DevicesInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"devicesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Dns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DnsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dnsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DnsOpts() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dnsOpts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DnsOptsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dnsOptsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DnsSearch() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dnsSearch",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DnsSearchInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dnsSearchInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Domainname() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainname",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) DomainnameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainnameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Entrypoint() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"entrypoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) EntrypointInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"entrypointInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Env() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) EnvInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"envInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ExitCode() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"exitCode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Gpus() *string {
	var returns *string
	_jsii_.Get(
		j,
		"gpus",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) GpusInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"gpusInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) GroupAdd() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groupAdd",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) GroupAddInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groupAddInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Healthcheck() ContainerHealthcheckOutputReference {
	var returns ContainerHealthcheckOutputReference
	_jsii_.Get(
		j,
		"healthcheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) HealthcheckInput() *ContainerHealthcheck {
	var returns *ContainerHealthcheck
	_jsii_.Get(
		j,
		"healthcheckInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Host() ContainerHostList {
	var returns ContainerHostList
	_jsii_.Get(
		j,
		"host",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) HostInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"hostInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Hostname() *string {
	var returns *string
	_jsii_.Get(
		j,
		"hostname",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) HostnameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"hostnameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) IdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"idInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Image() *string {
	var returns *string
	_jsii_.Get(
		j,
		"image",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ImageInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Init() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"init",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) InitInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"initInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) IpcMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipcMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) IpcModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipcModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Labels() ContainerLabelsList {
	var returns ContainerLabelsList
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LabelsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Lifecycle() *cdktf.TerraformResourceLifecycle {
	var returns *cdktf.TerraformResourceLifecycle
	_jsii_.Get(
		j,
		"lifecycle",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LogDriver() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logDriver",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LogDriverInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logDriverInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LogOpts() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"logOpts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LogOptsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"logOptsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Logs() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"logs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) LogsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"logsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MaxRetryCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxRetryCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MaxRetryCountInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxRetryCountInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Memory() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memory",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MemoryInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memoryInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MemorySwap() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memorySwap",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MemorySwapInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memorySwapInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Mounts() ContainerMountsList {
	var returns ContainerMountsList
	_jsii_.Get(
		j,
		"mounts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MountsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mountsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MustRun() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mustRun",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) MustRunInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mustRunInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"nameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NetworkData() ContainerNetworkDataList {
	var returns ContainerNetworkDataList
	_jsii_.Get(
		j,
		"networkData",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NetworkMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NetworkModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"networkModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NetworksAdvanced() ContainerNetworksAdvancedList {
	var returns ContainerNetworksAdvancedList
	_jsii_.Get(
		j,
		"networksAdvanced",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) NetworksAdvancedInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"networksAdvancedInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PidMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pidMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PidModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pidModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Ports() ContainerPortsList {
	var returns ContainerPortsList
	_jsii_.Get(
		j,
		"ports",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PortsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"portsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Privileged() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"privileged",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PrivilegedInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"privilegedInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Provider() cdktf.TerraformProvider {
	var returns cdktf.TerraformProvider
	_jsii_.Get(
		j,
		"provider",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Provisioners() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"provisioners",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PublishAllPorts() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"publishAllPorts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) PublishAllPortsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"publishAllPortsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ReadOnly() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnly",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ReadOnlyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"readOnlyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RemoveVolumes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"removeVolumes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RemoveVolumesInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"removeVolumesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Restart() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restart",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RestartInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restartInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Rm() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rm",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RmInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rmInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Runtime() *string {
	var returns *string
	_jsii_.Get(
		j,
		"runtime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) RuntimeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"runtimeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) SecurityOpts() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityOpts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) SecurityOptsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityOptsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ShmSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"shmSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) ShmSizeInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"shmSizeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Start() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"start",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StartInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"startInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StdinOpen() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stdinOpen",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StdinOpenInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stdinOpenInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StopSignal() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopSignal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StopSignalInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stopSignalInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StopTimeout() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"stopTimeout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StopTimeoutInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"stopTimeoutInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StorageOpts() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"storageOpts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) StorageOptsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"storageOptsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Sysctls() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"sysctls",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) SysctlsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"sysctlsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata {
	var returns *cdktf.TerraformProviderGeneratorMetadata
	_jsii_.Get(
		j,
		"terraformGeneratorMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) TerraformMetaArguments() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"terraformMetaArguments",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) TerraformResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Tmpfs() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"tmpfs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) TmpfsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"tmpfsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Tty() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"tty",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) TtyInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ttyInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Ulimit() ContainerUlimitList {
	var returns ContainerUlimitList
	_jsii_.Get(
		j,
		"ulimit",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) UlimitInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ulimitInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Upload() ContainerUploadList {
	var returns ContainerUploadList
	_jsii_.Get(
		j,
		"upload",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) UploadInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"uploadInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) User() *string {
	var returns *string
	_jsii_.Get(
		j,
		"user",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) UserInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) UsernsMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usernsMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) UsernsModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usernsModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Volumes() ContainerVolumesList {
	var returns ContainerVolumesList
	_jsii_.Get(
		j,
		"volumes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) VolumesInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"volumesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) Wait() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"wait",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) WaitInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"waitInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) WaitTimeout() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"waitTimeout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) WaitTimeoutInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"waitTimeoutInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) WorkingDir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"workingDir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Container) WorkingDirInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"workingDirInput",
		&returns,
	)
	return returns
}


// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container docker_container} Resource.
func NewContainer(scope constructs.Construct, id *string, config *ContainerConfig) Container {
	_init_.Initialize()

	if err := validateNewContainerParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Container{}

	_jsii_.Create(
		"docker.container.Container",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container docker_container} Resource.
func NewContainer_Override(c Container, scope constructs.Construct, id *string, config *ContainerConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.container.Container",
		[]interface{}{scope, id, config},
		c,
	)
}

func (j *jsiiProxy_Container)SetAttach(val interface{}) {
	if err := j.validateSetAttachParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"attach",
		val,
	)
}

func (j *jsiiProxy_Container)SetCgroupnsMode(val *string) {
	if err := j.validateSetCgroupnsModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cgroupnsMode",
		val,
	)
}

func (j *jsiiProxy_Container)SetCommand(val *[]*string) {
	if err := j.validateSetCommandParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"command",
		val,
	)
}

func (j *jsiiProxy_Container)SetConnection(val interface{}) {
	if err := j.validateSetConnectionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"connection",
		val,
	)
}

func (j *jsiiProxy_Container)SetContainerReadRefreshTimeoutMilliseconds(val *float64) {
	if err := j.validateSetContainerReadRefreshTimeoutMillisecondsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"containerReadRefreshTimeoutMilliseconds",
		val,
	)
}

func (j *jsiiProxy_Container)SetCount(val interface{}) {
	if err := j.validateSetCountParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"count",
		val,
	)
}

func (j *jsiiProxy_Container)SetCpuSet(val *string) {
	if err := j.validateSetCpuSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuSet",
		val,
	)
}

func (j *jsiiProxy_Container)SetCpuShares(val *float64) {
	if err := j.validateSetCpuSharesParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuShares",
		val,
	)
}

func (j *jsiiProxy_Container)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Container)SetDestroyGraceSeconds(val *float64) {
	if err := j.validateSetDestroyGraceSecondsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"destroyGraceSeconds",
		val,
	)
}

func (j *jsiiProxy_Container)SetDns(val *[]*string) {
	if err := j.validateSetDnsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dns",
		val,
	)
}

func (j *jsiiProxy_Container)SetDnsOpts(val *[]*string) {
	if err := j.validateSetDnsOptsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dnsOpts",
		val,
	)
}

func (j *jsiiProxy_Container)SetDnsSearch(val *[]*string) {
	if err := j.validateSetDnsSearchParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dnsSearch",
		val,
	)
}

func (j *jsiiProxy_Container)SetDomainname(val *string) {
	if err := j.validateSetDomainnameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"domainname",
		val,
	)
}

func (j *jsiiProxy_Container)SetEntrypoint(val *[]*string) {
	if err := j.validateSetEntrypointParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"entrypoint",
		val,
	)
}

func (j *jsiiProxy_Container)SetEnv(val *[]*string) {
	if err := j.validateSetEnvParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"env",
		val,
	)
}

func (j *jsiiProxy_Container)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Container)SetGpus(val *string) {
	if err := j.validateSetGpusParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"gpus",
		val,
	)
}

func (j *jsiiProxy_Container)SetGroupAdd(val *[]*string) {
	if err := j.validateSetGroupAddParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"groupAdd",
		val,
	)
}

func (j *jsiiProxy_Container)SetHostname(val *string) {
	if err := j.validateSetHostnameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"hostname",
		val,
	)
}

func (j *jsiiProxy_Container)SetId(val *string) {
	if err := j.validateSetIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"id",
		val,
	)
}

func (j *jsiiProxy_Container)SetImage(val *string) {
	if err := j.validateSetImageParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"image",
		val,
	)
}

func (j *jsiiProxy_Container)SetInit(val interface{}) {
	if err := j.validateSetInitParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"init",
		val,
	)
}

func (j *jsiiProxy_Container)SetIpcMode(val *string) {
	if err := j.validateSetIpcModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"ipcMode",
		val,
	)
}

func (j *jsiiProxy_Container)SetLifecycle(val *cdktf.TerraformResourceLifecycle) {
	if err := j.validateSetLifecycleParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lifecycle",
		val,
	)
}

func (j *jsiiProxy_Container)SetLogDriver(val *string) {
	if err := j.validateSetLogDriverParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"logDriver",
		val,
	)
}

func (j *jsiiProxy_Container)SetLogOpts(val *map[string]*string) {
	if err := j.validateSetLogOptsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"logOpts",
		val,
	)
}

func (j *jsiiProxy_Container)SetLogs(val interface{}) {
	if err := j.validateSetLogsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"logs",
		val,
	)
}

func (j *jsiiProxy_Container)SetMaxRetryCount(val *float64) {
	if err := j.validateSetMaxRetryCountParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"maxRetryCount",
		val,
	)
}

func (j *jsiiProxy_Container)SetMemory(val *float64) {
	if err := j.validateSetMemoryParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"memory",
		val,
	)
}

func (j *jsiiProxy_Container)SetMemorySwap(val *float64) {
	if err := j.validateSetMemorySwapParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"memorySwap",
		val,
	)
}

func (j *jsiiProxy_Container)SetMustRun(val interface{}) {
	if err := j.validateSetMustRunParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"mustRun",
		val,
	)
}

func (j *jsiiProxy_Container)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_Container)SetNetworkMode(val *string) {
	if err := j.validateSetNetworkModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"networkMode",
		val,
	)
}

func (j *jsiiProxy_Container)SetPidMode(val *string) {
	if err := j.validateSetPidModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"pidMode",
		val,
	)
}

func (j *jsiiProxy_Container)SetPrivileged(val interface{}) {
	if err := j.validateSetPrivilegedParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"privileged",
		val,
	)
}

func (j *jsiiProxy_Container)SetProvider(val cdktf.TerraformProvider) {
	_jsii_.Set(
		j,
		"provider",
		val,
	)
}

func (j *jsiiProxy_Container)SetProvisioners(val *[]interface{}) {
	if err := j.validateSetProvisionersParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"provisioners",
		val,
	)
}

func (j *jsiiProxy_Container)SetPublishAllPorts(val interface{}) {
	if err := j.validateSetPublishAllPortsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"publishAllPorts",
		val,
	)
}

func (j *jsiiProxy_Container)SetReadOnly(val interface{}) {
	if err := j.validateSetReadOnlyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"readOnly",
		val,
	)
}

func (j *jsiiProxy_Container)SetRemoveVolumes(val interface{}) {
	if err := j.validateSetRemoveVolumesParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"removeVolumes",
		val,
	)
}

func (j *jsiiProxy_Container)SetRestart(val *string) {
	if err := j.validateSetRestartParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"restart",
		val,
	)
}

func (j *jsiiProxy_Container)SetRm(val interface{}) {
	if err := j.validateSetRmParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"rm",
		val,
	)
}

func (j *jsiiProxy_Container)SetRuntime(val *string) {
	if err := j.validateSetRuntimeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"runtime",
		val,
	)
}

func (j *jsiiProxy_Container)SetSecurityOpts(val *[]*string) {
	if err := j.validateSetSecurityOptsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"securityOpts",
		val,
	)
}

func (j *jsiiProxy_Container)SetShmSize(val *float64) {
	if err := j.validateSetShmSizeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"shmSize",
		val,
	)
}

func (j *jsiiProxy_Container)SetStart(val interface{}) {
	if err := j.validateSetStartParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"start",
		val,
	)
}

func (j *jsiiProxy_Container)SetStdinOpen(val interface{}) {
	if err := j.validateSetStdinOpenParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stdinOpen",
		val,
	)
}

func (j *jsiiProxy_Container)SetStopSignal(val *string) {
	if err := j.validateSetStopSignalParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stopSignal",
		val,
	)
}

func (j *jsiiProxy_Container)SetStopTimeout(val *float64) {
	if err := j.validateSetStopTimeoutParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stopTimeout",
		val,
	)
}

func (j *jsiiProxy_Container)SetStorageOpts(val *map[string]*string) {
	if err := j.validateSetStorageOptsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"storageOpts",
		val,
	)
}

func (j *jsiiProxy_Container)SetSysctls(val *map[string]*string) {
	if err := j.validateSetSysctlsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"sysctls",
		val,
	)
}

func (j *jsiiProxy_Container)SetTmpfs(val *map[string]*string) {
	if err := j.validateSetTmpfsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"tmpfs",
		val,
	)
}

func (j *jsiiProxy_Container)SetTty(val interface{}) {
	if err := j.validateSetTtyParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"tty",
		val,
	)
}

func (j *jsiiProxy_Container)SetUser(val *string) {
	if err := j.validateSetUserParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"user",
		val,
	)
}

func (j *jsiiProxy_Container)SetUsernsMode(val *string) {
	if err := j.validateSetUsernsModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"usernsMode",
		val,
	)
}

func (j *jsiiProxy_Container)SetWait(val interface{}) {
	if err := j.validateSetWaitParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"wait",
		val,
	)
}

func (j *jsiiProxy_Container)SetWaitTimeout(val *float64) {
	if err := j.validateSetWaitTimeoutParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"waitTimeout",
		val,
	)
}

func (j *jsiiProxy_Container)SetWorkingDir(val *string) {
	if err := j.validateSetWorkingDirParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"workingDir",
		val,
	)
}

// Generates CDKTF code for importing a Container resource upon running "cdktf plan <stack-name>".
func Container_GenerateConfigForImport(scope constructs.Construct, importToId *string, importFromId *string, provider cdktf.TerraformProvider) cdktf.ImportableResource {
	_init_.Initialize()

	if err := validateContainer_GenerateConfigForImportParameters(scope, importToId, importFromId); err != nil {
		panic(err)
	}
	var returns cdktf.ImportableResource

	_jsii_.StaticInvoke(
		"docker.container.Container",
		"generateConfigForImport",
		[]interface{}{scope, importToId, importFromId, provider},
		&returns,
	)

	return returns
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
func Container_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateContainer_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.container.Container",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Container_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateContainer_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.container.Container",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Container_IsTerraformResource(x interface{}) *bool {
	_init_.Initialize()

	if err := validateContainer_IsTerraformResourceParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.container.Container",
		"isTerraformResource",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func Container_TfResourceType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"docker.container.Container",
		"tfResourceType",
		&returns,
	)
	return returns
}

func (c *jsiiProxy_Container) AddMoveTarget(moveTarget *string) {
	if err := c.validateAddMoveTargetParameters(moveTarget); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addMoveTarget",
		[]interface{}{moveTarget},
	)
}

func (c *jsiiProxy_Container) AddOverride(path *string, value interface{}) {
	if err := c.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (c *jsiiProxy_Container) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
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

func (c *jsiiProxy_Container) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
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

func (c *jsiiProxy_Container) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
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

func (c *jsiiProxy_Container) GetListAttribute(terraformAttribute *string) *[]*string {
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

func (c *jsiiProxy_Container) GetNumberAttribute(terraformAttribute *string) *float64 {
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

func (c *jsiiProxy_Container) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
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

func (c *jsiiProxy_Container) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
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

func (c *jsiiProxy_Container) GetStringAttribute(terraformAttribute *string) *string {
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

func (c *jsiiProxy_Container) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
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

func (c *jsiiProxy_Container) HasResourceMove() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"hasResourceMove",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) ImportFrom(id *string, provider cdktf.TerraformProvider) {
	if err := c.validateImportFromParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"importFrom",
		[]interface{}{id, provider},
	)
}

func (c *jsiiProxy_Container) InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := c.validateInterpolationForAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		c,
		"interpolationForAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) MoveFromId(id *string) {
	if err := c.validateMoveFromIdParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"moveFromId",
		[]interface{}{id},
	)
}

func (c *jsiiProxy_Container) MoveTo(moveTarget *string, index interface{}) {
	if err := c.validateMoveToParameters(moveTarget, index); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"moveTo",
		[]interface{}{moveTarget, index},
	)
}

func (c *jsiiProxy_Container) MoveToId(id *string) {
	if err := c.validateMoveToIdParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"moveToId",
		[]interface{}{id},
	)
}

func (c *jsiiProxy_Container) OverrideLogicalId(newLogicalId *string) {
	if err := c.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (c *jsiiProxy_Container) PutCapabilities(value *ContainerCapabilities) {
	if err := c.validatePutCapabilitiesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putCapabilities",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutDevices(value interface{}) {
	if err := c.validatePutDevicesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putDevices",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutHealthcheck(value *ContainerHealthcheck) {
	if err := c.validatePutHealthcheckParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putHealthcheck",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutHost(value interface{}) {
	if err := c.validatePutHostParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putHost",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutLabels(value interface{}) {
	if err := c.validatePutLabelsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putLabels",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutMounts(value interface{}) {
	if err := c.validatePutMountsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putMounts",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutNetworksAdvanced(value interface{}) {
	if err := c.validatePutNetworksAdvancedParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putNetworksAdvanced",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutPorts(value interface{}) {
	if err := c.validatePutPortsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putPorts",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutUlimit(value interface{}) {
	if err := c.validatePutUlimitParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putUlimit",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutUpload(value interface{}) {
	if err := c.validatePutUploadParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putUpload",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) PutVolumes(value interface{}) {
	if err := c.validatePutVolumesParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		c,
		"putVolumes",
		[]interface{}{value},
	)
}

func (c *jsiiProxy_Container) ResetAttach() {
	_jsii_.InvokeVoid(
		c,
		"resetAttach",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetCapabilities() {
	_jsii_.InvokeVoid(
		c,
		"resetCapabilities",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetCgroupnsMode() {
	_jsii_.InvokeVoid(
		c,
		"resetCgroupnsMode",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetCommand() {
	_jsii_.InvokeVoid(
		c,
		"resetCommand",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetContainerReadRefreshTimeoutMilliseconds() {
	_jsii_.InvokeVoid(
		c,
		"resetContainerReadRefreshTimeoutMilliseconds",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetCpuSet() {
	_jsii_.InvokeVoid(
		c,
		"resetCpuSet",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetCpuShares() {
	_jsii_.InvokeVoid(
		c,
		"resetCpuShares",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDestroyGraceSeconds() {
	_jsii_.InvokeVoid(
		c,
		"resetDestroyGraceSeconds",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDevices() {
	_jsii_.InvokeVoid(
		c,
		"resetDevices",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDns() {
	_jsii_.InvokeVoid(
		c,
		"resetDns",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDnsOpts() {
	_jsii_.InvokeVoid(
		c,
		"resetDnsOpts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDnsSearch() {
	_jsii_.InvokeVoid(
		c,
		"resetDnsSearch",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetDomainname() {
	_jsii_.InvokeVoid(
		c,
		"resetDomainname",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetEntrypoint() {
	_jsii_.InvokeVoid(
		c,
		"resetEntrypoint",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetEnv() {
	_jsii_.InvokeVoid(
		c,
		"resetEnv",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetGpus() {
	_jsii_.InvokeVoid(
		c,
		"resetGpus",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetGroupAdd() {
	_jsii_.InvokeVoid(
		c,
		"resetGroupAdd",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetHealthcheck() {
	_jsii_.InvokeVoid(
		c,
		"resetHealthcheck",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetHost() {
	_jsii_.InvokeVoid(
		c,
		"resetHost",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetHostname() {
	_jsii_.InvokeVoid(
		c,
		"resetHostname",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetId() {
	_jsii_.InvokeVoid(
		c,
		"resetId",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetInit() {
	_jsii_.InvokeVoid(
		c,
		"resetInit",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetIpcMode() {
	_jsii_.InvokeVoid(
		c,
		"resetIpcMode",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetLabels() {
	_jsii_.InvokeVoid(
		c,
		"resetLabels",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetLogDriver() {
	_jsii_.InvokeVoid(
		c,
		"resetLogDriver",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetLogOpts() {
	_jsii_.InvokeVoid(
		c,
		"resetLogOpts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetLogs() {
	_jsii_.InvokeVoid(
		c,
		"resetLogs",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetMaxRetryCount() {
	_jsii_.InvokeVoid(
		c,
		"resetMaxRetryCount",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetMemory() {
	_jsii_.InvokeVoid(
		c,
		"resetMemory",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetMemorySwap() {
	_jsii_.InvokeVoid(
		c,
		"resetMemorySwap",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetMounts() {
	_jsii_.InvokeVoid(
		c,
		"resetMounts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetMustRun() {
	_jsii_.InvokeVoid(
		c,
		"resetMustRun",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetNetworkMode() {
	_jsii_.InvokeVoid(
		c,
		"resetNetworkMode",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetNetworksAdvanced() {
	_jsii_.InvokeVoid(
		c,
		"resetNetworksAdvanced",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		c,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetPidMode() {
	_jsii_.InvokeVoid(
		c,
		"resetPidMode",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetPorts() {
	_jsii_.InvokeVoid(
		c,
		"resetPorts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetPrivileged() {
	_jsii_.InvokeVoid(
		c,
		"resetPrivileged",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetPublishAllPorts() {
	_jsii_.InvokeVoid(
		c,
		"resetPublishAllPorts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetReadOnly() {
	_jsii_.InvokeVoid(
		c,
		"resetReadOnly",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetRemoveVolumes() {
	_jsii_.InvokeVoid(
		c,
		"resetRemoveVolumes",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetRestart() {
	_jsii_.InvokeVoid(
		c,
		"resetRestart",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetRm() {
	_jsii_.InvokeVoid(
		c,
		"resetRm",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetRuntime() {
	_jsii_.InvokeVoid(
		c,
		"resetRuntime",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetSecurityOpts() {
	_jsii_.InvokeVoid(
		c,
		"resetSecurityOpts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetShmSize() {
	_jsii_.InvokeVoid(
		c,
		"resetShmSize",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetStart() {
	_jsii_.InvokeVoid(
		c,
		"resetStart",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetStdinOpen() {
	_jsii_.InvokeVoid(
		c,
		"resetStdinOpen",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetStopSignal() {
	_jsii_.InvokeVoid(
		c,
		"resetStopSignal",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetStopTimeout() {
	_jsii_.InvokeVoid(
		c,
		"resetStopTimeout",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetStorageOpts() {
	_jsii_.InvokeVoid(
		c,
		"resetStorageOpts",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetSysctls() {
	_jsii_.InvokeVoid(
		c,
		"resetSysctls",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetTmpfs() {
	_jsii_.InvokeVoid(
		c,
		"resetTmpfs",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetTty() {
	_jsii_.InvokeVoid(
		c,
		"resetTty",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetUlimit() {
	_jsii_.InvokeVoid(
		c,
		"resetUlimit",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetUpload() {
	_jsii_.InvokeVoid(
		c,
		"resetUpload",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetUser() {
	_jsii_.InvokeVoid(
		c,
		"resetUser",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetUsernsMode() {
	_jsii_.InvokeVoid(
		c,
		"resetUsernsMode",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetVolumes() {
	_jsii_.InvokeVoid(
		c,
		"resetVolumes",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetWait() {
	_jsii_.InvokeVoid(
		c,
		"resetWait",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetWaitTimeout() {
	_jsii_.InvokeVoid(
		c,
		"resetWaitTimeout",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) ResetWorkingDir() {
	_jsii_.InvokeVoid(
		c,
		"resetWorkingDir",
		nil, // no parameters
	)
}

func (c *jsiiProxy_Container) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Container) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

