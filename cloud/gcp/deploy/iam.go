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

package deploy

import (
	"regexp"

	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type GcpIamServiceAccount struct {
	pulumi.ResourceState

	Name           string
	ServiceAccount *serviceaccount.Account
}

type GcpIamServiceAccountArgs struct {
	AccountId string
}

// Create a new GCP IAM service account
func NewServiceAccount(ctx *pulumi.Context, name string, args *GcpIamServiceAccountArgs, opts ...pulumi.ResourceOption) (*GcpIamServiceAccount, error) {
	res := &GcpIamServiceAccount{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitricgcp:iam:GCPServiceAccount", name, res, opts...)
	if err != nil {
		return nil, err
	}

	invalidChars := regexp.MustCompile(`[^a-z0-9\-]`)
	args.AccountId = invalidChars.ReplaceAllString(args.AccountId, "-")

	// Create a random id
	acctRandId, err := random.NewRandomString(ctx, name+"-id", &random.RandomStringArgs{
		Length:  pulumi.Int(7),
		Upper:   pulumi.Bool(false),
		Number:  pulumi.Bool(false),
		Special: pulumi.Bool(false),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"name": name,
		}),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	res.ServiceAccount, err = serviceaccount.NewAccount(ctx, name, &serviceaccount.AccountArgs{
		AccountId: pulumi.Sprintf("%s-%s", utils.StringTrunc(args.AccountId, 30-8), acctRandId.Result),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	err = ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":             pulumi.String(res.Name),
		"serviceAccountId": res.ServiceAccount.AccountId,
	})

	return res, err
}
