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

package bucket

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type BucketConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The name of the bucket.
	//
	// This must be globally unique.
	BucketName *string `field:"required" json:"bucketName" yaml:"bucketName"`
	// The notification target configurations.
	NotificationTargets interface{} `field:"required" json:"notificationTargets" yaml:"notificationTargets"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The class of storage used to store the bucket's contents.
	//
	// This can be STANDARD, NEARLINE, COLDLINE, ARCHIVE, or MULTI_REGIONAL.
	// STANDARD.
	StorageClass *string `field:"optional" json:"storageClass" yaml:"storageClass"`
}
