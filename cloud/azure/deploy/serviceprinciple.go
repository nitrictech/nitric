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
	"github.com/pulumi/pulumi-azuread/sdk/v5/go/azuread"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ServicePrincipalArgs struct{}

type ServicePrincipal struct {
	pulumi.ResourceState
	Name               string
	ClientID           pulumi.StringOutput
	TenantID           pulumi.StringOutput
	ServicePrincipalId pulumi.StringOutput
	ClientSecret       pulumi.StringOutput
}

func NewServicePrincipal(ctx *pulumi.Context, name string, args *ServicePrincipalArgs, opts ...pulumi.ResourceOption) (*ServicePrincipal, error) {
	res := &ServicePrincipal{Name: name}

	err := ctx.RegisterComponentResource("nitric:principal:AzureAD", name, res, opts...)
	if err != nil {
		return nil, err
	}

	current, err := azuread.GetClientConfig(ctx, nil)
	if err != nil {
		return nil, err
	}

	// create an application per service principal
	appRoleId := pulumi.String("4962773b-9cdb-44cf-a8bf-237846a00ab7")

	app, err := azuread.NewApplication(ctx, ResourceName(ctx, name, ADApplicationRT), &azuread.ApplicationArgs{
		DisplayName: pulumi.String(name + "App"),
		Owners: pulumi.StringArray{
			pulumi.String(current.ObjectId),
		},
		AppRoles: azuread.ApplicationAppRoleTypeArray{
			&azuread.ApplicationAppRoleTypeArgs{
				AllowedMemberTypes: pulumi.StringArray{
					pulumi.String("Application"),
				},
				Description: pulumi.String("Enables webhook subscriptions to authenticate using this application"),
				DisplayName: pulumi.String("AzureEventGridSecureWebhookSubscriber"),
				Enabled:     pulumi.Bool(true),
				Id:          appRoleId,
				Value:       pulumi.String("4962773b-9cdb-44cf-a8bf-237846a00ab7"),
			},
		},
		// Tags:        common.Tags(ctx, a.stackID, name+"App"),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	sp, err := azuread.NewServicePrincipal(ctx, ResourceName(ctx, name, ADServicePrincipalRT), &azuread.ServicePrincipalArgs{
		ClientId: app.ClientId,
		Owners: pulumi.StringArray{
			pulumi.String(current.ObjectId),
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	res.TenantID = sp.ApplicationTenantId
	res.ServicePrincipalId = pulumi.StringOutput(sp.ID())

	_, err = azuread.NewAppRoleAssignment(ctx, name+"sub-role", &azuread.AppRoleAssignmentArgs{
		AppRoleId:         appRoleId,
		PrincipalObjectId: pulumi.String(current.ObjectId),
		ResourceObjectId:  sp.ID(),
	})
	if err != nil {
		return nil, err
	}

	spPwd, err := azuread.NewServicePrincipalPassword(ctx, ResourceName(ctx, name, ADServicePrincipalPasswordRT), &azuread.ServicePrincipalPasswordArgs{
		ServicePrincipalId: sp.ID().ToStringOutput(),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	res.ClientSecret = spPwd.Value

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":               pulumi.StringPtr(res.Name),
		"tenantID":           res.TenantID,
		"clientID":           res.ClientID,
		"clientSecret":       res.ClientSecret,
		"servicePrincipalId": res.ServicePrincipalId,
	})
}
