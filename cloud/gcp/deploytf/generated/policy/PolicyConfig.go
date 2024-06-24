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

package policy

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type PolicyConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The actions to apply to the policy.
	Actions *[]*string `field:"required" json:"actions" yaml:"actions"`
	// The IAM roles available to the policy.
	IamRoles interface{} `field:"required" json:"iamRoles" yaml:"iamRoles"`
	// The name of the resource.
	ResourceName *string `field:"required" json:"resourceName" yaml:"resourceName"`
	// The type of the resource (Bucket, Secret, KeyValueStore, Queue).
	ResourceType *string `field:"required" json:"resourceType" yaml:"resourceType"`
	// The service account to apply the policy to.
	ServiceAccountEmail *string `field:"required" json:"serviceAccountEmail" yaml:"serviceAccountEmail"`
}
