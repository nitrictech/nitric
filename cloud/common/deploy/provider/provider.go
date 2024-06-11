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

package provider

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

type NitricPulumiProvider interface {
	// Init - Initialize the provider with the given attributes, prior to any resource creation or Pulumi Context creation
	Init(attributes map[string]interface{}) error
	// Pre - Called prior to any resource creation, after the Pulumi Context has been established
	Pre(ctx *pulumi.Context, resources []*pulumix.NitricPulumiResource[any]) error
	// Config - Return the Pulumi ConfigMap for the provider
	Config() (auto.ConfigMap, error)

	// Order - Return the order that resources should be deployed in.
	// The order of resources is important as some resources depend on others.
	// Changing the default order is not recommended unless you know what you are doing.
	Order(resources []*deploymentspb.Resource) []*deploymentspb.Resource

	// Api - Deploy an API Gateway
	Api(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Api) error
	// Http - Deploy a HTTP Proxy
	Http(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Http) error
	// Bucket - Deploy a Storage Bucket
	Bucket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Bucket) error
	// Service - Deploy an service (Service)
	Service(ctx *pulumi.Context, parent pulumi.Resource, name string, config *pulumix.NitricPulumiServiceConfig, runtimeProvider RuntimeProvider) error
	// Topic - Deploy a Pub/Sub Topic
	Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error
	// Queue - Deploy a Queue
	Queue(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Queue) error
	// Secret - Deploy a Secret
	Secret(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Secret) error
	// Schedule - Deploy a Schedule
	Schedule(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Schedule) error
	// Websocket - Deploy a Websocket Gateway
	Websocket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Websocket) error
	// Policy - Deploy a Policy
	Policy(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Policy) error
	// KeyValueStore - Deploy a Key Value Store
	KeyValueStore(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.KeyValueStore) error
	// SqlDatabase - Deploy a SQL Database
	SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error

	// Post - Called after all resources have been created, before the Pulumi Context is concluded
	Post(ctx *pulumi.Context) error

	// Result - Last method to be called, return the result of the deployment to be printed to stdout
	Result(ctx *pulumi.Context) (pulumi.StringOutput, error)
}

// NitricDefaultOrder - Partial implementation of NitricPulumiProvider which implements the standard resource deployment order
type NitricDefaultOrder struct{}

func just(all []*deploymentspb.Resource, only resourcespb.ResourceType) []*deploymentspb.Resource {
	return lo.Filter(all, func(item *deploymentspb.Resource, index int) bool {
		return item.Id.Type == only
	})
}

// Order - the default resource deployment order
// By default deploy services (services) first, other resources typically depend on them
// e.g. topics may need to know about services in order to setup subscriptions.
func (*NitricDefaultOrder) Order(resources []*deploymentspb.Resource) []*deploymentspb.Resource {
	typeOrder := []resourcespb.ResourceType{
		resourcespb.ResourceType_SqlDatabase,
		resourcespb.ResourceType_Service,
		resourcespb.ResourceType_Secret,
		resourcespb.ResourceType_Queue,
		resourcespb.ResourceType_Topic,
		resourcespb.ResourceType_Bucket,
		resourcespb.ResourceType_KeyValueStore,
		resourcespb.ResourceType_Api,
		resourcespb.ResourceType_Websocket,
		resourcespb.ResourceType_Schedule,
		resourcespb.ResourceType_Http,
		resourcespb.ResourceType_Policy,
	}

	sorted := []*deploymentspb.Resource{}
	for _, resourceType := range typeOrder {
		sorted = append(sorted, just(resources, resourceType)...)
	}

	return sorted
}
