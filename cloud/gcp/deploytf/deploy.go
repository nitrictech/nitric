// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploytf

import (
	"embed"

	"github.com/aws/jsii-runtime-go"
	dockerprovider "github.com/cdktf/cdktf-provider-docker-go/docker/v11/provider"
	"github.com/cdktf/cdktf-provider-google-go/google/v14/datagoogleclientconfig"
	gcpprovider "github.com/cdktf/cdktf-provider-google-go/google/v14/provider"
	gcpbetaprovider "github.com/cdktf/cdktf-provider-googlebeta-go/googlebeta/v14/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/gcp/common"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/schedule"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/secret"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/service"
	tfstack "github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/topic"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/websocket"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VpcConfig struct {
	Network     string   `mapstructure:"network" json:"network"`
	Subnet      string   `mapstructure:"subnet" json:"subnet"`
	NetworkTags []string `mapstructure:"network-tags" json:"network_tags"`
	AllTraffic  bool     `mapstructure:"all-traffic" json:"all_traffic"`
}

type NitricGcpTerraformProvider struct {
	*deploy.CommonStackDetails
	Stack tfstack.Stack

	serviceIngress  bool
	requiresKvStore bool

	// CmekEnabled - Enable Customer Managed Encryption Keys
	cmekEnabled bool
	vpcConfig   *VpcConfig

	GcpConfig      *common.GcpConfig
	Apis           map[string]api.Api
	Buckets        map[string]bucket.Bucket
	Topics         map[string]topic.Topic
	Schedules      map[string]schedule.Schedule
	Services       map[string]service.Service
	Secrets        map[string]secret.Secret
	Queues         map[string]queue.Queue
	KeyValueStores map[string]keyvalue.Keyvalue
	Websockets     map[string]websocket.Websocket
	RawAttributes  map[string]interface{}

	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricGcpTerraformProvider)(nil)

func (a *NitricGcpTerraformProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	a.GcpConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	a.RawAttributes = attributes

	var ok bool
	a.cmekEnabled, ok = a.RawAttributes["cmek"].(bool)
	if !ok {
		a.cmekEnabled = false
	}

	a.vpcConfig = nil
	vpcConfig, ok := a.RawAttributes["vpc"].(map[string]interface{})
	if ok {
		a.vpcConfig = &VpcConfig{}
		mapstructure.Decode(vpcConfig, a.vpcConfig)
	}

	serviceIngress, ok := a.RawAttributes["internal-service-ingress"].(bool)
	if ok {
		a.serviceIngress = serviceIngress
	}

	return nil
}

// embed the modules directory here
//
//go:embed .nitric/modules/**/*
var modules embed.FS

func (a *NitricGcpTerraformProvider) RequiredProviders() map[string]interface{} {
	return map[string]interface{}{
		"google": map[string]string{
			"source":  "hashicorp/google",
			"version": "~> 6.12.0",
		},
		"google-beta": map[string]string{
			"source":  "hashicorp/google-beta",
			"version": "~> 6.12.0",
		},
	}
}

func (a *NitricGcpTerraformProvider) CdkTfModules() ([]provider.ModuleDirectory, error) {
	return []provider.ModuleDirectory{
		{
			ParentDir: ".nitric/modules",
			Modules:   modules,
		},
	}, nil
}

func (a *NitricGcpTerraformProvider) prepareGcpProviders(stack cdktf.TerraformStack) {
	impersonateSa, impersonateOk := a.RawAttributes["impersonate"].(string)

	tfRegion := cdktf.NewTerraformVariable(stack, jsii.String("region"), &cdktf.TerraformVariableConfig{
		Type:        jsii.String("string"),
		Default:     jsii.String(a.Region),
		Description: jsii.String("The GCP region to deploy resources to"),
	})

	if impersonateSa != "" && impersonateOk {
		gcpprovider.NewGoogleProvider(stack, jsii.String("gcp"), &gcpprovider.GoogleProviderConfig{
			Region:                    tfRegion.StringValue(),
			Project:                   jsii.String(a.GcpConfig.ProjectId),
			ImpersonateServiceAccount: jsii.String(impersonateSa),
		})

		gcpbetaprovider.NewGoogleBetaProvider(stack, jsii.String("gcp_beta"), &gcpbetaprovider.GoogleBetaProviderConfig{
			Region:                    tfRegion.StringValue(),
			Project:                   jsii.String(a.GcpConfig.ProjectId),
			ImpersonateServiceAccount: jsii.String(impersonateSa),
		})
	} else {
		gcpprovider.NewGoogleProvider(stack, jsii.String("gcp"), &gcpprovider.GoogleProviderConfig{
			Region:  tfRegion.StringValue(),
			Project: jsii.String(a.GcpConfig.ProjectId),
		})

		gcpbetaprovider.NewGoogleBetaProvider(stack, jsii.String("gcp_beta"), &gcpbetaprovider.GoogleBetaProviderConfig{
			Region:  tfRegion.StringValue(),
			Project: jsii.String(a.GcpConfig.ProjectId),
		})
	}
}

func (a *NitricGcpTerraformProvider) Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error {
	a.prepareGcpProviders(stack)

	googleConf := datagoogleclientconfig.NewDataGoogleClientConfig(stack, jsii.String("gcp_client_config"), &datagoogleclientconfig.DataGoogleClientConfigConfig{})

	var registryAuths []dockerprovider.DockerProviderRegistryAuth = []dockerprovider.DockerProviderRegistryAuth{
		{
			Address:  jsii.Sprintf("%s-docker.pkg.dev", a.Region),
			Username: jsii.String("oauth2accesstoken"),
			Password: googleConf.AccessToken(),
		},
	}

	dockerprovider.NewDockerProvider(stack, jsii.String("docker"), &dockerprovider.DockerProviderConfig{
		RegistryAuth: registryAuths,
	})

	// if resources has any kv stores, make sure kv is enabled for the stack
	for _, resource := range resources {
		if resource.Id.GetType() == resourcespb.ResourceType_KeyValueStore {
			a.requiresKvStore = true
			break
		}
	}

	a.Stack = tfstack.NewStack(stack, jsii.String("stack"), &tfstack.StackConfig{
		Location:         jsii.String(a.Region),
		StackName:        jsii.String(a.StackName),
		CmekEnabled:      jsii.Bool(a.cmekEnabled),
		FirestoreEnabled: jsii.Bool(a.requiresKvStore),
	})

	return nil
}

func (a *NitricGcpTerraformProvider) Post(stack cdktf.TerraformStack) error {
	// write terraform outputs

	cdktf.NewTerraformOutput(stack, jsii.Sprintf("stack-output"), &cdktf.TerraformOutputConfig{
		Value: a.Stack,
	})

	// loop over all the resources and create outputs for them
	for name, api := range a.Apis {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     api,
		})
	}

	for name, bucket := range a.Buckets {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     bucket,
		})
	}

	for name, topic := range a.Topics {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     topic,
		})
	}

	for name, schedule := range a.Schedules {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     schedule,
		})
	}

	for name, service := range a.Services {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     service,
		})
	}

	for name, secret := range a.Secrets {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     secret,
		})
	}

	for name, queue := range a.Queues {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     queue,
		})
	}

	return nil
}

// // Post - Called after all resources have been created, before the Pulumi Context is concluded
// Post(stack cdktf.TerraformStack) error

func NewNitricGcpProvider() *NitricGcpTerraformProvider {
	return &NitricGcpTerraformProvider{
		Apis:           make(map[string]api.Api),
		Buckets:        make(map[string]bucket.Bucket),
		Services:       make(map[string]service.Service),
		Topics:         make(map[string]topic.Topic),
		Schedules:      make(map[string]schedule.Schedule),
		Secrets:        make(map[string]secret.Secret),
		Queues:         make(map[string]queue.Queue),
		KeyValueStores: make(map[string]keyvalue.Keyvalue),
		Websockets:     make(map[string]websocket.Websocket),
	}
}
