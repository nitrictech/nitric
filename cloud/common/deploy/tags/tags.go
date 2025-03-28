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

package tags

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
)

// Tags generates standard resource tags used to map nitric resources on to deployed resources
func Tags(stackID string, resourceName string, resourceType resources.ResourceType) map[string]string {
	return map[string]string{
		// Locate the unique stack by the presence of the key and the resource by its name
		GetResourceNameKey(stackID): resourceName,
		// Identifies the nitric resource type that led to the creation of the cloud resource.
		GetResourceTypeKey(stackID): string(resourceType),
	}
}

// GetResourceNameKey returns the key used to retrieve a resource's stack specific name from its tags.
func GetResourceNameKey(stackID string) string {
	if stackID == "" {
		panic("blank stack ID, resource mapping isn't possible")
	}
	return fmt.Sprintf("x-nitric-%s-name", stackID)
}

// GetResourceTypeKey returns the key used to retrieve a resource's stack specific type from its tags.
func GetResourceTypeKey(stackID string) string {
	if stackID == "" {
		panic("blank stack ID, resource mapping isn't possible")
	}
	return fmt.Sprintf("x-nitric-%s-type", stackID)
}

type TagsConfig struct {
	Tags map[string]string `mapstructure:"tags"`
}

type GetGlobalTagsFunc func() map[string]string

func CreateGlobalTagsFromAttributes(attributes map[string]interface{}) (GetGlobalTagsFunc, error) {
	config := new(TagsConfig)

	err := mapstructure.Decode(attributes, config)
	if err != nil {
		return nil, err
	}

	return func() map[string]string {
		return config.Tags
	}, nil
}
