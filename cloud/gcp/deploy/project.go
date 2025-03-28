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

package deploy

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Project struct {
	pulumi.ResourceState

	Name     string
	Services []*projects.Service
}

type ProjectArgs struct {
	ProjectId     string
	ProjectNumber string
}

var requiredServices = []string{
	// Enable IAM
	"iam.googleapis.com",
	// Enable cloud run
	"run.googleapis.com",
	// Enable pubsub
	"pubsub.googleapis.com",
	// Enable cloud scheduler
	"cloudscheduler.googleapis.com",
	// Enable cloud scheduler
	"storage.googleapis.com",
	// Enable Compute API (Networking/Load Balancing)
	"compute.googleapis.com",
	// Enable Artifact Registry API and Container Registry API
	"containerregistry.googleapis.com",
	"artifactregistry.googleapis.com",
	// Enable firestore API
	"firestore.googleapis.com",
	// Enable ApiGateway API
	"apigateway.googleapis.com",
	// Enable SecretManager API
	"secretmanager.googleapis.com",
	// Enable Cloud Tasks API
	"cloudtasks.googleapis.com",
	// Enable monitoring API
	"monitoring.googleapis.com",
	// Enable Service Management API
	"servicemanagement.googleapis.com",
	// Enable Service Control API
	"servicecontrol.googleapis.com",
	// Enable Cloud SQL API
	"sql-component.googleapis.com",
	// Enable VPC Access API
	"vpcaccess.googleapis.com",
	// Enable Cloud Build API
	"cloudbuild.googleapis.com",
	// Enable Batch API
	"batch.googleapis.com",
}

// Creates a new GCP Project
func NewProject(ctx *pulumi.Context, name string, args *ProjectArgs, opts ...pulumi.ResourceOption) (*Project, error) {
	res := &Project{
		Name:     name,
		Services: []*projects.Service{},
	}

	err := ctx.RegisterComponentResource("ntiricgcp:project:GcpProject", name, res, opts...)
	if err != nil {
		return nil, err
	}

	deps := []pulumi.Resource{}

	// Enable all the required services
	for _, serv := range requiredServices {
		s, err := projects.NewService(ctx, serv+"-enabled", &projects.ServiceArgs{
			DisableDependentServices: pulumi.Bool(true),
			DisableOnDestroy:         pulumi.Bool(false),
			Project:                  pulumi.String(args.ProjectId),
			Service:                  pulumi.String(serv),
		})
		if err != nil {
			return nil, err
		}

		res.Services = append(res.Services, s)
		deps = append(deps, s)
	}

	// Add ServiceAccount Token Creator Role to the default pubsub gservice account
	// services-{projectNumber}@gcp-sa-pubsub.iam.gserviceaccount.com
	serviceAccount := pulumi.Sprintf("serviceAccount:service-%s@gcp-sa-pubsub.iam.gserviceaccount.com", args.ProjectNumber)

	_, err = projects.NewIAMMember(ctx, "pubsub-token-creator", &projects.IAMMemberArgs{
		Role:    pulumi.String("roles/iam.serviceAccountTokenCreator"),
		Member:  serviceAccount,
		Project: pulumi.String(args.ProjectId),
		// Only create this once the google managed service account is available
	}, pulumi.Parent(res), pulumi.DependsOn(deps))
	if err != nil {
		return nil, err
	}

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name": pulumi.StringPtr(res.Name),
		// "services": res.Services,
	})
}
