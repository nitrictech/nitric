// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type ContainerAppArgs struct {
	ResourceGroupName             pulumi.StringInput
	Location                      pulumi.StringInput
	StackID                       string
	SubscriptionID                pulumi.StringInput
	Registry                      *containerregistry.Registry
	RegistryUser                  pulumi.StringPtrInput
	RegistryPass                  pulumi.StringPtrInput
	ManagedEnv                    *app.ManagedEnvironment
	Env                           app.EnvironmentVarArray
	ImageUri                      pulumi.StringInput
	Service                       *deploy.Service
	ManagedIdentityID             pulumi.StringOutput
	MongoDatabaseName             pulumi.StringInput
	MongoDatabaseConnectionString pulumi.StringInput
	Config                        AzureContainerAppsConfig
	Schedules                     []*deploy.Resource
}

type ContainerApp struct {
	pulumi.ResourceState

	Name       string
	hostUrl    *pulumi.StringOutput
	Sp         *ServicePrincipal
	App        *app.ContainerApp
	EventToken pulumi.StringOutput
}

// HostUrl - Returns the HostURL of the application
// this will also ensure that the application has been successfully deployed
func (c *ContainerApp) HostUrl() (pulumi.StringOutput, error) {
	if c.hostUrl == nil {
		// Set the hostUrl from the App FQDN
		hostUrl := pulumi.All(c.App.LatestRevisionFqdn, c.EventToken).ApplyT(func(args []interface{}) (string, error) {
			// Get the full URL of the deployed container
			fqdn := args[0].(string)
			token := args[1].(string)

			// This callback in mainly used to get the hostURL for subscriptions
			// so ensuring that the deployed gateway is configured to validate subscription requests
			// doesn't hurt
			hostUrl := fmt.Sprintf("https://%s", fqdn)
			healthcheckUrl := fmt.Sprintf("%s/%s/x-nitric-topic/test", hostUrl, token)

			hCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			// Poll the URL until the host has started.
			for {
				// Provide data in the expected shape. The content is current not important.
				empty := ""
				dummyEvgt := eventgrid.Event{
					ID:          &empty,
					Data:        &empty,
					EventType:   &empty,
					Subject:     &empty,
					DataVersion: &empty,
				}

				jsonStr, err := dummyEvgt.MarshalJSON()
				if err != nil {
					return "", err
				}

				body := bytes.NewBuffer(jsonStr)

				req, err := http.NewRequestWithContext(hCtx, "POST", healthcheckUrl, body)
				if err != nil {
					return "", err
				}

				// Ensure that the gateway is configured to accept subscription validation
				req.Header.Set("aeg-event-type", "SubscriptionValidation")
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				resp, err := client.Do(req)
				if err == nil {
					_ = resp.Body.Close()
					break
				}
			}

			return hostUrl, nil
		}).(pulumi.StringOutput)

		c.hostUrl = &hostUrl
	}

	return *c.hostUrl, nil
}

// Built in role definitions for Azure
// See below URL for mapping
// https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
var RoleDefinitions = map[string]string{
	"KVSecretsOfficer":    "b86a8fe4-44ce-4948-aee5-eccb2c155cd7",
	"BlobDataContrib":     "ba92f5b4-2d11-453d-a403-e96b0029c9fe",
	"QueueDataContrib":    "974c5e8b-45b9-4653-ba55-5f855dd0fb88",
	"EventGridDataSender": "d5a91429-5739-47e2-a06b-3470a27159e7",
	// Access for locating resources
	"TagContributor": "4a9ae827-6dc8-4573-8ac7-8239d42aa03f",
}

func (p *NitricAzurePulumiProvider) Service(ctx *pulumi.Context, parent pulumi.Resource, name string, service *pulumix.NitricPulumiServiceConfig, runtime provider.RuntimeProvider) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.Provider(p.ContainerEnv.DockerProvider)}

	res := &ContainerApp{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitricazure:ContainerApp", name, res, opts...)
	if err != nil {
		return err
	}

	repositoryUrl := pulumi.Sprintf("%s/%s-%s-%s", p.ContainerEnv.Registry.LoginServer, p.ProjectName, name, "azure")

	image, err := image.NewImage(ctx, name, &image.ImageArgs{
		SourceImage:   service.GetImage().Uri,
		RepositoryUrl: repositoryUrl,
		Runtime:       runtime(),
	}, opts...)
	if err != nil {
		return err
	}

	if service.Type == "" {
		service.Type = "default"
	}

	serviceConfig := p.AzureConfig.Config[service.Type]

	if serviceConfig.ContainerApps == nil {
		return fmt.Errorf("invalid container app config type: %s", service.Type)
	}

	token, err := random.NewRandomPassword(ctx, res.Name+"-event-token", &random.RandomPasswordArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(32),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"name": name,
		}),
	})
	if err != nil {
		return errors.WithMessage(err, "service event token")
	}

	res.EventToken = token.Result

	// the service principal's named doesn't need to be unique from the container app, so we reuse it.
	principal, err := NewServicePrincipal(ctx, name, &ServicePrincipalArgs{}, pulumi.Parent(res))
	if err != nil {
		return err
	}
	p.Principals[resourcespb.ResourceType_Service][name] = principal
	res.Sp = principal

	scope := pulumi.Sprintf("subscriptions/%s/resourceGroups/%s", p.ClientConfig.SubscriptionId, p.ResourceGroup.Name)

	// Assign roles to the new SP
	for defName, id := range RoleDefinitions {
		_ = ctx.Log.Info("Assignment "+ResourceName(ctx, name+defName, AssignmentRT)+" roleDef "+id, &pulumi.LogArgs{Ephemeral: true})

		_, err = authorization.NewRoleAssignment(ctx, ResourceName(ctx, name+defName, AssignmentRT), &authorization.RoleAssignmentArgs{
			PrincipalId:      res.Sp.ServicePrincipalId,
			PrincipalType:    pulumi.StringPtr("ServicePrincipal"),
			RoleDefinitionId: pulumi.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", p.ClientConfig.SubscriptionId, id),
			Scope:            scope,
		}, pulumi.Parent(res))
		if err != nil {
			return err
		}
	}

	env := app.EnvironmentVarArray{
		app.EnvironmentVarArgs{
			Name:  pulumi.String("EVENT_TOKEN"),
			Value: res.EventToken,
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String(resource.NITRIC_STACK_ID),
			Value: pulumi.String(p.StackId),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprint(service.Workers)),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_SUBSCRIPTION_ID"),
			Value: pulumi.String(p.ClientConfig.SubscriptionId),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_RESOURCE_GROUP"),
			Value: p.ResourceGroup.Name,
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_CLIENT_ID"),
			SecretRef: pulumi.String("client-id"),
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_TENANT_ID"),
			SecretRef: pulumi.String("tenant-id"),
		},
		app.EnvironmentVarArgs{
			Name:      pulumi.String("AZURE_CLIENT_SECRET"),
			SecretRef: pulumi.String("client-secret"),
		},
		app.EnvironmentVarArgs{
			Name:  pulumi.String("TOLERATE_MISSING_SERVICES"),
			Value: pulumi.String("true"),
		},
	}

	if p.DatabaseServer != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name: pulumi.String("NITRIC_DATABASE_BASE_URL"),
			Value: pulumi.Sprintf("postgres://%s:%s@%s:%s", "nitric", p.DbMasterPassword.Result,
				p.DatabaseServer.FullyQualifiedDomainName, "5432"),
		})
	}

	for k, v := range service.Env() {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String(k),
			Value: v,
		})
	}

	env = append(env, p.ContainerEnv.Env...)

	// if len(args.Env) > 0 {
	// 	env = append(env, args.Env...)
	// }

	resOpts := []pulumi.ResourceOption{pulumi.Parent(res)}
	for _, db := range p.SqlMigrations {
		resOpts = append(resOpts, pulumi.DependsOn([]pulumi.Resource{db}))
	}

	//	If this instance contains a schedule set the minimum instances to 1
	// schedules rely on the Dapr Runtime to trigger the function, without a running instance the Dapr Runtime will not execute, so the schedule won't trigger.
	_, schedulesFound := lo.Find(p.resources, func(item *pulumix.NitricPulumiResource[any]) bool {
		switch t := item.Config.(type) {
		case *deploymentspb.Schedule:
			return t.Target.GetService() == name
		}

		return false
	})

	minReplicas := serviceConfig.ContainerApps.MinReplicas
	if schedulesFound {
		minReplicas = lo.Max([]int{minReplicas, 1})
	}

	appName := ResourceName(ctx, name, ContainerAppRT)

	res.App, err = app.NewContainerApp(ctx, appName, &app.ContainerAppArgs{
		ResourceGroupName:    p.ResourceGroup.Name,
		Location:             p.ResourceGroup.Location,
		ManagedEnvironmentId: p.ContainerEnv.ManagedEnv.ID(),
		Configuration: app.ConfigurationArgs{
			ActiveRevisionsMode: pulumi.String("Single"),
			Ingress: app.IngressArgs{
				External:   pulumi.BoolPtr(true),
				TargetPort: pulumi.Int(9001),
			},
			Registries: app.RegistryCredentialsArray{
				app.RegistryCredentialsArgs{
					Server:            p.ContainerEnv.Registry.LoginServer,
					Username:          p.ContainerEnv.RegistryUser,
					PasswordSecretRef: pulumi.String("pwd"),
				},
			},
			Dapr: &app.DaprArgs{
				AppId:       pulumi.String(appName),
				AppPort:     pulumi.Int(9001),
				AppProtocol: pulumi.String("http"),
				Enabled:     pulumi.Bool(true),
			},
			Secrets: app.SecretArray{
				app.SecretArgs{
					Name:  pulumi.String("pwd"),
					Value: p.ContainerEnv.RegistryPass,
				},
				app.SecretArgs{
					Name:  pulumi.String("client-id"),
					Value: res.Sp.ClientID,
				},
				app.SecretArgs{
					Name:  pulumi.String("tenant-id"),
					Value: res.Sp.TenantID,
				},
				app.SecretArgs{
					Name:  pulumi.String("client-secret"),
					Value: res.Sp.ClientSecret,
				},
			},
		},
		Tags: pulumi.ToStringMap(common.Tags(p.StackId, name, resources.Service)),
		Template: app.TemplateArgs{
			Scale: app.ScaleArgs{
				MaxReplicas: pulumi.Int(serviceConfig.ContainerApps.MaxReplicas),
				MinReplicas: pulumi.Int(minReplicas),
			},
			Containers: app.ContainerArray{
				app.ContainerArgs{
					Name:  pulumi.String("myapp"),
					Image: image.URI(),
					Resources: app.ContainerResourcesArgs{
						Cpu:    pulumi.Float64(serviceConfig.ContainerApps.Cpu),
						Memory: pulumi.Sprintf("%.2fGi", serviceConfig.ContainerApps.Memory),
					},
					Env: env,
				},
			},
		},
	}, resOpts...)
	if err != nil {
		return err
	}

	authName := fmt.Sprintf("%s-auth", appName)

	_, err = app.NewContainerAppsAuthConfig(ctx, authName, &app.ContainerAppsAuthConfigArgs{
		AuthConfigName:   pulumi.String("current"),
		ContainerAppName: res.App.Name,
		GlobalValidation: &app.GlobalValidationArgs{
			UnauthenticatedClientAction: app.UnauthenticatedClientActionV2Return401,
		},
		IdentityProviders: &app.IdentityProvidersArgs{
			AzureActiveDirectory: &app.AzureActiveDirectoryArgs{
				Enabled: pulumi.Bool(true),
				Registration: &app.AzureActiveDirectoryRegistrationArgs{
					ClientId:                res.Sp.ClientID,
					ClientSecretSettingName: pulumi.String("client-secret"),
					OpenIdIssuer:            pulumi.Sprintf("https://sts.windows.net/%s/v2.0", res.Sp.TenantID),
				},
				Validation: &app.AzureActiveDirectoryValidationArgs{
					AllowedAudiences: pulumi.StringArray{p.ContainerEnv.ManagedUser.ClientId},
				},
			},
		},
		Platform: &app.AuthPlatformArgs{
			Enabled: pulumi.Bool(true),
		},
		ResourceGroupName: p.ResourceGroup.Name,
	}, pulumi.Parent(res.App))
	if err != nil {
		return err
	}

	err = ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":         pulumi.StringPtr(res.Name),
		"containerApp": res.App,
	})

	p.ContainerApps[name] = res

	return err
}
