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
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"

	goruntime "runtime"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/common/deploy/env"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricTerraformProvider interface {
	// Init - Initialize the provider with the given attributes, prior to any resource creation or Pulumi Context creation
	Init(attributes map[string]interface{}) error
	// Pre - Called prior to any resource creation, after the Pulumi Context has been established
	Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error

	// CdkTfModules - Return the relative parent directory (root golang packed) and embedded modules directory
	CdkTfModules() (string, fs.FS, error)

	// RequiredProviders - Return a list of required providers for this provider
	RequiredProviders() map[string]interface{}

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
	// SqlDatabase - Deploy a SQL Database
	SqlDatabase(stack cdktf.TerraformStack, name string, config *deploymentspb.SqlDatabase) error

	// Post - Called after all resources have been created, before the Pulumi Context is concluded
	Post(stack cdktf.TerraformStack) error
}

type TerraformProviderServer struct {
	provider NitricTerraformProvider
	runtime  RuntimeProvider
}

func (s *TerraformProviderServer) Up(req *deploymentspb.DeploymentUpRequest, stream deploymentspb.Deployment_UpServer) error {
	if beta, err := env.BETA_PROVIDERS.Bool(); err != nil || !beta {
		return status.Error(codes.FailedPrecondition, "Nitric terraform providers are currently in beta, please add beta-providers to the preview field of your nitric.yaml to enable")
	}

	return createTerraformStackForNitricProvider(req, s.provider, s.runtime)
}

func (s *TerraformProviderServer) Down(req *deploymentspb.DeploymentDownRequest, stream deploymentspb.Deployment_DownServer) error {
	if beta, err := env.BETA_PROVIDERS.Bool(); err != nil || !beta {
		return status.Error(codes.FailedPrecondition, "Nitric terraform providers are currently in beta, please add beta-providers to the preview field of your nitric.yaml to enable")
	}

	return status.Error(codes.Unimplemented, "Down not implemented for Terraform providers, please run terraform destroy against your stack state")
}

func NewTerraformProviderServer(provider NitricTerraformProvider, runtime RuntimeProvider) *TerraformProviderServer {
	return &TerraformProviderServer{
		provider: provider,
		runtime:  runtime,
	}
}

func createTerraformStackForNitricProvider(req *deploymentspb.DeploymentUpRequest, nitricProvider NitricTerraformProvider, runtime RuntimeProvider) (err error) {
	defer func() {
		if r := recover(); r != nil {
			b := make([]byte, 2048) // adjust buffer size to be larger than expected stack
			n := goruntime.Stack(b, false)
			s := string(b[:n])
			err = fmt.Errorf("panic: %v [%s]", r, s)
		}
	}()

	projectName, stackName, err := stackAndProjectFromAttributes(req.Attributes.AsMap())
	if err != nil {
		return err
	}

	fullStackName := fmt.Sprintf("%s-%s", projectName, stackName)

	parentDir, modules, err := nitricProvider.CdkTfModules()
	if err != nil {
		return err
	}

	// modules dir
	modulesDir := filepath.Join(parentDir)

	err = os.MkdirAll(modulesDir, 0o750)
	if err != nil {
		return err
	}
	// cleanup the modules when we're done
	// NOTE: Its importent to ensure that the modules are written to a temporary directory like .nitric
	defer os.RemoveAll(modulesDir)

	err = fs.WalkDir(modules, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			data, err := modules.Open(path)
			if err != nil {
				return err
			}
			defer data.Close()

			//#nosec G304 -- path unpacking known modules embedded fs
			out, err := os.Create(path)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, data)
			return err
		}

		return os.MkdirAll(path, 0o750)
	})
	if err != nil {
		return err
	}

	appCtx := map[string]interface{}{
		"cdktfRelativeModules": []string{filepath.Join(modulesDir)},
		// Ensure static output
		"cdktfStaticModuleAssetHash": "nitric_modules",
	}

	attributesMap := req.Attributes.AsMap()
	outdir, ok := attributesMap["outdir"].(string)
	if !ok {
		outdir = "cdktf.out"
	}

	// TODO: This setting currently doesn't work
	// actually needs an env variable to be set, but output is broken
	hcl, ok := attributesMap["hcl"].(bool)
	if !ok {
		hcl = false
	}

	app := cdktf.NewApp(&cdktf.AppConfig{
		HclOutput: jsii.Bool(hcl),
		Outdir:    jsii.String(outdir),
		Context:   &appCtx,
	})

	err = nitricProvider.Init(attributesMap)
	if err != nil {
		return err
	}

	stack := cdktf.NewTerraformStack(app, &fullStackName)

	stack.AddOverride(jsii.String("terraform.required_providers"), nitricProvider.RequiredProviders())

	// The code that defines your stack goes here
	resources := nitricProvider.Order(req.Spec.Resources)

	// TODO: Ideally this would be configured via a NewBackend for type safety
	// instead allowing for arbitrary map overrides that map directly to the backend
	// property of any terraform based stack file
	backend, ok := attributesMap["backend"].(map[string]interface{})
	// if ok {
	// 	stack.AddOverride(jsii.String("terraform.backend"), nil)
	// 	// Do backend as map based override for now
	// 	stack.AddOverride(jsii.String("terraform.backend"), backend)
	// }

	if ok && len(backend) > 1 {
		logger.Fatalf("Only one backend is supported, found %d", len(backend))
	}

	if ok && len(backend) == 1 {
		backendType := lo.Keys(backend)[0]
		config := backend[backendType].(map[string]interface{})
		jsonMap, err := json.Marshal(config)
		if err != nil {
			logger.Fatalf("Failed to serialize backend config %v", err)
			return err
		}

		switch backendType {
		case "gcs":
			gcsConfig := &cdktf.GcsBackendConfig{}
			// serialize the backend config
			err := json.Unmarshal(jsonMap, gcsConfig)
			if err != nil {
				return err
			}
			cdktf.NewGcsBackend(stack, gcsConfig)
		case "s3":
			s3Config := &cdktf.S3BackendConfig{}
			// serialize the backend config
			err := json.Unmarshal(jsonMap, s3Config)
			if err != nil {
				return err
			}
			cdktf.NewS3Backend(stack, s3Config)
		default:
			logger.Fatalf("Unsupported backend type %s", backendType)
		}
	}

	err = nitricProvider.Pre(stack, resources)
	if err != nil {
		return err
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
		case *deploymentspb.Resource_SqlDatabase:
			err = nitricProvider.SqlDatabase(stack, res.Id.Name, t.SqlDatabase)
		}
		if err != nil {
			return err
		}
	}

	err = nitricProvider.Post(stack)
	if err != nil {
		return err
	}

	app.ToString()

	// result, err := nitricProvider.Result(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	app.Synth()

	return nil
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
