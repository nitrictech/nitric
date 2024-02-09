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
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricGcpPulumiProvider struct {
	stackId       string
	projectName   string
	stackName     string
	fullStackName string

	config *GcpConfig
	region string

	delayQueue      *cloudtasks.Queue
	authToken       *oauth2.Token
	baseComputeRole *projects.IAMCustomRole

	project            *Project
	apiGateways        map[string]*apigateway.Gateway
	httpProxies        map[string]*apigateway.Gateway
	cloudRunServices   map[string]*NitricCloudRunService
	buckets            map[string]*storage.Bucket
	topics             map[string]*pubsub.Topic
	queues             map[string]*pubsub.Topic
	queueSubscriptions map[string]*pubsub.Subscription
	secrets            map[string]*secretmanager.Secret

	provider.NitricDefaultOrder
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-gcp
var runtime []byte

var _ provider.NitricPulumiProvider = (*NitricGcpPulumiProvider)(nil)

const pulumiGcpVersion = "6.67.0"

func (a *NitricGcpPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"gcp:region":     auto.ConfigValue{Value: a.region},
		"gcp:project":    auto.ConfigValue{Value: a.config.ProjectId},
		"gcp:version":    auto.ConfigValue{Value: pulumiGcpVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricGcpPulumiProvider) WithDefaultResourceOptions(opts ...pulumi.ResourceOption) []pulumi.ResourceOption {
	defaultOptions := []pulumi.ResourceOption{
		pulumi.DependsOn([]pulumi.Resource{a.project}),
	}

	return append(defaultOptions, opts...)
}

func (a *NitricGcpPulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	region, ok := attributes["region"].(string)
	if !ok {
		return fmt.Errorf("Missing region attribute")
	}

	a.region = region

	a.config, err = ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	var isString bool

	iProject, hasProject := attributes["project"]
	a.projectName, isString = iProject.(string)
	if !hasProject || !isString || a.projectName == "" {
		// need a valid project name
		return fmt.Errorf("project is not set or invalid")
	}

	iStack, hasStack := attributes["stack"]
	a.stackName, isString = iStack.(string)
	if !hasStack || !isString || a.stackName == "" {
		// need a valid stack name
		return fmt.Errorf("stack is not set or invalid")
	}

	// Backwards compatible stack name
	// The existing providers in the CLI
	// Use the combined project and stack name
	a.fullStackName = fmt.Sprintf("%s-%s", a.projectName, a.stackName)

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

func (a *NitricGcpPulumiProvider) Pre(ctx *pulumi.Context, resources []*deploymentspb.Resource) error {
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

	a.stackId = <-stackIdChan

	project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
		ProjectId: &a.config.ProjectId,
	}, nil)
	if err != nil {
		return err
	}

	a.project, err = NewProject(ctx, "project", &ProjectArgs{
		ProjectId:     a.config.ProjectId,
		ProjectNumber: project.Number,
	})
	if err != nil {
		return err
	}

	a.delayQueue, err = cloudtasks.NewQueue(ctx, "delay-queue", &cloudtasks.QueueArgs{
		Location: pulumi.String(a.region),
	})
	if err != nil {
		return err
	}

	// Deploy all services
	a.authToken, err = getGCPToken(ctx)
	if err != nil {
		return err
	}

	baseCustomRoleId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-base-role", a.fullStackName), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": a.fullStackName,
		}),
	})
	if err != nil {
		return errors.WithMessage(err, "base customRole id")
	}

	a.baseComputeRole, err = projects.NewIAMCustomRole(ctx, "base-role", &projects.IAMCustomRoleArgs{
		Title:       pulumi.String(a.fullStackName + "-functions-base-role"),
		Permissions: pulumi.ToStringArray(baseComputePermissions),
		RoleId:      baseCustomRoleId.ID(),
	})
	if err != nil {
		return errors.WithMessage(err, "base customRole")
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
	if len(a.apiGateways) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.apiGateways {
			outputs = append(outputs, pulumi.Sprintf("%s: https://%s", apiName, api.DefaultHostname))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.httpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.httpProxies {
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
		httpProxies:        make(map[string]*apigateway.Gateway),
		apiGateways:        make(map[string]*apigateway.Gateway),
		cloudRunServices:   make(map[string]*NitricCloudRunService),
		buckets:            make(map[string]*storage.Bucket),
		topics:             make(map[string]*pubsub.Topic),
		queues:             make(map[string]*pubsub.Topic),
		queueSubscriptions: make(map[string]*pubsub.Subscription),
		secrets:            make(map[string]*secretmanager.Secret),
	}
}
