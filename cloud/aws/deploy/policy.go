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
	"crypto/md5" //#nosec G501 -- md5 used only to produce a unique ID from non-sensistive information (policy IDs)
	"encoding/hex"
	"encoding/json"
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

func md5Hash(b []byte) string {
	hasher := md5.New() //#nosec G401 -- md5 used only to produce a unique ID from non-sensistive information (policy IDs)
	hasher.Write(b)

	return hex.EncodeToString(hasher.Sum(nil))
}

var awsActionsMap map[resourcespb.Action][]string = map[resourcespb.Action][]string{
	resourcespb.Action_BucketFileList: {
		"s3:ListBucket",
	},
	resourcespb.Action_BucketFileGet: {
		"s3:GetObject",
	},
	resourcespb.Action_BucketFilePut: {
		"s3:PutObject",
	},
	resourcespb.Action_BucketFileDelete: {
		"s3:DeleteObject",
	},
	resourcespb.Action_TopicPublish: {
		"sns:GetTopicAttributes",
		"sns:Publish",
		"states:StartExecution",
		"states:StateSyncExecution",
	},
	resourcespb.Action_KeyValueStoreRead: {
		"dynamodb:GetItem",
		"dynamodb:BatchGetItem",
	},
	resourcespb.Action_KeyValueStoreWrite: {
		"dynamodb:UpdateItem",
		"dynamodb:PutItem",
	},
	resourcespb.Action_KeyValueStoreDelete: {
		"dynamodb:DeleteItem",
	},
	resourcespb.Action_SecretAccess: {
		"secretsmanager:GetSecretValue",
	},
	resourcespb.Action_SecretPut: {
		"secretsmanager:PutSecretValue",
	},
	resourcespb.Action_WebsocketManage: {
		"execute-api:ManageConnections",
	},
	resourcespb.Action_QueueEnqueue: {
		"sqs:SendMessage",
		"sqs:GetQueueAttributes",
		"sqs:GetQueueUrl",
		"sqs:ListQueueTags",
	},
	resourcespb.Action_QueueDequeue: {
		"sqs:ReceiveMessage",
		"sqs:DeleteMessage",
		"sqs:GetQueueAttributes",
		"sqs:GetQueueUrl",
		"sqs:ListQueueTags",
	},
}

func actionsToAwsActions(actions []resourcespb.Action) []string {
	awsActions := make([]string, 0)

	for _, a := range actions {
		awsActions = append(awsActions, awsActionsMap[a]...)
	}

	awsActions = lo.Uniq(awsActions)

	return awsActions
}

// // discover the arn of a deployed resource
func (a *NitricAwsPulumiProvider) arnForResource(resource *deploymentspb.Resource) ([]interface{}, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Bucket:
		if b, ok := a.buckets[resource.Id.Name]; ok {
			return []interface{}{b.Arn, pulumi.Sprintf("%s/*", b.Arn)}, nil
		}
	case resourcespb.ResourceType_Topic:
		if t, ok := a.topics[resource.Id.Name]; ok {
			return []interface{}{t.sns.Arn, t.sfn.Arn}, nil
		}
	case resourcespb.ResourceType_Queue:
		if q, ok := a.queues[resource.Id.Name]; ok {
			return []interface{}{q.Arn}, nil
		}
	case resourcespb.ResourceType_KeyValueStore:
		if c, ok := a.keyValueStores[resource.Id.Name]; ok {
			return []interface{}{c.Arn}, nil
		}
	case resourcespb.ResourceType_Secret:
		if s, ok := a.secrets[resource.Id.Name]; ok {
			return []interface{}{s.Arn}, nil
		}
	case resourcespb.ResourceType_Websocket:
		if w, ok := a.websockets[resource.Id.Name]; ok {
			return []interface{}{pulumi.Sprintf("%s/*", w.ExecutionArn)}, nil
		}
	default:
		return nil, fmt.Errorf(
			"invalid resource type: %s. Did you mean to define it as a principal?", resource.Id.Type)
	}

	return nil, fmt.Errorf("unable to find %s named %s in AWS provider resource cache", resource.Id.Type, resource.Id.Name)
}

func (a *NitricAwsPulumiProvider) roleForPrincipal(resource *deploymentspb.Resource) (*iam.Role, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Service:
		if f, ok := a.lambdaRoles[resource.Id.Name]; ok {
			return f, nil
		}
	default:
		return nil, fmt.Errorf("could not find role for principal: %+v", resource)
	}

	return nil, fmt.Errorf("could not find role for principal: %+v", resource)
}

func (a *NitricAwsPulumiProvider) Policy(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Policy) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// Get Actions
	actions := actionsToAwsActions(config.Actions)

	// Get Targets
	targetArns := make([]interface{}, 0, len(config.Resources))

	for _, res := range config.Resources {
		if arn, err := a.arnForResource(res); err == nil {
			targetArns = append(targetArns, arn...)
		} else {
			return fmt.Errorf("failed to create policy, unable to determine resource ARN: %w", err)
		}
	}

	// Get principal roles
	// We're collecting roles here to ensure all defined principals are valid before proceeding
	principalRoles := make(map[string]*iam.Role)

	for _, princ := range config.Principals {
		if role, err := a.roleForPrincipal(princ); err == nil {
			if princ.Id.Type != resourcespb.ResourceType_Service {
				return fmt.Errorf("invalid principal type: %s. Only services can be principals", princ.Id.Type)
			}

			principalRoles[princ.Id.Name] = role
		} else {
			return err
		}
	}

	serialPolicy, err := json.Marshal(config)
	if err != nil {
		return err
	}

	policyJson := pulumi.All(targetArns...).ApplyT(func(args []interface{}) (string, error) {
		arns := make([]string, 0, len(args))

		for _, iArn := range args {
			arn, ok := iArn.(string)
			if !ok {
				return "", fmt.Errorf("input not a string: %T %v", arn, arn)
			}

			arns = append(arns, arn)
		}

		jsonb, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action":   actions,
					"Effect":   "Allow",
					"Resource": arns,
				},
			},
		})
		if err != nil {
			return "", err
		}

		return string(jsonb), nil
	})

	// create role policy for each role
	for k, r := range principalRoles {
		// Role policies require a unique name
		// Use a hash of the policy document to help create a unique name
		policyName := fmt.Sprintf("%s-%s", k, md5Hash(serialPolicy))

		_, err := iam.NewRolePolicy(ctx, policyName, &iam.RolePolicyArgs{
			Role:   r.ID(),
			Policy: policyJson,
		}, opts...)
		if err != nil {
			return err
		}
	}

	return nil
}
