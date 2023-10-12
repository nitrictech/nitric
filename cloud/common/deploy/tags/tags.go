// Copyright 2021 Nitric Pty Ltd.
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

package tags

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Tags(ctx *pulumi.Context, stackID string, name string) map[string]string {
	return map[string]string{
		"x-nitric-project":    ctx.Project(),
		"x-nitric-stack":      stackID,
		"x-nitric-stack-name": ctx.Stack(),
		"x-nitric-name":       name,
		// New tag for resource identification
		// Locate the unique stack by the presence of the key and the resource by its name
		GetResourceNameKey(stackID): name,
	}
}

// GetResourceNameKey returns the key used to retrieve a resource's stack specific name from its tags.
func GetResourceNameKey(stackID string) string {
	return fmt.Sprintf("x-nitric-stack-%s", stackID)
}
