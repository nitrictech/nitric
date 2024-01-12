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

package utils

import (
	"context"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type logger struct {
	t *testing.T
}

func (l *logger) Debug(msg string, args *pulumi.LogArgs) error {
	l.t.Log(msg)

	return nil
}

func (l *logger) Info(msg string, args *pulumi.LogArgs) error  { return l.Debug(msg, args) }
func (l *logger) Warn(msg string, args *pulumi.LogArgs) error  { return l.Debug(msg, args) }
func (l *logger) Error(msg string, args *pulumi.LogArgs) error { return l.Debug(msg, args) }

func Test_resourceName(t *testing.T) {
	tests := []struct {
		name       string
		nameSuffix string
		project    string
		stackName  string
		rt         ResourceType
		want       string
	}{
		{
			name:       "nothing changed",
			nameSuffix: "foo",
			project:    "stack",
			stackName:  "stack-dep",
			rt:         ContainerAppRT,
			want:       "foo-app",
		},
		{
			name:       "long stack name",
			nameSuffix: "foo",
			project:    "looooooooooooooooooooong",
			stackName:  "looooooooooooooooooooong-dep",
			rt:         StorageAccountRT,
			want:       "loooooodepst",
		},
		{
			name:       "long deployment name",
			nameSuffix: "foo",
			project:    "stack",
			stackName:  "stack-deploooooooooooooooooy",
			rt:         StorageAccountRT,
			want:       "stackdeplooost",
		},
		{
			name:       "camelCase name",
			nameSuffix: "important",
			project:    "stack",
			stackName:  "stack-deploy",
			rt:         KeyVaultRT,
			want:       "stDeKv",
		},
		{
			name:       "camelCase name with hyphen",
			nameSuffix: "important",
			project:    "stack",
			stackName:  "stack--deploy",
			rt:         KeyVaultRT,
			want:       "stDeKv",
		},
		{
			name:       "long stack and deployment name",
			nameSuffix: "foo",
			project:    "staaaaaaaaaaaaaaack",
			stackName:  "staaaaaaaaaaaaaaack-deploooooooooooooooooy",
			rt:         StorageAccountRT,
			want:       "staaaaadeplooost",
		},
		{
			name:       "containerApp non alphanumeric",
			nameSuffix: "foo&*)-ok",
			project:    "stack",
			stackName:  "stack-deploy",
			rt:         ContainerAppRT,
			want:       "foo-ok-app",
		},
		{
			name:       "resourcegroup non alphanumeric, with hypen",
			nameSuffix: "foo&*)-ok",
			project:    "stack",
			stackName:  "stack-deploy-one",
			rt:         ResourceGroupRT,
			want:       "stack-deploy-one-rg",
		},
		{
			name:       "mongocollection non alphanumeric, with Upper",
			nameSuffix: "*foo&)OK",
			project:    "stack",
			stackName:  "stack-deploy-one",
			rt:         MongoCollectionRT,
			want:       "fooOKColl",
		},
		{
			name:       "storage acct always fits in",
			nameSuffix: "",
			project:    "stack123456789",
			stackName:  "stack123456789-deploy123456789",
			rt:         StorageAccountRT,
			want:       "stack12deploy1st",
		},
		{
			name:       "overall too long",
			nameSuffix: "wow-this-is-long-isn't-it? could-be-longer-tho",
			project:    "stack123456789",
			stackName:  "stack123456789-deploy123456789",
			rt:         EventGridRT,
			want:       "wow-this-is-evgt",
		},
		{
			name:       "first char not a letter",
			nameSuffix: "14-red-frog",
			project:    "stack",
			stackName:  "stack-deploy",
			rt:         ApiOperationPolicyRT,
			want:       "red-frog-api-op-pol",
		},
		{
			name:       "last char a hypen",
			nameSuffix: "red-",
			project:    "stack",
			stackName:  "stack-deploy",
			rt:         ApiOperationPolicyRT,
			want:       "red-api-op-pol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx, err := pulumi.NewContext(context.Background(), pulumi.RunInfo{
				DryRun:  true,
				Project: tt.project,
				Stack:   tt.stackName,
			})
			if err != nil {
				t.Error(err)
			}

			testCtx.Log = &logger{t: t}

			if got := ResourceName(testCtx, tt.nameSuffix, tt.rt); got != tt.want {
				t.Errorf("resourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
