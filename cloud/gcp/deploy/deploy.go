// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "embed"

	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
	gcpsecretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/gcp/common"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricGcpPulumiProvider struct {
	*deploy.CommonStackDetails

	StackId   string
	GcpConfig *common.GcpConfig

	DelayQueue      *cloudtasks.Queue
	AuthToken       *oauth2.Token
	BaseComputeRole *projects.IAMCustomRole
	DockerProvider  *docker.Provider

	SecretManagerClient *gcpsecretmanager.Client

	Project            *Project
	ApiGateways        map[string]*apigateway.Gateway
	HttpProxies        map[string]*apigateway.Gateway
	CloudRunServices   map[string]*NitricCloudRunService
	Buckets            map[string]*storage.Bucket
	Topics             map[string]*pubsub.Topic
	Queues             map[string]*pubsub.Topic
	QueueSubscriptions map[string]*pubsub.Subscription
	Secrets            map[string]*secretmanager.Secret

	provider.NitricDefaultOrder
}

var _ provider.NitricPulumiProvider = (*NitricGcpPulumiProvider)(nil)

const pulumiGcpVersion = "6.67.0"

func (a *NitricGcpPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"gcp:region":     auto.ConfigValue{Value: a.Region},
		"gcp:project":    auto.ConfigValue{Value: a.GcpConfig.ProjectId},
		"gcp:version":    auto.ConfigValue{Value: pulumiGcpVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricGcpPulumiProvider) WithDefaultResourceOptions(opts ...pulumi.ResourceOption) []pulumi.ResourceOption {
	defaultOptions := []pulumi.ResourceOption{
		pulumi.DependsOn([]pulumi.Resource{a.Project}),
	}

	return append(defaultOptions, opts...)
}

func (a *NitricGcpPulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	a.GcpConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	a.SecretManagerClient, err = gcpsecretmanager.NewClient(context.Background())
	if err != nil {
		return err
	}

	return nil
}

var baseComputePermissions []string = []string{
	"storage.buckets.list",
	"storage.buckets.get",
	"cloudtasks.queues.get",
	"cloudtasks.tasks.create",
	"cloudtrace.traces.patch",
	"monitoring.timeSeries.create",
	// permission for blob signing
	// this is safe as only permissions this account has are delegated
	"iam.serviceAccounts.signBlob",
	// Basic list permissions
	"pubsub.topics.list",
	"pubsub.topics.get",
	"pubsub.snapshots.list",
	"pubsub.subscriptions.get",
	"resourcemanager.projects.get",
	"secretmanager.secrets.list",
	"apigateway.gateways.list",

	// telemetry
	"monitoring.metricDescriptors.create",
	"monitoring.metricDescriptors.get",
	"monitoring.metricDescriptors.list",
	"monitoring.monitoredResourceDescriptors.get",
	"monitoring.monitoredResourceDescriptors.list",
	"monitoring.timeSeries.create",
}

func (a *NitricGcpPulumiProvider) Pre(ctx *pulumi.Context, resources []*pulumix.NitricPulumiResource[any]) error {
	// make our random stackId
	stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Upper:   pulumi.Bool(false),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": ctx.Stack(),
		}),
	})
	if err != nil {
		return err
	}

	stackIdChan := make(chan string)
	pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.Result).ApplyT(func(id string) string {
		stackIdChan <- id
		return id
	})

	a.StackId = <-stackIdChan

	project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
		ProjectId: &a.GcpConfig.ProjectId,
	}, nil)
	if err != nil {
		return err
	}

	a.Project, err = NewProject(ctx, "project", &ProjectArgs{
		ProjectId:     a.GcpConfig.ProjectId,
		ProjectNumber: project.Number,
	})
	if err != nil {
		return err
	}

	a.DelayQueue, err = cloudtasks.NewQueue(ctx, "delay-queue", &cloudtasks.QueueArgs{
		Location: pulumi.String(a.Region),
	})
	if err != nil {
		return err
	}

	// Deploy all services
	a.AuthToken, err = getGCPToken(ctx)
	if err != nil {
		return err
	}

	baseCustomRoleId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-base-role", a.FullStackName), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": a.FullStackName,
		}),
	})
	if err != nil {
		return errors.WithMessage(err, "base customRole id")
	}

	a.BaseComputeRole, err = projects.NewIAMCustomRole(ctx, "base-role", &projects.IAMCustomRoleArgs{
		Title:       pulumi.String(a.FullStackName + "-functions-base-role"),
		Permissions: pulumi.ToStringArray(baseComputePermissions),
		RoleId:      baseCustomRoleId.ID(),
	})
	if err != nil {
		return errors.WithMessage(err, "base customRole")
	}

	// Check if a key value store exists, if so get/create a (default) firestore database
	kvStoreExists := lo.SomeBy(resources, func(res *pulumix.NitricPulumiResource[any]) bool {
		_, ok := res.Config.(*deploymentspb.Resource_KeyValueStore)
		return ok
	})

	if kvStoreExists {
		err := createFirestoreDatabase(ctx, *project.ProjectId, a.Region)
		if err != nil {
			return err
		}
	}

	a.DockerProvider, err = docker.NewProvider(ctx, "docker-auth-provider", &docker.ProviderArgs{
		RegistryAuth: &docker.ProviderRegistryAuthArray{
			docker.ProviderRegistryAuthArgs{
				Username: pulumi.String("oauth2accesstoken"),
				Password: pulumi.String(a.AuthToken.AccessToken),
				Address:  pulumi.String("https://gcr.io"),
			},
			docker.ProviderRegistryAuthArgs{
				Address:  pulumi.Sprintf("%s-docker.pkg.dev", a.Region),
				Username: pulumi.String("oauth2accesstoken"),
				Password: pulumi.String(a.AuthToken.AccessToken),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getGCPToken(ctx *pulumi.Context) (*oauth2.Token, error) {
	// If the user is attempting to impersonate a gcp service account using pulumi using the GOOGLE_IMPERSONATE_SERVICE_ACCOUNT env var
	// Read more: (https://www.pulumi.com/registry/packages/gcp/installation-configuration/#configuration-reference)
	targetSA := os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT")

	var token *oauth2.Token

	if targetSA != "" {
		service, err := iamcredentials.NewService(ctx.Context())
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Unable to impersonate service account: %s", targetSA))
		}

		accessToken, err := service.Projects.ServiceAccounts.GenerateAccessToken(fmt.Sprintf("projects/-/serviceAccounts/%s", targetSA), &iamcredentials.GenerateAccessTokenRequest{
			Scope: []string{
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/trace.append",
			},
		}).Do()
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Unable to impersonate service account: %s", targetSA))
		}

		if accessToken == nil {
			return nil, fmt.Errorf("unable to impersonate service account")
		}

		token = &oauth2.Token{AccessToken: accessToken.AccessToken}
	}

	if token == nil {
		creds, err := google.FindDefaultCredentialsWithParams(ctx.Context(), google.CredentialsParams{
			Scopes: []string{
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/trace.append",
			},
		})
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to find credentials, try 'gcloud auth application-default login'")
		}

		token, err = creds.TokenSource.Token()
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to acquire token source")
		}
	}

	return token, nil
}

func (a *NitricGcpPulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricGcpPulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add APIs outputs
	if len(a.ApiGateways) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.ApiGateways {
			outputs = append(outputs, pulumi.Sprintf("%s: https://%s", apiName, api.DefaultHostname))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.HttpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.HttpProxies {
			outputs = append(outputs, pulumi.Sprintf("%s: https://%s", proxyName, proxy.DefaultHostname))
		}
	}

	output, ok := pulumi.All(outputs...).ApplyT(func(deets []interface{}) string {
		stringyOutputs := make([]string, len(deets))
		for i, d := range deets {
			stringyOutputs[i] = d.(string)
		}

		return strings.Join(stringyOutputs, "\n")
	}).(pulumi.StringOutput)

	if !ok {
		return pulumi.StringOutput{}, fmt.Errorf("Failed to generate pulumi output")
	}

	return output, nil
}

func NewNitricGcpProvider() *NitricGcpPulumiProvider {
	return &NitricGcpPulumiProvider{
		HttpProxies:        make(map[string]*apigateway.Gateway),
		ApiGateways:        make(map[string]*apigateway.Gateway),
		CloudRunServices:   make(map[string]*NitricCloudRunService),
		Buckets:            make(map[string]*storage.Bucket),
		Topics:             make(map[string]*pubsub.Topic),
		Queues:             make(map[string]*pubsub.Topic),
		QueueSubscriptions: make(map[string]*pubsub.Subscription),
		Secrets:            make(map[string]*secretmanager.Secret),
	}
}

func createFirestoreDatabase(ctx *pulumi.Context, projectId string, location string) error {
	fsAdminClient, err := apiv1.NewFirestoreAdminClient(context.TODO())
	if err != nil {
		return err
	}

	defaultDb, _ := fsAdminClient.GetDatabase(context.TODO(), &adminpb.GetDatabaseRequest{
		Name: fmt.Sprintf("projects/%s/databases/(default)", projectId),
	})

	defaultFirestoreId := pulumi.ID("(default)")

	if defaultDb != nil {
		_, err = firestore.GetDatabase(ctx, "default", defaultFirestoreId, nil)
		if err != nil {
			return err
		}
	} else {
		_, err = firestore.NewDatabase(ctx, "default", &firestore.DatabaseArgs{
			Name:                     defaultFirestoreId,
			AppEngineIntegrationMode: pulumi.String("DISABLED"),
			LocationId:               pulumi.String(location),
			Type:                     pulumi.String("FIRESTORE_NATIVE"),
		}, pulumi.RetainOnDelete(true))
		if err != nil {
			return err
		}
	}

	return nil
}
