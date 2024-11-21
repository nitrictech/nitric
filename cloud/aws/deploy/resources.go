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
	"encoding/json"

	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAwsPulumiProvider) resourcesStore(ctx *pulumi.Context) error {
	// Build the AWS resource index from the provider information
	// This will be used to store the ARNs/Identifiers of all resources created by the stack
	bucketNameMap := pulumi.StringMap{}
	for name, bucket := range a.Buckets {
		bucketNameMap[name] = bucket.Bucket
	}

	apiArnMap := pulumi.StringMap{}
	for name, api := range a.Apis {
		apiArnMap[name] = api.Arn
	}

	apiEndpointMap := pulumi.StringMap{}
	for name, api := range a.Apis {
		apiEndpointMap[name] = api.ApiEndpoint
	}

	websocketArnMap := pulumi.StringMap{}
	for name, api := range a.Websockets {
		apiArnMap[name] = api.Arn
	}

	websocketEndpointMap := pulumi.StringMap{}
	for name, api := range a.Websockets {
		apiEndpointMap[name] = api.ApiEndpoint
	}

	topicArnMap := pulumi.StringMap{}
	stateMachArnMap := pulumi.StringMap{}
	for name, topic := range a.Topics {
		topicArnMap[name] = topic.sns.Arn
		stateMachArnMap[name] = topic.sfn.Arn
	}

	kvStoreArnMap := pulumi.StringMap{}
	for name, kvStore := range a.KeyValueStores {
		kvStoreArnMap[name] = kvStore.Arn
	}

	queueArnMap := pulumi.StringMap{}
	for name, queue := range a.Queues {
		queueArnMap[name] = queue.Arn
	}

	secretsArnMap := pulumi.StringMap{}
	for name, secret := range a.Secrets {
		secretsArnMap[name] = secret.Arn
	}

	httpProxyArnMap := pulumi.StringMap{}
	for name, proxy := range a.HttpProxies {
		httpProxyArnMap[name] = proxy.Arn
	}

	httpProxyEndpointMap := pulumi.StringMap{}
	for name, proxy := range a.HttpProxies {
		httpProxyEndpointMap[name] = proxy.ApiEndpoint
	}

	// Build the index from the provider information
	resourceIndexJson := pulumi.All(
		bucketNameMap,
		apiArnMap,
		apiEndpointMap,
		websocketArnMap,
		websocketEndpointMap,
		topicArnMap,
		kvStoreArnMap,
		queueArnMap,
		secretsArnMap,
		stateMachArnMap,
		httpProxyArnMap,
		httpProxyEndpointMap,
	).ApplyT(func(args []interface{}) (string, error) {
		bucketNameMap := args[0].(map[string]string)
		apiArnMap := args[1].(map[string]string)
		apiEndpointMap := args[2].(map[string]string)
		websocketArnMap := args[3].(map[string]string)
		websocketEndpointMap := args[4].(map[string]string)
		topicArnMap := args[5].(map[string]string)
		kvStoreArnMap := args[6].(map[string]string)
		queueArnMap := args[7].(map[string]string)
		secretsArnMap := args[8].(map[string]string)
		stateMachArnMap := args[9].(map[string]string)
		httpProxyArnMap := args[10].(map[string]string)
		httpProxyEndpointMap := args[11].(map[string]string)

		index := common.NewResourceIndex()
		for name, bucket := range bucketNameMap {
			index.Buckets[name] = bucket
		}

		for name, arn := range apiArnMap {
			index.Apis[name] = common.ApiGateway{
				Arn:      arn,
				Endpoint: apiEndpointMap[name],
			}
		}

		for name, arn := range websocketArnMap {
			index.Websockets[name] = common.ApiGateway{
				Arn:      arn,
				Endpoint: websocketEndpointMap[name],
			}
		}

		for name, arn := range httpProxyArnMap {
			index.HttpProxies[name] = common.ApiGateway{
				Arn:      arn,
				Endpoint: httpProxyEndpointMap[name],
			}
		}

		for name, arn := range topicArnMap {
			index.Topics[name] = common.Topic{
				Arn:             arn,
				StateMachineArn: stateMachArnMap[name],
			}
		}

		for name, arn := range kvStoreArnMap {
			index.KvStores[name] = arn
		}

		for name, arn := range queueArnMap {
			index.Queues[name] = arn
		}

		for name, arn := range secretsArnMap {
			index.Secrets[name] = arn
		}

		indexJson, err := json.Marshal(index)
		if err != nil {
			return "", err
		}

		return string(indexJson), nil
	}).(pulumi.StringOutput)

	_, err := ssm.NewParameter(ctx, "nitric-resource-index", &ssm.ParameterArgs{
		// Create a deterministic name for the resource index
		Name:     pulumi.Sprintf("/nitric/%s/resource-index", a.StackId),
		DataType: pulumi.String("text"),
		Type:     pulumi.String("String"),
		// Store the nitric resource index serialized as a JSON string
		Value: resourceIndexJson,
	})
	if err != nil {
		return err
	}

	// Create a role that gives read access to the above parameter
	policy, err := iam.NewPolicy(ctx, "nitric-index-policy", &iam.PolicyArgs{Policy: pulumi.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "ssm:GetParameter",
						"Resource": "arn:aws:ssm:*:*:parameter/nitric/%s/resource-index"
					}
				]
			}`, a.StackId)},
	)
	if err != nil {
		return err
	}

	for name, role := range a.LambdaRoles {
		_, err = iam.NewRolePolicyAttachment(ctx, name+"NitricIndexAccess", &iam.RolePolicyAttachmentArgs{
			PolicyArn: policy.Arn,
			Role:      role.ID(),
		})
		if err != nil {
			return err
		}
	}

	for name, role := range a.BatchRoles {
		_, err = iam.NewRolePolicyAttachment(ctx, name+"NitricIndexAccess", &iam.RolePolicyAttachmentArgs{
			PolicyArn: policy.Arn,
			Role:      role.ID(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
