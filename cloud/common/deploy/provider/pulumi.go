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
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/env"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc"
)

type PulumiProviderServer struct {
	provider NitricPulumiProvider
	runtime  RuntimeProvider
}

func NewPulumiProviderServer(provider NitricPulumiProvider, runtime RuntimeProvider) *PulumiProviderServer {
	return &PulumiProviderServer{
		provider: provider,
		runtime:  runtime,
	}
}

const resultCtxKey = "nitric:stack:result"

func nitricResourceToPulumiResource(res *deploymentspb.Resource) *pulumix.NitricPulumiResource[any] {
	switch t := res.Config.(type) {
	case *deploymentspb.Resource_Service:
		return &pulumix.NitricPulumiResource[any]{
			Id: res.Id,
			Config: &pulumix.NitricPulumiServiceConfig{
				Service: t.Service,
			},
		}
	default:
		return &pulumix.NitricPulumiResource[any]{
			Id:     res.Id,
			Config: res.Config,
		}
	}
}

func createPulumiProgramForNitricProvider(req *deploymentspb.DeploymentUpRequest, nitricProvider NitricPulumiProvider, runtime RuntimeProvider) func(*pulumi.Context) error {
	return func(ctx *pulumi.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				err = fmt.Errorf("recovered panic: %+v\n Stack: %s", r, stack)
			}
		}()

		// Need to convert the Nitric resources to Pulumi resources, this will allow us to extend their configurations with pulumi inputs/outputs
		pulumiResources := make([]*pulumix.NitricPulumiResource[any], 0, len(req.Spec.Resources))
		for _, res := range nitricProvider.Order(req.Spec.Resources) {
			pulumiResources = append(pulumiResources, nitricResourceToPulumiResource(res))
		}

		// Pre-deployment Hook, used for validation, extension, spec modification, etc.
		err = nitricProvider.Pre(ctx, pulumiResources)
		if err != nil {
			return err
		}

		for _, res := range pulumiResources {
			parent, err := pulumix.ParentResourceFromResourceId(ctx, res.Id)
			if err != nil {
				return err
			}

			switch t := res.Config.(type) {
			case *deploymentspb.Resource_SqlDatabase:
				err = nitricProvider.SqlDatabase(ctx, parent, res.Id.Name, t.SqlDatabase)
			case *pulumix.NitricPulumiServiceConfig:
				err = nitricProvider.Service(ctx, parent, res.Id.Name, t, runtime)
			case *deploymentspb.Resource_Secret:
				err = nitricProvider.Secret(ctx, parent, res.Id.Name, t.Secret)
			case *deploymentspb.Resource_Topic:
				err = nitricProvider.Topic(ctx, parent, res.Id.Name, t.Topic)
			case *deploymentspb.Resource_Queue:
				err = nitricProvider.Queue(ctx, parent, res.Id.Name, t.Queue)
			case *deploymentspb.Resource_Bucket:
				err = nitricProvider.Bucket(ctx, parent, res.Id.Name, t.Bucket)
			case *deploymentspb.Resource_Api:
				err = nitricProvider.Api(ctx, parent, res.Id.Name, t.Api)
			case *deploymentspb.Resource_Websocket:
				err = nitricProvider.Websocket(ctx, parent, res.Id.Name, t.Websocket)
			case *deploymentspb.Resource_Schedule:
				err = nitricProvider.Schedule(ctx, parent, res.Id.Name, t.Schedule)
			case *deploymentspb.Resource_Policy:
				err = nitricProvider.Policy(ctx, parent, res.Id.Name, t.Policy)
			case *deploymentspb.Resource_Http:
				err = nitricProvider.Http(ctx, parent, res.Id.Name, t.Http)
			case *deploymentspb.Resource_KeyValueStore:
				err = nitricProvider.KeyValueStore(ctx, parent, res.Id.Name, t.KeyValueStore)
			}
			if err != nil {
				return err
			}
		}

		err = nitricProvider.Post(ctx)
		if err != nil {
			return err
		}

		result, err := nitricProvider.Result(ctx)
		if err != nil {
			return err
		}

		ctx.Export(resultCtxKey, result)

		// Validate extract and whatever else
		return nil
	}
}

func stackAndProjectFromAttributes(attributesMap map[string]interface{}) (string, string, error) {
	projectName, ok := attributesMap["project"].(string)
	if !ok {
		return "", "", fmt.Errorf("")
	}

	stackName, ok := attributesMap["stack"].(string)
	if !ok {
		return "", "", fmt.Errorf("")
	}

	return projectName, stackName, nil
}

type pulumiError struct {
	Stdout  string
	Stderr  string
	Message string
}

func (p pulumiError) Error() string {
	parts := []string{}
	if p.Message != "" {
		parts = append(parts, p.Message)
	}
	if p.Stdout != "" {
		parts = append(parts, p.Stdout)
	}
	if p.Stderr != "" {
		parts = append(parts, p.Stderr)
	}
	return strings.Join(parts, "\n")
}

func parsePulumiError(err error) error {
	splits := []string{"code: ", "stdout: ", "stderr: "}
	splitIndexes := make([]int, len(splits))

	for i, split := range splits {
		splitIndexes[i] = strings.Index(err.Error(), split)
		if splitIndexes[i] == -1 {
			return nil
		}
	}

	// we have all the splits, we can parse the error
	pe := pulumiError{
		Message: err.Error()[:splitIndexes[0]],
		Stdout:  err.Error()[splitIndexes[1]:splitIndexes[2]],
		Stderr:  err.Error()[splitIndexes[2]:],
	}

	return pe
}

// Up - automatically called by the Nitric CLI via the `up` command
func (s *PulumiProviderServer) Up(req *deploymentspb.DeploymentUpRequest, stream deploymentspb.Deployment_UpServer) error {
	projectName, stackName, err := stackAndProjectFromAttributes(req.Attributes.AsMap())
	if err != nil {
		return err
	}

	attributesMap := req.Attributes.AsMap()

	err = s.provider.Init(attributesMap)
	if err != nil {
		return err
	}

	pulumiProgram := createPulumiProgramForNitricProvider(req, s.provider, s.runtime)

	autoStack, err := auto.UpsertStackInlineSource(context.TODO(), fmt.Sprintf("%s-%s", projectName, stackName), projectName, pulumiProgram)
	if err != nil {
		return err
	}

	pulumiEventsChan := make(chan events.EngineEvent)

	go func() {
		// output the stream
		_ = pulumix.StreamPulumiUpEngineEvents(stream, pulumiEventsChan)
	}()

	config, err := s.provider.Config()
	if err != nil {
		return err
	}

	err = autoStack.SetAllConfig(context.TODO(), config)
	if err != nil {
		return err
	}

	refresh, ok := attributesMap["refresh"].(bool)

	if ok && refresh {
		logger.Info("refreshing pulumi stack")
		_, err = autoStack.Refresh(context.TODO())
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	}

	result, err := autoStack.Up(context.TODO(), optup.EventStreams(pulumiEventsChan))
	if err != nil {
		// Check for common Pulumi 'autoError' types
		if auto.IsConcurrentUpdateError(err) {
			if pe := parsePulumiError(err); pe != nil {
				err = pe
			}
			err = fmt.Errorf("the pulumi stack file is locked.\nThis occurs when a previous deployment is still in progress or was interrupted.\n%w", err)
		} else if auto.IsSelectStack404Error(err) {
			err = fmt.Errorf("stack not found. %w", err)
		} else if auto.IsCreateStack409Error(err) {
			err = fmt.Errorf("failed to create Pulumi stack, this may be a bug in nitric. Seek help https://github.com/nitrictech/nitric/issues\n%w", err)
		} else if auto.IsCompilationError(err) {
			err = fmt.Errorf("failed to compile Pulumi program, this may be a bug in your chosen provider or with nitric. Seek help https://github.com/nitrictech/nitric/issues\n%w", err)
		} else if pe := parsePulumiError(err); pe != nil {
			err = pe
		}

		return err
	}

	resultStr, ok := result.Outputs[resultCtxKey].Value.(string)
	if !ok {
		resultStr = ""
	}

	err = stream.Send(&deploymentspb.DeploymentUpEvent{
		Content: &deploymentspb.DeploymentUpEvent_Result{
			Result: &deploymentspb.UpResult{
				Content: &deploymentspb.UpResult_Text{
					Text: resultStr,
				},
			},
		},
	})

	return err
}

// Down - automatically called by the Nitric CLI via the `down` command
func (s *PulumiProviderServer) Down(req *deploymentspb.DeploymentDownRequest, stream deploymentspb.Deployment_DownServer) error {
	projectName, stackName, err := stackAndProjectFromAttributes(req.Attributes.AsMap())
	if err != nil {
		return err
	}

	// run down on the stack
	err = s.provider.Init(req.Attributes.AsMap())
	if err != nil {
		return err
	}

	stack, err := auto.UpsertStackInlineSource(context.TODO(), fmt.Sprintf("%s-%s", projectName, stackName), projectName, nil)
	if err != nil {
		return err
	}

	pulumiEventsChan := make(chan events.EngineEvent)

	go func() {
		_ = pulumix.StreamPulumiDownEngineEvents(stream, pulumiEventsChan)
	}()

	config, err := s.provider.Config()
	if err != nil {
		return err
	}

	err = stack.SetAllConfig(context.TODO(), config)
	if err != nil {
		return err
	}

	_, err = stack.Refresh(context.TODO())
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	_, err = stack.Destroy(context.TODO(), optdestroy.EventStreams(pulumiEventsChan))
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

// Start - starts the Nitric Provider gRPC server, making it callable by the Nitric CLI during deployments.
func (s *PulumiProviderServer) Start() {
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
