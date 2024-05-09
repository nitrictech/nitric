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
	"fmt"
	"net"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/common/deploy/env"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricTerraformProvider interface {
	// Init - Initialize the provider with the given attributes, prior to any resource creation or Pulumi Context creation
	Init(attributes map[string]interface{}) error
	// Pre - Called prior to any resource creation, after the Pulumi Context has been established
	Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error

	// Order - Return the order that resources should be deployed in.
	// The order of resources is important as some resources depend on others.
	// Changing the default order is not recommended unless you know what you are doing.
	Order(resources []*deploymentspb.Resource) []*deploymentspb.Resource

	// Api - Deploy an API Gateway
	Api(tack cdktf.TerraformStack, name string, config *deploymentspb.Api) error
	// Http - Deploy a HTTP Proxy
	Http(tack cdktf.TerraformStack, name string, config *deploymentspb.Http) error
	// Bucket - Deploy a Storage Bucket
	Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error
	// Service - Deploy an service (Service)
	Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider RuntimeProvider) error
	// Topic - Deploy a Pub/Sub Topic
	Topic(stack cdktf.TerraformStack, name string, config *deploymentspb.Topic) error
	// Queue - Deploy a Queue
	Queue(stack cdktf.TerraformStack, name string, config *deploymentspb.Queue) error
	// Secret - Deploy a Secret
	Secret(stack cdktf.TerraformStack, name string, config *deploymentspb.Secret) error
	// Schedule - Deploy a Schedule
	Schedule(stack cdktf.TerraformStack, name string, config *deploymentspb.Schedule) error
	// Websocket - Deploy a Websocket Gateway
	Websocket(stack cdktf.TerraformStack, name string, config *deploymentspb.Websocket) error
	// Policy - Deploy a Policy
	Policy(stack cdktf.TerraformStack, name string, config *deploymentspb.Policy) error
	// KeyValueStore - Deploy a Key Value Store
	KeyValueStore(stack cdktf.TerraformStack, name string, config *deploymentspb.KeyValueStore) error

	// Post - Called after all resources have been created, before the Pulumi Context is concluded
	Post(stack cdktf.TerraformStack) error

	// Result - Last method to be called, return the result of the deployment to be printed to stdout
	// Result(stack cdktf.TerraformStack) (string, error)
}

type TerraformProviderServer struct {
	provider NitricTerraformProvider
	runtime  RuntimeProvider
}

func (s *TerraformProviderServer) Up(req *deploymentspb.DeploymentUpRequest, stream deploymentspb.Deployment_UpServer) error {
	app, err := createTerraformStackForNitricProvider(req, s.provider, s.runtime)
	if err != nil {
		return err
	}

	// Sythesize the Terraform stack
	app.Synth()

	return nil
}

func (s *TerraformProviderServer) Down(req *deploymentspb.DeploymentDownRequest, stream deploymentspb.Deployment_DownServer) error {
	return status.Error(codes.Unimplemented, "Down not implemented for Terraform providers, please run terraform destroy against your stack state")
}

func NewTerraformProviderServer(provider NitricTerraformProvider, runtime RuntimeProvider) *TerraformProviderServer {
	return &TerraformProviderServer{
		provider: provider,
		runtime:  runtime,
	}
}

func createTerraformStackForNitricProvider(req *deploymentspb.DeploymentUpRequest, nitricProvider NitricTerraformProvider, runtime RuntimeProvider) (cdktf.App, error) {
	projectName, stackName, err := stackAndProjectFromAttributes(req.Attributes.AsMap())
	if err != nil {
		return nil, err
	}

	fullStackName := fmt.Sprintf("%s-%s", projectName, stackName)

	app := cdktf.NewApp(nil)

	stack := cdktf.NewTerraformStack(app, &fullStackName)

	// The code that defines your stack goes here
	resources := nitricProvider.Order(req.Spec.Resources)

	err = nitricProvider.Pre(stack, resources)
	if err != nil {
		return nil, err
	}

	for _, res := range resources {
		switch t := res.Config.(type) {
		case *deploymentspb.Resource_Service:
			err = nitricProvider.Service(stack, res.Id.Name, t.Service, runtime)
		case *deploymentspb.Resource_Secret:
			err = nitricProvider.Secret(stack, res.Id.Name, t.Secret)
		case *deploymentspb.Resource_Topic:
			err = nitricProvider.Topic(stack, res.Id.Name, t.Topic)
		case *deploymentspb.Resource_Queue:
			err = nitricProvider.Queue(stack, res.Id.Name, t.Queue)
		case *deploymentspb.Resource_Bucket:
			err = nitricProvider.Bucket(stack, res.Id.Name, t.Bucket)
		case *deploymentspb.Resource_Api:
			err = nitricProvider.Api(stack, res.Id.Name, t.Api)
		case *deploymentspb.Resource_Websocket:
			err = nitricProvider.Websocket(stack, res.Id.Name, t.Websocket)
		case *deploymentspb.Resource_Schedule:
			err = nitricProvider.Schedule(stack, res.Id.Name, t.Schedule)
		case *deploymentspb.Resource_Policy:
			err = nitricProvider.Policy(stack, res.Id.Name, t.Policy)
		case *deploymentspb.Resource_Http:
			err = nitricProvider.Http(stack, res.Id.Name, t.Http)
		case *deploymentspb.Resource_KeyValueStore:
			err = nitricProvider.KeyValueStore(stack, res.Id.Name, t.KeyValueStore)
		}
		if err != nil {
			return nil, err
		}
	}

	err = nitricProvider.Post(stack)
	if err != nil {
		return nil, err
	}

	// result, err := nitricProvider.Result(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return app, nil
}

func (s *TerraformProviderServer) Start() {
	port := env.PORT.String()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("error listening on port %s %v", port, err)
	}

	srv := grpc.NewServer()

	deploymentspb.RegisterDeploymentServer(srv, s)

	fmt.Printf("Deployment server started on %s\n", lis.Addr().String())
	err = srv.Serve(lis)
	if err != nil {
		logger.Fatalf("error serving requests %v", err)
	}
}
