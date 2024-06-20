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
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/policy"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

func (a *NitricGcpTerraformProvider) Policy(stack cdktf.TerraformStack, name string, config *deploymentspb.Policy) error {
	stringActions := []string{}
	for _, action := range config.Actions {
		stringActions = append(stringActions, v1.Action_name[int32(action)])
	}
	for _, principal := range config.Principals {
		if principal.Id.Type != v1.ResourceType_Service {
			return fmt.Errorf("non-service principals are not supported")
		}
		service, ok := a.Services[principal.Id.Name]
		if !ok {
			return fmt.Errorf("could not find service %s", principal.Id.Name)
		}

		for _, resource := range config.Resources {
			memberName := fmt.Sprintf("%s-%s", principal.Id.Name, resource.Id.Name)

			var concreteResourceName *string = nil

			switch resource.Id.Type {
			case v1.ResourceType_Bucket:
				if concreteResource, ok := a.Buckets[resource.Id.Name]; ok {
					concreteResourceName = concreteResource.NameOutput()
				} else {
					return fmt.Errorf("could not find bucket %s", resource.Id.Name)
				}
			case v1.ResourceType_Queue:
				if concreteResource, ok := a.Queues[resource.Id.Name]; ok {
					concreteResourceName = concreteResource.NameOutput()
				} else {
					return fmt.Errorf("could not find queue %s", resource.Id.Name)
				}
			case v1.ResourceType_Secret:
				if concreteResource, ok := a.Secrets[resource.Id.Name]; ok {
					concreteResourceName = concreteResource.NameOutput()
				} else {
					return fmt.Errorf("could not find secret %s", resource.Id.Name)
				}
			case v1.ResourceType_KeyValueStore:
				// Resource name isn't used for KeyValue stores, so it can be blank
				concreteResourceName = jsii.String("")
			case v1.ResourceType_Topic:
				if concreteResource, ok := a.Topics[resource.Id.Name]; ok {
					concreteResourceName = concreteResource.TopicNameOutput()
				} else {
					return fmt.Errorf("could not find topic %s", resource.Id.Name)
				}
			}

			resourceType := resource.Id.Type.String()

			policy.NewPolicy(stack, jsii.String(memberName), &policy.PolicyConfig{
				Actions:             jsii.Strings(stringActions...),
				ResourceName:        concreteResourceName,
				ResourceType:        jsii.String(resourceType),
				ServiceAccountEmail: service.ServiceAccountEmailOutput(),
				IamRoles:            cdktf.Token_AsAny(a.Stack.IamRolesOutput()),
			})
		}
	}

	return nil
}
