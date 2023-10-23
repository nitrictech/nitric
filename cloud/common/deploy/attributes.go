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

import "fmt"

type CommonStackDetails struct {
	Project       string
	FullStackName string
	Stack         string
	Region        string
}

// Read nitric attributes from the provided deployment attributes
func CommonStackDetailsFromAttributes(attributes map[string]interface{}) (*CommonStackDetails, error) {
	iProject, hasProject := attributes["project"]
	project, isString := iProject.(string)
	if !hasProject || !isString || project == "" {
		// need a valid project name
		return nil, fmt.Errorf("project is not set or invalid")
	}

	iStack, hasStack := attributes["stack"]
	stack, isString := iStack.(string)
	if !hasStack || !isString || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("stack is not set or invalid")
	}

	iRegion, hasRegion := attributes["region"]
	region, isString := iRegion.(string)
	if !hasRegion || !isString || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("region is not set or invalid")
	}
	// Backwards compatible stack name
	// The existing providers in the CLI
	// Use the combined project and stack name
	fullStackName := fmt.Sprintf("%s-%s", project, stack)

	return &CommonStackDetails{
		Project:       project,
		FullStackName: fullStackName,
		Region:        region,
		Stack:         stack,
	}, nil
}
