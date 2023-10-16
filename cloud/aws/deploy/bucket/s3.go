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

package bucket

import (
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type S3Bucket struct {
	pulumi.ResourceState
	S3   *s3.Bucket
	Name string
}

type S3BucketArgs struct {
	StackID string
	Bucket  *v1.Bucket
}

// NewS3Bucket creates new S3 Buckets
func NewS3Bucket(ctx *pulumi.Context, name string, args *S3BucketArgs, opts ...pulumi.ResourceOption) (*S3Bucket, error) {
	res := &S3Bucket{
		Name: name,
	}
	err := ctx.RegisterComponentResource("nitric:bucket:AwsS3Bucket", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.S3, err = s3.NewBucket(ctx, name, &s3.BucketArgs{
		Tags: pulumi.ToStringMap(common.Tags(args.StackID, name, resources.Bucket)),
	})
	if err != nil {
		return nil, errors.WithMessage(err, "s3 bucket "+name)
	}

	return res, nil
}
