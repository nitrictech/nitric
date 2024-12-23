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
	"encoding/json"
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/parameter"
)

func (a *NitricAwsTerraformProvider) ResourcesStore(stack cdktf.TerraformStack, accessRoleNames []string) error {
	index := common.NewResourceIndex()

	for name, bucket := range a.Buckets {
		index.Buckets[name] = *bucket.BucketArnOutput()
	}

	for name, api := range a.Apis {
		index.Apis[name] = common.ApiGateway{
			Arn:      *api.ArnOutput(),
			Endpoint: *api.EndpointOutput(),
		}
	}

	for name, ws := range a.Websockets {
		index.Websockets[name] = common.ApiGateway{
			Arn:      *ws.WebsocketArnOutput(),
			Endpoint: *ws.EndpointOutput(),
		}
	}

	for name, proxy := range a.HttpProxies {
		index.HttpProxies[name] = common.ApiGateway{
			Arn:      *proxy.ArnOutput(),
			Endpoint: *proxy.EndpointOutput(),
		}
	}

	for name, topic := range a.Topics {
		index.Topics[name] = common.Topic{
			Arn:             *topic.TopicArnOutput(),
			StateMachineArn: *topic.SfnArnOutput(),
		}
	}

	for name, kv := range a.KeyValueStores {
		index.KvStores[name] = *kv.KvArnOutput()
	}

	for name, queue := range a.Queues {
		index.Queues[name] = *queue.QueueArnOutput()
	}

	for name, secret := range a.Secrets {
		index.Secrets[name] = *secret.SecretArnOutput()
	}

	indexJson, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal resource index: %w", err)
	}

	parameter.NewParameter(stack, jsii.String("nitric_resources"), &parameter.ParameterConfig{
		ParameterName:   jsii.Sprintf("/nitric/%s/resource-index", *a.Stack.StackIdOutput()),
		ParameterValue:  jsii.String(string(indexJson)),
		AccessRoleNames: jsii.Strings(accessRoleNames...),
	})

	return nil
}
