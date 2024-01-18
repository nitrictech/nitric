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

package topic

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	nitricresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Topics
type AzureEventGridTopic struct {
	pulumi.ResourceState

	Name          string
	Topic         *eventgrid.Topic
	ResourceGroup *resources.ResourceGroup
}

type AzureEventGridTopicArgs struct {
	StackID       string
	ResourceGroup *resources.ResourceGroup
}

func NewAzureEventGridTopic(ctx *pulumi.Context, name string, args *AzureEventGridTopicArgs, opts ...pulumi.ResourceOption) (*AzureEventGridTopic, error) {
	res := &AzureEventGridTopic{
		Name:          name,
		ResourceGroup: args.ResourceGroup,
	}

	err := ctx.RegisterComponentResource("nitric:topic:AzureEventGridTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Topic, err = eventgrid.NewTopic(ctx, utils.ResourceName(ctx, res.Name, utils.EventGridRT), &eventgrid.TopicArgs{
		ResourceGroupName: args.ResourceGroup.Name,
		Location:          args.ResourceGroup.Location,
		Tags:              pulumi.ToStringMap(common.Tags(args.StackID, res.Name, nitricresources.Topic)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
