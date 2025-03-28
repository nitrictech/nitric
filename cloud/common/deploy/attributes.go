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
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type CommonStackDetails struct {
	ProjectName   string            `mapstructure:"project"`
	FullStackName string            `mapstructure:"-"`
	StackName     string            `mapstructure:"stack"`
	Region        string            `mapstructure:"region"`
	Tags          map[string]string `mapstructure:"tags"`
}

// Read nitric attributes from the provided deployment attributes
func CommonStackDetailsFromAttributes(attributes map[string]interface{}) (*CommonStackDetails, error) {
	commonStackDetails := new(CommonStackDetails)

	err := mapstructure.Decode(attributes, commonStackDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attributes: %w", err)
	}

	if commonStackDetails.ProjectName == "" {
		return nil, fmt.Errorf("project is not set or invalid")
	}

	if commonStackDetails.StackName == "" {
		return nil, fmt.Errorf("stack is not set or invalid")
	}

	if commonStackDetails.Region == "" {
		return nil, fmt.Errorf("region is not set or invalid")
	}

	if commonStackDetails.Tags == nil {
		commonStackDetails.Tags = make(map[string]string)
	}

	// Backwards compatible stack name
	// The existing providers in the CLI
	// Use the combined project and stack name
	fullStackName := fmt.Sprintf("%s-%s", commonStackDetails.ProjectName, commonStackDetails.StackName)

	commonStackDetails.FullStackName = fullStackName

	return commonStackDetails, nil
}
